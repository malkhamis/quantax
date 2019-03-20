// Package calc defines interfaces for various tax-related calculators
package calc

// TaxCalculator is used to calculate payable tax. Constructors for types
// implementing this interface should typically accept 'Finances' as argument
type TaxCalculator interface {
	// Calc returns the payable amount of tax set/initialized in this calculator
	Calc() float64
	// Update sets the financial numbers which the tax will be calculated for in
	// subsequent calls to Calc(). Users may call this method to set the financial
	// numbers to anything other than what the calculator was initialized with
	Update(Finances)
}

// ChildBenefitFormula is a formula used to calculate the child benefit amount
// from the given family finances and children
type ChildBenefitFormula func(FamilyFinances, Children) float64

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
	UpdateBeneficiaries(Children)
	// UpdateRule sets the rule for calculating the amount of benefits for child
	// given family finances. Users may call this method to set the rule to
	// anything other than what the calculator was initialized with
	UpdateForumla(ChildBenefitFormula)
}
