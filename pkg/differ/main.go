package differ

import (
	"fmt"
)

func Diff(spec map[string]interface{}, changedSpec map[string]interface{}) []error {
	var result []error
	changesPaths := changedSpec["paths"].(map[string]interface{})

	for url, node := range spec["paths"].(map[string]interface{}) {
		changedNode := getNode(changesPaths, url)

		if changedNode == nil {
			result = append(result, fmt.Errorf("resource %v mustn't be removed", url))
			continue
		}

		for method, methodNode := range node.(map[string]interface{}) {
			resource := methodNode.(map[string]interface{})
			changedMethod := getNode(changedNode, method)

			if changedMethod == nil {
				result = append(result, fmt.Errorf("%v method of %v path mustn't be removed", method, url))
				continue
			}

			params := resource["parameters"].([]interface{})
			changedParams := changedMethod["parameters"].([]interface{})

			for _, p := range params {
				param := p.(map[string]interface{})
				changedParam := findParam(changedParams, param["name"].(string))

				paramRequired := getRequiredProp(param)
				changedParamRequired := getRequiredProp(changedParam)

				if changedParam == nil && paramRequired {
					result = append(result, fmt.Errorf("required param %v deleted", param["name"].(string)))
				}

				if !paramRequired && changedParamRequired {
					result = append(result, fmt.Errorf("param %v mustn't be required because it wasn't be required", param["name"].(string)))
				}

				if changedParam != nil && param != nil {
					paramType := getTypeProp(param)
					changedParamType := getTypeProp(changedParam)

					if paramType != changedParamType {
						result = append(result, fmt.Errorf("param %v mustn't change type from %v to %v", param["name"].(string), paramType, changedParamType))
					}
				}
			}

			for _, p := range changedParams {
				changedParam := p.(map[string]interface{})
				if findParam(params, changedParam["name"].(string)) == nil && getRequiredProp(changedParam) {
					result = append(result, fmt.Errorf("new required param %v mustn't be added", changedParam["name"].(string)))
				}
			}
		}

	}

	return result
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

func getTypeProp(node map[string]interface{}) string {
	value, ok := node["type"]

	if ok {
		return value.(string)
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
