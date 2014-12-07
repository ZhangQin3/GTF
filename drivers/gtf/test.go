package gtf

import (
	"fmt"
	"gtf/log"
	"reflect"
)

type tcDefinition struct {
	tcid          string
	tcDescription string
	tcApplicable  string
	tcRequirement string //The requirements that a TestCase needs to test environment or settings.
	a             bool   //The applicability of the defined testcase is satisfied
	r             bool   //The requirements of the defined testcase are satisfied
}

const (
	Overridable    bool = true
	NonOverridable bool = false
)

type Test struct{ DDD string }

func (t *Test) DefineCase(tcid, tcDescr string) *tcDefinition {
	tcdef := &tcDefinition{tcid: tcid, tcDescription: tcDescr}
	tcDefinitions[tcid] = tcdef
	return tcdef
}

func (tcdef *tcDefinition) R(req string) {
	tcdef.tcRequirement = req
	tcdef.calcReqirement()
	tcdef.a = true
}

func (tcdef *tcDefinition) A(app string) {
	tcdef.tcApplicable = app
	tcdef.calcApplicable()
	tcdef.r = true
}

// TODO priority
func (tcdef *tcDefinition) P(app string) {

}

// TODO feature list
func (tcdef *tcDefinition) F(app string) {

}

func (tcdef *tcDefinition) RA(req, app string) {
	tcdef.tcRequirement = req
	tcdef.tcApplicable = app
	tcdef.calcReqirement()
	tcdef.calcApplicable()
}

func (tcdef *tcDefinition) calcApplicable() {
	//TODO: the calculation of the applicalibity according to tcdef.tcApplicable
	//The method should be independent from the feature tested.
	tcdef.a = true
}

func (tcdef *tcDefinition) calcReqirement() {
	//TODO: the calculation of the applicalibity according to tcdef.tcRequirement
	//The method should be independent from the feature tested.
	tcdef.r = true
}

func (t *Test) SetParam(param string, value interface{}, overridable ...bool) {
	if len(overridable) > 1 {
		panic("The overrideable parameter needs only ONE argument.")
	}

	if _, ok := testParams[param]; !ok || (len(overridable) == 1 && overridable[0] == NonOverridable) {
		testParams[param] = value
	}
}

// Called by TestCaseProcedure in ths testcase scripts to run real tests.
func (t *Test) ExecuteTestCase(f interface{}, tcid string, params ...interface{}) {
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
	tc := newTestCase(f, tcid, &params)
	tc.runTcMethod()
}

// ExecStep exemine if the (first) return of the func f matchs the string expect.
// The expect string may be: "string", "regexp", "glob string", [num1, num2], [num1,num2), [num,),
// {elem1, elem2, elem3,}, exp1||exp2||exp3
func (t *Test) ExecStep(expect interface{}, f interface{}, params ...interface{}) {
	var tcmParams []reflect.Value
	sf := reflect.ValueOf(f)
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
