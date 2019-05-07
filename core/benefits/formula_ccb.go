package benefits

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

var _ ChildBenefitFormula = (*CCBMaxReducer)(nil)

// CCBMaxReducer computes Canada Child Benefits as a function of income, number
// of children, and children's ages. The formula calculates the maximum
// entitlement for all children, then the max is reduced based on the income,
// where reduction is calculated according to multi-tier, rated brackets
type CCBMaxReducer struct {
	// the [min, max] dollar amounts for given age groups (bound-inclusive)
	BeneficiaryClasses []AgeGroupBenefits
	// Reducers are used to map amount-reducing formulas to child count,
	// where the index of the formula represents the number of children.
	// If the number of children is greater than the number of formulas,
	// the last formula is used
	Reducers []core.WeightedBrackets
}

// Apply returns the total annual benefits for the children given the net income
func (mr *CCBMaxReducer) Apply(netIncome float64, children ...*human.Person) float64 {

	childCount := getChildCount(children)
	if childCount == 0 {
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
	reduction := mr.reducerFormula(childCount).Apply(netIncome)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (mr *CCBMaxReducer) Validate() error {

	for _, ageGroupBenefit := range mr.BeneficiaryClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if len(mr.Reducers) < 1 {
		return ErrNoFormula
	}

	for _, formula := range mr.Reducers {

		if formula == nil {
			return ErrNoFormula
		}

		if err := formula.Validate(); err != nil {
			return errors.Wrap(err, "invalid reducer")
		}

	}

	return nil
}

// Clone returns a copy of this instance
func (mr *CCBMaxReducer) Clone() ChildBenefitFormula {

	if mr == nil {
		return nil
	}

	clone := &CCBMaxReducer{}

	if mr.Reducers != nil {
		clone.Reducers = make([]core.WeightedBrackets, len(mr.Reducers))
		for i, reducer := range mr.Reducers {
			clone.Reducers[i] = reducer.Clone()
		}
	}

	if mr.BeneficiaryClasses != nil {
		clone.BeneficiaryClasses = make([]AgeGroupBenefits, len(mr.BeneficiaryClasses))
		copy(clone.BeneficiaryClasses, mr.BeneficiaryClasses)
	}

	return clone
}

// reducerFormula returns the reduction formula based on the child count
func (mr *CCBMaxReducer) reducerFormula(childCount int) core.WeightedBrackets {

	var reducerIndex int

	if childCount >= len(mr.Reducers) {
		reducerIndex = len(mr.Reducers)
	} else {
		reducerIndex = childCount
	}

	reducerIndex--
	return mr.Reducers[reducerIndex]
}
