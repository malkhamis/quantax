// Package tax provides implementations for core.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.TaxCalculator = (*Aggregator)(nil)

// Aggregator is used to aggregate payable tax from multiple tax calculators
type Aggregator struct {
	calculators []core.TaxCalculator
}

// NewAggregator returns a new tax aggregator for the given tax calculators
func NewAggregator(c0, c1 core.TaxCalculator, extras ...core.TaxCalculator) (*Aggregator, error) {

	cAgg := &Aggregator{
		calculators: make([]core.TaxCalculator, 0, len(extras)+2),
	}

	for i, c := range append([]core.TaxCalculator{c0, c1}, extras...) {
		if c == nil {
			return nil, errors.Wrapf(ErrNoCalc, "index %d: invalid calculator", i)
		}
		cAgg.calculators = append(cAgg.calculators, c)
	}

	return cAgg, nil

}

// TaxPayable returns the sum of payable tax from the underlying calculators
func (agg *Aggregator) TaxPayable() (float64, []core.TaxCredit) {

	var (
		taxAgg float64
		crAgg  []core.TaxCredit
	)

	for _, c := range agg.calculators {
		taxPayable, credits := c.TaxPayable()
		taxAgg += taxPayable
		crAgg = append(crAgg, credits...)
	}

	return taxAgg, crAgg
}

// SetFinances sets the given finances in all underlying tax calculators
func (agg *Aggregator) SetFinances(f *finance.IndividualFinances) {
	for _, c := range agg.calculators {
		c.SetFinances(f)
	}
}

// SetCredits sets the given credits in all underlying tax calculators
func (agg *Aggregator) SetCredits(credits []core.TaxCredit) {
	for _, c := range agg.calculators {
		c.SetCredits(credits)
	}
}
