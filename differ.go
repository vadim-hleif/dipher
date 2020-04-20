// Nice package which allows view difference between json objects
// and also fail build if by some rules
package differ

import "reflect"

func Diff(obj interface{}, obj2 interface{}) map[string]interface{} {
	if (reflect.TypeOf(obj).Kind() != reflect.Map || reflect.TypeOf(obj2).Kind() != reflect.Map) ||
		(reflect.TypeOf(obj) != reflect.TypeOf(obj2)) {
		// unsupported yet
		return nil
	}

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map:
		original := reflect.ValueOf(obj)
		another := reflect.ValueOf(obj2)

		result := diff(original, another)
		secondPart := diff(another, original)

		for key, value := range secondPart {
			result[key] = value
		}

		if len(result) == 0 {
			return nil
		} else {
			return result
		}

	default:
		return nil
	}
}

func diff(original reflect.Value, another reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}
	it := original.MapRange()

	for it.Next() {
		key := it.Key()
		value := it.Value()

		v := another.MapIndex(key)
		if !v.IsValid() || value.Interface() != v.Interface() {
			result[key.String()] = value.Interface()
		}
	}

	return result
}