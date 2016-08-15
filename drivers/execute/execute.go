package main

import (
	test_verify_csvdata "gtf/scripts/test_verify_csvdata"
)

import (
	"gtf/drivers/gtf"
)

func main() {
	gtf.TestSuiteSchema.TestScripts = map[string]interface{}{
		`test_verify_csvdata`: new(test_verify_csvdata.Test),
	}

	gtf.TestSuiteSchema.Repetitions = map[string]int{
		`test_verify_csvdata`: 1,
	}

	gtf.GtfMain()
}
