// Nice package which allows view difference between json objects
// and also fail build if by some rules
package pkg

import (
	. "differ/pkg/report"
	. "differ/pkg/reporter"
	"github.com/google/go-cmp/cmp"
)

func Diff(obj interface{}, obj2 interface{}) []Report {
	var r DiffReporter
	cmp.Equal(obj, obj2, cmp.Reporter(&r))

	return r.Reports
}
