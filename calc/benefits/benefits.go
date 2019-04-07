// Package benefits implements benefit calculators' interfaces defined by
// package calc
package benefits

import (
	"errors"

	"github.com/malkhamis/quantax/calc/human"
)

// Sentinel errors that can be wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
	ErrNoIncCalc = errors.New("no income calculator given")
	ErrNoCalc    = errors.New("no benefit calculator given")
)

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(netIncome float64, children ...human.Person) float64
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() ChildBenefitFormula
}
