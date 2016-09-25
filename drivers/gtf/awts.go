package gtf

import (
	"fmt"
	"gtf/drivers/common"
	"gtf/drivers/csv"
	"gtf/drivers/log"
	tsuite "gtf/testsuites/tsuite"
	"io"
	"strconv"
	"time"
)

var currentAWScript *awScript

/* Contains the data for each test script */
type awScript struct {
	fileName     string /* test script file name without suffix(.go). */
	acctionwords *csv.Actionwords
	tSuite       *tsuite.TSuite
	tcDefs       map[string]*tcDef
	logger       *log.Logger /* logger for each test script.  */
	startTime    time.Time
	// endTime   time.Time
}

func newAWScript(fileName string, ts *tsuite.TSuite) *awScript {
	aws := csv.NewActionwords(common.AWFilesDir + fileName)
	s := &awScript{fileName: fileName, acctionwords: aws}
	s.initLogger()
	s.tSuite = ts
	currentAWScript = s
	return s
}

func (s *awScript) initLogger() {
	logFile := s.fileName + "." + strconv.FormatInt(time.Now().Unix(), 10) + ".html"
	common.CopyFile(logFile, `..\src\gtf\drivers\log\tmpl\header.html`)

	s.logger = log.NewLogger(logFile, `..\src\gtf\drivers\log\tmpl\t1.tmpl`)
}

/* Log a test script information in the report file. */
func (s *awScript) logHeader() {
	s.startTime = time.Now()
	data := log.TestScriptHdr{
		s.startTime.String(),
		s.fileName,
	}
	s.logger.Output("LOG_HEADER", log.LOnlyFile, data)
}
func (s *awScript) logTailer() {
	end := time.Now()
	data := log.TestScriptTlr{
		end.String(),
		fmt.Sprintf("%.2f", end.Sub(s.startTime).Minutes()),
	}
	s.logger.Output("LOG_TAILER", log.LOnlyFile, data)
}

func (s *awScript) setup() {
	s.logHeader()
	/* Initialize TestParams from testsuite SuiteParams. */
	TestParams = s.tSuite.SuiteParams
	s.tSuite.CaseSetup()
}

func (s *awScript) cleanup() {
	s.tSuite.CaseCleanup()
	s.logTailer()
}

func (s *awScript) runTestCases() {
	for _, b := range s.acctionwords.Blocks() {
		s.DefineTestCase(b)
	}

	for i := 0; i < TestSuiteSchema.Repetitions[s.fileName]; i++ {
		// for {
		// 	h, rs, err := s.acctionwords.Read()
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	log.Info(h)
		// 	log.Warning(rs)
		// 	log.Info(err)
		// }
		s.ExecuteAwTestCase()
	}
}

func (s *awScript) DefineTestCase(b *csv.Block) *tcDef {
	if s.tcDefs == nil {
		s.tcDefs = make(map[string]*tcDef)
	}

	tcid := b.GetTCID()
	d := &tcDef{tcid: tcid, description: b.GetDescription()}
	s.tcDefs[tcid] = d
	return d
}

func (s *awScript) ExecuteAwTestCase() {
	h, rs, err := s.acctionwords.Read()
	if err == io.EOF {
		return
	}

	tcid := h[0]
	if tcDef, ok := s.tcDefs[tcid]; ok {
		if !tcDef.CalculateAppliability() {
			return
		}
	} else {
		panic("The testcase: " + tcid + " is not defined.")
	}

	tc := newAwTestcase(tcid)
	defer func() {
		tc.logResult()
	}()

	tc.logHeader()
	tc.runAwTestcase(h, rs)

	// recur current method to execute all the test cases in the csv file
	s.ExecuteAwTestCase()
}
