package gtf

import (
	"gtf/log"
	"io/ioutil"
	"reflect"
	"regexp"
	"runtime"
	"time"
)

/* Contains the data for each testcase in a specfic test script. */
type testCase struct {
	tcID          string
	description   string
	method        reflect.Value    /* reflect method of testcase method. */
	methodParams  *[]reflect.Value /* params of the testcase method. */
	methodName    string           /* method name of the testcase method. */
	cleaupMethod  string           /* method name , if any, called if the testcase method ends normally, to clean up the test environment. */
	onCrashMethod string           /* method name , if any, called if the testcase method crashed. */
	testScript    *testScript
}

// tcTestLogicMethod is the real test method with test logic
// tcid is the first parameter of the method tcTestLogicMethod
// params is other parameter(s), if any, of the method tcTestLogicMethod
func newTestCase(tcTestLogicMethod interface{}, tcid string, params *[]interface{}) *testCase {
	var tcMParams []reflect.Value
	var tc testCase
	tc.testScript = currentTestScript
	tp := reflect.ValueOf(tcTestLogicMethod)
	_, funcName := getFunctionName(tp)

	if tp.Kind() != reflect.Func {
		panic("The first param of the gtf.Execute must be a testcase method!")
	}
	tc.method = tp
	tc.methodName = funcName
	tc.tcID = tcid
	tc.description = tcDefinitions[tcid].description

	paramsLen := len(*params)
	tcMParams = make([]reflect.Value, paramsLen+1)
	tcMParams[0] = reflect.ValueOf(tcid)
	if paramsLen != 0 {
		for i := 0; i < paramsLen; i++ {
			tcMParams[i+1] = reflect.ValueOf((*params)[i])
		}
	}
	tc.methodParams = &tcMParams
	return &tc
}

func (tc *testCase) runTcMethod() {
	/* Catch exeptions in the test method body, if any, in the test method. */
	var cleanupCalledFlag bool = false
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			logStackTrace(currentTestScript.logger, buf)

			/* Call testcase cleanup on crash methed if testcase method of cleanup method panics. */
			if !cleanupCalledFlag {
				/* In case the same following two line are not executed after tc.tcMethod.Call(*tc.tcMParams). */
				logHorizonLine(currentTestScript.logger)
				currentTestScript.logger.GenerateStep("PostTest", "PostTest")
			}
			tc.callCleanupOnCrashMethod()
		}
	}()
	/* Add PRE-FIRST-STEP in case error occurs before first step. */
	currentTestScript.logger.GenerateStep("PreTest", "PreTest")
	/* Call testcase method. */
	tc.method.Call(*tc.methodParams)
	/* Call testcase cleanup method if there is not panic in the procedure of testcase method
	   if there is panic in the testcase method the cleanup method will not be called.*/
	cleanupCalledFlag = true
	logHorizonLine(currentTestScript.logger)
	currentTestScript.logger.GenerateStep("PostTest", "PostTest")
	tc.callCleanupMethod()
}

func (tc *testCase) callCleanupOnCrashMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			logStackTrace(currentTestScript.logger, buf)
		}
	}()
	tc.callMethod("CleanupOnCrash")
}

func (tc *testCase) callCleanupMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			logStackTrace(currentTestScript.logger, buf)
		}
	}()
	tc.callMethod("Cleanup")
}

func (tc *testCase) callMethod(method string) {
	mc := tc.methodName + method
	m := tc.testScript.tTest.MethodByName(mc)
	if m.Kind() != reflect.Func {
		log.Warningf("The %s method %s is NOT definded.", method, mc)
	} else {
		m.Call(nil)
	}
}

func getFunctionName(rv reflect.Value) (string, string) {
	qualifiedFuncName := runtime.FuncForPC(rv.Pointer()).Name()

	reg := regexp.MustCompile(`(\w+).*?\w+.(\w+)·fm`)
	matchs := reg.FindStringSubmatch(qualifiedFuncName)
	if matchs == nil {
		panic("The qualified function name: " + qualifiedFuncName + `does NOT match the regexp: \w+.*?\w+.(\w)·fm.`)
	}
	return matchs[1], matchs[2]
}

func clearTcSteps(l *log.Logger) {
	l.Steps = l.Steps[0:0]
}

func (tc *testCase) logHeader() {
	tc.testScript.logger.Output("TC_HEADING",
		log.LOnlyFile,
		log.TestcaseHdrInfo{
			tc.tcID,
			time.Now().Format("2006-01-02 15:04:05"),
			tc.description,
		})
}

func (tc *testCase) logResult() {
	var faildSteps string
	var logger = tc.testScript.logger
	defer clearTcSteps(logger)

	for _, step := range logger.Steps {
		if step.IsFailed {
			faildSteps = faildSteps + "{" + step.Index + "} "
		}
	}
	if faildSteps != "" {
		logFailedSteps(logger, faildSteps)
	}
	tcSummaryResult := generateTcResultSummary(logger, tc.tcID, tc.description, faildSteps)

	/* TODO: enhance it if possible. */
	logger.CloseFile()
	logFileContent, err := ioutil.ReadFile(logger.FileName())
	if err != nil {
		panic(err)
	}
	regexpInsertTcSummary := regexp.MustCompile(`<div style="display:none">hide</div>`)
	logFileContent = regexpInsertTcSummary.ReplaceAll(logFileContent, tcSummaryResult.Bytes())
	ioutil.WriteFile(logger.FileName(), logFileContent, 0666)
	logger.ReopenFile()
}
