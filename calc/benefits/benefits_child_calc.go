package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"

	"github.com/pkg/errors"
)

// ChildBenfitCalculator is used to calculate recievable child benefits for
// families with dependent children. This type implements the following
// interface: 'calc.ChildBenefitCalculator'
type ChildBenfitCalculator struct {
	children         []human.Person
	formula          ChildBenefitFormula
	incomeCalculator calc.IncomeCalculator
}

// compile-time check for interface implementation
var _ calc.ChildBenefitCalculator = (*ChildBenfitCalculator)(nil)

// NewChildBenefitCalculator returns a new child benefit calculator for the
// given formula and the income calculator
func NewChildBenefitCalculator(formula ChildBenefitFormula, incCalc calc.IncomeCalculator) (*ChildBenfitCalculator, error) {

	if formula == nil {
		return nil, ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	if incCalc == nil {
		return nil, ErrNoIncCalc
	}

	c := &ChildBenfitCalculator{
		formula:          formula.Clone(),
		incomeCalculator: incCalc,
	}
	return c, nil
}

// Calc returns the recievable amount of child benefits
func (c *ChildBenfitCalculator) Calc(finances finance.IncomeDeductor) float64 {

	netIncome := c.incomeCalculator.NetIncome(finances)
	benefits := c.formula.Apply(netIncome, c.children...)
	return benefits
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (c *ChildBenfitCalculator) SetBeneficiaries(children ...human.Person) {
	c.children = make([]human.Person, len(children))
	copy(c.children, children)
}