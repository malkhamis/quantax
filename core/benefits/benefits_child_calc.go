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
	formula          ChildBenefitFormula
	incomeCalculator core.IncomeCalculator
	children         []*human.Person
	finances         core.HouseholdFinances
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

// BenefitRecievable returns the recievable amount of child benefits
func (c *ChildBenfitCalculator) BenefitRecievable() float64 {

	netIncome := c.householdNetIncome()
	benefits := c.formula.Apply(netIncome, c.children...)
	return benefits
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to BenefitRecievable()
func (c *ChildBenfitCalculator) SetBeneficiaries(children []*human.Person) {
	c.children = children
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil, a non-nil, empty finances is set. It must
// be noted that this function also calls SetFinances on the income calculator
func (c *ChildBenfitCalculator) SetFinances(finances core.HouseholdFinances) {
	if finances == nil {
		finances = core.NewHouseholdFinancesNop()
	}

	c.finances = finances
}

// householdNetIncome calculates the net income of the stored household finances
func (c *ChildBenfitCalculator) householdNetIncome() float64 {

	c.incomeCalculator.SetFinances(c.finances.SpouseA())
	netIncome := c.incomeCalculator.NetIncome()

	c.incomeCalculator.SetFinances(c.finances.SpouseB())
	netIncome += c.incomeCalculator.NetIncome()

	return netIncome
}
