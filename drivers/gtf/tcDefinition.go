package gtf

type tcDefinition struct {
	tcid          string
	tcDescription string
	tcApplicable  string
	tcRequirement string //The requirements that a TestCase needs to test environment or settings.
	a             bool   //The applicability of the defined testcase is satisfied
	r             bool   //The requirements of the defined testcase are satisfied
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
