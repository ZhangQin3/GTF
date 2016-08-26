package gtf

import (
	"fmt"
	"gtf/drivers/common"
	"gtf/drivers/csv"
	"gtf/drivers/log"
	tsuite "gtf/testsuites/tsuite"
	// "reflect"
	"strconv"
	"time"
)

var currentAWScript *actionwordScript

/* Contains the data for each test script */
type actionwordScript struct {
	fileName  string /* test script file name without suffix(.go). */
	data      *csv.Data
	tSuite    *tsuite.TSuite
	logger    *log.Logger /* logger for each test script.  */
	startTime time.Time
	// endTime   time.Time
}

func newAWScript(fileName string, ts *tsuite.TSuite) *actionwordScript {
	d := csv.NewData(common.AWFilesDir + fileName)
	s := &actionwordScript{fileName: fileName, data: d}
	s.initLogger()
	s.tSuite = ts
	currentAWScript = s
	return s
}

func (s *actionwordScript) initLogger() {
	logFile := s.fileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(logFile, `..\src\gtf\drivers\log\tmpl\header.html`)

	s.logger = log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}

/* Log a test script information in the report file. */
func (s *actionwordScript) logHeader() {
	s.startTime = time.Now()
	data := log.TestScriptHdr{
		s.startTime.String(),
		s.fileName,
	}
	s.logger.Output("LOG_HEADER", log.LOnlyFile, data)
}
func (s *actionwordScript) logTailer() {
	end := time.Now()
	data := log.TestScriptTlr{
		end.String(),
		fmt.Sprintf("%.2f", end.Sub(s.startTime).Minutes()),
	}
	s.logger.Output("LOG_TAILER", log.LOnlyFile, data)
}

// /* fieldName is the field name of Test in test.go, tTest promotes them. */
// func (s *actionwordScript) tTestField(fieldName string) reflect.Value {
// 	return s.tTest.Elem().FieldByName(fieldName)
// }

// func (s *actionwordScript) tcDefField(tcid, fieldName string) string {
// 	return s.tTest.Elem().FieldByName("tcDefs").MapIndex(reflect.ValueOf(tcid)).Elem().FieldByName(fieldName).String()
// }

func (s *actionwordScript) setup() {
	s.logHeader()
	/* Initialize TestParams from testsuite SuiteParams. */
	TestParams = s.tSuite.SuiteParams
	s.tSuite.CaseSetup()
}

func (s *actionwordScript) cleanup() {
	s.tSuite.CaseCleanup()
	s.logTailer()
}

func (s *actionwordScript) runTestCases() {
	/* Execute TestCaseProcedure, in the method TestCaseProcedure the function ExecuteTestCase
	   will be called to execute each test procedure for each testcase via executing Test.ExecuteTestCase method. */
	for i := 0; i < TestSuiteSchema.Repetitions[s.fileName]; i++ {
		log.Warning(s.data.ReadRecord())
	}
}
