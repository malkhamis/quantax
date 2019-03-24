package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type AgeGroupBenefits struct {
	AgesInMonths    calc.AgeRange
	AmountsPerMonth calc.Bracket
}

func NewAgeGroupBenefits(ages calc.AgeRange, minmaxAmnts calc.Bracket) (AgeGroupBenefits, error) {

	benf := AgeGroupBenefits{
		AgesInMonths:    ages.Clone(),
		AmountsPerMonth: minmaxAmnts.Clone(),
	}

	return benf, benf.Validate()
}

func (g AgeGroupBenefits) IsInAgeGroup(child calc.Person) bool {

	geqMinAge := child.AgeMonths >= g.AgesInMonths.Min()
	leqMaxAge := child.AgeMonths <= g.AgesInMonths.Max()
	return geqMinAge && leqMaxAge
}

func (g AgeGroupBenefits) Validate() error {

	err := g.AmountsPerMonth.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid bracket")
	}

	err = g.AgesInMonths.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid age range")
	}

	return nil
}
