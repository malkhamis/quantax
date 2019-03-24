package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitFormula = (*CCBFormula)(nil)

// CCBFormula computes Canada Child Benefits amounts for child(ren)
type CCBFormula struct {
	// the [min, max] dollar amounts for given age groups
	BenefitClasses []AgeGroupBenefits
	// the sub-formulas to reduce the maximum benefits. Step numbers
	// indicate the number of children, where zero means no children
	BenefitReducer Reducer
}

// Apply returns a 12-month payment schedule for the children given the income
func (cbf *CCBFormula) Apply(income float64, first Child, others ...Child) float64 {

	maxBenefits := cbf.maxAnnualAmount(first)
	for _, child := range others {
		maxBenefits += cbf.maxAnnualAmount(child)
	}

	minBenefits := cbf.minAnnualAmount(first)
	for _, child := range others {
		minBenefits += cbf.minAnnualAmount(child)
	}

	childCount := float64(len(others) + 1)
	reduction := cbf.BenefitReducer.Reduce(income, childCount)

	reducedBenefits := maxBenefits - reduction
	if reducedBenefits < minBenefits {
		return minBenefits
	}

	return reducedBenefits
}

func (cbf *CCBFormula) Validate() error {

	for _, ageGroupBenefit := range cbf.BenefitClasses {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if err := cbf.BenefitReducer.Validate(); err != nil {
		return errors.Wrap(err, "invalid step reducers")
	}

	return nil
}

// maxAnnualAmount returns the maximum annual benefits for the given child
func (cbf *CCBFormula) maxAnnualAmount(c Child) float64 {

	child := c.Clone()
	maxPayments := make(Payments, 12)

	for month := range maxPayments {
		for _, ageGroup := range cbf.BenefitClasses {

			if ageGroup.IsInAgeGroup(child) {
				maxPayments[month] += ageGroup.Amounts.Upper()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit classes
		}
		child.AgeMonths++
	}

	return maxPayments.Total()
}

// minAnnualAmount returns the minimum annual benefits for the given child
func (cbf *CCBFormula) minAnnualAmount(c Child) float64 {

	child := c.Clone()
	minPayments := make(Payments, 12)

	for month := range minPayments {
		for _, ageGroup := range cbf.BenefitClasses {

			if ageGroup.IsInAgeGroup(child) {
				minPayments[month] += ageGroup.Amounts.Lower()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit classes
		}
		child.AgeMonths++
	}

	return minPayments.Total()
}
