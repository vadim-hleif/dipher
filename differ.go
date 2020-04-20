// Nice tool which allows view difference between json objects
//and also fail build if by some rules
package differ

func Diff(obj interface{}, obj2 interface{}) interface{} {
	original, first := obj.(map[string]string)
	another, second := obj2.(map[string]string)
	if !first || !second {
		return nil
	}

	result := diff(another, original)
	secondPart := diff(original, another)

	for key, value := range secondPart {
		result[key] = value
	}

	if len(result) == 0 {
		return nil
	} else {
		return result
	}
}

func diff(original map[string]string, another map[string]string) map[string]string {
	result := map[string]string{}

	for key, value := range another {
		if v, ok := original[key]; ok {
			if value != v {
				result[key] = value
			}
		} else {
			result[key] = value
		}
	}

	return result
}