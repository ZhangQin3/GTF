package main

import (
	test_verify_arrs_web `gtf/scripts/test_verify_arrs_web`
)

import (
	"gtf/drivers/gtf"
)

func main() {
	gtf.TestSuiteSchema.TestScripts = map[string]interface{}{
		`test_verify_arrs_web`: new(test_verify_arrs_web.Test),
	}

	gtf.TestSuiteSchema.Repetitions = map[string]int{
		`test_verify_arrs_web`: 1,
	}

	gtf.GtfMain()
}
