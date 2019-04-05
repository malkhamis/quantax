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
	// AdjustableIncomeSources returns the income sources that trigger
	// income source and payable tax adjustments
	AdjustableIncomeSources() map[finance.IncomeSource]Adjuster
	// AdjustableDeductionSources returns the deduction sources that
	// trigger deduction source and payable tax adjustments
	AdjustableDeductionSources() map[finance.DeductionSource]Adjuster
	// Clone returns a copy of this formula
	Clone() Formula
	// Validate checks if the formula is valid for use
	Validate() error
}

// CanadianFormula is used to calculate Canadian federal and provincial taxes
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
	// IncomeAdjusters specify income sources that should be adjusted
	// to accurately calculate the tax. The formula uses these fields
	// to communicate to the client about how net income should be
	// calculated, but the formula itself won't use them at all for
	// calculating the payable tax
	IncomeAdjusters map[finance.IncomeSource]Adjuster
	// DeductionAdjusters specify deduction sources that should be
	// adjusted to accurately calculate the tax. The formula uses these
	// fields to communicate to the client about how net income should
	// be calculated, but the formula itself won't use them at all for
	// calculating the payable tax
	DeductionAdjusters map[finance.DeductionSource]Adjuster
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

// AdjustableIncomeSources returns the adjuster of a given income source. This
// formula expects the client to call the adjuster before calling Apply()
func (ct *CanadianFormula) AdjustableIncomeSources() map[finance.IncomeSource]Adjuster {
	return ct.IncomeAdjusters
}

// AdjustableDeductionSources returns the adjuster of a given deduction source.
// This formula expects the client to call the adjuster before calling Apply()
func (ct *CanadianFormula) AdjustableDeductionSources() map[finance.DeductionSource]Adjuster {
	return ct.DeductionAdjusters
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

	if ct.IncomeAdjusters != nil {
		clone.IncomeAdjusters = make(map[finance.IncomeSource]Adjuster)
		for source, adjuster := range ct.IncomeAdjusters {
			clone.IncomeAdjusters[source] = adjuster.Clone()
		}
	}

	if ct.DeductionAdjusters != nil {
		clone.DeductionAdjusters = make(map[finance.DeductionSource]Adjuster)
		for source, adjuster := range ct.DeductionAdjusters {
			clone.DeductionAdjusters[source] = adjuster.Clone()
		}
	}

	return clone
}

// Validate ensures that this formula is valid for use
func (ct *CanadianFormula) Validate() error {
	return ct.WeightedBrackets.Validate()
}
