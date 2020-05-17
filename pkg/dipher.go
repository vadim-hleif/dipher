package pkg

import "fmt"

type dipher struct {
	specV1 map[string]interface{}
	specV2 map[string]interface{}
}

func (dipher *dipher) diff() []Report {
	errs := make([]Report, 0)

	pathsV2 := dipher.specV2["paths"].(map[string]interface{})

	for url, urlNodeV1 := range dipher.specV1["paths"].(map[string]interface{}) {
		urlNodeV2 := getNode(pathsV2, url)

		if urlNodeV2 == nil {
			errs = append(errs, Report{
				Err:      fmt.Errorf("resource %v mustn't be removed", url),
				JSONPath: "$",
			})
			continue
		}

		for methodV1, m := range urlNodeV1.(map[string]interface{}) {
			methodPath := fmt.Sprintf("$.%v.%v", url, methodV1)

			methodNodeV1 := m.(map[string]interface{})
			methodNodeV2 := getNode(urlNodeV2, methodV1)

			if methodNodeV2 == nil {
				errs = append(errs, Report{
					Err:      fmt.Errorf("%v method of %v path mustn't be removed", methodV1, url),
					JSONPath: methodPath,
				})
				continue
			}

			paramsV1 := methodNodeV1["parameters"].([]interface{})
			paramsV2 := methodNodeV2["parameters"].([]interface{})

			for _, p := range paramsV1 {
				localParamPath := fmt.Sprintf("%v.parameters", methodPath)

				paramV1 := p.(map[string]interface{})
				paramV2 := findParam(paramsV2, paramV1["name"].(string))

				typeV1, _ := getTypeProp(paramV1)
				switch typeV1 {
				case "reference", "object":
					errs = append(errs, toReports(dipher.compareParams(paramV1, paramV2), localParamPath)...)
				default:
					errs = append(errs,
						toReports(compareSimpleParams(paramV1, paramV2), localParamPath)...)
				}
			}

			for _, p := range paramsV2 {
				paramV2 := p.(map[string]interface{})

				if findParam(paramsV1, paramV2["name"].(string)) == nil && getRequiredProp(paramV2) {
					errs = append(errs, Report{
						Err:      fmt.Errorf("new required param %v mustn't be added", paramV2["name"].(string)),
						JSONPath: fmt.Sprintf("%v.parameters", methodPath),
					})
				}
			}

			responsesV1 := methodNodeV1["responses"].(map[string]interface{})
			responsesV2 := methodNodeV2["responses"].(map[string]interface{})

			responsesPath := fmt.Sprintf("%v.responses", methodPath)

			for code, c := range responsesV1 {
				responseV1 := c.(map[string]interface{})
				responseV2, ok := responsesV2[code].(map[string]interface{})

				if !ok {
					errs = append(errs, Report{
						Err:      fmt.Errorf("response with code %v mustn't be removed", code),
						JSONPath: responsesPath,
					})
				} else {
					errs = append(errs,
						toReports(dipher.compareResponse(responseV1, responseV2), responsesPath)...)
				}
			}
		}
	}

	return errs
}
