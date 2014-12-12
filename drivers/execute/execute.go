package main

import (
	test_verify_web "gtf/scripts/test_verify_web"
)

import (
	"gtf"
)

func main() {
	gtf.TestSuiteSchema.TestScripts = map[string]interface{}{
		"test_verify_web": new(test_verify_web.Test),
	}

	gtf.TestSuiteSchema.Repetitions = map[string]int{
		"test_verify_web": 1,
	}

	gtf.GtfMain()
}
