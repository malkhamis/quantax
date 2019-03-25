package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

// AgeGroupBenefits represents the min/max benefit amounts for an age group
type AgeGroupBenefits struct {
	AgesMonths      calc.AgeRange
	AmountsPerMonth calc.Bracket
}

// NewAgeGroupBenefits returns a new age group benefit instance. If the given
// arguments are invalid, an error is returned
func NewAgeGroupBenefits(ages calc.AgeRange, minmaxAmnts calc.Bracket) (AgeGroupBenefits, error) {

	benf := AgeGroupBenefits{
		AgesMonths:      ages.Clone(),
		AmountsPerMonth: minmaxAmnts.Clone(),
	}

	return benf, benf.Validate()
}

// IsInAgeGroup returns true if the age of the given person is with the range
// of this group's age range
func (g AgeGroupBenefits) IsInAgeGroup(person calc.Person) bool {

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
