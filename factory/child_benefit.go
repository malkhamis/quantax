package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
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
	formulas := make([]benefits.ChildBenefitFormula, len(regions))
	for i, r := range regions {

		convertedRegion, ok := knownRegions[r]
		if !ok {
			calcFactory.setFailingConstructor(
				errors.Wrapf(ErrRegionNotExist, "child benefit region %q", r),
			)
			return calcFactory
		}

		foundFormula, err := history.GetChildBenefitFormula(year, convertedRegion)
		if err != nil {
			calcFactory.setFailingConstructor(
				errors.Wrapf(err, "child benefit formula for region %q", r),
			)
			return calcFactory
		}

		formulas[i] = foundFormula
	}

	calcFactory.initConstructor(formulas...)
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
func (f *ChildBenefitFactory) initConstructor(formulas ...benefits.ChildBenefitFormula) {

	switch {

	case len(formulas) == 0:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			return benefits.NewChildBenefitCalculator(nil)
		}

	case len(formulas) == 1:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			return benefits.NewChildBenefitCalculator(formulas[0])
		}

	default:
		f.newCalculator = func() (calc.ChildBenefitCalculator, error) {
			return benefits.NewChildBenefitAggregator(formulas[0], formulas[1], formulas[2:]...)
		}
	}
}
