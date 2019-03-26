package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitFormula = (*MaxReducerFormula)(nil)

// MaxReducerFormula computes Canada Child Benefits amounts for child(ren)
type MaxReducerFormula struct {
	// the [min, max] dollar amounts for given age groups
	BenefitClasses []AgeGroupBenefits
	// the sub-formulas to reduce the maximum benefits. Step numbers
	// indicate the number of children, where zero means no children
	BenefitReducer Reducer
}

// Apply returns the total annual benefits for the children given the income
func (mrf *MaxReducerFormula) Apply(income float64, first calc.Person, others ...calc.Person) float64 {

	maxBenefits := mrf.maxAnnualAmount(first)
	for _, child := range others {
		maxBenefits += mrf.maxAnnualAmount(child)
	}

	minBenefits := mrf.minAnnualAmount(first)
	for _, child := range others {
		minBenefits += mrf.minAnnualAmount(child)
	}

	childCount := float64(len(others) + 1)
	reduction := mrf.BenefitReducer.Reduce(income, childCount)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (mrf *MaxReducerFormula) Validate() error {

	for _, ageGroupBenefit := range mrf.BenefitClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if err := mrf.BenefitReducer.Validate(); err != nil {
		return errors.Wrap(err, "invalid step reducers")
	}

	return nil
}

// maxAnnualAmount returns the maximum annual benefits for the given child
func (mrf *MaxReducerFormula) maxAnnualAmount(child calc.Person) float64 {

	maxPayments := make(payments, 12)

	for month := range maxPayments {
		for _, ageGroup := range mrf.BenefitClasses {

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
func (mrf *MaxReducerFormula) minAnnualAmount(child calc.Person) float64 {

	minPayments := make(payments, 12)

	for month := range minPayments {
		for _, ageGroup := range mrf.BenefitClasses {

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
