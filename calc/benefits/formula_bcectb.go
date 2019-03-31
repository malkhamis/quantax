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
	// IncomeType the income calculation method this formula expects
	IncomeType finance.IncomeType
}

// Apply returns the total annual benefits for the children given the income
func (mr *BCECTBMaxReducer) Apply(income float64, children ...human.Person) float64 {

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
	reduction := float64(childCount) * mr.ReducerFormula.Apply(income)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

// IncomeCalcMethod returns the type of income this formula expects
func (mr *BCECTBMaxReducer) IncomeCalcMethod() finance.IncomeType {
	// TODO: legislation is not explicit about it
	return mr.IncomeType
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
		BeneficiaryClasses: make([]AgeGroupBenefits, len(mr.BeneficiaryClasses)),
		ReducerFormula:     mr.ReducerFormula.Clone(),
	}

	copy(clone.BeneficiaryClasses, mr.BeneficiaryClasses)

	return clone
}
