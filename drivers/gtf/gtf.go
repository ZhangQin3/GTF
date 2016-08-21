package gtf

import (
	tsuite "gtf/testsuites/tsuite"
	"strings"
)

type testSuiteSchema struct {
	TestScripts map[string]interface{}
	Repetitions map[string]int
}

var (
	TestSuiteSchema testSuiteSchema
	TestParams      = make(map[string]interface{}) /* TestParams uased to store params inherited from testsuite and set from the current test script. */
)

func GtfMain() {
	ts := suiteSetup()
	runTestScripts(ts)
	ts.SuiteCleanup()
}

func suiteSetup() *tsuite.TSuite {
	ts := new(tsuite.TSuite)
	ts.SetSuiteParams()
	ts.SuitSetup()
	return ts
}

func runTestScripts(ts *tsuite.TSuite) {
	/*  fileName, tTest := "test_verify_test", new(test_verify_test.Test) or
	    fileName, tTest := "test_verify_test", "csv") */
	for fileName, tTest := range TestSuiteSchema.TestScripts {
		str, ok := tTest.(string)

		if ok {
			// run actionword testscript
			if strings.EqualFold(str, "csv") {
				s := newAWScript(fileName, ts)
				s.setup()
				s.runTestCases()
				s.cleanup()
			} else {
				panic("Only support csv file as actionword script.")
			}
		} else {
			// run standard testscript
			s := newTestScript(fileName, tTest, ts)
			s.setup()
			s.runTestCases()
			s.cleanup()
		}
	}
}
