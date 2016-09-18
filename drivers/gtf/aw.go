package gtf

import (
	"gtf/library/github"
	"reflect"
)

func runTestcase(heading []string, records [][]string) {
	var p interface{}
	if heading[2] == "github" {
		p = github.Github{}
	}

	for _, record := range records {
		p_typ := reflect.TypeOf(p)

	}
}
