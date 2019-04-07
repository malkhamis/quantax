// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*Aggregator)(nil)

// Aggregator is used to aggregate payable tax from multiple tax calculators
type Aggregator struct {
	calculators []calc.TaxCalculator
}

// NewAggregator returns a new tax aggregator for the given tax calculators
func NewAggregator(c0, c1 calc.TaxCalculator, extras ...calc.TaxCalculator) (*Aggregator, error) {

	cAgg := &Aggregator{
		calculators: make([]calc.TaxCalculator, 0, len(extras)+2),
	}

	for i, c := range append([]calc.TaxCalculator{c0, c1}, extras...) {
		if c == nil {
			return nil, errors.Wrapf(ErrNoCalc, "index %d: invalid calculator", i)
		}
		cAgg.calculators = append(cAgg.calculators, c)
	}

	return cAgg, nil

}

// Calc computes the tax on the taxable amount set in this calculator
func (agg *Aggregator) Calc(finances *finance.IndividualFinances) float64 {

	var payableTax float64
	for _, c := range agg.calculators {
		payableTax += c.Calc(finances)
	}
	return payableTax
}
