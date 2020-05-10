package differ

import (
	"fmt"
	"strings"
)

// Diff returns slice with errors if any breaking changes was founded
//
// returns empty array if there aren't any errors
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

				typeV1, _ := getTypeProp(paramV1)
				switch typeV1 {
				case "reference", "object":
					errs = compareObjectParams(paramV1, paramV2, specV1, specV2, errs)
				default:
					errs = comparePrimitiveParams(paramV1, paramV2, errs)
				}
			}

			for _, p := range paramsV2 {
				paramV2 := p.(map[string]interface{})
				if findParam(paramsV1, paramV2["name"].(string)) == nil && getRequiredProp(paramV2) {
					errs = append(errs, fmt.Errorf("new required param %v mustn't be added", paramV2["name"].(string)))
				}
			}

			responsesV1 := methodNodeV1["responses"].(map[string]interface{})
			responsesV2 := methodNodeV2["responses"].(map[string]interface{})

			for code, c := range responsesV1 {
				responseV1 := c.(map[string]interface{})
				responseV2, ok := responsesV2[code].(map[string]interface{})

				if !ok {
					errs = append(errs, fmt.Errorf("response with code %v mustn't be removed", code))
				} else {
					errs = compareResponseObjects(responseV1, responseV2, specV1, specV2, errs)
				}
			}
		}
	}

	return errs
}

func comparePrimitiveParams(paramV1 map[string]interface{}, paramV2 map[string]interface{}, errs []error) []error {
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

	typeV1, _ := getTypeProp(paramV1)
	typeV2, _ := getTypeProp(paramV2)

	if typeV1 != typeV2 {
		errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", paramV1["name"].(string), typeV1, typeV2))
	}

	return errs
}

func compareResponseObjects(responseV1 map[string]interface{}, responseV2 map[string]interface{},
	specV1 map[string]interface{}, specV2 map[string]interface{}, errs []error) []error {
	schemaV1 := getNode(responseV1, "schema")
	if schemaV1 == nil {
		schemaV1 = responseV1
	}
	schemaV2 := getNode(responseV2, "schema")
	if schemaV2 == nil {
		schemaV2 = responseV2
	}

	typeV1, _ := getTypeProp(responseV1)
	if typeV1 == "reference" {
		return compareResponseObjects(getModelByRef(responseV1, specV1), getModelByRef(responseV2, specV2),
			specV1, specV2, errs)
	}

	pV2 := getNode(schemaV2, "properties")

	for nameV1, p := range getNode(schemaV1, "properties") {
		propsV1 := p.(map[string]interface{})
		propsV2, ok := pV2[nameV1].(map[string]interface{})

		if ok {
			typeV1, _ := getTypeProp(propsV1)
			typeV2, _ := getTypeProp(propsV2)

			if typeV1 == "reference" {
				errs = compareResponseObjects(getModelByRef(propsV1, specV1), getModelByRef(propsV2, specV2),
					specV1, specV2, errs)
			}

			if typeV1 != typeV2 {
				errs = append(errs, fmt.Errorf("response field %v mustn't change type from %v to %v", nameV1, typeV1, typeV2))
			}
		} else {
			errs = append(errs, fmt.Errorf("response field %v mustn't be deleted", nameV1))
		}

	}
	return errs
}

func compareObjectParams(paramV1 map[string]interface{}, paramV2 map[string]interface{},
	specV1 map[string]interface{}, specV2 map[string]interface{}, errs []error) []error {
	schemaV1 := getNode(paramV1, "schema")
	if schemaV1 == nil {
		schemaV1 = paramV1
	}
	schemaV2 := getNode(paramV2, "schema")
	if schemaV2 == nil {
		schemaV2 = paramV2
	}

	typeV1, _ := getTypeProp(paramV1)
	if typeV1 == "reference" {
		return compareObjectParams(getModelByRef(paramV1, specV1), getModelByRef(paramV2, specV2),
			specV1, specV2, errs)
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

	for nameV1, p := range getNode(schemaV1, "properties") {
		propsV1 := p.(map[string]interface{})
		propsV2, ok := pV2[nameV1].(map[string]interface{})
		if ok {
			typeV1, _ := getTypeProp(propsV1)
			typeV2, _ := getTypeProp(propsV2)

			if typeV1 == "reference" {
				errs = compareObjectParams(getModelByRef(propsV1, specV1), getModelByRef(propsV2, specV2),
					specV1, specV2, errs)
			}

			if typeV1 != typeV2 {
				errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", nameV1, typeV1, typeV2))
			}

			enumV1 := getEnum(propsV1)
			enumV2 := getEnum(propsV2)
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
	if _, isArray := getTypeProp(node); isArray {
		items := node["items"].(map[string]interface{})

		if enum, ok := items["enum"]; ok {
			return enum.([]interface{})
		}
	}

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

func getTypeProp(node map[string]interface{}) (string, bool) {
	value, ok := node["type"]
	var elType string

	if ok {
		elType = value.(string)
	}

	value, ok = node["schema"]
	if ok {
		schema := value.(map[string]interface{})
		value, ok := schema["type"]
		if ok {
			elType = value.(string)
		} else if _, ok := schema["$ref"]; ok {
			elType = "reference"
		}
	}

	if _, ok := node["$ref"]; ok {
		elType = "reference"
	}

	var isArray bool
	if elType == "array" {
		isArray = true
		items := node["items"].(map[string]interface{})

		t, ok := items["type"]

		if ok {
			elType = t.(string)
		} else if _, ok := items["$ref"]; ok {
			elType = "reference"
		}
	}

	return elType, isArray
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
	schema, ok := node["schema"].(map[string]interface{})
	var ref string

	if _, isArray := getTypeProp(node); isArray {
		ref = node["items"].(map[string]interface{})["$ref"].(string)
	} else if ok {
		ref = schema["$ref"].(string)
	} else {
		ref = node["$ref"].(string)
	}

	paths := strings.Split(ref, "/")

	var currentNode = spec
	for i := 1; i < len(paths); i++ {
		currentNode = currentNode[paths[i]].(map[string]interface{})
	}

	return currentNode

}
