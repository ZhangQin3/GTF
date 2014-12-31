package gtf

import tsuite "gtf/testsuites/tsuite"

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
	ts.SuiteTeardown()
}

func suiteSetup() *tsuite.TSuite {
	ts := new(tsuite.TSuite)
	ts.SetSuiteParams()
	ts.SuitSetup()
	return ts
}

func runTestScripts(ts *tsuite.TSuite) {
	/*  fileName, tTest := "test_verify_test", new(test_verify_test.Test) */
	for fileName, tTest := range TestSuiteSchema.TestScripts {
		s := newTestScript(fileName, tTest, ts)
		s.setup()
		s.runTestCases()
		s.cleanup()
	}
}
