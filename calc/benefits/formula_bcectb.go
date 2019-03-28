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
		maxBenefits += mr.maxAnnualAmount(child)
	}

	var minBenefits float64
	for _, child := range children {
		minBenefits += mr.minAnnualAmount(child)
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

// maxAnnualAmount returns the maximum annual benefits for the given child
func (mr *BCECTBMaxReducer) maxAnnualAmount(child calc.Person) float64 {

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
func (mr *BCECTBMaxReducer) minAnnualAmount(child calc.Person) float64 {

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
