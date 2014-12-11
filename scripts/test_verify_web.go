package verify

import (
	"fmt"
	"gtf"
	"gtf/log"
	"webgui"
)

type Test struct{ gtf.Test }

// Set Test Params
func (t *Test) SetTestParams() {
	fmt.Println("Print From <<SetTestParams>>.")
	t.SetParam("EnableRouter", 1, gtf.NonOverridable)
	t.SetParam("CpeIpVersion", 1)
}

func (t *Test) CaseDefinitions() {
	log.Error("=====================>")
	t.DefineCase("tcid001", "this is my first case.").R("BB")
	t.DefineCase("tcid002", "this is my seco  case.").A("D30")
	t.DefineCase("tcid003", "this is my third case.").RA("BB", "D30")
}

// --------------- Testcase Procedure ---------------
func (t *Test) TestCaseProcedure() {
	t.ExecuteTestCase(t.VerifyPrototype, "tcid001")
}

// --------------- Test Procedure ---------------
func (t *Test) VerifyPrototype(tcid string) {
	log.Step(".", "Login the WebGui.")
	// p := new(webgui.GWTWLoginPage)
	p := webgui.OPenLoginPage()
	p.UserName().SetText("technician")
	p.Passord().SetText("T!m3W4rn3rC4bl3")
	p.Apply().Click()

	// p.WanSetup().Click()

	webgui.OpenWanSetup(p.WD)
	// k, v := p.WanSetup().Text()
	// k, v = p.WanSetup().TagName()
	// fmt.Println("------------>>>>>>>>>>>>>>>>>>>>>>=====", k, v)
	// webgui.OpenWanSetup(p.WD)
	// p.HostName().SetText("text")

	// e := webgui.OPenLoginPage()
	// e.UserName().SetText("technician")
	// e.Passord().SetText("T!m3W4rn3rC4bl3")
	// e.Apply().Click()
	// p.HostName().SetText("text")
}

func (t *Test) stepTest(str string) string {
	log.Log("INFO", str)
	log.Log("ERROR", "Enter an error args.", "messge a gain")
	log.Debug("Debug form stepTest")
	log.Error("dddddddddd", "eeeeeeeeeeeeeeee", "fffffffffffffffffff")
	log.Warning("This is a warning")

	// panic("panic in the test method.")
	return "OK"
}

// --------------- Test Cleanup ---------------
func (t *Test) VerifyPrototypeCleanupOnCrash() {
	log.Error("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<=====")
	fmt.Println("Print From <<VerifyPrototypeOnCrash>>. With Para")
}

func (t *Test) VerifyPrototypeCleanup() {
	log.Error("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<=====")
	fmt.Println("Print From <<VerifyPrototypeCleanup>>. With Para")
}

func (t *Test) TestCaseProcedureCleanup() {
	log.Log("in TestCaseProcedureCleanup")
}
