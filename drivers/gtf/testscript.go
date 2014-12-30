package gtf

import (
	"gtf/drivers/common"
	"gtf/drivers/log"
	"reflect"
	"strconv"
	"time"
)

/* Contains the data for each test script */
type testScript struct {
	fileName string         /* test script file name without suffix(.go). */
	tTest    *reflect.Value /* the the pointer to the instance of the Test struct in the test script */
	logger   *log.Logger    /* logger for each test script.  */
}

func newTestScript(fileName string, tTest interface{}) *testScript {
	t := reflect.ValueOf(tTest)
	s := &testScript{fileName: fileName, tTest: &t}
	s.initLogger()
	return s
}

func (s *testScript) initLogger() {
	logFile := s.fileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(`..\src\gtf\drivers\log\tmpl\header.html`, logFile)

	s.logger = log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}

/* Log a test script information in the report file. */
func (s *testScript) logHeader() {
	s.logger.Output("SCRIPT_HEADING",
		log.LOnlyFile,
		log.TestScriptHdr{
			time.Now().String(),
			s.fileName,
		})
}
func (s *testScript) logTailer() {
	s.logger.Output("SCRIPT_TAIL", log.LOnlyFile, nil)
}

/* fieldName is the field name of Test in test.go, tTest promotes them. */
func (s *testScript) tTestField(fieldName string) reflect.Value {
	return s.tTest.Elem().FieldByName(fieldName)
}

func (s *testScript) tcDefField(tcid, fieldName string) string {
	return s.tTest.Elem().FieldByName("tcDefs").MapIndex(reflect.ValueOf(tcid)).Elem().FieldByName(fieldName).String()
}
