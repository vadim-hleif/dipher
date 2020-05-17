package pkg

import "strings"

// compareAndApply compares every elements from source every element from target
// and calls cb func with element, that missed is target but exist is source
func compareAndApply(source []interface{}, target []interface{}, cb func(name interface{})) {
	for _, elem := range source {
		var exist bool

		for _, targetElem := range target {
			if elem == targetElem {
				exist = true
				break
			}
		}

		if !exist {
			cb(elem)
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

func toReports(err []error, localParamPath string) []Report {
	temp := make([]Report, len(err))

	for i := 0; i < len(err); i++ {
		temp[i] = Report{
			Err:      err[i],
			JSONPath: localParamPath,
		}
	}
	return temp
}
