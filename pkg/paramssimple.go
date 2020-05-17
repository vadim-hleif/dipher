package pkg

import "fmt"

func compareSimpleParams(paramV1 map[string]interface{}, paramV2 map[string]interface{}) []error {
	errs := make([]error, 0)
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

	if enumV1 == nil && enumV2 != nil {
		errs = append(errs, fmt.Errorf("param %v mustn't have enum", paramV1["name"].(string)))
	}

	compareAndApply(enumV1, enumV2, func(name interface{}) {
		errs = append(errs, fmt.Errorf("param %v mustn't remove value %v from enum", paramV1["name"].(string), name))
	})

	typeV1, _ := getTypeProp(paramV1)
	typeV2, _ := getTypeProp(paramV2)

	if typeV1 != typeV2 {
		errs = append(errs, fmt.Errorf("param %v mustn't change type from %v to %v", paramV1["name"].(string), typeV1, typeV2))
	}

	return errs
}
