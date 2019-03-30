package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

// ChildBenfitCalculator is used to calculate recievable child benefits for
// families with dependent children. This type implements the following
// interface: 'calc.ChildBenefitCalculator'
type ChildBenfitCalculator struct {
	children []calc.Person
	formula  ChildBenefitFormula
}

// compile-time check for interface implementation
var _ calc.ChildBenefitCalculator = (*ChildBenfitCalculator)(nil)

// NewChildBenefitCalculator returns a new child benefit calculator for the
// given formula
func NewChildBenefitCalculator(formula ChildBenefitFormula) (*ChildBenfitCalculator, error) {

	if formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	return &ChildBenfitCalculator{formula: formula.Clone()}, nil
}

// Calc returns the recievable amount of child benefits
func (c *ChildBenfitCalculator) Calc(finances calc.FamilyFinances) float64 {

	netIncome := c.formula.IncomeCalcMethod().Calc(finances)
	benefits := c.formula.Apply(netIncome, c.children...)
	return benefits
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (c *ChildBenfitCalculator) SetBeneficiaries(children ...calc.Person) {
	c.children = make([]calc.Person, len(children))
	copy(c.children, children)
}