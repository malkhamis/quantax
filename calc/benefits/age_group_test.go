package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestAgeGroupBenefits_IsInAgeGroup(t *testing.T) {

	agb, err := NewAgeGroupBenefits(calc.AgeRange{10, 20}, calc.Bracket{0, 100})
	if err != nil {
		t.Fatal(err)
	}

	if !agb.IsInAgeGroup(calc.Person{AgeMonths: 15}) {
		t.Errorf("expected a 5 month old person to be in age group")
	}
	if !agb.IsInAgeGroup(calc.Person{AgeMonths: 10}) {
		t.Errorf("expected a 10 month old person to be in age group")
	}
	if !agb.IsInAgeGroup(calc.Person{AgeMonths: 20}) {
		t.Errorf("expected a 20 month old person to be in age group")
	}

	if agb.IsInAgeGroup(calc.Person{AgeMonths: 9}) {
		t.Errorf("expected a 9 month old person to not be in age group")
	}
	if agb.IsInAgeGroup(calc.Person{AgeMonths: 21}) {
		t.Errorf("expected a 21 month old person to not be in age group")
	}
}

func TestAgeGroupBenefits_Validate(t *testing.T) {

	agb := AgeGroupBenefits{
		AgesMonths:      calc.AgeRange{20, 10},
		AmountsPerMonth: calc.Bracket{0, 100},
	}

	err := agb.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}

	agb = AgeGroupBenefits{
		AgesMonths:      calc.AgeRange{0, 10},
		AmountsPerMonth: calc.Bracket{200, 100},
	}

	err = agb.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}
}
