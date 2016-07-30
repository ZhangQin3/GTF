package git

import (
	"gtf/drivers/gtf"
	"gtf/drivers/log"
	"gtf/library/github"
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

/* ===================== Testcase Procedure ===================== */
func (t *Test) TestCaseProcedure() {
	t.ExecuteTestCase(t.VerifyPrototype, "tcid001")
}

/* ===================== Test Procedure ===================== */
func (t *Test) VerifyPrototype(tcid string) {
	log.Step(1, "Open github.")
	p := github.OpenGithub()

	log.Step(2, "Login github.")
	p.SignIn().Click()
	p.UserName().SetText("goautomation")
	p.Password().SetText("0web.driver")
	p.Signin().Click()

	log.Step(3, "Create a new project")
	p.NewProject().Click()

	log.Step(4, "Logout")
	p.Profile().Click()
	p.Logout().Click()

	// p.UserName().SetText("ddddddddd")

	// log.Step(2, "Input user name.")
	// p.UserName().SetText("technician")

	// b, e := p.UserName().DoesExist()
	// log.Warning("----------- ", b, e)
	// x, y := p.ScreenShot()
	// log.Warning("------+++++----- ", x, y)

	// log.Step(3, "Input user password.")
	// c, f := p.Passord().DoesExist()
	// log.Warning("----------- ", c, f)

	// p.Passord().SetText("T!m3W4rn3rC4bl3")

	// log.Step(4, "Apply inputs.")
	// p.Apply().Click()

	// // p.WanSetup().Click()

	// log.Step(5, "Goto wan setup page.")
	// github.OpenWanSetup(p)
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

/* ===================== Test Cleanup ===================== */
func (t *Test) VerifyPrototypeCleanupOnCrash() {

}

func (t *Test) VerifyPrototypeCleanup() {

}

func (t *Test) TestCaseProcedureCleanup() {

}
