package tsuite

import (
	"fmt"
	suite "gtf/suite"
)

type TSuite struct{ suite.Suite }

/* Setting suite level parameters */
func (ts *TSuite) SetSuiteParams() {
	ts.SetParam("NewSWVersion", "1.1.1.1")
	ts.SetParam("OldSWVersion", "1.1.1.2")
}

func (ts *TSuite) SuiteSchema() {
	ts.Schema = `
       {"script": "test_verify_test.go","repetitions":1}
       {"script": "test_verify_second.go","repetitions":2}
	`

	// {"script":"test.go", "repetitions":10, "other":"Test example."}
	// {"script":"test_verify_webs.go"}
	// {"script":"auto.go"}
	// {"script":"test_verify_web.go"}
}

// run on the beginning of the test suite.
func (ts *TSuite) SuitSetup() {
	fmt.Println("--------------->>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

// run on the end on the test suite.
func (ts *TSuite) SuiteTeardown() {
	fmt.Println(">>>>>>>>>>>>>>>>>---------------------------")
}

// run on the beginning of every testcase.
func (ts *TSuite) CaseSetup() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

// run on the end of every testcase.
func (ts *TSuite) CaseTeardown() {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
