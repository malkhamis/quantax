package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ ChildBenefitFormula = (*CCBMaxReducer)(nil)

// CCBMaxReducer computes Canada Child Benefits as a function of adjusted
// family net income (AFNI), number of children, and children's ages. The
// formula calculates the maximum entitlement for all children, then the max
// is reduced based on the income, where reduction is calculated according to
// multi-tier, rated brackets
type CCBMaxReducer struct {
	// the [min, max] dollar amounts for given age groups (bound-inclusive)
	BenefitClasses []AgeGroupBenefits
	// Reducers are used to map amount-reducing formulas to child count,
	// where the index of the formula represents the number of children.
	// If the number of children is greater than the number of formulas,
	// the last formula is used
	Reducers []calc.WeightedBracketFormula
}

// Apply returns the total annual benefits for the children given the income
func (mr *CCBMaxReducer) Apply(income float64, children ...calc.Person) float64 {

	if len(children) == 0 {
		return 0.0
	}

	var maxBenefits float64
	for _, child := range children {
		maxBenefits += mr.maxAnnualAmount(child)
	}

	var minBenefits float64
	for _, child := range children {
		minBenefits += mr.minAnnualAmount(child)
	}

	childCount := len(children)
	reduction := mr.reducerFormula(childCount).Apply(income)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

// IncomeCalcMethod returns the type of income this formula expects
func (mr *CCBMaxReducer) IncomeCalcMethod() IncomeType {
	return AFNI
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (mr *CCBMaxReducer) Validate() error {

	for _, ageGroupBenefit := range mr.BenefitClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if len(mr.Reducers) < 1 {
		return calc.ErrNoFormula
	}

	for _, formula := range mr.Reducers {

		if formula == nil {
			return calc.ErrNoFormula
		}

		if err := formula.Validate(); err != nil {
			return errors.Wrap(err, "invalid reducer")
		}

	}

	return nil
}

// Clone returns a copy of this instance
func (mr *CCBMaxReducer) Clone() ChildBenefitFormula {

	clone := &CCBMaxReducer{
		BenefitClasses: make([]AgeGroupBenefits, len(mr.BenefitClasses)),
		Reducers:       make([]calc.WeightedBracketFormula, len(mr.Reducers)),
	}

	copy(clone.BenefitClasses, mr.BenefitClasses)

	for i, reducer := range mr.Reducers {
		clone.Reducers[i] = reducer.Clone()
	}

	return clone
}

// maxAnnualAmount returns the maximum annual benefits for the given child
func (mr *CCBMaxReducer) maxAnnualAmount(child calc.Person) float64 {

	maxPayments := make(payments, 12)

	for month := range maxPayments {
		for _, ageGroup := range mr.BenefitClasses {

			if ageGroup.IsInAgeGroup(child) {
				maxPayments[month] += ageGroup.AmountsPerMonth.Upper()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit classes
		}
		child.AgeMonths++
	}

	return maxPayments.Total()
}

// minAnnualAmount returns the minimum annual benefits for the given child
func (mr *CCBMaxReducer) minAnnualAmount(child calc.Person) float64 {

	minPayments := make(payments, 12)

	for month := range minPayments {
		for _, ageGroup := range mr.BenefitClasses {

			if ageGroup.IsInAgeGroup(child) {
				minPayments[month] += ageGroup.AmountsPerMonth.Lower()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit age groups
		}
		child.AgeMonths++
	}

	return minPayments.Total()
}

// reducerFormula returns the reduction formula based on the child count
func (mr *CCBMaxReducer) reducerFormula(childCount int) calc.WeightedBracketFormula {

	var reducerIndex int

	if childCount >= len(mr.Reducers) {
		reducerIndex = len(mr.Reducers)
	} else {
		reducerIndex = childCount
	}

	reducerIndex--
	return mr.Reducers[reducerIndex]

}