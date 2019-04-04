package benefits

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"

	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ ChildBenefitFormula = (*BCECTBMaxReducer)(nil)

// BCECTBMaxReducer computes British Columbia Early Childhood Tax Benefits
// as a function of income, number of children, and children's ages. The
// formula calculates the maximum entitlement for all children, then the max
// is reduced based on the income, where reduction is calculated according to
// to a single rated-bracket and is amplified by the number of children
type BCECTBMaxReducer struct {
	// the [min, max] dollar amounts for given age groups (bound-inclusive)
	BeneficiaryClasses []AgeGroupBenefits
	// Reducer is the sub-formula used to reduce the maximum benefits
	ReducerFormula finance.WeightedBrackets
	// ExcludedIncome is the income sources which this formula does not
	// expect to be part of the income passed for calculating benefits.
	// The formula uses these fields to communicate to the client about
	// how net income should be calculated, but the formula itself won't
	// use them at all for calculating child benefit amount
	ExcludedIncome []finance.IncomeSource
	// ExcludedDeductions is the deduction sources which this formula
	// does not expect to be part of the income passed for calculating
	// benefits. The formula uses these fields to communicate to the
	// client about how net income should be calculated, but the formula
	// itself won't use them at all for calculating child benefit amount
	ExcludedDeductions []finance.DeductionSource
}

// Apply returns the total annual benefits for the children given the net
// income. It is up to the client to calculate the net income appropriately
// by checking excluded income and deduction sources through calling methods
// 'ExcludedIncomeSources()' and 'ExcludedDeductionSources()'
func (mr *BCECTBMaxReducer) Apply(netIncome float64, children ...human.Person) float64 {

	if len(children) == 0 {
		return 0.0
	}

	var maxBenefits float64
	for _, child := range children {
		maxBenefits += multiAgeGroupBenefits(
			mr.BeneficiaryClasses,
		).MaxAnnualAmount(
			child,
		)
	}

	var minBenefits float64
	for _, child := range children {
		minBenefits += multiAgeGroupBenefits(
			mr.BeneficiaryClasses,
		).MinAnnualAmount(
			child,
		)
	}

	childCount := len(children)
	reduction := float64(childCount) * mr.ReducerFormula.Apply(netIncome)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

// ExcludedIncomeSources returns the net income sources which this formula expects
// to not be part of the income passed to Apply()
func (mr *BCECTBMaxReducer) ExcludedIncomeSources() []finance.IncomeSource {
	return mr.ExcludedIncome
}

// ExcludedDeductionSources returns the income sources which this formula
// expects to not be part of the net income passed to Apply()
func (mr *BCECTBMaxReducer) ExcludedDeductionSources() []finance.DeductionSource {
	return mr.ExcludedDeductions
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (mr *BCECTBMaxReducer) Validate() error {

	for _, ageGroupBenefit := range mr.BeneficiaryClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if mr.ReducerFormula == nil {
		return ErrNoFormula
	}

	if err := mr.ReducerFormula.Validate(); err != nil {
		return errors.Wrap(err, "invalid reducer formula")
	}

	return nil
}

// Clone returns a copy of this instance
func (mr *BCECTBMaxReducer) Clone() ChildBenefitFormula {

	clone := &BCECTBMaxReducer{
		ReducerFormula: mr.ReducerFormula.Clone(),
	}

	if mr.BeneficiaryClasses != nil {
		clone.BeneficiaryClasses = make([]AgeGroupBenefits, len(mr.BeneficiaryClasses))
		copy(clone.BeneficiaryClasses, mr.BeneficiaryClasses)
	}

	if mr.ExcludedIncome != nil {
		clone.ExcludedIncome = make([]finance.IncomeSource, len(mr.ExcludedIncome))
		copy(clone.ExcludedIncome, mr.ExcludedIncome)
	}

	if mr.ExcludedDeductions != nil {
		clone.ExcludedDeductions = make([]finance.DeductionSource, len(mr.ExcludedDeductions))
		copy(clone.ExcludedDeductions, mr.ExcludedDeductions)
	}

	return clone
}
