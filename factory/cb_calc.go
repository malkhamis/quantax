package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/history"
)

// ChildBenefitCalcFactory is a type used to conveniently create child benefit
// calculators
type ChildBenefitCalcFactory struct {
	year     uint
	region   history.Jurisdiction
	children []calc.Person
	formula  benefits.ChildBenefitFormula
}

// NewChildBenefitCalcFactoryreturns a new child benefit calculator factory
// from the given options for the given children
func NewChildBenefitCalcFactory(opts Options, children ...calc.Person) (*ChildBenefitCalcFactory, error) {

	region, ok := knownRegions[opts.Region]
	if !ok {
		return nil, ErrRegionNotExist
	}

	foundFormula, err := history.GetChildBenefitFormula(opts.Year, region)
	if err != nil {
		return nil, err
	}

	c := &ChildBenefitCalcFactory{
		year:     opts.Year,
		region:   region,
		children: children,
		formula:  foundFormula,
	}

	return c, nil
}

// NewCalculator creates a new child benefit calculator that is configured with
// the the parameters/options set in this factory
func (f *ChildBenefitCalcFactory) NewCalculator(finances calc.FamilyFinances) (calc.ChildBenefitCalculator, error) {
	return benefits.NewCBCalculator(f.formula, finances, f.children...)
}
