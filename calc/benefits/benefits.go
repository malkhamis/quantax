// Package benefits provides implementation for various benefit interfaces
// as defined in package calc
package benefits

import "github.com/malkhamis/quantax/calc"

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(income float64, first calc.Person, others ...calc.Person) float64
	// Validate checks if the formula is valid for use
	Validate() error
	// TODO: embed Formula
}
