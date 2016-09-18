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

func (d *tcDef) R(requirement string) *tcDef {
	d.requirement = requirement
	return d
}

func (d *tcDef) A(applicable string) *tcDef {
	d.applicable = applicable
	return d
}

func (d *tcDef) RA(requirement, applicable string) *tcDef {
	d.requirement = requirement
	d.applicable = applicable
	return d
}

func (d *tcDef) Description() string {
	return d.description
}

// TODO priority
func (d *tcDef) P(priority string) *tcDef {
	d.priority = priority
	return d
}

// TODO feature list
func (d *tcDef) F(feature string) *tcDef {
	d.feature = feature
	return d
}

func (d *tcDef) calculateApplicable() bool {
	//TODO: the calculation of the applicalibity according to def.tcApplicable
	//The method should be independent from the feature tested.
	return true
}

func (d *tcDef) calculateReqirement() bool {
	//TODO: the calculation of the applicalibity according to def.tcRequirement
	//The method should be independent from the feature tested.
	return true
}

func (d *tcDef) CalculateAppliability() bool {
	if !d.calculateApplicable() {
		fmt.Println("[ERROR] The testcase is not applicable.")
		return false
	}
	if !d.calculateReqirement() {
		fmt.Println("[ERROR] The testcase's requirements are not be satisfied.")
		return false
	}
	return true
}
