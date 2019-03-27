package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestCCBMaxReducer_MinAnnualAmount(t *testing.T) {

	formula := &CCBMaxReducer{
		BenefitClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 10},
				AmountsPerMonth: calc.Bracket{50, 100},
			},
			{
				AgesMonths:      calc.AgeRange{11, 20},
				AmountsPerMonth: calc.Bracket{25, 50},
			},
		},
	}

	expected := float64(50*6) + float64(25*6)
	actual := formula.minAnnualAmount(calc.Person{AgeMonths: 5})
	if actual != expected {
		t.Errorf(
			"expected a 5 month old child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

}

func TestMaxReducer_MaxAnnualAmount(t *testing.T) {

	formula := CCBMaxReducer{
		BenefitClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 10},
				AmountsPerMonth: calc.Bracket{50, 100},
			},
			{
				AgesMonths:      calc.AgeRange{11, 20},
				AmountsPerMonth: calc.Bracket{25, 50},
			},
		},
	}

	expected := float64(100*6) + float64(50*6)
	actual := formula.maxAnnualAmount(calc.Person{AgeMonths: 5})
	if actual != expected {
		t.Errorf(
			"expected a 5 month old child to be entitled to %.2f, got %.2f",
			expected, actual,
		)
	}

}

func TestCCBMaxReducer_Validate_InvalidAgeRanges(t *testing.T) {

	formula := CCBMaxReducer{
		BenefitClasses: []AgeGroupBenefits{
			AgeGroupBenefits{
				AgesMonths:      calc.AgeRange{10, 0},
				AmountsPerMonth: calc.Bracket{0, 55},
			},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}
}

func TestCCBMaxReducer_Validate_NilFormula(t *testing.T) {

	formula := CCBMaxReducer{
		BenefitClasses: nil,
		Reducers:       nil,
	}

	err := formula.Validate()
	if errors.Cause(err) != calc.ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}
}

func TestCCBMaxReducer_Validate_InvalidFormula(t *testing.T) {

	formula := CCBMaxReducer{
		Reducers: []calc.WeightedBracketFormula{
			{0.0132: calc.Bracket{100000, 1}},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}
}
