package gtf

import (
	"fmt"
	"gtf/log"
	"reflect"
)

type paramsFlag bool

const (
	Overridable    paramsFlag = true
	NonOverridable paramsFlag = false
)

type Test struct{ DDD string }

func (t *Test) DefineCase(tcid, tcDescr string) *tcDefinition {
	tcdef := &tcDefinition{tcid: tcid, tcDescription: tcDescr}
	tcDefinitions[tcid] = tcdef
	return tcdef
}

func (t *Test) SetParam(param string, value interface{}, overridable ...paramsFlag) {
	if len(overridable) > 1 {
		panic("The overrideable parameter needs only ONE argument.")
	}

	if _, ok := testParams[param]; !ok || (len(overridable) == 1 && overridable[0] == NonOverridable) {
		testParams[param] = value
	}
}

// Called by TestCaseProcedure in ths testcase scripts to run real tests,
// tcTestLogicMethod is the real test method with test logic
// tcid is the first parameter of the method tcTestLogicMethod
// params is other parameter(s), if any, of the method tcTestLogicMethod
func (t *Test) ExecuteTestCase(tcTestLogicMethod interface{}, tcid string, params ...interface{}) {
	defer func() {
		logTcResult(currentTestScript.logger, tcid, tcDefinitions[tcid].tcDescription)
	}()

	fmt.Println("[STEP] FROM GTF's <<ExecuteTestCase>>")
	if tcdef, ok := tcDefinitions[tcid]; ok {
		if !tcdef.a {
			fmt.Println("[ERROR] The testcase is not applicable.")
			return
		}
		if !tcdef.r {
			fmt.Println("[ERROR] The testcase's requirements are not be satisfied.")
			return
		}
	} else {
		// The testcase is not defined.
		fmt.Println("[ERROR] The testcase is not defined.")
		return
	}
	logTcHearder(currentTestScript.logger, tcid, tcDefinitions[tcid].tcDescription)
	tc := newTestCase(tcTestLogicMethod, tcid, &params)
	tc.runTcMethod()
}

// ExecStep exemine if the (first) return of the func f matchs the string expect.
// The expect string may be: "string", "regexp", "glob string", [num1, num2], [num1,num2), [num,),
// {elem1, elem2, elem3,}, exp1||exp2||exp3
// testStepLogicMethod is the test logic method of a step
// params are  parameter(s), if any, of the method testStepLogicMethod
func (t *Test) ExecStep(expect interface{}, testStepLogicMethod interface{}, params ...interface{}) {
	var tcmParams []reflect.Value
	sf := reflect.ValueOf(testStepLogicMethod)
	if sf.Kind() != reflect.Func {
		panic("the step func mast be a function!")
	}

	paramsLen := len(params)
	if paramsLen != 0 {
		tcmParams = make([]reflect.Value, paramsLen)
		for i := 0; i < paramsLen; i++ {
			tcmParams[i] = reflect.ValueOf((params)[i])
		}
	}

	ret := sf.Call(tcmParams)
	if len(ret) == 0 {
		panic("It seems the step func does NOT return any value, so should not be called by ExecStep func.")
	}

	switch expect.(type) {
	case string:
		switch ret[0].Type().String() {
		case "string":
			var exp string = expect.(string)

			log.Logf("Expected result for the step: %s", exp)
			if exp == ret[0].String() {
				log.Logf("Actual result for the step: %s", ret[0].String())
			} else {
				log.Errorf("Actual result for the step: %s", ret[0].String())
			}
		case "int":
			fmt.Println("--->", "int", ret[0].Type().String())
		}
	case int:
		fmt.Println("--->", ret[0].Type())
	}
	fmt.Println(expect, ret[0].Type())
}
