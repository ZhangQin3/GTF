package gtf

import (
	"gtf/log"
	"reflect"
	"regexp"
	"runtime"
)

/* Contains the data for each testcase in a specfic test script. */
type testCase struct {
	tcMethod      reflect.Value    /* reflect method of testcase method. */
	tcMParams     *[]reflect.Value /* params of the testcase method. */
	tcMName       string           /* method name of the testcase method. */
	methodCleaup  string           /* method name , if any, called if the testcase method ends normally, to clean up the test environment. */
	methodOnCrash string           /* method name , if any, called if the testcase method crashed. */
	tstScript     *testScript
}

func newTestCase(f interface{}, tcid string, params *[]interface{}) *testCase {
	var tcMParams []reflect.Value
	var tc testCase
	tc.tstScript = tstScript
	tp := reflect.ValueOf(f)
	_, funcName := getFunctionName(tp)

	if tp.Kind() != reflect.Func {
		panic("The first param of the gtf.Execute must be a testcase method!")
	}
	tc.tcMName = funcName
	tc.tcMethod = tp

	paramsLen := len(*params)
	tcMParams = make([]reflect.Value, paramsLen+1)
	tcMParams[0] = reflect.ValueOf(tcid)
	if paramsLen != 0 {
		for i := 0; i < paramsLen; i++ {
			tcMParams[i+1] = reflect.ValueOf((*params)[i])
		}
	}
	tc.tcMParams = &tcMParams
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
			logStack(tstScript.logger, buf)

			/* Call testcase cleanup on crash methed if testcase method of cleanup method panics. */
			if !cleanupCalledFlag {
				/* In case the same following two line are not executed after tc.tcMethod.Call(*tc.tcMParams). */
				logHorizon(tstScript.logger)
				tstScript.logger.GenStp("POST-TEST", "POST-TEST")
			}
			tc.callCleanupOnCrashMethod()
		}
	}()
	/* Add PRE-FIRST-STEP in case error occurs before first step. */
	tstScript.logger.GenStp("PRE-FIRST-STEP", "PRE-FIRST-STEP")
	/* Call testcase method. */
	tc.tcMethod.Call(*tc.tcMParams)
	/* Call testcase cleanup method if there is not panic in the procedure of testcase method
	   if there is panic in the testcase method the cleanup method will not be called.*/
	cleanupCalledFlag = true
	logHorizon(tstScript.logger)
	tstScript.logger.GenStp("POST-TEST", "POST-TEST")
	tc.callCleanupMethod()
}

func (tc *testCase) callCleanupOnCrashMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			logStack(tstScript.logger, buf)
		}
	}()
	tc.callCMethod("CleanupOnCrash")
}

func (tc *testCase) callCleanupMethod() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			logStack(tstScript.logger, buf)
		}
	}()
	tc.callCMethod("Cleanup")
}

func (tc *testCase) callCMethod(method string) {
	mc := tc.tcMName + method
	m := tc.tstScript.tTest.MethodByName(mc)
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
