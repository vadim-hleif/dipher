package differ

import (
	"fmt"
)

func Diff(specV1 map[string]interface{}, specV2 map[string]interface{}) []error {
	var errs []error
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

				isParamV1Required := getRequiredProp(paramV1)
				isParamV2Required := getRequiredProp(paramV2)

				if paramV2 == nil && isParamV1Required {
					errs = append(errs, fmt.Errorf("required param %v deleted", paramV1["name"].(string)))
				}

				if !isParamV1Required && isParamV2Required {
					errs = append(errs, fmt.Errorf("param %v mustn't be required because it wasn't be required", paramV1["name"].(string)))
				}

				if paramV2 != nil && paramV1 != nil {
					paramType := getTypeProp(paramV1)
					changedParamType := getTypeProp(paramV2)

					if paramType != changedParamType {
						errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", paramV1["name"].(string), paramType, changedParamType))
					}

					enumV1 := getEnum(paramV1)
					enumV2 := getEnum(paramV2)

					if enumV2 == nil && enumV1 == nil {
						continue
					}

					if enumV1 == nil && enumV2 != nil {
						errs = append(errs, fmt.Errorf("param %v mustn't have enum", paramV1["name"].(string)))
					}

					for _, elV1 := range enumV1 {
						var exist bool
						for _, elV2 := range enumV2 {
							if elV1 == elV2 {
								exist = true
								break
							}
						}

						if !exist {
							errs = append(errs, fmt.Errorf("param %v mustn't remove value %v from enum", paramV1["name"].(string), elV1))
						}
					}
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

	return ""
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
