package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type AgeGroupBenefits struct {
	Ages    calc.AgeRange
	Amounts calc.Bracket
}

func NewAgeGroupBenefits(ages calc.AgeRange, minmaxAmnts calc.Bracket) (AgeGroupBenefits, error) {

	benf := AgeGroupBenefits{
		Ages:    ages.Clone(),
		Amounts: minmaxAmnts.Clone(),
	}

	return benf, benf.Validate()
}

func (g AgeGroupBenefits) IsInAgeGroup(c Child) bool {
	return c.AgeMonths >= g.Ages.Min() && c.AgeMonths <= g.Ages.Max()
}

func (g AgeGroupBenefits) Validate() error {

	err := g.Amounts.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid bracket")
	}

	err = g.Ages.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid age range")
	}

	return nil
}
