package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitCalculator = (*CBCalculator)(nil)

// CBCalculator is used to calculate recievable child benefits for families
// with dependent children. This type implements 'calc.ChildBenefitCalculator'
type CBCalculator struct {
	children []calc.Person
	formula  ChildBenefitFormula
	finances calc.FamilyFinances
}

// NewCBCalculator returns a new child benefit calculator
func NewCBCalculator(formula ChildBenefitFormula, finances calc.FamilyFinances, children ...calc.Person) (*CBCalculator, error) {

	if formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	cbc := &CBCalculator{formula: formula.Clone()}
	cbc.UpdateFinances(finances)
	cbc.UpdateBeneficiaries(children...)
	return cbc, nil
}

// Calc returns the recievable amount of child benefits
func (c *CBCalculator) Calc() float64 {

	netIncome := c.formula.IncomeCalcMethod().Calc(c.finances)
	benefits := c.formula.Apply(netIncome, c.children...)
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
func (c *CBCalculator) UpdateBeneficiaries(children ...calc.Person) {
	c.children = append(c.children, children...)
}
