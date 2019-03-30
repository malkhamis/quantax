// Package benefits implements benefit calculators' interfaces defined by
// package calc
package benefits

import (
	"errors"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
)

// Sentinel errors that can be wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
)

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(income float64, children ...calc.Person) float64
	// IncomeCalcMethod returns the method of calculating the income
	IncomeCalcMethod() finance.IncomeType
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() ChildBenefitFormula
}
