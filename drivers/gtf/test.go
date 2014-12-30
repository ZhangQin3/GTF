// The test.go defines a Test struct which is used as an embedded field in test scripts to promote its methods for test scripts,
// the method of Test struct will be called in test scripts.
package gtf

import (
	"fmt"
	"gtf/drivers/log"
	"reflect"
)

type paramFlag bool

const (
	Overridable    paramFlag = true // Overridable means the test parameters can be overrided by test suite parameters with same name.
	NonOverridable paramFlag = false
)

type Test struct {
	tcDefs       map[string]*tcDef
	DemoVariable string
}

// flag parameter is just an optional param, not variadic.
func (t *Test) SetParam(param string, value interface{}, flag ...paramFlag) {
	if len(flag) > 1 {
		panic("The overrideable parameter is just an optional param, not variadic.")
	}

	if _, ok := TestParams[param]; !ok || (len(flag) == 1 && flag[0] == Overridable) {
		TestParams[param] = value
	}
}

func (t *Test) DefineCase(tcid, description string) *tcDef {
	if t.tcDefs == nil {
		t.tcDefs = make(map[string]*tcDef)
	}
	tcDef := &tcDef{tcid: tcid, description: description}
	t.tcDefs[tcid] = tcDef
	return tcDef
}

// Called by TestCaseProcedure in ths testcase scripts to run real tests,
// testLogicMethod is the real test method with test logic
// tcid is the first parameter of the method testLogicMethod
// params is other parameter(s), if any, of the method testLogicMethod
func (t *Test) ExecuteTestCase(testLogicMethod interface{}, tcid string, params ...interface{}) {
	tc := newTestCase(testLogicMethod, tcid, &params)
	defer func() {
		tc.logResult()
	}()

	if tcDef, ok := t.tcDefs[tcid]; ok {
		if !tcDef.CalculateAppliability() {
			return
		}
	} else {
		fmt.Println("[ERROR] The testcase is not defined.")
		return
	}

	tc.logHeader()
	tc.runTcMethod()
}

// ExecStep exemine if the (first) return of the func f matchs the string expect.
// The expect string may be: "string", "regexp", "glob string", [num1, num2], [num1,num2), [num,),
// {elem1, elem2, elem3,}, exp1||exp2||exp3
// stepLogicMethod is the test logic method of a step
// params are  parameter(s), if any, of the method stepLogicMethod
func (t *Test) ExecStep(expected interface{}, stepLogicMethod interface{}, params ...interface{}) {
	var tcmParams []reflect.Value
	stepMethod := reflect.ValueOf(stepLogicMethod)
	if stepMethod.Kind() != reflect.Func {
		panic("the step func mast be a function!")
	}

	paramsLen := len(params)
	if paramsLen != 0 {
		tcmParams = make([]reflect.Value, paramsLen)
		for i := 0; i < paramsLen; i++ {
			tcmParams[i] = reflect.ValueOf((params)[i])
		}
	}

	ret := stepMethod.Call(tcmParams)
	if len(ret) == 0 {
		panic("It seems the step func does NOT return any value, so should not be called by ExecStep func.")
	}

	switch expected.(type) {
	case string:
		switch ret[0].Type().String() {
		case "string":
			var exp string = expected.(string)

			log.Info("Expected result for the step: ", exp)
			if exp == ret[0].String() {
				log.Info("Actual result for the step: ", ret[0])
			} else {
				log.Error("Actual result for the step: ", ret[0])
			}
		case "int":
			fmt.Println("--->", "int", ret[0].Type().String())
		}
	case int:
		fmt.Println("--->", ret[0].Type())
	}
	fmt.Println(expected, ret[0].Type())
}
