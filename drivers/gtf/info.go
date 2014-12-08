// These functions wraps gtf/log to used only in the gtf package.
package gtf

import (
	"bytes"
	"fmt"
	"gtf/log"
	"io/ioutil"
	"regexp"
	"time"
)

func clearTcSteps(l *log.Logger) {
	l.Steps = l.Steps[0:0]
}

/* Log a test script information in the report file. */
func logTsHeader(logger *log.Logger, pkgName string) {
	logger.Output("TS_HEADING",
		log.LOnlyFile,
		log.TsHeaderInfo{
			time.Now().String(),
			pkgName,
		})
}

func logTcHeader(logger *log.Logger, tcid, tcDescr string) {
	logger.Output("TC_HEADING",
		log.LOnlyFile,
		log.TcHeaderInfo{
			tcid,
			time.Now().Format("2006-01-02 15:04:05"),
			tcDescr,
		})
}

func logTcResult(logger *log.Logger, tcid, tcDescription string) {
	var FaildSteps string
	defer clearTcSteps(logger)

	for _, step := range logger.Steps {
		if step.IsFailed {
			FaildSteps = FaildSteps + "{" + step.StepIndex + "} "
		}
	}
	if FaildSteps != "" {
		logFailAtTail(logger, FaildSteps)
	}
	tcSummaryResult := generateTcResultSummary(logger, tcid, tcDescription, FaildSteps)

	/* TODO: enhance it if possible. */
	logger.CloseFile()
	logFileContent, err := ioutil.ReadFile(logger.GetFileName())
	if err != nil {
		panic(err)
	}
	regexpInsertTcSummary := regexp.MustCompile(`<div style="display:none">hide</div>`)
	logFileContent = regexpInsertTcSummary.ReplaceAll(logFileContent, tcSummaryResult.Bytes())
	ioutil.WriteFile(logger.GetFileName(), logFileContent, 0666)
	logger.ReopenFile()
}

func logHorizon(logger *log.Logger) {
	logger.Output("HORIZON", log.LOnlyFile, nil)
}

func logStack(logger *log.Logger, buf []byte) {
	logger.Output("PANIC", log.LFileAndConsole, fmt.Sprintf("%s\n", buf))
}

func logFailAtTail(logger *log.Logger, v ...interface{}) {
	logger.Output("D_FAIL", log.LFileAndConsole, fmt.Sprint(v...))
}

/* Here only input failedStps, if failedStps != "" indicates there is some error happened. */
func generateTcResultSummary(logger *log.Logger, tcid, tcDescr, failedStps string) bytes.Buffer {
	data := log.TcResultSummary{
		tcid,
		tcDescr,
		failedStps == "",
		failedStps,
	}
	var buf bytes.Buffer
	if err := logger.GetTemplate().ExecuteTemplate(&buf, "RESULT_SUMMARY", data); err != nil {
		panic(err)
	}
	return buf
}
