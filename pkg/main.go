// Package pkg provides api which can found breaking changes in swagger 2.0 specifications
package pkg

// Report contains jsonPath and errors which related to that jsonPath
type Report struct {
	Err      error
	JSONPath string
}

// Diff returns slice with errors if any breaking changes was founded
//
// returns empty array if there aren't any errors
func Diff(specV1, specV2 map[string]interface{}) []Report {
	dipher := dipher{
		specV1:  specV1,
		specV2:  specV2,
		reports: make([]Report, 0),
	}

	return dipher.diff()
}
