// Package tax provides implementations for core.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
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
	ErrInvalidTaxInfo      = errors.New("invalid tax info") // TODO better name
)

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(netIncome float64) float64
	// Clone returns a copy of this formula
	Clone() Formula
	// TaxInfo: TODO
	TaxInfo() core.TaxInfo
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
	// TaxInfo: TODO
	TaxInfo() core.TaxInfo
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

	taxInfoTF := cfg.TaxFormula.TaxInfo()
	taxInfoCTF := cfg.ContraTaxFormula.TaxInfo()
	if taxInfoTF != taxInfoCTF {
		return errors.Wrapf(ErrInvalidTaxInfo, "%v != %v", taxInfoTF, taxInfoCTF)
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
