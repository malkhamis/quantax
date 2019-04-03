// package calc defines interfaces for various tax-related calculators
package calc

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
)

// TaxCalculator is used to calculate payable tax on individual earnings
type TaxCalculator interface {
	// Calc returns the payable amount of tax on the income in this calculator
	Calc(*finance.IndividualFinances) float64
}

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children.
type ChildBenefitCalculator interface {
	// Calc returns the recievable amount of child benefits given family finances
	// and the children set in the calculator
	Calc(finance.IncomeDeductor) float64
	// SetBeneficiaries sets the children which the calculator will compute the
	// benefits for in subsequent calls to Calc()
	SetBeneficiaries(...human.Person)
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
	SetFinances(*finance.IndividualFinances)
}
