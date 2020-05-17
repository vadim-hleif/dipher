package pkg

// Report contains jsonPath and error under that
type Report struct {
	Err      error
	JSONPath string
}

// Diff returns slice with errors if any breaking changes was founded
//
// returns empty array if there aren't any errors
func Diff(specV1 map[string]interface{}, specV2 map[string]interface{}) []Report {
	dipher := dipher{
		specV1: specV1,
		specV2: specV2,
	}

	return dipher.diff()
}
