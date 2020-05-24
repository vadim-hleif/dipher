package pkg

import (
	"encoding/json"
	"fmt"
)

// Use Diff to compare the same swagger specifications.
// Returns empty slice because there are no breaking changes
func ExampleDiff_sameSpecs() {
	// swagger specification
	specV1Json := []byte(`
{
  "swagger": "2.0",
  "paths": {
    "/pet": {
      "post": {
        "parameters": [
          {
            "name": "create-request",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "id",
                "name"
              ],
              "properties": {
                "id": {
                  "type": "integer"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "alex",
                    "john",
                    "tom"
                  ]
                },
                "age": {
                  "type": "integer"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`)

	// parse spec
	var specV1 map[string]interface{}
	_ = json.Unmarshal(specV1Json, &specV1)

	fmt.Println(Diff(specV1, specV1))
	// Output: []
}

// Use Diff to compare the different swagger specifications.
// Returns slice with errors because there are some breaking changes
func ExampleDiff_specsWithBreakingChanges() {
	// first swagger specification
	specV1Json := `
{
  "swagger": "2.0",
  "paths": {
    "/pet": {
      "post": {
        "parameters": [
          {
            "name": "create-request",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "id",
                "name"
              ],
              "properties": {
                "id": {
                  "type": "integer"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "alex",
                    "john",
                    "tom"
                  ]
                },
                "age": {
                  "type": "integer"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`

	// parse spec
	var specV1 map[string]interface{}
	_ = json.Unmarshal([]byte(specV1Json), &specV1)

	// second swagger specification with deleted id param in request body
	specV2Json := `
{
  "swagger": "2.0",
  "paths": {
    "/pet": {
      "post": {
        "parameters": [
          {
            "name": "create-request",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "name"
              ],
              "properties": {
                "name": {
                  "type": "string",
                  "enum": [
                    "alex",
                    "john",
                    "tom"
                  ]
                },
                "age": {
                  "type": "integer"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`

	// parse spec
	var specV2 map[string]interface{}
	_ = json.Unmarshal([]byte(specV2Json), &specV2)

	// print reports
	reports := Diff(specV1, specV2)
	for _, report := range reports {
		fmt.Println(report.JSONPath)
		fmt.Println("\t", report.Err)
	}

	// Output:
	// $./pet.post.parameters
	// 	 required param id mustn't be deleted
}
