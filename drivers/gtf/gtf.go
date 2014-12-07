package gtf

import (
	"fmt"
	"gtf/drivers/common"
	"gtf/log"
	tsuite "gtf/testsuites/tsuite"
	"os"
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
)

/* Contains the data for each test script */
type testScript struct {
	scriptName  string         /* test script name without suffix. */
	tTest       *reflect.Value /* the the pointer to the instance of the Test struct in the test script */
	logger      *log.Logger    /* logger for each test script.  */
	logFileName string         /* log file name. */
}

/* Global variable(s) NOT exported. */
var (
	tstScript  *testScript
	tcDefs     = make(map[string]*tcDefinition) /* The testcase defined in the method CaseDefinitions in the test script, the key is string tcid.. */
	testParams = make(map[string]interface{})   /* The map uased to lay params inherited from testsuite and set from the testcase.. */
)

func GtfMain() {
	var ts = initTestSuite()

	/*  scriptFileName, tTest := "test_verify_test", new(test_verify_test.Test) */
	for scriptFileName, tTest := range TestSuiteSchema.TestScripts {
		initTestScript(scriptFileName, tTest, ts)
		if err := runTests(scriptFileName); err != nil {
			continue /*Jump out to execute the next script. */
		}
		cleanupTestScript(ts)
	}

	/* Call suite SuiteTeardown function. */
	ts.SuiteTeardown()
}

func initTestSuite() *tsuite.TSuite {
	var ts = new(tsuite.TSuite)
	ts.SetSuiteParams()
	ts.SuitSetup()
	return ts
}

func initTestScript(scriptFileName string, tTest interface{}, ts *tsuite.TSuite) {
	tstScript = nil
	logFileName, logger := initTestLogFile(scriptFileName)
	Test := reflect.ValueOf(tTest)
	tstScript = &testScript{scriptName: scriptFileName, tTest: &Test, logger: logger, logFileName: logFileName}

	/* Initialize test execution params from the testsuite Params. */
	testParams = ts.SuiteParams
	ts.CaseSetup()

	logTsHearder(tstScript.logger, scriptFileName)
}

func runTests(scriptFileName string) (err error) {
	tstScript.tTest.MethodByName("SetTestParams").Call(nil)
	if csDef := tstScript.tTest.MethodByName("CaseDefinitions"); csDef.IsValid() {
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
		tp := tstScript.tTest.MethodByName("TestCaseProcedure")
		tp.Call(nil)
	}
	return nil
}

func cleanupTestScript(ts *tsuite.TSuite) {
	/* Call test script level Cleanup method. */
	tcpCleanup := tstScript.tTest.MethodByName("TestCaseProcedureCleanup")
	if tcpCleanup.Kind() == reflect.Func {
		tcpCleanup.Call(nil)
	}

	ts.CaseTeardown()
}

func initTestLogFile(testFileName string) (string, *log.Logger) {
	logFileName := testFileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	copy := "copy"
	if os.Getenv("OSTYPE") == "linux" {
		copy = `cp`
		common.ExecOSCmd(`%s ../src/gtf/drivers/log/tmpl/header.html %s`, copy, logFileName)
	} else {
		common.ExecOSCmd(`%s ..\src\gtf\drivers\log\tmpl\header.html %s`, copy, logFileName)
	}
	logger := log.NewLogger(logFileName, `../src/gtf/drivers/log/tmpl/t1.tmpl`)

	return logFileName, logger
}
