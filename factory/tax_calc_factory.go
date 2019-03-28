// Package factory provides functions that creates financial calculators
package factory

import (
	"errors"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
)

// TaxCalcFactory is a type used to conveniently create tax calculators
type TaxCalcFactory struct {
	year    uint
	region  history.Jurisdiction
	formula tax.Formula
}

// NewTaxCalcFactory returns a new tax calculator factory from the given options
func NewTaxCalcFactory(opts Options) (*TaxCalcFactory, error) {

	region, ok := knownRegions[opts.Region]
	if !ok {
		return nil, errors.New("unkown region")
	}

	foundFormula, err := history.GetTaxFormula(opts.Year, region)
	if err != nil {
		return nil, err
	}

	c := &TaxCalcFactory{
		year:    opts.Year,
		region:  region,
		formula: foundFormula,
	}

	return c, nil
}

// NewCalculator creates a new tax calculator that is configured with the the
// parameters/options set in this factory
func (f *TaxCalcFactory) NewCalculator(finances calc.IndividualFinances) (calc.TaxCalculator, error) {
	return tax.NewCalculator(finances, f.formula)
}
