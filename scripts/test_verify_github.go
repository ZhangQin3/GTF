package git

import (
	"gtf/drivers/assert"
	"gtf/drivers/gtf"
	"gtf/drivers/log"
	"gtf/library/github"
	"se"
)

type Test struct{ gtf.Test }

/* Set Test Params */
func (t *Test) SetTestParams() {
	t.SetParam("EnableRouter", 1, gtf.NonOverridable)
	t.SetParam("CpeIpVersion", 1)
}

func (t *Test) CaseDefinitions() {
	t.DefineTestCase("tcid001", "Verify login to the github.")
	t.DefineTestCase("tcid002", "Create a github project.")
}

/* ============================================================= */
func (t *Test) TestCaseProcedure() {
	t.ExecuteTestCase(t.VerifyLogin, "tcid001")
	// t.ExecuteTestCase(t.VerifyCreateProject, "tcid002")
}

func (t *Test) VerifyLogin(tcid string) {
	log.Step(1, "Open github.")
	p, e := github.OpenGithub(se.Browsername("chrome"))
	assert.Nil(e)
	defer p.Close()

	// log.DoCatch(ttt, 1, 2)

	log.Step(2, "Login github.")
	e = p.SignIn().Click()
	log.DoPanic(e)
	p.UserName().SetText("goautomation", se.PreClear)
	p.Password().SetText("0web.driver")
	p.Signin().Click()

	log.Step(3, "Logout")
	p.Profile().Click()
	p.Logout().Click()
}

func ttt(a, b int) int {
	c := a + b
	log.Info("dddddddddd", c)

	return c
}

func (t *Test) VerifyCreateProject(tcid string) {
	log.Step(1, "Open github.")
	p, _ := github.OpenGithub()
	defer p.Close()

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
func (t *Test) VerifyLoginCleanupOnCrash() {
	log.Info("VerifyLoginCleanupOnCrash is called")
}

func (t *Test) VerifyLoginCleanup() {
	log.Info("VerifyLoginCleanup is called")
}

func (t *Test) TestCaseProcedureCleanup() {
	log.Info("TestCaseProcedureCleanup is called")
}
