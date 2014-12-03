package main

import (
	test_verify_second "gtf/scripts/test_verify_second"
	test_verify_test "gtf/scripts/test_verify_test"
)

import (
	"gtf"
)

func main() {
	gtf.Tss.TestScripts = map[string]interface{}{
		"test_verify_test":   new(test_verify_test.Test),
		"test_verify_second": new(test_verify_second.Test),
	}

	gtf.Tss.Repetitions = map[string]int{
		"test_verify_test":   1,
		"test_verify_second": 2,
	}

	gtf.GtfMain()
}
