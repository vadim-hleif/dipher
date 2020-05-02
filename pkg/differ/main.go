package differ

import (
	"fmt"
	"strings"
)

func Diff(specV1 map[string]interface{}, specV2 map[string]interface{}) []error {
	errs := make([]error, 0)
	pathsV2 := specV2["paths"].(map[string]interface{})

	for url, urlNodeV1 := range specV1["paths"].(map[string]interface{}) {
		urlNodeV2 := getNode(pathsV2, url)

		if urlNodeV2 == nil {
			errs = append(errs, fmt.Errorf("resource %v mustn't be removed", url))
			continue
		}

		for methodV1, m := range urlNodeV1.(map[string]interface{}) {
			methodNodeV1 := m.(map[string]interface{})
			methodNodeV2 := getNode(urlNodeV2, methodV1)

			if methodNodeV2 == nil {
				errs = append(errs, fmt.Errorf("%v method of %v path mustn't be removed", methodV1, url))
				continue
			}

			paramsV1 := methodNodeV1["parameters"].([]interface{})
			paramsV2 := methodNodeV2["parameters"].([]interface{})

			for _, p := range paramsV1 {
				paramV1 := p.(map[string]interface{})
				paramV2 := findParam(paramsV2, paramV1["name"].(string))

				typeV1 := getTypeProp(paramV1)
				typeV2 := getTypeProp(paramV2)

				switch typeV1 {
				case "reference":
					errs = validateParam(typeV1, typeV2, getModelByRef(paramV1, specV1), getModelByRef(paramV2, specV2), errs)
				case "object":
					errs = validateParam(typeV1, typeV2, paramV1, paramV2, errs)
				default:
					errs = validateParamPrimitive(paramV1, paramV2, errs, typeV1, typeV2)
				}

			}

			for _, p := range paramsV2 {
				paramV2 := p.(map[string]interface{})
				if findParam(paramsV1, paramV2["name"].(string)) == nil && getRequiredProp(paramV2) {
					errs = append(errs, fmt.Errorf("new required param %v mustn't be added", paramV2["name"].(string)))
				}
			}
		}

	}

	return errs
}

func validateParamPrimitive(paramV1 map[string]interface{}, paramV2 map[string]interface{}, errs []error, typeV1 string, typeV2 string) []error {
	isParamV1Required := getRequiredProp(paramV1)
	isParamV2Required := getRequiredProp(paramV2)

	if paramV2 == nil && isParamV1Required {
		return append(errs, fmt.Errorf("required param %v mustn't be deleted", paramV1["name"].(string)))

	}

	if !isParamV1Required && isParamV2Required {
		errs = append(errs, fmt.Errorf("param %v mustn't be required because it wasn't be required", paramV1["name"].(string)))
	}

	enumV1 := getEnum(paramV1)
	enumV2 := getEnum(paramV2)

	if !(enumV2 == nil && enumV1 == nil) {
		if enumV1 == nil && enumV2 != nil {
			errs = append(errs, fmt.Errorf("param %v mustn't have enum", paramV1["name"].(string)))
		}

		compareAndApply(enumV1, enumV2, func(name interface{}) {
			errs = append(errs, fmt.Errorf("param %v mustn't remove value %v from enum", paramV1["name"].(string), name))
		})
	}

	if typeV1 != typeV2 {
		errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", paramV1["name"].(string), typeV1, typeV2))
	}

	return errs
}

func validateParam(typeV1 string, typeV2 string, paramV1 map[string]interface{}, paramV2 map[string]interface{}, errs []error) []error {
	if typeV1 != typeV2 {
		errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", paramV1["name"].(string), typeV1, typeV2))
	}

	schemaV1 := getNode(paramV1, "schema")
	if schemaV1 == nil {
		schemaV1 = paramV1
	}
	schemaV2 := getNode(paramV2, "schema")
	if schemaV2 == nil {
		schemaV2 = paramV2
	}

	requiredPropsV1 := getRequiredProps(schemaV1)
	requiredPropsV2 := getRequiredProps(schemaV2)

	compareAndApply(requiredPropsV2, requiredPropsV1, func(name interface{}) {
		errs = append(errs, fmt.Errorf("param %v mustn't be required because it wasn't be required", name))
	})

	compareAndApply(requiredPropsV1, requiredPropsV2, func(name interface{}) {
		errs = append(errs, fmt.Errorf("required param %v mustn't be deleted", name))
	})

	pV2 := getNode(schemaV2, "properties")

	for nameV1, propsV1 := range getNode(schemaV1, "properties") {
		propsV2, ok := pV2[nameV1]
		if ok {
			typeV1 := getTypeProp(propsV1.(map[string]interface{}))
			typeV2 := getTypeProp(propsV2.(map[string]interface{}))

			if typeV1 != typeV2 {
				errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", nameV1, typeV1, typeV2))
			}

			enumV1 := getEnum(propsV1.(map[string]interface{}))
			enumV2 := getEnum(propsV2.(map[string]interface{}))
			if enumV2 == nil && enumV1 == nil {
				continue
			}

			if enumV1 == nil && enumV2 != nil {
				errs = append(errs, fmt.Errorf("param %v mustn't have enum", nameV1))
			}

			compareAndApply(enumV1, enumV2, func(name interface{}) {
				errs = append(errs, fmt.Errorf("param %v mustn't remove value %v from enum", nameV1, name))
			})
		}

	}
	return errs
}

func compareAndApply(sliceV1 []interface{}, sliceV2 []interface{}, cb func(name interface{})) {
	for _, elV1 := range sliceV1 {
		var exist bool
		for _, elV2 := range sliceV2 {
			if elV1 == elV2 {
				exist = true
				break
			}
		}

		if !exist {
			cb(elV1)
		}
	}
}

func getNode(node map[string]interface{}, key string) map[string]interface{} {
	value, ok := node[key]

	if ok {
		return value.(map[string]interface{})
	}

	return nil
}

func getRequiredProp(node map[string]interface{}) bool {
	value, ok := node["required"]

	if ok {
		return value.(bool)
	}

	return false
}

func getRequiredProps(node map[string]interface{}) []interface{} {
	value, ok := node["required"]

	if ok {
		return value.([]interface{})
	}

	return nil
}

func getEnum(node map[string]interface{}) []interface{} {
	value, ok := node["enum"]

	if ok {
		return value.([]interface{})
	}

	value, ok = node["schema"]
	if ok {
		schema := value.(map[string]interface{})
		value, ok := schema["enum"]
		if ok {
			return value.([]interface{})
		}
	}

	return nil
}

func getTypeProp(node map[string]interface{}) string {
	value, ok := node["type"]

	if ok {
		return value.(string)
	}

	value, ok = node["schema"]
	if ok {
		schema := value.(map[string]interface{})
		value, ok := schema["type"]
		if ok {
			return value.(string)
		}
	}

	return "reference"
}

func findParam(in []interface{}, name string) map[string]interface{} {
	for _, p := range in {
		param := p.(map[string]interface{})
		if param["name"].(string) == name {
			return param
		}
	}

	return nil
}

func getModelByRef(node map[string]interface{}, spec map[string]interface{}) map[string]interface{} {
	schema := node["schema"].(map[string]interface{})
	ref := schema["ref"].(string)
	//#/definitions/Pet

	paths := strings.Split(ref, "/")

	var currentNode = spec
	for i := 1; i < len(paths); i++ {
		currentNode = currentNode[paths[i]].(map[string]interface{})
	}

	return currentNode

}
