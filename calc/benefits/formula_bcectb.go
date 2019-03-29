package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ ChildBenefitFormula = (*BCECTBMaxReducer)(nil)

// BCECTBMaxReducer computes British Columbia Early Childhood Tax Benefits
// as a function of income, number of children, and children's ages. The
// formula calculates the maximum entitlement for all children, then the max
// is reduced based on the income, where reduction is calculated according to
// to a single rated-bracket and is amplified by the number of children
type BCECTBMaxReducer struct {
	// the [min, max] dollar amounts for given age groups (bound-inclusive)
	BenefitClasses []AgeGroupBenefits
	// Reducer is the sub-formula used to reduce the maximum benefits
	ReducerFormula calc.WeightedBracketFormula
}

// Apply returns the total annual benefits for the children given the income
func (mr *BCECTBMaxReducer) Apply(income float64, children ...calc.Person) float64 {

	if len(children) == 0 {
		return 0.0
	}

	var maxBenefits float64
	for _, child := range children {
		maxBenefits += multiAgeGroupBenefits(
			mr.BenefitClasses,
		).MaxAnnualAmount(
			child,
		)
	}

	var minBenefits float64
	for _, child := range children {
		minBenefits += multiAgeGroupBenefits(
			mr.BenefitClasses,
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
func (mr *BCECTBMaxReducer) IncomeCalcMethod() IncomeType {
	// TODO: legislation is not explicit about it
	return AFNI
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (mr *BCECTBMaxReducer) Validate() error {

	for _, ageGroupBenefit := range mr.BenefitClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if mr.ReducerFormula == nil {
		return calc.ErrNoFormula
	}

	if err := mr.ReducerFormula.Validate(); err != nil {
		return errors.Wrap(err, "invalid reducer formula")
	}

	return nil
}

// Clone returns a copy of this instance
func (mr *BCECTBMaxReducer) Clone() ChildBenefitFormula {

	clone := &BCECTBMaxReducer{
		BenefitClasses: make([]AgeGroupBenefits, len(mr.BenefitClasses)),
		ReducerFormula: mr.ReducerFormula.Clone(),
	}

	copy(clone.BenefitClasses, mr.BenefitClasses)

	return clone
}
