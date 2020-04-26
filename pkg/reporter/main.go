package reporter

import (
	"bytes"
	"differ/pkg/report"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"strings"
)

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	Path    cmp.Path
	Reports report.Reports
}

// Report custom logic to aggregate difference between jsons
func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		var result bytes.Buffer
		for _, step := range r.Path {
			node := step.String()
			r := regexp.MustCompile("\\[\"(.*)\"]")
			match := r.FindStringSubmatch(node)
			if len(match) > 0 {
				result.WriteString(match[1])
				result.WriteString(".")
			}
		}
		result.Truncate(result.Len() - 1)

		oldProp, newProp := r.Path.Last().Values()
		var diff string
		if oldProp.IsValid() && !newProp.IsValid() {
			diff = "removed"
		} else if !oldProp.IsValid() && newProp.IsValid() {
			diff = "added"
		}

		if oldProp.IsValid() && newProp.IsValid() {
			if oldProp.Type() != newProp.Type() {
				diff = "type_changed"
			} else if oldProp.Interface() != newProp.Interface() {
				diff = "value_changed"
			}
		}
		var actualValue interface{}
		if newProp.IsValid() {
			actualValue = newProp.Interface()
		}

		r.Reports = append(r.Reports, report.Report{
			JSONPath:    result.String(),
			Diff:        diff,
			ActualValue: actualValue,
		})
	}
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.Path = append(r.Path, ps)
}

func (r *DiffReporter) PopStep() {
	r.Path = r.Path[:len(r.Path)-1]
}

func (r *DiffReporter) String() string {
	var result []string
	for _, v := range r.Reports {
		result = append(result, v.String())
	}
	return strings.Join(result, "\n")
}
