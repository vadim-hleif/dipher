package pkg

import (
	"strings"
)

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
	value, ok := node["enum"]

	if ok {
		return value.([]interface{})
	}

	value, ok = node["schema"]
	_, isArray, _ := getMetadata(node)

	if ok {
		node = value.(map[string]interface{})
	}

	if isArray {
		node = node["items"].(map[string]interface{})
	}

	value, ok = node["enum"]
	if ok {
		return value.([]interface{})
	}

	return nil
}

func getMetadata(node map[string]interface{}) (string, bool, string) {
	value, ok := node["schema"]
	if ok {
		node = value.(map[string]interface{})
	}

	var elType string
	t, ok := node["type"]

	if ok {
		elType = t.(string)
	}

	isArray := elType == "array"
	if isArray {
		node = node["items"].(map[string]interface{})

		t, ok := node["type"]

		if ok {
			elType = t.(string)
		}
	}

	var ref string
	i, ok := node["$ref"]
	if ok {
		ref = i.(string)
		elType = "reference"
	}

	return elType, isArray, ref
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

func getModelByRef(ref string, spec map[string]interface{}) map[string]interface{} {
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
