// Package calc defines interfaces for various tax-related calculators
package calc

type Formula interface {
	// Name returns the name of this formula
	Name() string
	// Validate checks if the formula is valid for use
	Validate() error
}

// BasicFormula takes a single numeric input and produces a numeric output
type BasicFormula interface {
	// Apply applies the formula on the given numeric parameter
	Apply(param float64) float64
	// Validate checks if the formula is valid for use
	Validate() error
	// TODO: embed Formula
}

// TaxFormula computes payable taxes on individual income
type TaxFormula interface {
	BasicFormula
}

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(income float64, first Person, others ...Person) float64
	// Validate checks if the formula is valid for use
	Validate() error
	// TODO: embed Formula
}

// TaxCalculator is used to calculate payable tax.
type TaxCalculator interface {
	// Calc returns the payable amount of tax on the income in this calculator
	Calc(taxCredits ...float64) float64
	// UpdateFinances sets the financial numbers which the tax will be calculated
	// for in subsequent calls to Calc(). Users may call this method to set the
	// financial numbers to anything other than what the calculator was
	// initialized with
	UpdateFinances(IndividualFinances)
	// UpdateFormula sets the tax formula used in this calculator. Users may call
	// this method to set the formula to anything other than what the calculator
	// was initialized with
	UpdateFormula(TaxFormula) error
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
	UpdateBeneficiaries(Person, ...Person)
	// UpdateForumla sets the formula for calculating the amount of benefits for
	// children given family finances. Users may call this method to set the
	// formula to anything other than what the calculator was initialized with
	UpdateForumla(ChildBenefitFormula) error
}
