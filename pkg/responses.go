package pkg

import (
	"fmt"
)

func (dipher *dipher) compareResponse(responseV1 map[string]interface{}, responseV2 map[string]interface{}) []error {
	errs := make([]error, 0)

	schemaV1 := getNode(responseV1, "schema")
	if schemaV1 == nil {
		schemaV1 = responseV1
	}
	schemaV2 := getNode(responseV2, "schema")
	if schemaV2 == nil {
		schemaV2 = responseV2
	}

	typeV1, _, refV1 := getMetadata(responseV1)
	if typeV1 == "reference" {
		_, _, refV2 := getMetadata(responseV2)
		modelV1 := getModelByRef(refV1, dipher.specV1)
		modelV2 := getModelByRef(refV2, dipher.specV2)

		errs = append(errs, dipher.compareResponse(modelV1, modelV2)...)
	}

	pV2 := getNode(schemaV2, "properties")

	for nameV1, p := range getNode(schemaV1, "properties") {
		propsV1 := p.(map[string]interface{})
		propsV2, ok := pV2[nameV1].(map[string]interface{})

		if ok {
			typeV1, _, refV1 := getMetadata(propsV1)
			typeV2, _, refV2 := getMetadata(propsV2)

			if typeV1 == "reference" {
				modelV1 := getModelByRef(refV1, dipher.specV1)
				modelV2 := getModelByRef(refV2, dipher.specV2)

				if !dipher.containsDef(refV1) {
					dipher.addDefs(refV1, refV2)
					errs = append(errs, dipher.compareResponse(modelV1, modelV2)...)
				}

			}

			if typeV1 != typeV2 {
				errs = append(errs, fmt.Errorf("response field %v mustn't change type from %v to %v", nameV1, typeV1, typeV2))
			}
		} else {
			errs = append(errs, fmt.Errorf("response field %v mustn't be deleted", nameV1))
		}

	}

	dipher.dropRefs()
	return errs
}
