package rrsp

import "github.com/malkhamis/quantax/calc/finance"

// compile-time check for interface implementation
var _ Formula = (*MaxCapper)(nil)

// MaxCapper computes the maximum acquired RRSP contribution on earned income.
// The formula works by computing a percentage of earned income, where the
// result is capped by a maximum value
type MaxCapper struct {
	Rate       float64            // the percentage of earned income
	Cap        float64            // the maximum contributable amount
	IncomeType finance.IncomeType // The method of calculating income
}

// Contribution returns the max contribution room acquired given then income
func (mc *MaxCapper) Contribution(income float64) float64 {

	contribution := mc.Rate * income
	if contribution > mc.Cap {
		return mc.Cap
	}
	return contribution
}

// IncomeCalcMethod returns the method of calculating the income
func (mc *MaxCapper) IncomeCalcMethod() finance.IncomeType {
	return mc.IncomeType
}

// Validate checks if the formula is valid for use
func (mc *MaxCapper) Validate() error {
	return nil
}

// Clone returns a copy of the formula
func (mc *MaxCapper) Clone() Formula {
	clone := *mc
	return &clone
}
