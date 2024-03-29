package benefits

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

// AgeGroupBenefits represents the min/max benefit amounts for an age group
type AgeGroupBenefits struct {
	AgesMonths      human.AgeRange
	AmountsPerMonth core.Bracket
}

type multiAgeGroupBenefits []AgeGroupBenefits

func (ma multiAgeGroupBenefits) MaxAnnualAmount(child *human.Person) float64 {

	if child == nil {
		return 0.0
	}

	var (
		childClone = *child
		max        float64
	)

	for range make([]struct{}, 12) {
		for _, ageGroup := range ma {

			if ageGroup.IsInAgeGroup(&childClone) {
				max += ageGroup.AmountsPerMonth.Upper()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit classes
		}
		childClone.AgeMonths++
	}

	return max

}

func (ma multiAgeGroupBenefits) MinAnnualAmount(child *human.Person) float64 {

	if child == nil {
		return 0.0
	}

	var (
		childClone = *child
		min        float64
	)

	for range make([]struct{}, 12) {
		for _, ageGroup := range ma {

			if ageGroup.IsInAgeGroup(&childClone) {
				min += ageGroup.AmountsPerMonth.Lower()
			}
			// we still want to loop in case the child
			// belongs to multiple benefit classes
		}
		childClone.AgeMonths++
	}

	return min

}

// NewAgeGroupBenefits returns a new age group benefit instance. The age range
// is expected to be in months (not years). If the given arguments are invalid,
// an error is returned
func NewAgeGroupBenefits(ages human.AgeRange, minmaxAmnts core.Bracket) (AgeGroupBenefits, error) {

	benf := AgeGroupBenefits{
		AgesMonths:      ages,
		AmountsPerMonth: minmaxAmnts,
	}

	return benf, benf.Validate()
}

// IsInAgeGroup returns true if the age of the given person is with the range
// of this group's age range
func (g AgeGroupBenefits) IsInAgeGroup(person *human.Person) bool {

	if person == nil {
		return false
	}

	geqMinAge := person.AgeMonths >= g.AgesMonths.Min()
	leqMaxAge := person.AgeMonths <= g.AgesMonths.Max()
	return geqMinAge && leqMaxAge
}

// Validate ensures that this instance is valid for use. Users need to call this
// method before use only if the instance was manually created/modified
func (g AgeGroupBenefits) Validate() error {

	err := g.AmountsPerMonth.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid bracket")
	}

	err = g.AgesMonths.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid age range")
	}

	return nil
}
