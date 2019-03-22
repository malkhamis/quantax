// Package calc defines interfaces for various tax-related calculators
package calc

// Formula represents a method for doing calculations on numeric parameters
type Formula interface {
	// Apply applies the formula on numeric parameters and returns the results
	Apply() float64
	// Validate checks if the formula is valid for use
	Validate() error
}

// TaxFormula represents a method for calculating tax on income
type TaxFormula interface {
	// Update sets the financial numbers which the tax will be calculated for.
	// Users may call this method to set the financial numbers to anything other
	// than what the tax formula was initialized with
	Update(IndividualFinances)
	Formula
}

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Update sets the financial numbers which benefits will be calculated for.
	// Users may call this method to set the financial numbers to anything other
	// than what the child benefit formula was initialized with
	Update(FamilyFinances)
	Formula
}

// TaxCalculator is used to calculate payable tax. Constructors for types
// implementing this interface should typically accept 'Finances' as argument
type TaxCalculator interface {
	// Calc returns the payable amount of tax set/initialized in this calculator
	Calc(taxCredits ...float64) float64
	// UpdateFinances sets the financial numbers which the tax will be calculated
	// for in subsequent calls to Calc(). Users may call this method to set the
	// financial numbers to anything other than what the calculator was
	// initialized with
	UpdateFinances(IndividualFinances)
	// UpdateFormula sets the tax formula used in this calculator. Users may call
	// this method to set the formula to anything other than what the calculator
	// was initialized with
	UpdateFormula(TaxFormula)
}

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children. Constructors for types implementing this
// interface should typically accept 'FamilyFinances' as argument
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
	UpdateBeneficiaries([]Person)
	// UpdateForumla sets the formula for calculating the amount of benefits for
	// children given family finances. Users may call this method to set the
	// formula to anything other than what the calculator was initialized with
	UpdateForumla(ChildBenefitFormula)
}
