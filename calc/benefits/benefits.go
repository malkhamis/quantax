package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type (
	// ChildBenefitFormula represents a method for calculating child benefits
	ChildBenefitFormula interface {
		// Apply returns the sum of benefits for all beneficiaries
		Apply(income float64, children ...calc.Person) float64
		// IncomeCalcMethod returns the method of calculating the income
		IncomeCalcMethod() IncomeType
		// Validate checks if the formula is valid for use
		Validate() error
		// Clone returns a copy of the formula
		Clone() ChildBenefitFormula
	}

	// Calculator is used to calculate recievable child benefits for families
	// with dependent children. This type implements 'calc.ChildBenefitCalculator'
	Calculator struct {
		children []calc.Person
		formula  ChildBenefitFormula
	}
)

// compile-time check for interface implementation
var _ calc.ChildBenefitCalculator = (*Calculator)(nil)

// NewCalculator returns a new child benefit calculator
func NewCalculator(formula ChildBenefitFormula) (*Calculator, error) {

	if formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	return &Calculator{formula: formula.Clone()}, nil
}

// Calc returns the recievable amount of child benefits
func (c *Calculator) Calc(finances calc.FamilyFinances) float64 {

	netIncome := c.formula.IncomeCalcMethod().Calc(finances)
	benefits := c.formula.Apply(netIncome, c.children...)
	return benefits
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (c *Calculator) SetBeneficiaries(children ...calc.Person) {
	c.children = make([]calc.Person, len(children))
	copy(c.children, children)
}
