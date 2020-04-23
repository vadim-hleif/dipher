// Nice package which allows view difference between json objects
// and also fail build if by some rules
package differ

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"strings"
)

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path  cmp.Path
	paths []string
}

type Report struct {
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

		r.paths = append(r.paths, result.String())
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.paths, "\n")
}

func Diff(obj interface{}, obj2 interface{}) []string {
	var r DiffReporter
	cmp.Equal(obj, obj2, cmp.Reporter(&r))

	return r.paths
}
