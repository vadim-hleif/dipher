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

	expected := []string{"name", "other-name"}

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

	expected := []string{"name"}

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

	expected := []string{"name.second-level"}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}
