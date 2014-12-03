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
	var faildStps string

	defer clearTcSteps(logger)

	for _, v := range logger.Steps {
		if (*v).StpFailed {
			faildStps = faildStps + "{" + (*v).StpIdx + "} "
		}
	}
	if faildStps != "" {
		logFailAtTail(logger, faildStps)
	}
	buf := genTcDataToTbl(logger, tcid, tcDescription, faildStps)

	/* TODO: enhance it if possible. */
	log.CloseLogFile()
	fileContent, err := ioutil.ReadFile(logger.GetLogFileName())
	if err != nil {
		panic(err)
	}
	insert := regexp.MustCompile(`<div style="display:none">hide</div>`)
	fileContent = insert.ReplaceAll(fileContent, buf.Bytes())
	ioutil.WriteFile(logger.GetLogFileName(), fileContent, 0666)
	log.ReopenLogFile()
}

func logTcHearder(logger *log.Logger, tcid, tcDescr string) {
	logger.Output("TC_HEADING",
		log.LOnlyFile,
		log.TcHdrInfo{
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
		log.TsHdrInfo{
			time.Now().String(),
			pkgName,
		})
}

func logStack(logger *log.Logger, buf []byte) {
	logger.Output("PANIC", log.LFileConsole, fmt.Sprintf("%s\n", buf))
}

func logFailAtTail(logger *log.Logger, v ...interface{}) {
	logger.Output("D_FAIL", log.LFileConsole, fmt.Sprint(v...))
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
	if err := logger.GetLogTmpl().ExecuteTemplate(&buf, "REPORT_TBL", data); err != nil {
		panic(err)
	}

	return buf
}
