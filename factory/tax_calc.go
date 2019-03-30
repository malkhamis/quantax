// Package factory provides functions that creates financial calculators
package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
)

// TaxCalcFactory is a type used to conveniently create tax calculators
type TaxCalcFactory struct {
	formulas []tax.Formula
}

// NewTaxCalcFactory returns a new tax calculator factory from the given
// params. If extra regions are specified, the returned calculator aggregates
// the taxes for all the given regions
func NewTaxCalcFactory(year uint, region Region, extras ...Region) (*TaxCalcFactory, error) {

	allRegions := []Region{region}
	allRegions = append(allRegions, extras...)

	c := &TaxCalcFactory{
		formulas: make([]tax.Formula, len(extras)+1),
	}

	for i, r := range allRegions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			return nil, ErrRegionNotExist
		}

		foundFormula, err := history.GetTaxFormula(year, convertedRegion)
		if err != nil {
			return nil, err
		}

		c.formulas[i] = foundFormula
	}

	return c, nil
}

// NewCalculator creates a new tax calculator that is configured with the the
// parameters/options set in this factory
func (f *TaxCalcFactory) NewCalculator(finances calc.IndividualFinances) (calc.TaxCalculator, error) {

	if len(f.formulas) == 0 {
		return nil, ErrFactoryNotInit
	}

	if len(f.formulas) == 1 {
		return tax.NewCalculator(finances, f.formulas[0])
	}

	return tax.NewCalculatorAgg(finances, f.formulas)
}
