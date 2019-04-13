// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula           = errors.New("no formula given/set")
	ErrNoContraFormula     = errors.New("no contra-formula given/set")
	ErrNoCreditor          = errors.New("no creditor given/set")
	ErrDupCreditSource     = errors.New("duplicates are not allowed")
	ErrUnknownCreditSource = errors.New("unknown credit source")
	ErrNoIncCalc           = errors.New("no income calculator given")
	ErrNoCalc              = errors.New("no benefit calculator given")
)

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(netIncome float64) float64
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

// ContraFormula computes reduction on payable taxes for the given finances
type ContraFormula interface {
	// Apply applies the contra-formula and returns a slice of Credits that is
	// sorted in a priority-of-use sequence, where the first item has the highest
	// priority of use before the next item
	Apply(finances *finance.IndividualFinances, netIncome float64) []Credits
	// Clone returns a copy of this contra-formula
	Clone() ContraFormula
	// Validate checks if the formula is valid for use
	Validate() error
}

// CalcConfig is used to pass configurations to create new tax calculator
type CalcConfig struct {
	IncomeCalc       calc.IncomeCalculator
	TaxFormula       Formula
	ContraTaxFormula ContraFormula
}

// validate checks if the configurations are valid for use by calc constructors
func (cfg CalcConfig) validate() error {

	if cfg.TaxFormula == nil {
		return ErrNoFormula
	}

	err := cfg.TaxFormula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula")
	}

	if cfg.ContraTaxFormula == nil {
		return ErrNoContraFormula
	}

	err = cfg.ContraTaxFormula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid contra-formula")
	}

	if cfg.IncomeCalc == nil {
		return ErrNoIncCalc
	}

	return nil
}
