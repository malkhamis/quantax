package tax

import "github.com/malkhamis/quantax/calc/finance"

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(income float64) float64
	// TODO
	ExcludedNetIncomeSources() ([]finance.IncomeSource, []finance.DeductionSource)
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

type CanadianFormula struct {
	// TODO
	ExcludedIncome []finance.IncomeSource
	// TODO
	ExcludedDeductions []finance.DeductionSource
	finance.WeightedBrackets
}

// Apply applies this formula on the given income and returns the payable tax
func (ct *CanadianFormula) Apply(income float64) float64 {
	return ct.WeightedBrackets.Apply(income)
}

// TODO
func (ct *CanadianFormula) ExcludedNetIncomeSources() ([]finance.IncomeSource, []finance.DeductionSource) {
	return ct.ExcludedIncome, ct.ExcludedDeductions
}

// Clone returns a copy of this formula
func (ct *CanadianFormula) Clone() Formula {

	clone := &CanadianFormula{
		WeightedBrackets: ct.WeightedBrackets.Clone(),
	}

	if ct.ExcludedIncome != nil {
		clone.ExcludedIncome = make([]finance.IncomeSource, len(ct.ExcludedIncome))
		copy(clone.ExcludedIncome, ct.ExcludedIncome)
	}

	if ct.ExcludedDeductions != nil {
		clone.ExcludedDeductions = make([]finance.DeductionSource, len(ct.ExcludedDeductions))
		copy(clone.ExcludedDeductions, ct.ExcludedDeductions)
	}

	return clone
}

// Validate ensures that this formula is valid for use
func (ct *CanadianFormula) Validate() error {
	return ct.WeightedBrackets.Validate()
}
