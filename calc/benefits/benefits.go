package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

const maxChildAgeMonths = 12 * 18

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

func (b AgeGroupBenefits) Validate() error {

	err := b.Amounts.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid bracket")
	}

	err = b.Ages.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid age range")
	}

	return nil
}
