// These functions wraps gtf/log to used only in the gtf package.
package gtf

import (
	"bytes"
	"fmt"
	"gtf/log"
)

func logHorizonLine(logger *log.Logger) {
	logger.Output("HORIZON", log.LOnlyFile, nil)
}

func logStackTrace(logger *log.Logger, buf []byte) {
	logger.Output("PANIC", log.LFileAndConsole, fmt.Sprintf("%s\n", buf))
}

func logFailedSteps(logger *log.Logger, v ...interface{}) {
	logger.Output("D_FAIL", log.LFileAndConsole, fmt.Sprint(v...))
}

/* Here only input failedStps, if failedStps != "" indicates there is some error happened. */
func generateTcResultSummary(logger *log.Logger, tcid, tcDescription, failedSteps string) bytes.Buffer {
	data := log.TestcaseResultSummary{
		tcid,
		tcDescription,
		failedSteps == "",
		failedSteps,
	}
	var buf bytes.Buffer
	if err := logger.GetTemplate().ExecuteTemplate(&buf, "RESULT_SUMMARY", data); err != nil {
		panic(err)
	}
	return buf
}
