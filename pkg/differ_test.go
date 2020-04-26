package pkg

import (
	. "differ/pkg/report"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestDiff_sameDocuments(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{})

	if result != nil {
		t.Error("Should return nil if there are no differences, but returned: ", result)
	}
}

func TestDiff_differentKeys(t *testing.T) {
	result := Diff(map[string]string{
		"name": "John",
	}, map[string]string{
		"other-name": "John",
	})
	eraseActualValues(result)

	expected := []Report{{
		JSONPath: "name",
		Diff:     "removed",
	}, {
		JSONPath: "other-name",
		Diff:     "added",
	}}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}

func TestDiff_sameKeys(t *testing.T) {
	result := Diff(map[string]string{
		"name": "John",
	}, map[string]string{
		"name": "Other name",
	})
	eraseActualValues(result)

	expected := []Report{{
		JSONPath: "name",
		Diff:     "value_changed",
	}}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}

func TestDiff_nestedDifferentValues(t *testing.T) {
	result := Diff(
		map[string]interface{}{
			"name": map[string]interface{}{
				"second-level": "value",
			},
		}, map[string]interface{}{
			"name": map[string]interface{}{
				"second-level": "value2",
			},
		})
	eraseActualValues(result)

	expected := []Report{{
		JSONPath: "name.second-level",
		Diff:     "value_changed",
	}}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}

func TestDiff_realSwagger(t *testing.T) {
	file, _ := ioutil.ReadFile("../old-swagger.json")
	var oldSwagger interface{}
	_ = json.Unmarshal(file, &oldSwagger)

	file, _ = ioutil.ReadFile("../new-swagger.json")
	var newSwagger interface{}
	_ = json.Unmarshal(file, &newSwagger)

	result := Diff(oldSwagger, newSwagger)
	eraseActualValues(result)

	expected := []Report{{
		JSONPath: "paths./http_test/test_get/{pathParam}",
		Diff:     "removed",
	}, {
		JSONPath: "paths./http_test/test_get2/{pathParam}",
		Diff:     "added",
	}, {
		JSONPath: "paths./proxy-config/api/v1/proxy-config/environments.get.parameters",
		Diff:     "added",
	}}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}

func eraseActualValues(r []Report) {
	for i := range r {
		r[i].ActualValue = nil
	}
}
