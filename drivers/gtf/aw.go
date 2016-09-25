package gtf

import (
	"fmt"
	"gtf/library/github"
	"reflect"
)

func executeAwTestCase(heading []string, records [][]string) {
	var p interface{}
	if heading[2] == "github" {
		fmt.Println("------------", heading[2])
		p = &github.Github{}
	}

	for _, record := range records {
		p_typ := reflect.TypeOf(p)

		fmt.Println(p_typ.MethodByName("OpenURL"))

		fmt.Println(p_typ.NumMethod())
		fmt.Println(record)
	}
}
