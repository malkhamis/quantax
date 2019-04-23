package factory

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/income"
	"github.com/malkhamis/quantax/core/tax"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

// TaxFactory is a type used to conveniently create tax calculators
type TaxFactory struct {
	newCalculator func() (core.TaxCalculator, error)
}

// NewTaxFactory returns a new tax calculator factory from the given params. If
// multiple regions are specified, the returned calculator aggregates the taxes
// for all the given regions
// TODO: make it mandatory to give one region
func NewTaxFactory(year uint, regions ...core.Region) *TaxFactory {

	calcFactory := &TaxFactory{}
	allParams := make([]history.TaxParams, len(regions))
	for i, region := range regions {

		foundParams, err := history.GetTaxParams(year, region)
		if err != nil {
			calcFactory.setFailingConstructor(
				errors.Wrapf(err, "tax formula for region %q", region),
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
func (f *TaxFactory) NewCalculator() (core.TaxCalculator, error) {
	if f.newCalculator == nil {
		return nil, ErrFactoryNotInit
	}
	return f.newCalculator()
}

// setFailingConstructor makes calls to NewCalculator returns nil, err
func (f *TaxFactory) setFailingConstructor(err error) {
	f.newCalculator = func() (core.TaxCalculator, error) {
		return nil, errors.Wrap(err, "tax factory error")
	}
}

// initConstructor initializes this factory's 'newCalculator' function from the
// given formulas
func (f *TaxFactory) initConstructor(allParams ...history.TaxParams) {

	switch len(allParams) {
	case 0:
		f.newCalculator = func() (core.TaxCalculator, error) {
			return tax.NewCalculator(tax.CalcConfig{})
		}

	case 1:
		f.newCalculator = func() (core.TaxCalculator, error) {
			incomeCalc, err := income.NewCalculator(allParams[0].IncomeRecipe)
			if err != nil {
				return nil, errors.Wrap(err, "error creating income calculator")
			}
			cfg := tax.CalcConfig{
				IncomeCalc:       incomeCalc,
				TaxFormula:       allParams[0].Formula,
				ContraTaxFormula: allParams[0].ContraFormula,
			}
			return tax.NewCalculator(cfg)
		}

	default:
		f.newCalculator = func() (core.TaxCalculator, error) {
			taxCalcs := make([]core.TaxCalculator, len(allParams))

			for i, p := range allParams {
				incomeCalc, err := income.NewCalculator(p.IncomeRecipe)
				if err != nil {
					return nil, errors.Wrap(err, "error creating income calculator")
				}
				cfg := tax.CalcConfig{
					IncomeCalc:       incomeCalc,
					TaxFormula:       p.Formula,
					ContraTaxFormula: p.ContraFormula,
				}
				taxCalcs[i], err = tax.NewCalculator(cfg)
				if err != nil {
					return nil, errors.Wrap(err, "error creating child benefit calculator")
				}
			}
			return tax.NewAggregator(taxCalcs[0], taxCalcs[1], taxCalcs[2:]...)
		}
	}

}
