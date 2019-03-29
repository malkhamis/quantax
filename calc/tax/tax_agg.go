// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
)

// compile-time check for interface implementation
var (
	_ calc.TaxCalculator = (*CalculatorAgg)(nil)
)

// Calculator is used to calculate payable tax for individuals
type CalculatorAgg struct {
	calculators []*Calculator
	finances    calc.IndividualFinances
}

// NewCalculatorAgg returns a new tax calculator for the given financial numbers
// and tax formulas. The returned calculator will calculate the sum of taxes for
// using the given formulas for the given finances
func NewCalculatorAgg(finances calc.IndividualFinances, formulas []Formula) (*CalculatorAgg, error) {

	if len(formulas) < 1 {
		return nil, calc.ErrNoFormula
	}

	cAgg := &CalculatorAgg{
		calculators: make([]*Calculator, len(formulas)),
	}

	for i, formula := range formulas {
		c, err := NewCalculator(finances, formula)
		if err != nil {
			return nil, err
		}
		cAgg.calculators[i] = c
	}

	return cAgg, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (agg *CalculatorAgg) Calc() float64 {

	var payableTax float64
	for _, c := range agg.calculators {
		payableTax += c.Calc()
	}
	return payableTax
}

// Update sets the financial numbers which the tax will be calculated for
func (agg *CalculatorAgg) UpdateFinances(newFinances calc.IndividualFinances) {
	for _, c := range agg.calculators {
		c.UpdateFinances(newFinances)
	}
}
