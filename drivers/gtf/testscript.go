package gtf

import (
	"gtf/drivers/common"
	"gtf/log"
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

func newTestScript(scriptFileName string, tTest interface{}) *testScript {
	test := reflect.ValueOf(tTest)
	testscript := &testScript{fileName: scriptFileName, tTest: &test}
	testscript.initLogger()
	return testscript
}

func (s *testScript) initLogger() {
	logFile := s.fileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(`..\src\gtf\drivers\log\tmpl\header.html`, logFile)

	s.logger = log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}

/* Log a test script information in the report file. */
func (s *testScript) logHeader() {
	s.logger.Output("TS_HEADING",
		log.LOnlyFile,
		log.TestScriptHdrInfo{
			time.Now().String(),
			s.fileName,
		})
}
