package gtf

import (
	"bytes"
	"errors"
	"fmt"
	"gtf/drivers/log"
	"io/ioutil"
	// "reflect"
	"regexp"
	"runtime"
	"time"
)

/* Contains the data for each testcase in a specfic test script. */
type awTestcase struct {
	tcid          string
	description   string
	testScript    *awScript
	cleaupMethod  string /* method name , if any, called if the testcase method ends normally, to clean up the test environment. */
	onCrashMethod string /* method name , if any, called if the testcase method crashed. */
	startTime     time.Time
	endTime       time.Time
}

func newAwTestcase(tcid string) *awTestcase {
	descr := currentAWScript.tcDefs[tcid].Description()
	return &awTestcase{testScript: currentAWScript, tcid: tcid, description: descr}
}

func (tc *awTestcase) runAwTestcase(heading []string, records [][]string) (err error) {
	/* Catch exeptions in the test method body, if any, in the test method. */
	err = errors.New("Panic in runTcMethod.")
	var flagCleanupCalled bool = false
	var l = tc.testScript.logger
	defer func() {
		if err := recover(); err != nil {
			var b bytes.Buffer
			b.WriteString(fmt.Sprintf("%s\n%s", err, "===========================-------------------------"))

			var buf []byte = make([]byte, 3072)
			runtime.Stack(buf, true)
			if err := l.GetTemplate().ExecuteTemplate(&b, "PANIC", fmt.Sprintf("%s", buf)); err != nil {
				log.Error(err)
				panic(err)
			}
			l.LabelError()
			tc.logStackTrace(b.Bytes())

			/* Call testcase cleanup on crash methed if testcase method of cleanup method panics. */
			if !flagCleanupCalled {
				/* In case the same following two line are not executed after tc.tcMethod.Call(*tc.tcMParams). */
				tc.logHorizonLine()
				l.GenerateStep("PostTest", "PostTest")
			}
			tc.callOnCrashMethod()
		}
	}()
	/* Add PreTest in case error occurs before first step. */
	l.GenerateStep("PreTest", "PreTest")

	/* Call testcase method. */
	fmt.Println(heading)
	fmt.Println(records)

	/* Call testcase cleanup method if there is not panic in the procedure of testcase method
	   if there is panic in the testcase method the cleanup method will not be called.*/
	flagCleanupCalled = true
	tc.logHorizonLine()
	l.GenerateStep("PostTest", "PostTest")
	tc.callCleanupMethod()

	return nil
}

func (tc *awTestcase) callOnCrashMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			tc.logStackTrace(buf)
		}
	}()
	tc.callMethod("CleanupOnCrash")
}

func (tc *awTestcase) callCleanupMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			tc.logStackTrace(buf)
		}
	}()
	tc.callMethod("Cleanup")
}

func (tc *awTestcase) callMethod(method string) {
	// mc := tc.methodName + method
	// m := tc.testScript.tTest.MethodByName(mc)
	// if m.Kind() != reflect.Func {
	// 	log.Warningf("The %s method %s is NOT definded.", method, mc)
	// } else {
	// 	m.Call(nil)
	// }
}

// func getFunctionName(rv reflect.Value) (string, string) {
// 	qualifiedFuncName := runtime.FuncForPC(rv.Pointer()).Name()

// 	reg := regexp.MustCompile(`(\w+).*?\w+\.\(\*Test\)\.(\w+)-fm`)
// 	matchs := reg.FindStringSubmatch(qualifiedFuncName)
// 	if matchs == nil {
// 		panic("The qualified function name: " + qualifiedFuncName + ` does NOT match the regexp: (\w+).*?\w+\.\(\*Test\)\.(\w+)-fm.`)
// 	}
// 	return matchs[1], matchs[2]
// }

func (tc *awTestcase) logResult() {
	var faildSteps string
	var l = tc.testScript.logger
	defer func() {
		l.CleanupSteps()
	}()

	for _, step := range l.Steps {
		if step.IsFailed {
			faildSteps = faildSteps + "{" + step.Index + "} "
		}
	}
	if faildSteps != "" {
		tc.logFailedSteps(faildSteps)
	}
	summary := tc.testResultSummary(faildSteps)

	/* TODO: enhance it if possible. */
	l.CloseFile()
	content, err := ioutil.ReadFile(l.FileName())
	if err != nil {
		panic(err)
	}
	regexpTcSummary := regexp.MustCompile(`<div style="display:none">hide</div>`)
	content = regexpTcSummary.ReplaceAll(content, summary.Bytes())
	ioutil.WriteFile(l.FileName(), content, 0666)
	l.ReopenFile()
}

func (tc *awTestcase) logHeader() {
	start := time.Now()
	tc.startTime = start
	tc.testScript.logger.Output("TC_HEADER",
		log.LOnlyFile,
		log.TestcaseHdr{
			tc.tcid,
			start.Format("2006-01-02 15:04:05"),
			tc.description,
			start.UnixNano(),
		})
}

/* Here only input failedStps, if failedStps != "" indicates there is some error happened. */
func (tc *awTestcase) testResultSummary(failedSteps string) bytes.Buffer {
	end := time.Now()
	tc.endTime = end
	duration := end.Sub(tc.startTime)
	data := log.TestResultSummary{
		tc.tcid,
		tc.description,
		failedSteps == "",
		failedSteps,
		tc.startTime.UnixNano(),
		fmt.Sprintf("%.2f", duration.Minutes()),
	}
	var buf bytes.Buffer
	if err := tc.testScript.logger.GetTemplate().ExecuteTemplate(&buf, "RESULT_SUMMARY", data); err != nil {
		panic(err)
	}
	return buf
}

func (tc *awTestcase) logFailedSteps(failedSteps string) {
	tc.testScript.logger.Output("D_FAIL", log.LFileAndConsole, failedSteps)
}

func (tc *awTestcase) logStackTrace(buf []byte) {
	var l = tc.testScript.logger
	t := l.PanicTime()
	if t == 0 {
		l.Output("PANIC", log.LFileAndConsole, fmt.Sprintf("%s\n", buf))
	} else {
		l.ZeroPanicTime()
		/* TODO: enhance it if possible. */
		l.CloseFile()
		content, err := ioutil.ReadFile(l.FileName())
		if err != nil {
			panic(err)
		}
		regexpTcSummary := regexp.MustCompile(fmt.Sprintf("panic_here_%d", t))
		content = regexpTcSummary.ReplaceAll(content, buf)
		ioutil.WriteFile(l.FileName(), content, 0666)
		l.ReopenFile()
	}
}

func (tc *awTestcase) logHorizonLine() {
	tc.testScript.logger.Output("HORIZON", log.LOnlyFile, nil)
}
