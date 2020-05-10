package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

func ReadSpec(path string) map[string]interface{} {
	file, _ := ioutil.ReadFile(path)
	var spec map[string]interface{}
	err := json.Unmarshal(file, &spec)

	if err != nil {
		log.Fatalf("path %v not found", path)
	}

	return spec
}

func MakeReport(errs []error) string {
	var report strings.Builder

	for _, difference := range errs {
		report.WriteString("\n")
		report.WriteString(difference.Error())
		report.WriteString("\n")
	}

	return report.String()
}
