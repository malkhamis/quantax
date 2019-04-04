package rrsp

import "github.com/malkhamis/quantax/calc/finance"

// Formula computes the max contribution room acquired for the given income
type Formula interface {
	// Contribution returns the max contribution room acquired
	// given then income
	Contribution(income float64) float64
	// AllowedIncomeSources returns the sources that add to the
	// RRSP contribution room
	AllowedIncomeSources() []finance.IncomeSource
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() Formula
}
