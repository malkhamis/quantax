// Package calc defines interfaces for various tax-related calculators
package calc

import "github.com/malkhamis/quantax/calc/finance"

// TaxCalculator is used to calculate payable tax.
type TaxCalculator interface {
	// Calc returns the payable amount of tax on the income in this calculator
	Calc(finance.IndividualFinances) float64
}

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children.
type ChildBenefitCalculator interface {
	// Calc returns the recievable amount of child benefits given family finances
	// and the children set in the calculator
	Calc(finance.FamilyFinances) float64
	// SetBeneficiaries sets the children which the calculator will compute the
	// benefits for in subsequent calls to Calc()
	SetBeneficiaries(...Person)
}
