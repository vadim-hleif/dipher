// Nice package which allows view difference between json objects
// and also fail build if by some rules
package report

import (
	"bytes"
	"fmt"
)

type Report struct {
	JsonPath string
	Diff     string
}

type Reports []Report

func (r Reports) String() string {
	var result bytes.Buffer
	for _, report := range r {
		result.WriteString(report.String())
		result.WriteString("\n")
	}
	return result.String()
}

func (r *Report) String() string {
	return fmt.Sprintf("%v\n\t%v", r.JsonPath, r.Diff)
}
