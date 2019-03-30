// Package factory provides functions that creates financial calculators
package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/history"
)

// ChildBenefitCalcFactory is a type used to conveniently create child benefit
// calculators
type ChildBenefitCalcFactory struct {
	formulas []benefits.ChildBenefitFormula
}

// NewChildBenefitCalcFactory returns a new child benefit calculator factory
// from the given params. If extra regions are specified, the returned
// calculator aggregates the benefits for all beneficiaries
func NewChildBenefitCalcFactory(year uint, region Region, extras ...Region) (*ChildBenefitCalcFactory, error) {

	allRegions := []Region{region}
	allRegions = append(allRegions, extras...)

	c := &ChildBenefitCalcFactory{
		formulas: make([]benefits.ChildBenefitFormula, len(extras)+1),
	}

	for i, r := range allRegions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			return nil, ErrRegionNotExist
		}

		foundFormula, err := history.GetChildBenefitFormula(year, convertedRegion)
		if err != nil {
			return nil, err
		}

		c.formulas[i] = foundFormula
	}

	return c, nil
}

// NewCalculator creates a new child benefit calculator that is configured with
// the the parameters/options set in this factory
func (f *ChildBenefitCalcFactory) NewCalculator() (calc.ChildBenefitCalculator, error) {

	if len(f.formulas) == 0 {
		return nil, ErrFactoryNotInit
	}

	if len(f.formulas) == 1 {
		return benefits.NewChildBenefitCalculator(f.formulas[0])
	}

	return benefits.NewChildBenefitAggregator(f.formulas[0], f.formulas[1], f.formulas[2:]...)
}
