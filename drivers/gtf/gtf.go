package gtf

import (
	tsuite "gtf/testsuites/tsuite"
)

type testSuiteSchema struct {
	TestScripts map[string]interface{}
	Repetitions map[string]int
}

var (
	TestSuiteSchema testSuiteSchema
	TestParams      = make(map[string]interface{}) /* TestParams uased to store params inherited from testsuite and set from the current test script. */
	currentScript   *testScript
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

// func initTestScript(scriptFileName string, tTest interface{}, ts *tsuite.TSuite) {
// 	currentScript = newTestScript(scriptFileName, tTest)
// 	currentScript.logHeader()

// 	/* Initialize TestParams from testsuite SuiteParams. */
// 	TestParams = ts.SuiteParams
// 	ts.CaseSetup()
// }

// func cleanupTestScript(ts *tsuite.TSuite) {
// 	/* Call test script level Cleanup method. */
// 	tcpCleanup := currentScript.tTest.MethodByName("TestCaseProcedureCleanup")
// 	if tcpCleanup.Kind() == reflect.Func {
// 		tcpCleanup.Call(nil)
// 	}

// 	ts.CaseTeardown()
// 	currentScript.logTailer()
// }

func runTestScripts(ts *tsuite.TSuite) {
	/*  fileName, tTest := "test_verify_test", new(test_verify_test.Test) */
	for fileName, tTest := range TestSuiteSchema.TestScripts {
		s = newTestScript(fileName, tTest)
		s.setup(ts)
		if err := s.runTestCases(); err != nil {
			continue /*Jump out to execute the next script. */
		}
		s.cleanup(ts)
	}
}

// func runTestCases(scriptFileName string) (err error) {
// 	currentScript.tTest.MethodByName("SetTestParams").Call(nil)
// 	if tcDef := currentScript.tTest.MethodByName("CaseDefinitions"); tcDef.IsValid() {
// 		/* The global variable tcDefs will be filled here. */
// 		tcDef.Call(nil)
// 	} else {
// 		/* None testcase is defined. Log a message in the log file, and stop execute the testscript. */
// 		log.Error("[ERROR] No testcase defined in the script.")
// 		return fmt.Errorf("Jump out to execute the next script.")
// 	}

// 	/* Execute TestCaseProcedure, in the method TestCaseProcedure the function ExecuteTestCase
// 	   will be called to execute each test procedure for each testcase via executing Test.ExecuteTestCase method. */
// 	for i := 0; i < TestSuiteSchema.Repetitions[scriptFileName]; i++ {
// 		tp := currentScript.tTest.MethodByName("TestCaseProcedure")
// 		tp.Call(nil)
// 	}
// 	return nil
// }
