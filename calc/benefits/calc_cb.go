package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitCalculator = (*CBCalculator)(nil)

type Person = calc.Person
type FamilyFinances = calc.FamilyFinances
type ChildBenefitFormula = calc.ChildBenefitFormula

// ConfigCCB are the parameters for creating a new child benefit calculator
type ConfigCCB struct {
	Finances FamilyFinances
	Formula  ChildBenefitFormula
}

// CBCalculator is used to calculate recievable child benefits for families
// with dependent children. This type implements 'calc.ChildBenefitCalculator'
type CBCalculator struct {
	children []Person
	formula  ChildBenefitFormula
	finances FamilyFinances
}

// NewCBCalculator returns a new child benefit calculator
func NewCBCalculator(cfg ConfigCCB, child Person, others ...Person) (*CBCalculator, error) {

	cbc := &CBCalculator{}
	cbc.UpdateFinances(cfg.Finances)
	cbc.UpdateBeneficiaries(child, others...)
	err := cbc.UpdateForumla(cfg.Formula)
	return cbc, err
}

// Calc returns the recievable amount of child benefits
func (c *CBCalculator) Calc() float64 {

	benefits := c.formula.Apply(
		c.finances.NetIncome(),
		c.children[0], // calc has at least one child
		c.children[1:]...,
	)

	return benefits
}

// UpdateFinances sets the financial numbers which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set financial
// numbers to anything other than what the calculator was initialized with
func (c *CBCalculator) UpdateFinances(newFinances FamilyFinances) {
	c.finances = newFinances
}

// UpdateBeneficiary sets the child which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set beneficiary
// to anything other than what the calculator was initialized with
func (c *CBCalculator) UpdateBeneficiaries(child Person, others ...Person) {

	c.children = []Person{child.Clone()}
	for _, otherChild := range others {
		c.children = append(c.children, otherChild.Clone())
	}
}

// UpdateForumla sets the formula for calculating the amount of benefits for
// children given family finances. Users may call this method to set the
// formula to anything other than what the calculator was initialized with
func (c *CBCalculator) UpdateForumla(newFormula ChildBenefitFormula) error {

	if newFormula == nil {
		return calc.ErrNoFormula
	}

	err := newFormula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula")
	}

	c.formula = newFormula
	return nil
}
