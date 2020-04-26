// Package pkg allows view difference between json objects
// and also fail build if by some rules
package pkg

import (
	"differ/pkg/report"
	"differ/pkg/reporter"
	"github.com/google/go-cmp/cmp"
)

// Diff return diffs between two jsons
func Diff(obj interface{}, obj2 interface{}) []report.Report {
	var r reporter.DiffReporter
	cmp.Equal(obj, obj2, cmp.Reporter(&r))

	return r.Reports
}
