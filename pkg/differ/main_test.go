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

	file, _ := ioutil.ReadFile("specV1.json")
	var spec map[string]interface{}
	_ = json.Unmarshal(file, &spec)

	file, _ = ioutil.ReadFile("specV2.json")
	var changedSpec map[string]interface{}
	_ = json.Unmarshal(file, &changedSpec)

	errs := Diff(spec, changedSpec)

	if errs == nil {
		t.Error("Should produce errs")
	}

	expected := []error{
		errors.New("param id mustn't have enum"),
		errors.New("param name mustn't remove value alex from enum"),
		errors.New("param age mustn't change type from integer to string"),
		errors.New("required param name deleted"),
		errors.New("param age mustn't be required because it wasn't be required"),
		errors.New("get method of /pet/findByTags path mustn't be removed"),
		errors.New("resource /pet/findByStatus mustn't be removed"),
		errors.New("param additionalMetadata mustn't change type from string to integer"),
		errors.New("param petId mustn't be required because it wasn't be required"),
		errors.New("required param body deleted"),
		errors.New("new required param required-param mustn't be added"),
		errors.New("param sort-without-enum mustn't have enum"),
		errors.New("param missed-enum-value mustn't remove value desc from enum"),
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

	fmt.Println("____________________")
	fmt.Println("Expected errors:")
	for _, err := range expected {
		fmt.Println(err)
	}
	fmt.Println("____________________")

}
