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
	// SetBeneficiaries sets the children which the calculator will compute the
	// benefits for in subsequent calls to Calc()
	SetBeneficiaries(...human.Person)
	// SetFinances makes subsequent calculations based on the given finances
	SetFinances(Financer)
}

// RRSPCalculator is used to calculate recievable or payable tax on transactions
// related to Registered Retirement Saving Plan (RRSP) accounts
type RRSPCalculator interface {
	// TaxPaid calculates the tax payable upon withdrawal
	TaxPaid(withdrawal float64) float64
	// TaxRefund calculates the refundable tax upon deposit/contribution
	TaxRefund(contribution float64) (float64, error)
	// ContributionEarned calculates the newly acquired contribution room
	ContributionEarned() float64
	// SetFinances makes subsequent calculations based on the given finances
	SetFinances(*IndividualFinances)
}

// TaxCalculator is used to calculate payable tax on earnings
type TaxCalculator interface {
	// TaxPayable returns the payable amount of tax for the set finances.
	// The tax credit represent any amount owed to the tax payer without
	// implications for how they might be used.
	TaxPayable() (float64, []TaxCredit)
	// SetFinances stores the given financial data in the underlying tax
	// calculator. Subsequent calls to other functions are based on the
	// the given finances. Changes to the given finances after calling
	// this function should affect future calculations
	SetFinances(*IndividualFinances)
	// SetCredits stores the given credits in the underlying tax calculator.
	// Subsequent calls to other functions will be influenced by the given tax
	// credits. Treatment of given credits is implementation-specific. Ideally,
	// These given credits are originated by the same tax calculator.
	SetCredits([]TaxCredit)
}

// TaxCredit represents an amount that is owed to the tax payer
type TaxCredit interface {
	// Amount is the amount owed to tax payer
	Amount() float64
	// Source describes the source of this tax credit
	Source() string
}
