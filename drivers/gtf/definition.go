package gtf

import (
	"fmt"
)

// test case definition
type tcDef struct {
	tcid        string
	description string
	applicable  string
	requirement string // The requirements that a testcase needs to test environment or settings.
	priority    string // The priority of the testcase.
	feature     string // The feature(s) the testcase belongs to.
}

func (def *tcDef) R(requirement string) *tcDef {
	def.requirement = requirement
	return def
}

func (def *tcDef) A(applicable string) *tcDef {
	def.applicable = applicable
	return def
}

func (def *tcDef) RA(requirement, applicable string) *tcDef {
	def.requirement = requirement
	def.applicable = applicable
	return def
}

// TODO priority
func (def *tcDef) P(priority string) *tcDef {
	def.priority = priority
	return def
}

// TODO feature list
func (def *tcDef) F(feature string) *tcDef {
	def.feature = feature
	return def
}

func (def *tcDef) calculateApplicable() bool {
	//TODO: the calculation of the applicalibity according to def.tcApplicable
	//The method should be independent from the feature tested.
	return true
}

func (def *tcDef) calculateReqirement() bool {
	//TODO: the calculation of the applicalibity according to def.tcRequirement
	//The method should be independent from the feature tested.
	return true
}

func (def *tcDef) CalculateAppliability() bool {
	if !def.calculateApplicable() {
		fmt.Println("[ERROR] The testcase is not applicable.")
		return false
	}
	if !def.calculateReqirement() {
		fmt.Println("[ERROR] The testcase's requirements are not be satisfied.")
		return false
	}
	return true
}
