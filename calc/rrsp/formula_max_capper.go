package rrsp

import "github.com/malkhamis/quantax/calc/finance"

// compile-time check for interface implementation
var _ Formula = (*MaxCapper)(nil)

// MaxCapper computes the maximum acquired RRSP contribution on earned income.
// The formula works by computing a percentage of earned income, where the
// result is capped by a maximum value
type MaxCapper struct {
	// the percentage of earned income
	Rate float64
	// the maximum contributable amount
	Cap float64
	// sources that add to contribution room
	IncomeSources []finance.IncomeSource
	// affected income source when making withdrawal
	IncomeSourceForWithdrawal finance.IncomeSource
	// affected deduction source when making contribution
	DeductionSourceForContribution finance.DeductionSource
}

// ContributionEarned returns the max contribution room acquired given the net
// income. It is up to the client to calculate the net income appropriately by
// checking allowed income sources through calling 'AllowedIncomeSources()'
func (mc *MaxCapper) ContributionEarned(netIncome float64) float64 {

	contribution := mc.Rate * netIncome
	if contribution > mc.Cap {
		return mc.Cap
	}
	return contribution
}

// AllowedIncomeSources returns the sources which this formula expects as part
// of the net income when calculating the contribution
func (mc *MaxCapper) AllowedIncomeSources() []finance.IncomeSource {
	return mc.IncomeSources
}

// TargetSourceForWithdrawl returns the affected income source when making a
// withdrawal from an RRSP account
func (mc *MaxCapper) TargetSourceForWithdrawl() finance.IncomeSource {
	return mc.IncomeSourceForWithdrawal
}

// TargetSourceForContribution returns the affected deducion source when making
// contribution to an RRSP account
func (mc *MaxCapper) TargetSourceForContribution() finance.DeductionSource {
	return mc.DeductionSourceForContribution
}

// Validate checks if the formula is valid for use
func (mc *MaxCapper) Validate() error {
	return nil
}

// Clone returns a copy of the formula
func (mc *MaxCapper) Clone() Formula {

	clone := *mc

	if mc.IncomeSources != nil {
		clone.IncomeSources = make([]finance.IncomeSource, len(mc.IncomeSources))
		copy(clone.IncomeSources, mc.IncomeSources)
	}

	return &clone
}
