package tsuite

import (
	"gtf/drivers/suite"
)

type TSuite struct{ suite.Suite }

/* Setting suite level parameters */
func (ts *TSuite) SetSuiteParams() {
	ts.SetParam("NewSWVersion", "1.1.1.2")
	ts.SetParam("OldSWVersion", "1.1.1.1")
}

func (ts *TSuite) SuiteScheme() {
	ts.Scheme = `
       // {"script": "test_verify_test.go","repetitions":1}
       // {"script": "test_verify_second.go","repetitions":2}
       // {"script": "test_verify_csvdata.go"}
       {"script": "verify_github.csv"}
       // {"script": "test_verify_github.go"}
	`
}

/* run on the beginning of the test suite. */
func (ts *TSuite) SuitSetup() {

}

/* run on the end on the test suite. */
func (ts *TSuite) SuiteCleanup() {

}

/* run on the beginning of every testcase. */
func (ts *TSuite) CaseSetup() {

}

/* run on the end of every testcase. */
func (ts *TSuite) CaseCleanup() {

}
