// package core defines interfaces for various tax-related calculators. It also
// contains implementations of said interfaces as sub-packages. In short, this
// package and its subpackages provide the basic building blocks that can be
// used to conduct quantitative tax analysis or as standalone
package core

import "github.com/malkhamis/quantax/core/human"

// IncomeCalculator is used to calculate income used for benefits and tax
type IncomeCalculator interface {
	// TotalIncome returns the total adjusted income
	TotalIncome() float64
	// TotalDeductions returns total adjusted deductions
	TotalDeductions() float64
	// NetIncome returns total income less total deduction
	NetIncome() float64
	// SetFinances makes subsequent calculations based on the given finances
	SetFinances(Financer)
}

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children.
type ChildBenefitCalculator interface {
	// Calc returns the recievable amount of child benefits for the given
	// finances and the children set in the calculator
	Calc() float64
	// SetFinances makes subsequent calculations based on the given finances
	SetFinances(HouseholdFinances)
	// SetBeneficiaries sets the children which the calculator will compute the
	// benefits for in subsequent calls to Calc()
	SetBeneficiaries(...*human.Person)
}

// RRSPCalculator is used to calculate recievable or payable tax on transactions
// related to Registered Retirement Saving Plan (RRSP) accounts
type RRSPCalculator interface {
	// TaxPaid calculates the tax payable on the amount withdrawn. It should
	// return the tax credits after the withdrawal
	TaxPaid(withdrawal float64) (float64, []TaxCredit)
	// TaxRefund calculates the refundable tax upon deposit/contribution. It
	// should return the tax credits after the contribution
	TaxRefund(contribution float64) (float64, []TaxCredit)
	// ContributionEarned calculates the newly acquired contribution room
	ContributionEarned() float64
	// SetFinances stores the given financial data in the underlying tax
	// calculator. Subsequent calls to other functions are based on the
	// the given finances. Changes to the given finances after calling
	// this function should affect future calculations
	SetFinances(HouseholdFinances, []TaxCredit)
	// SetDependents sets the dependents which the calculator might use for tax-
	// related calculations
	SetDependents([]*human.Person)
	// SetTargetSpouseA makes subsequent calls based on spouse A finances
	SetTargetSpouseA()
	// SetTargetSpouseB makes subsequent calls based on spouse B finances
	SetTargetSpouseB()
}

// TaxCalculator is used to calculate payable tax on earnings
type TaxCalculator interface {
	// TaxPayable returns the payable amount of tax for the set finances.
	// The tax credit represent any amount owed to the tax payer without
	// implications for how they might be used.
	TaxPayable() (spouseA, spouseB float64, combinedCredits []TaxCredit)
	// SetFinances stores the given financial data in the underlying tax
	// calculator. Subsequent calls to other functions are based on the
	// the given finances. Changes to the given finances after calling
	// this function should affect future calculations
	SetFinances(HouseholdFinances, []TaxCredit)
	// SetDependents sets the dependents which the calculator might use for tax-
	// related calculations
	SetDependents([]*human.Person)
	// TaxYear returns the tax year of the calculator
	Year() uint
	// Regions is the tax regions of the calculator. The underlying implementation
	// may compute the taxes for multiple regions
	Regions() []Region
}

// TaxCredit represents an amount that is owed to the tax payer
type TaxCredit interface {
	// SetAmounts sets the initial, used, and remaining abouts of this tax credit
	SetAmounts(initial, used, remaining float64)
	// Amounts returns the tax credit amounts
	Amounts() (initial, used, remaining float64)
	// Rule returns the rule of how this tax credit can be used
	Rule() CreditRule
	// ReferenceFinancer returns the owner of this tax credit
	ReferenceFinancer() Financer
	// Source is the financial source of this tax credit. If the credit is
	// not associated with a specific source, it should return 'SrcNone'
	Source() FinancialSource
	// ShallowCopy returns a copy of this instance, where reference financer
	// is the same as the original instance
	ShallowCopy() TaxCredit
	// Year returns the tax year from which the tax credit was calculated
	Year() uint
	// Region returns the tax region for which the tax credit was calculated
	Region() Region
	// Description is a short description for the reason of the tax credit
	Description() string
}

// CreditRuleType is an enum type for recognized methods of using tax credits
type CreditRuleType int

const (
	// unknown/uninitialized
	_ CreditRuleType = iota
	// CrRuleTypeCashable indicates cashable credits
	// that may result in negative payable tax amount
	CrRuleTypeCashable
	// CrRuleTypeCanCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero
	// and the remaining balance may be carried forward
	// to the future
	CrRuleTypeCanCarryForward
	// CrRuleTypeNotCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero.
	// Unused balance cannot be carried forward to the future
	CrRuleTypeNotCarryForward
)

// CreditRule associate a credit source with a method of using its credit amount
type CreditRule struct {
	// the name of the credit source
	CrSource string
	// the way of using the credit source
	Type CreditRuleType
}
