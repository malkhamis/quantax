package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

func TestAgeGroupBenefits_IsInAgeGroup(t *testing.T) {

	agb, err := NewAgeGroupBenefits(human.AgeRange{10, 20}, core.Bracket{0, 100})
	if err != nil {
		t.Fatal(err)
	}

	if agb.IsInAgeGroup(nil) {
		t.Errorf("expected nil person to not be in age group")
	}

	if !agb.IsInAgeGroup(&human.Person{AgeMonths: 15}) {
		t.Errorf("expected a 5 month old person to be in age group")
	}
	if !agb.IsInAgeGroup(&human.Person{AgeMonths: 10}) {
		t.Errorf("expected a 10 month old person to be in age group")
	}
	if !agb.IsInAgeGroup(&human.Person{AgeMonths: 20}) {
		t.Errorf("expected a 20 month old person to be in age group")
	}

	if agb.IsInAgeGroup(&human.Person{AgeMonths: 9}) {
		t.Errorf("expected a 9 month old person to not be in age group")
	}
	if agb.IsInAgeGroup(&human.Person{AgeMonths: 21}) {
		t.Errorf("expected a 21 month old person to not be in age group")
	}
}

func TestMultiAgeGroupBenefits_MinAnnualAmount(t *testing.T) {

	agb := multiAgeGroupBenefits{
		{
			AgesMonths:      human.AgeRange{0, 10},
			AmountsPerMonth: core.Bracket{50, 100},
		},
		{
			AgesMonths:      human.AgeRange{11, 20},
			AmountsPerMonth: core.Bracket{25, 50},
		},
	}

	expected := 0.0
	actual := agb.MinAnnualAmount(nil)
	if actual != expected {
		t.Errorf(
			"expected a a nil child  child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

	expected = float64(50*6) + float64(25*6)
	actual = agb.MinAnnualAmount(&human.Person{AgeMonths: 5})
	if actual != expected {
		t.Errorf(
			"expected a 5 month old child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

}

func TestMultiAgeGroupBenefits_MaxAnnualAmount(t *testing.T) {

	agb := multiAgeGroupBenefits{
		{
			AgesMonths:      human.AgeRange{0, 10},
			AmountsPerMonth: core.Bracket{50, 100},
		},
		{
			AgesMonths:      human.AgeRange{11, 20},
			AmountsPerMonth: core.Bracket{25, 50},
		},
	}

	expected := 0.0
	actual := agb.MaxAnnualAmount(nil)
	if actual != expected {
		t.Errorf(
			"expected a 5 month old child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

	expected = float64(100*6) + float64(50*6)
	actual = agb.MaxAnnualAmount(&human.Person{AgeMonths: 5})
	if actual != expected {
		t.Errorf(
			"expected a 5 month old child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

}

func TestAgeGroupBenefits_Validate(t *testing.T) {

	agb := AgeGroupBenefits{
		AgesMonths:      human.AgeRange{20, 10},
		AmountsPerMonth: core.Bracket{0, 100},
	}

	err := agb.Validate()
	if errors.Cause(err) != human.ErrInvalidAgeRange {
		t.Errorf("unexpected error\nwant: %v\n got: %v", human.ErrInvalidAgeRange, err)
	}

	agb = AgeGroupBenefits{
		AgesMonths:      human.AgeRange{0, 10},
		AmountsPerMonth: core.Bracket{200, 100},
	}

	err = agb.Validate()
	if errors.Cause(err) != core.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", core.ErrBoundsReversed, err)
	}
}
