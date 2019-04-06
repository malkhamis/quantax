package income

import "github.com/malkhamis/quantax/calc/finance"

// compile-time check for interface implementatino
var (
	_ Adjuster = CanadianCapitalGainAdjuster{}
)

// Adjuster is a type that adjusts any amount given finances. The type is
/// mainly used by income calculators to adjust specific income/deductions
type Adjuster interface {
	// Adjusted returns an adjusted amount from the given finances
	Adjusted(finance.IncomeDeductor) float64
	// Clone returns a copy of this adjuster
	Clone() Adjuster
}

// CanadianCapitalGainAdjuster returns the adjusted income for Canadian-sourced
// capital gain income
type CanadianCapitalGainAdjuster struct {
	// the proportion of the capital gain that is considered income
	Proportion float64
}

// Adjusted returns the adjusted amount of capital gain income by returning the
// adjusted amount as follows:
//  AdjustedAmount = (Proportion) x (Capital Gain Income)
// If the given finances is nil, it returns 0.0
func (cg CanadianCapitalGainAdjuster) Adjusted(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	income := finances.TotalIncome(finance.IncSrcCapitalGainCA)
	return income * cg.Proportion
}

// Clone returns a copy of this instance
func (cg CanadianCapitalGainAdjuster) Clone() Adjuster {
	return cg
}
