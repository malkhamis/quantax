package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

// TaxFactory is a type used to conveniently create tax calculators
type TaxFactory struct {
	newCalculator func() (calc.TaxCalculator, error)
}

// NewTaxFactory returns a new tax calculator factory from the given params. If
// multiple regions are specified, the returned calculator aggregates the taxes
// for all the given regions
func NewTaxFactory(year uint, regions ...Region) *TaxFactory {

	calcFactory := &TaxFactory{}
	formulas := make([]tax.Formula, len(regions))
	for i, r := range regions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			calcFactory.setFailingConstructor(
				errors.Wrapf(ErrRegionNotExist, "tax region %q", r),
			)
			return calcFactory
		}

		foundFormula, err := history.GetTaxFormula(year, convertedRegion)
		if err != nil {
			calcFactory.setFailingConstructor(
				errors.Wrapf(err, "tax formula for region %q", r),
			)
			return calcFactory
		}

		formulas[i] = foundFormula
	}

	calcFactory.initConstructor(formulas...)
	return calcFactory
}

// NewCalculator creates a new tax calculator that is configured with the params
// set in this factory
func (f *TaxFactory) NewCalculator() (calc.TaxCalculator, error) {
	if f.newCalculator == nil {
		return nil, ErrFactoryNotInit
	}
	return f.newCalculator()
}

// setFailingConstructor makes calls to NewCalculator returns nil, err
func (f *TaxFactory) setFailingConstructor(err error) {
	f.newCalculator = func() (calc.TaxCalculator, error) {
		return nil, errors.Wrap(err, "tax factory error")
	}
}

// initConstructor initializes this factory's 'newCalculator' function from the
// given formulas
func (f *TaxFactory) initConstructor(formulas ...tax.Formula) {

	switch len(formulas) {
	case 0:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			return tax.NewCalculator(nil)
		}

	case 1:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			return tax.NewCalculator(formulas[0])
		}

	default:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			return tax.NewAggregator(formulas[0], formulas[1], formulas[2:]...)
		}
	}
}
