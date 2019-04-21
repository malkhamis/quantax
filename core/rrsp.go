package core

import "github.com/malkhamis/quantax/core/finance"

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
