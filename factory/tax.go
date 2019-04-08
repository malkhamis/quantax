package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance/income"
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
	allParams := make([]history.TaxParams, len(regions))
	for i, r := range regions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			calcFactory.setFailingConstructor(
				errors.Wrapf(ErrRegionNotExist, "tax region %q", r),
			)
			return calcFactory
		}

		foundParams, err := history.GetTaxParams(year, convertedRegion)
		if err != nil {
			calcFactory.setFailingConstructor(
				errors.Wrapf(err, "tax formula for region %q", r),
			)
			return calcFactory
		}

		allParams[i] = foundParams
	}

	calcFactory.initConstructor(allParams...)
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
func (f *TaxFactory) initConstructor(allParams ...history.TaxParams) {

	switch len(allParams) {
	case 0:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			return tax.NewCalculator(nil, nil)
		}

	case 1:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			formula, incomeRecipe := allParams[0].Formula, allParams[0].IncomeRecipe
			incomeCalc, err := income.NewCalculator(incomeRecipe)
			if err != nil {
				return nil, errors.Wrap(err, "error creating income calculator")
			}
			return tax.NewCalculator(formula, incomeCalc)
		}

	default:
		f.newCalculator = func() (calc.TaxCalculator, error) {
			taxCalcs := make([]calc.TaxCalculator, len(allParams))

			for i, p := range allParams {
				formula, incomeRecipe := p.Formula, p.IncomeRecipe
				incomeCalc, err := income.NewCalculator(incomeRecipe)
				if err != nil {
					return nil, errors.Wrap(err, "error creating income calculator")
				}
				taxCalcs[i], err = tax.NewCalculator(formula, incomeCalc)
				if err != nil {
					return nil, errors.Wrap(err, "error creating child benefit calculator")
				}
			}
			return tax.NewAggregator(taxCalcs[0], taxCalcs[1], taxCalcs[2:]...)
		}
	}

}
