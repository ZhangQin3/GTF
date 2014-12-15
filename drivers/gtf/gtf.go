package gtf

import (
	"fmt"
	"gtf/drivers/common"
	"gtf/log"
	tsuite "gtf/testsuites/tsuite"
	"reflect"
	"strconv"
	"time"
)

type testSuiteSchema struct {
	TestScripts map[string]interface{}
	Repetitions map[string]int
}

/* Global variable(s) exported. */
var (
	TestSuiteSchema testSuiteSchema
	TestParams      = make(map[string]interface{}) /* The map uased to lay params inherited from testsuite and set from the testcase.. */
)

/* Contains the data for each test script */
type testScript struct {
	fileName string         /* test script file name without suffix(.go). */
	tTest    *reflect.Value /* the the pointer to the instance of the Test struct in the test script */
	logger   *log.Logger    /* logger for each test script.  */
}

/* Global variable(s) NOT exported. */
var (
	currentTestScript *testScript
	tcDefinitions     = make(map[string]*tcDefinition) /* The testcase defined in the method CaseDefinitions in the test script, the key is string tcid.. */
)

func GtfMain() {
	var ts = initTestSuite()
	runTestScripts(ts)
	ts.SuiteTeardown()
}

func initTestSuite() *tsuite.TSuite {
	var ts = new(tsuite.TSuite)
	ts.SetSuiteParams()
	ts.SuitSetup()
	return ts
}

func initTestScript(scriptFileName string, tTest interface{}, ts *tsuite.TSuite) {
	currentTestScript = newTestScript(scriptFileName, tTest)
	logTestScriptHeader(currentTestScript)

	/* Initialize test execution params from the testsuite Params. */
	TestParams = ts.SuiteParams
	ts.CaseSetup()
}

func cleanupTestScript(ts *tsuite.TSuite) {
	/* Call test script level Cleanup method. */
	tcpCleanup := currentTestScript.tTest.MethodByName("TestCaseProcedureCleanup")
	if tcpCleanup.Kind() == reflect.Func {
		tcpCleanup.Call(nil)
	}

	ts.CaseTeardown()
}

func runTestScripts(ts *tsuite.TSuite) {
	/*  scriptFileName, tTest := "test_verify_test", new(test_verify_test.Test) */
	for scriptFileName, tTest := range TestSuiteSchema.TestScripts {
		initTestScript(scriptFileName, tTest, ts)
		if err := runTestCases(scriptFileName); err != nil {
			continue /*Jump out to execute the next script. */
		}
		cleanupTestScript(ts)
	}
}

func runTestCases(scriptFileName string) (err error) {
	currentTestScript.tTest.MethodByName("SetTestParams").Call(nil)
	if csDef := currentTestScript.tTest.MethodByName("CaseDefinitions"); csDef.IsValid() {
		/* The global variable tcDefs will be filled here. */
		csDef.Call(nil)
	} else {
		/* None testcase is defined. Log a message in the log file, and stop execute the testscript. */
		log.Error("[ERROR] No testcase defined in the script.")
		return fmt.Errorf("Jump out to execute the next script.")
	}

	/* Execute TestCaseProcedure, in the method TestCaseProcedure the function ExecuteTestCase
	   will be called to execute each test procedure for each testcase via executing Test.ExecuteTestCase method. */
	for i := 0; i < TestSuiteSchema.Repetitions[scriptFileName]; i++ {
		tp := currentTestScript.tTest.MethodByName("TestCaseProcedure")
		tp.Call(nil)
	}
	return nil
}

func newTestScript(scriptFileName string, tTest interface{}) *testScript {
	logger := initTestScriptLogger(scriptFileName)
	test := reflect.ValueOf(tTest)
	return &testScript{fileName: scriptFileName, tTest: &test, logger: logger}
}

func initTestScriptLogger(testFileName string) *log.Logger {
	logFile := testFileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(`..\src\gtf\drivers\log\tmpl\header.html`, logFile)

	return log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}
