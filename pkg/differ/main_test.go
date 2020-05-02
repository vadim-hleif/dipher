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
	name      string
	specsPath string
	want      []error
}

var cases = []test{
	// params-object test cases
	{
		name:      "should detect removing enum value in any request object property",
		specsPath: "params-object/removed_value_from_old_enum",
		want: []error{
			errors.New("param name mustn't remove value alex from enum"),
		},
	},
	{
		name:      "should detect adding enum in any request object property",
		specsPath: "params-object/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param age mustn't have enum"),
		},
	},
	{
		name:      "should detect type changing in the any request object property",
		specsPath: "params-object/type_changing",
		want: []error{
			errors.New("param id mustn't change type from integer to string"),
		},
	},
	{
		name:      "should detect required property removing in the request object",
		specsPath: "params-object/required_property_deletion",
		want: []error{
			errors.New("required param id mustn't be deleted"),
		},
	},
	{
		name:      "should detect new required property in the request object",
		specsPath: "params-object/required_property_adding",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	},
	// end params-object test cases

	// primitive params test cases
	{
		name:      "should detect removing enum value in any request param",
		specsPath: "params/removed_value_from_old_enum",
		want: []error{
			errors.New("param sort mustn't remove value asc from enum"),
		},
	},
	{
		name:      "should detect adding enum in any request param",
		specsPath: "params/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param sort mustn't have enum"),
		},
	},
	{
		name:      "should detect type changing in the any request param",
		specsPath: "params/type_changing",
		want: []error{
			errors.New("param sort mustn't change type from string to integer"),
		},
	},
	{
		name:      "should detect required property removing in the any request param",
		specsPath: "params/required_param_deletion",
		want: []error{
			errors.New("required param sort mustn't be deleted"),
		},
	},
	{
		name:      "should detect new required param in the request",
		specsPath: "params/new_required_param",
		want: []error{
			errors.New("new required param filter mustn't be added"),
		},
	},
	{
		name:      "should detect marking old param as required",
		specsPath: "params/mark_old_param_as_required",
		want: []error{
			errors.New("param sort mustn't be required because it wasn't be required"),
		},
	},
	// end primitive params test cases

	// objects with definitions test cases
	{
		name:      "should detect new required property in the request object definition",
		specsPath: "definitions/required_property_adding",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	},
	{
		name:      "should detect required property removing in the request object definition",
		specsPath: "definitions/required_property_deletion",
		want: []error{
			errors.New("required param name mustn't be deleted"),
		},
	},
	{
		name:      "should detect type changing in the any request object property definition",
		specsPath: "definitions/type_changing",
		want: []error{
			errors.New("param id mustn't change type from integer to string"),
		},
	},
	{
		name:      "should detect adding enum in any request object property definition",
		specsPath: "definitions/adding_enum_to_old_value_without_enum",
		want: []error{
			errors.New("param name mustn't have enum"),
		},
	},
	{
		name:      "should detect removing enum value in any request object property definition",
		specsPath: "definitions/removed_value_from_old_enum",
		want: []error{
			errors.New("param name mustn't remove value alex from enum"),
		},
	},
	{
		name:      "should detect object definition in any path",
		specsPath: "definitions/custom_def_path",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
		},
	},
	{
		name:      "should compare models by different refs",
		specsPath: "definitions/same_models_with_different_names",
		want:      make([]error, 0),
	},
	{
		name:      "should detect diff in models with different refs",
		specsPath: "definitions/different_models_with_different_names",
		want: []error{
			errors.New("param age mustn't be required because it wasn't be required"),
			errors.New("param age mustn't have enum"),
			errors.New("param name mustn't change type from string to integer"),
			errors.New("param name mustn't remove value alex from enum"),
			errors.New("required param id mustn't be deleted"),
		},
	},
	// end objects with definitions test cases

	// resources and verbs test cases
	{
		name:      "should detect path removing",
		specsPath: "paths_and_verbs/resource_removing",
		want: []error{
			errors.New("resource /pet mustn't be removed"),
		},
	},
	{
		name:      "should detect method removing in the any path",
		specsPath: "paths_and_verbs/verb_removing",
		want: []error{
			errors.New("post method of /pet path mustn't be removed"),
		},
	},
	// end resources and verbs test cases
}

func TestDiff(t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
