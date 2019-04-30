// Package tax provides implementations for core.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula       = errors.New("no formula given/set")
	ErrNoContraFormula = errors.New("no contra-formula given/set")
	ErrNoIncCalc       = errors.New("no income calculator given")
	ErrInvalidTaxArg   = errors.New("invalid tax arguments")
	ErrNoCreditor      = errors.New("no creditor given/set")
	ErrDupCreditSource = errors.New("duplicate credit sources are not allowed")
)

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(netIncome float64) float64
	// Year is the tax year this contra formula is associated with
	Year() uint
	// Region is the tax region this contra formula is associated with
	Region() core.Region
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

// ContraFormula computes reduction on payable taxes for the given finances
type ContraFormula interface {
	// Apply applies the contra-formula and returns a slice of tax credits
	Apply(*TaxPayer) []*TaxCredit
	// FilterAndSort removes tax credits that are not recognized by this contra-
	// formula and sort the remaining items in a priority-of-use sequence, where
	// the first item has the highest priority of use before the next one
	FilterAndSort([]core.TaxCredit) []core.TaxCredit
	// Clone returns a copy of this contra-formula
	Clone() ContraFormula
	// Year is the tax year this contra formula is associated with
	Year() uint
	// Region is the tax region this contra formula is associated with
	Region() core.Region
	// Validate checks if the formula is valid for use
	Validate() error
}

// CalcConfig is used to pass configurations to create new tax calculator
type CalcConfig struct {
	IncomeCalc       core.IncomeCalculator
	TaxFormula       Formula
	ContraTaxFormula ContraFormula
}

// TODO: validate formual and contra formula are for the same region and year?
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

	if cfg.TaxFormula.Year() != cfg.ContraTaxFormula.Year() {
		return errors.Wrap(ErrInvalidTaxArg, "formula/contra-formula tax year mismatch")
	}

	if cfg.TaxFormula.Region() != cfg.ContraTaxFormula.Region() {
		return errors.Wrap(ErrInvalidTaxArg, "formula/contra-formula tax region mismatch")
	}

	if cfg.IncomeCalc == nil {
		return ErrNoIncCalc
	}

	return nil
}

// TaxPayer represents an individual who pays taxes
type TaxPayer struct {
	// the financial data of the subject tax payer
	Finances core.Financer
	// the net income for tax purposes
	NetIncome float64
	// the financial data of the tax payer's spouse
	SpouseFinances core.Financer
	// the net income of the spouse if applicable
	SpouseNetIncome float64
	// Dependents the dependents of the tax payer
	Dependents []*human.Person
}
