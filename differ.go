// Nice package which allows view difference between json objects
// and also fail build if by some rules
package differ

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"strings"
)

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path    cmp.Path
	reports []Report
}

type Report struct {
	jsonPath string
	diff     string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		var result bytes.Buffer
		for _, step := range r.path {
			node := step.String()
			r := regexp.MustCompile("\\[\"(.*)\"]")
			match := r.FindStringSubmatch(node)
			if len(match) > 0 {
				result.WriteString(match[1])
				result.WriteString(".")
			}
		}
		result.Truncate(result.Len() - 1)

		oldProp, newProp := r.path.Last().Values()
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

		r.reports = append(r.reports, Report{
			jsonPath: result.String(),
			diff:     diff,
		})
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	var result []string
	for _, v := range r.reports {
		result = append(result, v.String())
	}
	return strings.Join(result, "\n")
}

func (r *Report) String() string {
	return fmt.Sprint("#v \n\t #v", r.jsonPath, r.diff)
}

func Diff(obj interface{}, obj2 interface{}) []Report {
	var r DiffReporter
	cmp.Equal(obj, obj2, cmp.Reporter(&r))

	return r.reports
}
