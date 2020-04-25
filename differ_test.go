package differ

import (
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

	expected := []Report{{
		jsonPath: "name",
		diff:     "removed",
	}, {
		jsonPath: "other-name",
		diff:     "added",
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

	expected := []Report{{
		jsonPath: "name",
		diff:     "value_changed",
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

	expected := []Report{{
		jsonPath: "name.second-level",
		diff:     "value_changed",
	}}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}
