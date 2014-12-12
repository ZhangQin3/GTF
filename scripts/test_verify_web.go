package verify

import (
	"fmt"
	"gtf"
	"gtf/log"
	"webgui"
)

type Test struct{ gtf.Test }

/* Set Test Params */
func (t *Test) SetTestParams() {
	t.SetParam("EnableRouter", 1, gtf.NonOverridable)
	t.SetParam("CpeIpVersion", 1)
}

func (t *Test) CaseDefinitions() {
	t.DefineCase("tcid001", "this is my first case.").R("BB")
}

/* --------------- Testcase Procedure ----------- */
func (t *Test) TestCaseProcedure() {
	t.ExecuteTestCase(t.VerifyPrototype, "tcid001")
}

/* --------------- Test Procedure --------------- */
func (t *Test) VerifyPrototype(tcid string) {
	log.Step(".", "Login the WebGui.")
	// p := new(webgui.GWTWLoginPage)
	p := webgui.OpenLoginPage()
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

/* --------------- Test Cleanup --------------- */
func (t *Test) VerifyPrototypeCleanupOnCrash() {

}

func (t *Test) VerifyPrototypeCleanup() {

}

func (t *Test) TestCaseProcedureCleanup() {

}
