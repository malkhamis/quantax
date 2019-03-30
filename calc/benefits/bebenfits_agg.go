package benefits

import (
	"github.com/malkhamis/quantax/calc"
)

// CalculatorAgg is used to calculate recievable child benefits for families
// with dependent children. This type implements 'calc.ChildBenefitCalculator'
type CalculatorAgg struct {
	calculators []*Calculator
}

// compile-time check for interface implementation
var _ calc.ChildBenefitCalculator = (*CalculatorAgg)(nil)

// NewCalculatorAgg returns a new child benefit calculator which aggregate the
// benefits from each given formula
func NewCalculatorAgg(formula1, formula2 ChildBenefitFormula, extras ...ChildBenefitFormula) (*CalculatorAgg, error) {

	cAgg := &CalculatorAgg{
		calculators: make([]*Calculator, len(extras)+2),
	}

	for i, formula := range append([]ChildBenefitFormula{formula1, formula2}, extras...) {
		c, err := NewCalculator(formula)
		if err != nil {
			return nil, err
		}
		cAgg.calculators[i] = c
	}

	return cAgg, nil
}

// Calc returns the aggregate recievable amount of child benefits
func (agg *CalculatorAgg) Calc(finances calc.FamilyFinances) float64 {

	var total float64
	for _, c := range agg.calculators {
		total += c.Calc(finances)
	}
	return total
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (agg *CalculatorAgg) SetBeneficiaries(children ...calc.Person) {
	for _, c := range agg.calculators {
		c.SetBeneficiaries(children...)
	}
}
