package differ

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"sort"
	"testing"
)

type test struct {
	specsPath string
	want      []error
}

// REQUEST OBJECTS
func TestDiff_should_detect_removing_enum_value_in_any_request_object_property(t *testing.T) {
	runTest(t, test{
		specsPath: "params-object/removed_value_from_old_enum",
		want: []error{
			errors.New("param name mustn't remove value alex from enum"),
		},
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_object_roperty(t *testing.T) {
	runTest(t, test{
		specsPath: "params-object/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param age mustn't have enum"),
		},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_object_property(t *testing.T) {
	runTest(t, test{
		specsPath: "params-object/type_changing",
		want: []error{
			errors.New("param id mustn't change type from integer to string"),
		},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_request_object(t *testing.T) {
	runTest(t, test{
		specsPath: "params-object/required_property_deletion",
		want: []error{
			errors.New("required param id mustn't be deleted"),
		},
	})
}

func TestDiff_should_detect_new_required_property_in_the_request_object(t *testing.T) {
	runTest(t, test{
		specsPath: "params-object/required_property_adding",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	})
}

func TestDiff_should_detect_removing_enum_value_in_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "params/removed_value_from_old_enum",
		want: []error{
			errors.New("param sort mustn't remove value asc from enum"),
		},
	})
}

// END REQUEST OBJECTS

// REQUEST PARAMS
func TestDiff_should_detect_adding_enum_in_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "params/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param sort mustn't have enum"),
		},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "params/required_param_deletion",
		want: []error{
			errors.New("required param sort mustn't be deleted"),
		},
	})
}

func TestDiff_should_detect_new_required_param_in_the_request(t *testing.T) {
	runTest(t, test{
		specsPath: "params/new_required_param",
		want: []error{
			errors.New("new required param filter mustn't be added"),
		},
	})
}

func TestDiff_should_detect_marking_old_param_as_required(t *testing.T) {
	runTest(t, test{
		specsPath: "params/mark_old_param_as_required",
		want: []error{
			errors.New("param sort mustn't be required because it wasn't be required"),
		},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "params/type_changing",
		want: []error{
			errors.New("param sort mustn't change type from string to integer"),
		},
	})
}

// END REQUEST PARAMS

// REQUEST OBJECTS WITH DEFINITION
func TestDiff_should_detect_new_required_property_in_the_request_object_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/required_property_adding",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_request_object_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/required_property_deletion",
		want: []error{
			errors.New("required param name mustn't be deleted"),
		},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/type_changing",
		want: []error{
			errors.New("param id mustn't change type from integer to string"),
		},
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param name mustn't have enum"),
		},
	})
}

func TestDiff_should_detect_removing_enum_value_in_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/removed_value_from_old_enum",
		want: []error{
			errors.New("param name mustn't remove value alex from enum"),
		},
	})
}

func TestDiff_should_detect_object_definition_in_any_path(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/custom_def_path",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	})
}

func TestDiff_should_compare_models_by_different_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/same_models_with_different_names",
		want:      make([]error, 0),
	})
}

func TestDiff_should_detect_diff_in_models_with_different_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/different_models_with_different_names",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
			errors.New("param age mustn't have enum"),
			errors.New("param name mustn't change type from string to integer"),
			errors.New("param name mustn't remove value alex from enum"),
			errors.New("required param id mustn't be deleted"),
		},
	})
}

func TestDiff_should_detect_diff_in_nested_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "definitions/models_with_recursive_refs",
		want: []error{
			errors.New("param zipCode mustn't be required because it wasn't be required"),
			errors.New("param city mustn't remove value NY from enum"),
			errors.New("param name mustn't have enum"),
			errors.New("required param id mustn't be deleted"),
			errors.New("required param city mustn't be deleted"),
		},
	})
}

// END REQUEST OBJECTS WITH DEFINITION

// RESOURCES AND VERBS TEST CASES
func TestDiff_should_detect_path_removing(t *testing.T) {
	runTest(t, test{
		specsPath: "paths_and_verbs/resource_removing",
		want: []error{
			errors.New("resource /pet mustn't be removed"),
		},
	})
}

func TestDiff_should_detect_method_removing_in_the_any_path(t *testing.T) {
	runTest(t, test{
		specsPath: "paths_and_verbs/verb_removing",
		want: []error{
			errors.New("post method of /pet path mustn't be removed"),
		},
	})
}

// END RESOURCES AND VERBS TEST CASES

// test helper
func runTest(t *testing.T, tt test) {
	file, _ := ioutil.ReadFile("test-specs/" + tt.specsPath + "/V1.json")
	var specV1 map[string]interface{}
	_ = json.Unmarshal(file, &specV1)

	file, _ = ioutil.ReadFile("test-specs/" + tt.specsPath + "/V2.json")
	var specV2 map[string]interface{}
	_ = json.Unmarshal(file, &specV2)

	got := Diff(specV1, specV2)

	sort.Slice(got, func(i, j int) bool {
		return got[i].Error() < got[j].Error()
	})
	sort.Slice(tt.want, func(i, j int) bool {
		return tt.want[i].Error() < tt.want[j].Error()
	})

	//for i, err := range got {
	//	fmt.Println(err)
	//	fmt.Println(tt.want[i])
	//	fmt.Println(err.Error() == tt.want[i].Error())
	//}

	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("Diff() = %v, want %v", got, tt.want)
	}
}
