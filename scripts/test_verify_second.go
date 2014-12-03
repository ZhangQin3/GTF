package verify

import "gtf"
import "gtf/log"

type Test struct{ gtf.TEST }

/* Script level global variables */
var Str string

/* Setting script level parameters */
func (t *Test) SetTestParams() {
	t.SetParam("EnableRouter", 1, gtf.NonOverridable)
	t.SetParam("CpeIpVersion", 1)
}

func (t *Test) CaseDefinitions() {
	t.DefineCase("tcid001", "this is my first case.").R("BB")
	t.DefineCase("tcid002", "this is my second case.").A("D30")
	t.DefineCase("tcid003", "this is my third case.").RA("BB", "D30")
}

// --------------- Testcase Procedure ---------------
func (t *Test) TestCaseProcedure() {
	t.ExecuteTestCase(t.VerifyPrototype, "tcid001")
	t.ExecuteTestCase(t.VerifyPrototypeSec, "tcid002")
	t.ExecuteTestCase(t.VerifyPrototypeThird, "tcid003")
}

// --------------- Test Procedure ---------------
func (t *Test) VerifyPrototype(tcid string) {
	log.Step(1, "setup ----null---- router")
	t.ExecStep("OK", t.stepTest, "123")

	log.Step("aaa", "setup ----%s---- router", "ospf")
	t.ExecStep("OK", t.stepTest, "123")

	log.Step(".", "setup ----%s---- router", "ospf")
	t.ExecStep("OK", t.stepTest, "123")
}

func (t *Test) VerifyPrototypeSec(tcid string) {
	log.Step(1, "setup ----%s---- router", "ospf")
	t.ExecStep("123", t.stepTest, "1233")

	log.Step(2, "setup ----%s---- router", "rip2")
	t.ExecStep("OK", t.stepTest, "456")
}

func (t *Test) VerifyPrototypeThird(tcid string) {
	log.Step(1, "setup ----%s---- router", "ospf")
	t.ExecStep("123", t.stepTest, "123")

	log.Step(2, "setup ----%s---- router", "rip2")
	t.ExecStep("OK", t.stepTest, "OK")
}

func (t *Test) stepTest(str string) string {

	return str
}

// --------------- Test Cleanup ---------------
func (t *Test) VerifyPrototypeCleanupOnCrash() {
	log.Error("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<=====")
}

func (t *Test) VerifyPrototypeCleanup() {
	log.Error("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<=====")
}

func (t *Test) TestCaseProcedureCleanup() {
	log.Log("in TestCaseProcedureCleanup")
}
