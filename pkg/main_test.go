package pkg

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
	want      []Report
}

// REQUEST OBJECTS
func TestDiff_should_detect_removing_enum_value_in_any_request_object_property(t *testing.T) {
	runTest(t, test{
		specsPath: "request/params-object/removed_value_from_old_enum",
		want: []Report{{
			Err:      errors.New("param name mustn't remove value alex from enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_object_roperty(t *testing.T) {
	runTest(t, test{
		specsPath: "request/params-object/adding_enum_to_old_value_without_enum",
		want: []Report{{
			Err:      errors.New("param age mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_object_property(t *testing.T) {
	runTest(t, test{
		specsPath: "request/params-object/type_changing",
		want: []Report{{
			Err:      errors.New("param id mustn't change type from integer to string"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_request_object(t *testing.T) {
	runTest(t, test{
		specsPath: "request/params-object/required_property_deletion",
		want: []Report{{
			Err:      errors.New("required param id mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_new_required_property_in_the_request_object(t *testing.T) {
	runTest(t, test{
		specsPath: "request/params-object/required_property_adding",
		want: []Report{{
			Err:      errors.New("param age mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

// END REQUEST OBJECTS

// REQUEST PARAMS

func TestDiff_should_handle_case_with_schema_and_array_def_in_params(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/schema_with_array",
		want:      []Report{},
	})
}

func TestDiff_should_detect_removing_enum_value_in_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/removed_value_from_old_enum",
		want: []Report{{
			Err:      errors.New("param sort mustn't remove value asc from enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/adding_enum_to_old_value_without_enum",
		want: []Report{{
			Err:      errors.New("param sort mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/required_param_deletion",
		want: []Report{{
			Err:      errors.New("required param sort mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_new_required_param_in_the_request(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/new_required_param",
		want: []Report{{
			Err:      errors.New("new required param filter mustn't be added"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_marking_old_param_as_required(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/mark_old_param_as_required",
		want: []Report{{
			Err:      errors.New("param sort mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_param(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/type_changing",
		want: []Report{{
			Err:      errors.New("param sort mustn't change type from string to integer"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_type_chaning_in_array(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/type_changing_in_array",
		want: []Report{{
			Err:      errors.New("param sort mustn't change type from string to integer"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}
func TestDiff_should_detect_adding_enum_old_value_array_value_without_enum(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/adding_enum_to_old_value_array_value_without_enum",
		want: []Report{{
			Err:      errors.New("param sort mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_removing_value_from_old_array_enum(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/params/removed_value_from_old_array_enum",
		want: []Report{{
			Err:      errors.New("param sort mustn't remove value desc from enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

// END REQUEST PARAMS

// REQUEST OBJECTS WITH DEFINITION
func TestDiff_should_detect_new_required_property_in_the_request_object_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/required_property_adding",
		want: []Report{{
			Err:      errors.New("param age mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_required_property_removing_in_the_request_object_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/required_property_deletion",
		want: []Report{{
			Err:      errors.New("required param name mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_type_changing_in_the_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/type_changing",
		want: []Report{{
			Err:      errors.New("param id mustn't change type from integer to string"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/adding_enum_to_old_value_without_enum",
		want: []Report{{
			Err:      errors.New("param age mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_removing_enum_value_in_any_request_object_property_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/removed_value_from_old_enum",
		want: []Report{{
			Err:      errors.New("param name mustn't remove value alex from enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_object_definition_in_any_path(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/custom_def_path",
		want: []Report{{
			Err:      errors.New("param age mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_compare_models_by_different_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/same_models_with_different_names",
		want:      make([]Report, 0),
	})
}

func TestDiff_should_detect_adding_enum_in_any_request_object_property_definition_array(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/adding_enum_to_old_array_value_without_enum",
		want: []Report{{
			Err:      errors.New("param age mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_diff_in_models_with_different_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/different_models_with_different_names",
		want: []Report{{
			Err:      errors.New("param age mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param age mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param name mustn't change type from string to integer"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param name mustn't remove value alex from enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("required param id mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_diff_in_models_with_different_refs_array(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/different_models_with_different_names_array",
		want: []Report{{
			Err:      errors.New("param age mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param age mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param name mustn't change type from string to integer"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param name mustn't remove value alex from enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("required param id mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_detect_diff_in_nested_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/models_with_nested_refs",
		want: []Report{{
			Err:      errors.New("param zipCode mustn't be required because it wasn't be required"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param city mustn't remove value NY from enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("param name mustn't have enum"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("required param id mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}, {
			Err:      errors.New("required param city mustn't be deleted"),
			JSONPath: "$./pet.post.parameters",
		}},
	})
}

func TestDiff_should_handle_self_reference(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/models_with_self_ref",
		want:      []Report{},
	})
}

func TestDiff_should_handle_case_with_recursive_refs_between_chain_of_models(t *testing.T) {
	runTest(t, test{
		specsPath: "/request/definitions/two_models_with_recursive_refs",
		want:      []Report{},
	})
}

// END REQUEST OBJECTS WITH DEFINITION

// RESOURCES AND VERBS TEST CASES
func TestDiff_should_detect_path_removing(t *testing.T) {
	runTest(t, test{
		specsPath: "paths_and_verbs/resource_removing",
		want: []Report{{
			Err:      errors.New("resource /pet mustn't be removed"),
			JSONPath: "$",
		}},
	})
}

func TestDiff_should_detect_method_removing_in_the_any_path(t *testing.T) {
	runTest(t, test{
		specsPath: "paths_and_verbs/verb_removing",
		want: []Report{{
			Err:      errors.New("post method of /pet path mustn't be removed"),
			JSONPath: "$./pet.post",
		}},
	})
}

// END RESOURCES AND VERBS TEST CASES

// RESPONSE
func TestDiff_should_detect_code_node_removing_in_response(t *testing.T) {
	runTest(t, test{
		specsPath: "response/code_removing",
		want: []Report{{
			Err:      errors.New("response with code 200 mustn't be removed"),
			JSONPath: "$./pet.post.responses",
		}},
	})
}

func TestDiff_should_detect_type_changes_in_response_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "response/definitions/type_changing",
		want: []Report{{
			Err:      errors.New("response field id mustn't change type from string to integer"),
			JSONPath: "$./pet.post.responses",
		}},
	})
}

func TestDiff_should_detect_field_removing_in_response_definition(t *testing.T) {
	runTest(t, test{
		specsPath: "response/definitions/field_removing",
		want: []Report{{
			Err:      errors.New("response field id mustn't be deleted"),
			JSONPath: "$./pet.post.responses",
		}},
	})
}

func TestDiff_should_detect_diff_in_response_nested_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/models_with_recursive_refs",
		want: []Report{{
			Err:      errors.New("response field id mustn't be deleted"),
			JSONPath: "$./pet.post.responses",
		}},
	})
}

func TestDiff_should_detect_diff_in_response_models_with_different_refs(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/different_models_with_different_names",
		want: []Report{{
			Err:      errors.New("response field id mustn't be deleted"),
			JSONPath: "$./pet.post.responses",
		}, {
			Err:      errors.New("response field name mustn't change type from string to integer"),
			JSONPath: "$./pet.post.responses",
		}},
	})
}

func TestDiff_should_detect_diff_in_response_models_with_different_refs_array(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/different_models_with_different_names_array",
		want: []Report{{
			Err:      errors.New("response field id mustn't be deleted"),
			JSONPath: "$./pet.post.responses"},
			{
				Err:      errors.New("response field name mustn't change type from string to integer"),
				JSONPath: "$./pet.post.responses",
			}},
	})
}

func TestDiff_should_handle_case_with_schema_and_array_def(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/schema_with_array",
		want:      []Report{},
	})
}

func TestDiff_should_handle_self_reference_response(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/models_with_self_ref",
		want:      []Report{},
	})
}

func TestDiff_should_handle_case_with_recursive_refs_between_chain_of_models_response(t *testing.T) {
	runTest(t, test{
		specsPath: "/response/definitions/two_models_with_recursive_refs",
		want:      []Report{},
	})
}

// END RESPONSE

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
		return got[i].Err.Error() < got[j].Err.Error()
	})
	sort.Slice(tt.want, func(i, j int) bool {
		return tt.want[i].Err.Error() < tt.want[j].Err.Error()
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
