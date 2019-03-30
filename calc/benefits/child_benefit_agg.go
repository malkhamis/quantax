package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
)

// ChildBenfitAggregator is used to calculate recievable child benefits for
// families with dependent children. This type can aggregate the benefits from
// multiple child benefit formulas. This type implements the following
// interface: 'calc.ChildBenefitCalculator'
type ChildBenfitAggregator struct {
	calculators []*ChildBenfitCalculator
}

// compile-time check for interface implementation
var _ calc.ChildBenefitCalculator = (*ChildBenfitAggregator)(nil)

// NewChildBenfitAggregator returns a new child benefit calculator which can
// aggregate the benefits from each given formula
func NewChildBenefitAggregator(formula1, formula2 ChildBenefitFormula, extras ...ChildBenefitFormula) (*ChildBenfitAggregator, error) {

	cAgg := &ChildBenfitAggregator{
		calculators: make([]*ChildBenfitCalculator, len(extras)+2),
	}

	for i, formula := range append([]ChildBenefitFormula{formula1, formula2}, extras...) {
		c, err := NewChildBenefitCalculator(formula)
		if err != nil {
			return nil, err
		}
		cAgg.calculators[i] = c
	}

	return cAgg, nil
}

// Calc returns the aggregate recievable amount of child benefits
func (agg *ChildBenfitAggregator) Calc(finances finance.FamilyFinances) float64 {

	var total float64
	for _, c := range agg.calculators {
		total += c.Calc(finances)
	}
	return total
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (agg *ChildBenfitAggregator) SetBeneficiaries(children ...calc.Person) {
	for _, c := range agg.calculators {
		c.SetBeneficiaries(children...)
	}
}
