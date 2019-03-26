package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitCalculator = (*CBCalculator)(nil)

// ConfigCB are the parameters for creating a new child benefit calculator
type ConfigCB struct {
	Finances calc.FamilyFinances
	Formula  calc.ChildBenefitFormula
}

// CBCalculator is used to calculate recievable child benefits for families
// with dependent children. This type implements 'calc.ChildBenefitCalculator'
type CBCalculator struct {
	children []calc.Person
	formula  calc.ChildBenefitFormula
	finances calc.FamilyFinances
}

// NewCBCalculator returns a new child benefit calculator
func NewCBCalculator(cfg ConfigCB, child calc.Person, others ...calc.Person) (*CBCalculator, error) {

	if cfg.Formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := cfg.Formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	cbc := &CBCalculator{formula: cfg.Formula}
	cbc.UpdateFinances(cfg.Finances)
	cbc.UpdateBeneficiaries(child, others...)
	return cbc, nil
}

// Calc returns the recievable amount of child benefits
func (c *CBCalculator) Calc() float64 {

	netIncome := c.finances.Income() - c.finances.Deductions()

	benefits := c.formula.Apply(
		netIncome,
		c.children[0], // calc has at least one child
		c.children[1:]...,
	)

	return benefits
}

// UpdateFinances sets the financial numbers which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set financial
// numbers to anything other than what the calculator was initialized with
func (c *CBCalculator) UpdateFinances(newFinances calc.FamilyFinances) {
	c.finances = newFinances
}

// UpdateBeneficiary sets the child which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set beneficiary
// to anything other than what the calculator was initialized with
func (c *CBCalculator) UpdateBeneficiaries(child calc.Person, others ...calc.Person) {

	c.children = []calc.Person{child}
	for _, otherChild := range others {
		c.children = append(c.children, otherChild)
	}
}
