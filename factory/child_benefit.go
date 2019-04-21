package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/income"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

// ChildBenefitFactory is a type used to conveniently create child benefit
// calculators
type ChildBenefitFactory struct {
	newCalculator func() (calc.ChildBenefitCalculator, error)
}

// NewChildBenefitFactory returns a new child benefit calculator factory
// from the given params. If multiple regions are specified, the returned
// calculator aggregates the benefits for all beneficiaries
func NewChildBenefitFactory(year uint, regions ...Region) *ChildBenefitFactory {

	calcFactory := &ChildBenefitFactory{}
	allParams := make([]history.CBParams, len(regions))
	for i, r := range regions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			calcFactory.setFailingConstructor(
				errors.Wrapf(ErrRegionNotExist, "child benefit region %q", r),
			)
			return calcFactory
		}

		foundParams, err := history.GetChildBenefitParams(year, convertedRegion)
		if err != nil {
			calcFactory.setFailingConstructor(
				errors.Wrapf(err, "child benefit formula for region %q", r),
			)
			return calcFactory
		}

		allParams[i] = foundParams
	}

	calcFactory.initConstructor(allParams...)
	return calcFactory
}

// NewCalculator creates a new child benefit calculator that is configured with
// the params set in this factory
func (f *ChildBenefitFactory) NewCalculator() (calc.ChildBenefitCalculator, error) {
	if f.newCalculator == nil {
		return nil, ErrFactoryNotInit
	}
	return f.newCalculator()
}

// setFailingConstructor makes calls to NewCalculator returns nil, err
func (f *ChildBenefitFactory) setFailingConstructor(err error) {
	f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
		return nil, errors.Wrap(err, "child benefit factory error")
	}
}

// initConstructor initializes this factory's 'newCalculator' function from the
// given formulas
func (f *ChildBenefitFactory) initConstructor(allParams ...history.CBParams) {

	switch {

	case len(allParams) == 0:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			return benefits.NewChildBenefitCalculator(nil, nil)
		}

	case len(allParams) == 1:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			formula, incomeRecipe := allParams[0].Formula, allParams[0].IncomeRecipe
			incomeCalc, err := income.NewCalculator(incomeRecipe)
			if err != nil {
				return nil, errors.Wrap(err, "error creating income calculator")
			}
			return benefits.NewChildBenefitCalculator(formula, incomeCalc)
		}

	default:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			cbCalcs := make([]calc.ChildBenefitCalculator, len(allParams))

			for i, p := range allParams {
				formula, incomeRecipe := p.Formula, p.IncomeRecipe
				incomeCalc, err := income.NewCalculator(incomeRecipe)
				if err != nil {
					return nil, errors.Wrap(err, "error creating income calculator")
				}
				cbCalcs[i], err = benefits.NewChildBenefitCalculator(formula, incomeCalc)
				if err != nil {
					return nil, errors.Wrap(err, "error creating child benefit calculator")
				}
			}

			return benefits.NewChildBenefitAggregator(cbCalcs[0], cbCalcs[1], cbCalcs[2:]...)
		}
	}
}
