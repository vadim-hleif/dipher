package differ

import (
	"encoding/json"
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

	expected := map[string]string{
		"other-name": "John",
		"name":       "John",
	}
	expectedBytes, _ := json.Marshal(expected)
	resultBytes, _ := json.Marshal(result)

	if string(expectedBytes) != string(resultBytes) {
		t.Error("expected: ", expected, "got: ", result)
	}
}

func TestDiff_sameKeys(t *testing.T) {
	result := Diff(map[string]string{
		"name": "John",
	}, map[string]string{
		"name": "Other name",
	})

	expected := map[string]string{
		"name": "Other name",
	}
	expectedBytes, _ := json.Marshal(expected)
	resultBytes, _ := json.Marshal(result)

	if string(expectedBytes) != string(resultBytes) {
		t.Error("expected: ", string(expectedBytes), "got: ", string(resultBytes))
	}
}
