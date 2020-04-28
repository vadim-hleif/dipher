package differ

import (
	"fmt"
)

func Diff(spec map[string]interface{}, changedSpec map[string]interface{}) []error {
	var result []error
	changesPaths := changedSpec["paths"].(map[string]interface{})

	for path, node := range spec["paths"].(map[string]interface{}) {
		changedNode := changesPaths[path].(map[string]interface{})

		for path, method := range node.(map[string]interface{}) {
			resource := method.(map[string]interface{})
			changedResource := changedNode[path].(map[string]interface{})

			params := resource["parameters"].([]interface{})
			changedParams := changedResource["parameters"].([]interface{})

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
