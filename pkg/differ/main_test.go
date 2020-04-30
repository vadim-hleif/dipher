package differ

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"testing"
)

func TestDiff(t *testing.T) {

	file, _ := ioutil.ReadFile("spec.json")
	var spec map[string]interface{}
	_ = json.Unmarshal(file, &spec)

	file, _ = ioutil.ReadFile("spec-changed.json")
	var changedSpec map[string]interface{}
	_ = json.Unmarshal(file, &changedSpec)

	errs := Diff(spec, changedSpec)

	if errs == nil {
		t.Error("Should produce errs")
	}

	expected := []error{
		errors.New("get method of /pet/findByTags path mustn't be removed"),
		errors.New("resource /pet/findByStatus mustn't be removed"),
		errors.New("param additionalMetadata mustn't change type from string to integer"),
		errors.New("param petId mustn't be required because it wasn't be required"),
		errors.New("required param body deleted"),
		errors.New("new required param required-param mustn't be added"),
	}

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Error() < expected[j].Error()
	})
	sort.Slice(errs, func(i, j int) bool {
		return errs[i].Error() < errs[j].Error()
	})

	if !reflect.DeepEqual(errs, expected) {
		t.Error("expected: ", expected, "got: ", errs)
	}

	fmt.Println("____________________")
	fmt.Println("Actual errors:")
	for _, err := range errs {
		fmt.Println(err)
	}
	fmt.Println("____________________")

}
