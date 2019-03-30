// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*CalculatorAgg)(nil)

// Calculator is used to calculate payable tax for individuals
type CalculatorAgg struct {
	calculators []*Calculator
}

// NewCalculatorAgg returns a new tax calculator for the given financial numbers
// and tax formulas. The returned calculator will calculate the sum of taxes for
// using the given formulas for the given finances
func NewCalculatorAgg(formula1, formula2 Formula, extras ...Formula) (*CalculatorAgg, error) {

	cAgg := &CalculatorAgg{
		calculators: make([]*Calculator, len(extras)+2),
	}

	for i, formula := range append([]Formula{formula1, formula2}, extras...) {
		c, err := NewCalculator(formula)
		if err != nil {
			return nil, err
		}
		cAgg.calculators[i] = c
	}

	return cAgg, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (agg *CalculatorAgg) Calc(finances calc.IndividualFinances) float64 {

	var payableTax float64
	for _, c := range agg.calculators {
		payableTax += c.Calc(finances)
	}
	return payableTax
}
