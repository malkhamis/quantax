// Package tax provides implementations for calc.TaxCalculator
package tax

import (
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
