package gtf

import (
	"fmt"
)

type tcDefinition struct {
	tcid        string
	description string
	applicable  string
	requirement string // The requirements that a testcase needs to test environment or settings.
	priority    string // The priority of the testcase.
	feature     string // The feature(s) the testcase belongs to.
}

func (tcDef *tcDefinition) R(requirement string) *tcDefinition {
	tcDef.requirement = requirement
	return tcDef
}

func (tcDef *tcDefinition) A(applicable string) *tcDefinition {
	tcDef.applicable = applicable
	return tcDef
}

func (tcDef *tcDefinition) RA(requirement, applicable string) *tcDefinition {
	tcDef.requirement = requirement
	tcDef.applicable = applicable
	return tcDef
}

// TODO priority
func (tcDef *tcDefinition) P(priority string) *tcDefinition {
	tcDef.priority = priority
	return tcDef
}

// TODO feature list
func (tcDef *tcDefinition) F(feature string) *tcDefinition {
	tcDef.feature = feature
	return tcDef
}

func (tcDef *tcDefinition) calculateApplicable() bool {
	//TODO: the calculation of the applicalibity according to tcDef.tcApplicable
	//The method should be independent from the feature tested.
	return true
}

func (tcDef *tcDefinition) calculateReqirement() bool {
	//TODO: the calculation of the applicalibity according to tcDef.tcRequirement
	//The method should be independent from the feature tested.
	return true
}

func (tcDef *tcDefinition) CalculateAppliability() bool {
	if !tcDef.calculateApplicable() {
		fmt.Println("[ERROR] The testcase is not applicable.")
		return false
	}
	if !tcDef.calculateReqirement() {
		fmt.Println("[ERROR] The testcase's requirements are not be satisfied.")
		return false
	}
	return true
}
