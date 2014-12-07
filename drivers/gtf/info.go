package gtf

import (
	"bytes"
	"fmt"
	"gtf/log"
	"io/ioutil"
	"regexp"
	"time"
)

/* These functions wraps gtf/log to used only in the gtf package. */

func clearTcSteps(l *log.Logger) {
	l.Steps = l.Steps[0:0]
}

func logTCResult(logger *log.Logger, tcid, tcDescription string) {
	var FaildSteps string

	defer clearTcSteps(logger)

	for _, v := range logger.Steps {
		if (*v).IsFailed {
			FaildSteps = FaildSteps + "{" + (*v).StepIndex + "} "
		}
	}
	if FaildSteps != "" {
		logFailAtTail(logger, FaildSteps)
	}
	buf := genTcDataToTbl(logger, tcid, tcDescription, FaildSteps)

	/* TODO: enhance it if possible. */
	log.CloseFile()
	fileContent, err := ioutil.ReadFile(logger.GetFileName())
	if err != nil {
		panic(err)
	}
	insert := regexp.MustCompile(`<div style="display:none">hide</div>`)
	fileContent = insert.ReplaceAll(fileContent, buf.Bytes())
	ioutil.WriteFile(logger.GetFileName(), fileContent, 0666)
	log.ReopenFile()
}

func logTcHearder(logger *log.Logger, tcid, tcDescr string) {
	logger.Output("TC_HEADING",
		log.LOnlyFile,
		log.TcHeaderInfo{
			tcid,
			time.Now().Format("2006-01-02 15:04:05"),
			tcDescr,
		})
}

func logHorizon(logger *log.Logger) {
	logger.Output("HORIZON", log.LOnlyFile, nil)
}

/* Log a test script information in the report file. */
func logTsHearder(logger *log.Logger, pkgName string) {
	logger.Output("TS_HEADING",
		log.LOnlyFile,
		log.TsHeaderInfo{
			time.Now().String(),
			pkgName,
		})
}

func logStack(logger *log.Logger, buf []byte) {
	logger.Output("PANIC", log.LFileAndConsole, fmt.Sprintf("%s\n", buf))
}

func logFailAtTail(logger *log.Logger, v ...interface{}) {
	logger.Output("D_FAIL", log.LFileAndConsole, fmt.Sprint(v...))
}

/* Here only input failedStps, if failedStps != "" indicates there is some error happened. */
func genTcDataToTbl(logger *log.Logger, tcid, tcDescr, failedStps string) bytes.Buffer {
	data := log.TcResultToTbl{
		tcid,
		tcDescr,
		failedStps == "",
		failedStps,
	}
	var buf bytes.Buffer
	if err := logger.GetTemplate().ExecuteTemplate(&buf, "REPORT_TBL", data); err != nil {
		panic(err)
	}

	return buf
}
