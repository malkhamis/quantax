package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.TaxCalculator = (*Aggregator)(nil)

// Aggregator is used to aggregate payable tax from multiple tax calculators
type Aggregator struct {
	calculators []core.TaxCalculator
}

// NewAggregator returns a new tax aggregator for the given tax calculators. If
// one calculator is nil, it returns ErrNoCalc. If the tax years of all tax
// calculators do not match, it returns ErrTooManyYears
func NewAggregator(c0, c1 core.TaxCalculator, extras ...core.TaxCalculator) (*Aggregator, error) {

	cAgg := &Aggregator{
		calculators: make([]core.TaxCalculator, 0, len(extras)+2),
	}

	allCalcs := append([]core.TaxCalculator{c0, c1}, extras...)
	years := map[uint]struct{}{}

	for i, c := range allCalcs {

		if c == nil {
			return nil, errors.Wrapf(ErrNoCalc, "index %d: invalid calculator", i)
		}
		years[c.Year()] = struct{}{}
		cAgg.calculators = append(cAgg.calculators, c)
	}

	if len(years) != 1 {
		return cAgg, errors.Wrap(ErrTooManyYears, "all calculators' tax year must be the same")
	}

	return cAgg, nil

}

// Year returns the tax year for which this calculator is configured
func (agg *Aggregator) Year() uint {
	// it was  established that all calculators
	// have the same tax year at constructor
	return agg.calculators[0].Year()
}

// Regions returns the tax regions for which this calculator is configured
func (agg *Aggregator) Regions() []core.Region {

	regions := make([]core.Region, 0, len(agg.calculators))
	for _, c := range agg.calculators {
		regions = append(regions, c.Regions()...)
	}
	return regions
}

// SetFinances sets the given finances in all underlying tax calculators
func (agg *Aggregator) SetFinances(f core.HouseholdFinances, credits []core.TaxCredit) {
	for _, c := range agg.calculators {
		c.SetFinances(f, credits)
	}
}

// SetDependents sets the dependents which the calculator might use for tax-
// related calculations
func (agg *Aggregator) SetDependents(deps ...*human.Person) {
	for _, c := range agg.calculators {
		c.SetDependents(deps...)
	}
}

// TaxPayable returns the sum of payable tax from the underlying calculators
func (agg *Aggregator) TaxPayable() (spouseA, spouseB float64, unusedCredits []core.TaxCredit) {

	var (
		taxAggA float64
		taxAggB float64
		crAgg   []core.TaxCredit
	)

	for _, c := range agg.calculators {
		taxA, taxB, credits := c.TaxPayable()
		taxAggA += taxA
		taxAggB += taxB
		crAgg = append(crAgg, credits...)
	}

	return taxAggA, taxAggB, crAgg
}
