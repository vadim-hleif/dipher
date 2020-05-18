package pkg

import "fmt"

func (dipher *dipher) compareParams(paramV1 map[string]interface{}, paramV2 map[string]interface{}) []error {
	errs := make([]error, 0)

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
		return dipher.compareParams(getModelByRef(paramV1, dipher.specV1), getModelByRef(paramV2, dipher.specV2))
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

		if !ok {
			continue
		}

		typeV1, _ := getTypeProp(propsV1)
		typeV2, _ := getTypeProp(propsV2)

		if typeV1 == "reference" {
			errs = append(errs,
				dipher.compareParams(getModelByRef(propsV1, dipher.specV1), getModelByRef(propsV2, dipher.specV2))...)
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
	return errs
}
