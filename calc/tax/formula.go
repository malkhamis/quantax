package tax

import "github.com/malkhamis/quantax/calc/finance"

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(income float64) float64
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

type CanadianFormula finance.WeightedBrackets

// Apply applies this formula on the given income and returns the payable tax
func (ct CanadianFormula) Apply(income float64) float64 {
	return finance.WeightedBrackets(ct).Apply(income)
}

// Clone returns a copy of this formula
func (ct CanadianFormula) Clone() Formula {
	clone := finance.WeightedBrackets(ct).Clone()
	return CanadianFormula(clone)
}

// Validate ensures that this formula is valid for use
func (ct CanadianFormula) Validate() error {
	return finance.WeightedBrackets(ct).Validate()
}
