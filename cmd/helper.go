package cmd

import (
	"differ/pkg/differ"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

func readSpec(path string) map[string]interface{} {
	file, _ := ioutil.ReadFile(path)
	var spec map[string]interface{}
	err := json.Unmarshal(file, &spec)

	if err != nil {
		log.Fatalf("path %v not found", path)
	}

	return spec
}

func makeReport(errs []differ.Report) string {
	var report strings.Builder

	for jsonPath, errs := range toMap(errs) {
		report.WriteString("\n")
		report.WriteString(jsonPath)
		for _, err := range errs {
			report.WriteString("\n\t")
			report.WriteString(err.Error())
		}
	}

	return report.String()
}

func toMap(errs []differ.Report) map[string][]error {
	result := map[string][]error{}
	for _, value := range errs {
		_, ok := result[value.JSONPath]

		if !ok {
			result[value.JSONPath] = make([]error, 0)
		}

		result[value.JSONPath] = append(result[value.JSONPath], value.Err)
	}

	return result
}
