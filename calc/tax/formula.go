package tax

import "github.com/malkhamis/quantax/calc/finance"

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(netIncome float64) float64
	// ExcludedIncomeSources returns the income sources that should
	// be excluded from the expected total income for the formula
	ExcludedIncomeSources() []finance.IncomeSource
	// ExcludedDeductionSources returns the income sources that should
	// be excluded from the expected total deductions for the formula
	ExcludedDeductionSources() []finance.DeductionSource
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

type CanadianFormula struct {
	// ExcludedIncome is the income sources which this formula does not
	// expect to be part of the income passed for calculating the tax.
	// The formula uses these fields to communicate to the client about
	// how net income should be calculated, but the formula itself won't
	// use them at all for calculating the payable tax
	ExcludedIncome []finance.IncomeSource
	// ExcludedDeductions is the deduction sources which this formula
	// does not expect to be part of the income passed for calculating
	// the tax. The formula uses these fields to communicate to the
	// client about how net income should be calculated, but the formula
	// itself won't use them at all for calculating the payable tax
	ExcludedDeductions []finance.DeductionSource
	finance.WeightedBrackets
}

// Apply applies this formula on the given net income and returns the payable
// tax. It is up to the client to calculate the net income appropriately by
// checking excluded income and deduction sources through calling methods
// 'ExcludedIncomeSources()' and 'ExcludedDeductionSources()'
func (ct *CanadianFormula) Apply(netIncome float64) float64 {
	return ct.WeightedBrackets.Apply(netIncome)
}

// ExcludedIncomeSources returns the income sources which this formula expects
// to not be part of the net income passed to Apply()
func (ct *CanadianFormula) ExcludedIncomeSources() []finance.IncomeSource {
	return ct.ExcludedIncome
}

// ExcludedDeductionSources returns the income sources which this formula
// expects to not be part of the net income passed to Apply()
func (ct *CanadianFormula) ExcludedDeductionSources() []finance.DeductionSource {
	return ct.ExcludedDeductions
}

// Clone returns a copy of this formula
func (ct *CanadianFormula) Clone() Formula {

	if ct == nil {
		return nil
	}

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
