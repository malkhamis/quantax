// Package calc defines interfaces for various tax-related calculators
package calc

// TaxCalculator is used to calculate payable tax.
type TaxCalculator interface {
	// Calc returns the payable amount of tax on the income in this calculator
	Calc() float64
	// UpdateFinances sets the financial numbers which the tax will be calculated
	// for in subsequent calls to Calc(). Users may call this method to set the
	// financial numbers to anything other than what the calculator was
	// initialized with
	UpdateFinances(IndividualFinances)
}

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children.
type ChildBenefitCalculator interface {
	// Calc returns the recievable amount of child benefits
	Calc() float64
	// UpdateFinances sets the financial numbers which the calculator will use in
	// subsequent calls to Calc(). Users may call this method to set financial
	// numbers to anything other than what the calculator was initialized with
	UpdateFinances(FamilyFinances)
	// UpdateBeneficiary sets the child which the calculator will use in
	// subsequent calls to Calc(). Users may call this method to set beneficiary
	// to anything other than what the calculator was initialized with
	UpdateBeneficiaries(...Person)
}
