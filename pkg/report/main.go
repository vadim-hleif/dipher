package report

import (
	"bytes"
	"fmt"
)

// Report contains jsonPath and difference description
type Report struct {
	JSONPath    string
	Diff        string
	ActualValue interface{}
}

// Reports it's slice of Report
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
	return fmt.Sprintf("%v\n\t%v\n\t%v", r.JSONPath, r.Diff, r.ActualValue)
}
