// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*Aggregator)(nil)

// Aggregator is used to aggregate payable tax for individuals from multiple
// tax formulas
type Aggregator struct {
	calculators []*Calculator
}

// NewAggregator returns a new tax calculator for the given financial numbers
// and tax formulas. The returned calculator will calculate the sum of taxes
// using the given formulas for the given finances
func NewAggregator(formula1, formula2 Formula, extras ...Formula) (*Aggregator, error) {

	cAgg := &Aggregator{
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
func (agg *Aggregator) Calc(finances finance.IndividualFinances) float64 {

	var payableTax float64
	for _, c := range agg.calculators {
		payableTax += c.Calc(finances)
	}
	return payableTax
}
