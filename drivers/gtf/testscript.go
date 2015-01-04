package gtf

import (
	"gtf/drivers/common"
	"gtf/drivers/log"
	tsuite "gtf/testsuites/tsuite"
	"reflect"
	"strconv"
	"time"
)

var currentScript *testScript

/* Contains the data for each test script */
type testScript struct {
	fileName string         /* test script file name without suffix(.go). */
	tTest    *reflect.Value /* the the pointer to the instance of the Test struct in the test script */
	tSuite   *tsuite.TSuite
	logger   *log.Logger /* logger for each test script.  */
}

func newTestScript(fileName string, tTest interface{}, ts *tsuite.TSuite) *testScript {
	t := reflect.ValueOf(tTest)
	s := &testScript{fileName: fileName, tTest: &t}
	s.initLogger()
	s.tSuite = ts
	currentScript = s
	return s
}

func (s *testScript) initLogger() {
	logFile := s.fileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(`..\src\gtf\drivers\log\tmpl\header.html`, logFile)

	s.logger = log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}

/* Log a test script information in the report file. */
func (s *testScript) logHeader() {
	s.logger.Output("LOG_HEADER",
		log.LOnlyFile,
		log.TestScriptHdr{
			time.Now().String(),
			s.fileName,
		})
}
func (s *testScript) logTailer() {
	s.logger.Output("LOG_TAILER", log.LOnlyFile, nil)
}

/* fieldName is the field name of Test in test.go, tTest promotes them. */
func (s *testScript) tTestField(fieldName string) reflect.Value {
	return s.tTest.Elem().FieldByName(fieldName)
}

func (s *testScript) tcDefField(tcid, fieldName string) string {
	return s.tTest.Elem().FieldByName("tcDefs").MapIndex(reflect.ValueOf(tcid)).Elem().FieldByName(fieldName).String()
}

func (s *testScript) setup() {
	s.logHeader()

	/* Initialize TestParams from testsuite SuiteParams. */
	TestParams = s.tSuite.SuiteParams
	s.tSuite.CaseSetup()
}

func (s *testScript) cleanup() {
	/* Call test script level Cleanup method. */
	c := s.tTest.MethodByName("TestCaseProcedureCleanup")
	if c.Kind() == reflect.Func {
		c.Call(nil)
	}

	s.tSuite.CaseCleanup()
	s.logTailer()
}

func (s *testScript) runTestCases() {
	s.tTest.MethodByName("SetTestParams").Call(nil)
	if d := s.tTest.MethodByName("CaseDefinitions"); d.IsValid() {
		/* The global variable tcDefs will be filled here. */
		d.Call(nil)
	} else {
		/* None testcase is defined. Log a message in the log file, and stop execute the testscript. */
		log.Error("[ERROR] No testcase defined in the script.")
		return
	}

	/* Execute TestCaseProcedure, in the method TestCaseProcedure the function ExecuteTestCase
	   will be called to execute each test procedure for each testcase via executing Test.ExecuteTestCase method. */
	for i := 0; i < TestSuiteSchema.Repetitions[s.fileName]; i++ {
		p := s.tTest.MethodByName("TestCaseProcedure")
		p.Call(nil)
	}
}
