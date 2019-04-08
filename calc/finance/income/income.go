package income

import "github.com/malkhamis/quantax/calc/finance"

type Formula interface {
	// IncomeAdjusters returns the income sources that require adjustment before
	// incorporating them in the net income calculation
	IncomeAdjusters() map[finance.IncomeSource]Adjuster
	// DeductionAdjusters returns the deduction sources that require adjustment
	// before incorporating them in the net income
	DeductionAdjusters() map[finance.DeductionSource]Adjuster
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}
