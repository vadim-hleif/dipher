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

	expected := map[string]interface{}{
		"other-name": "John",
		"name":       "John",
	}

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

	expected := map[string]interface{}{
		"name": "Other name",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Error("expected: ", expected, "got: ", result)
	}
}
