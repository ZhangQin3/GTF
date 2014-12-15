// The test.go defines a Test struct which is used as an embedded field in test scripts to promote its methods for test scripts,
// the method of Test struct will be called in test scripts.
package gtf

import (
	"fmt"
	"gtf/log"
	"reflect"
)

type testParamsFlag bool

const (
	Overridable    testParamsFlag = true // Overridable means the test parameters can be overrided by test suite parameters with same name.
	NonOverridable testParamsFlag = false
)

var (
	tcDefinitions = make(map[string]*tcDefinition) /* The testcase defined in the method CaseDefinitions in the test script, the key is string tcid.. */
)

type Test struct{ DemoVariable string }

// overridable parameter is just an optional param, not variadic.
func (t *Test) SetParam(param string, value interface{}, overridable ...testParamsFlag) {
	if len(overridable) > 1 {
		panic("The overrideable parameter is just an optional param, not variadic.")
	}

	if _, ok := TestParams[param]; !ok || (len(overridable) == 1 && overridable[0] == Overridable) {
		TestParams[param] = value
	}
}

func (t *Test) DefineCase(tcid, description string) *tcDefinition {
	tcDef := &tcDefinition{tcid: tcid, description: description}
	tcDefinitions[tcid] = tcDef
	return tcDef
}

// Called by TestCaseProcedure in ths testcase scripts to run real tests,
// tcTestLogicMethod is the real test method with test logic
// tcid is the first parameter of the method tcTestLogicMethod
// params is other parameter(s), if any, of the method tcTestLogicMethod
func (t *Test) ExecuteTestCase(tcTestLogicMethod interface{}, tcid string, params ...interface{}) {
	tc := newTestCase(tcTestLogicMethod, tcid, &params)
	defer func() {
		// logTestCaseResult(currentTestScript.logger, tcid, tcDefinitions[tcid].description)
		tc.logResult()
	}()

	if tcDef, ok := tcDefinitions[tcid]; ok {
		if !tcDef.CalculateAppliability() {
			return
		}
	} else {
		fmt.Println("[ERROR] The testcase is not defined.")
		return
	}

	// logTestCaseHeader(currentTestScript.logger, tcid, tcDefinitions[tcid].description)
	tc.logHeader()

	// tc := newTestCase(tcTestLogicMethod, tcid, &params)
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
