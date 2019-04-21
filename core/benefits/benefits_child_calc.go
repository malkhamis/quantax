package benefits

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

// ChildBenfitCalculator is used to calculate recievable child benefits for
// families with dependent children. This type implements the following
// interface: 'core.ChildBenefitCalculator'
type ChildBenfitCalculator struct {
	children         []human.Person
	formula          ChildBenefitFormula
	incomeCalculator core.IncomeCalculator
}

// compile-time check for interface implementation
var _ core.ChildBenefitCalculator = (*ChildBenfitCalculator)(nil)

// NewChildBenefitCalculator returns a new child benefit calculator for the
// given formula and the income calculator
func NewChildBenefitCalculator(cfg CalcConfigCB) (*ChildBenfitCalculator, error) {

	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := &ChildBenfitCalculator{
		formula:          cfg.Formula.Clone(),
		incomeCalculator: cfg.IncomeCalc,
	}
	return c, nil
}

// Calc returns the recievable amount of child benefits
func (c *ChildBenfitCalculator) Calc(finances core.Financer) float64 {

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
