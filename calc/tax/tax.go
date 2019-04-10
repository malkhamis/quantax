// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
	ErrNoIncCalc = errors.New("no income calculator given")
	ErrNoCalc    = errors.New("no benefit calculator given")
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
	// Apply applies the contra-formula on the income and the set finances
	Apply(netIncome float64) map[CreditSource]Credit
	// Order returns the order in which credit sources are used to reduce tax.
	// Ideally, the returned slice is a superset of what Apply() might return
	Order() []CreditSource
	// SetFinances makes subsequent calculations based on the given finances.
	// Client Changes to the given finances are reflected on future calls
	SetFinances(*finance.IndividualFinances)
	// Clone returns a copy of this contra-formula
	Clone() ContraFormula
	// Validate checks if the formula is valid for use
	Validate() error
}

// Credit represent an amount that reduces payable tax
type Credit struct {
	IsRefundable bool    // if not used, it is paid back
	Amount       float64 // the credit amount (owed)
}

// CreditSource represents a source of tax credits
type CreditSource int

// credit sources recognized by this package
const (
	BeginCanadaCreditSources CreditSource = iota
	PersonalAmountCA
	EligibleDividendsCA
	EndCanadaCreditSources
)
