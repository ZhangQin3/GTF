package main

import (
	test_verify_github "gtf/scripts/test_verify_github"
)

import (
	"gtf/drivers/gtf"
)

func main() {
	gtf.TestSuiteSchema.TestScripts = map[string]interface{}{
		`test_verify_github`: new(test_verify_github.Test),
	}

	gtf.TestSuiteSchema.Repetitions = map[string]int{
		`test_verify_github`: 1,
	}

	gtf.GtfMain()
}
