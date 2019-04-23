package factory

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/rrsp"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

// RRSPFactoryConfig is used to pass configs for creating new RRSP factory
type RRSPFactoryConfig struct {
	Year       uint
	RRSPRegion core.Region
	TaxRegions []core.Region
}

// RRSPFactory is a type used to conveniently create RRSP calculators
type RRSPFactory struct {
	newCalculator func() (core.RRSPCalculator, error)
	taxFactory    *TaxFactory
}

// RRSPFactory returns a new RRSP calculator factory from the given config.
// If more than a single tax region is specified, the underlying RRSP calculator
// will use a tax aggregator
func NewRRSPFactory(config RRSPFactoryConfig) *RRSPFactory {

	calcFactory := &RRSPFactory{
		taxFactory: NewTaxFactory(config.Year, config.TaxRegions...),
	}

	foundParams, err := history.GetRRSPParams(config.Year, config.RRSPRegion)
	if err != nil {
		calcFactory.setFailingConstructor(errors.Wrap(err, "RRSP formula"))
		return calcFactory
	}

	calcFactory.initConstructor(foundParams)
	return calcFactory
}

// NewCalculator creates a new RRSP calculator that is configured with params
// set in this factory
func (f *RRSPFactory) NewCalculator() (core.RRSPCalculator, error) {
	if f.newCalculator == nil {
		return nil, ErrFactoryNotInit
	}
	return f.newCalculator()
}

// setFailingConstructor makes calls to NewCalculator returns nil, wrapped(err)
func (f *RRSPFactory) setFailingConstructor(err error) {
	f.newCalculator = func() (core.RRSPCalculator, error) {
		return nil, errors.Wrap(err, "RRSP factory error")
	}
}

// initConstructor initializes this factory's 'newCalculator' function from the
// given RRSP formula and this factory's internal tax calculator factory
func (f *RRSPFactory) initConstructor(params history.RRSPParams) {

	f.newCalculator = func() (core.RRSPCalculator, error) {
		taxCalc, err := f.taxFactory.NewCalculator()
		if err != nil {
			return nil, err
		}
		cfg := rrsp.CalcConfig{params.Formula, taxCalc}
		return rrsp.NewCalculator(cfg)
	}

}
