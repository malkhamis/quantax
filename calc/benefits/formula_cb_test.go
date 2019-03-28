package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestStepReducer_Apply(t *testing.T) {

	childCount1 := calc.WeightedBracketFormula{
		0.000: calc.Bracket{0, 10000},
		0.030: calc.Bracket{10000, 50000},
		0.070: calc.Bracket{50000, math.Inf(1)},
	}
	childCount2 := calc.WeightedBracketFormula{
		0.000: calc.Bracket{0, 10000},
		0.050: calc.Bracket{10000, 50000},
		0.100: calc.Bracket{50000, math.Inf(1)},
	}

	mr := &CCBMaxReducer{
		Reducers: []calc.WeightedBracketFormula{childCount1, childCount2},
		BenefitClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 11},
				AmountsPerMonth: calc.Bracket{0, 500},
			},
			{
				AgesMonths:      calc.AgeRange{12, 23},
				AmountsPerMonth: calc.Bracket{0, 250},
			},
		},
	}

	child1, child2 := calc.Person{AgeMonths: 0}, calc.Person{AgeMonths: 6}
	max := (12.0 * 500) + (6.0*500 + 6.0*250)

	income := 100000.0
	expected := max - (0.050 * 40000) - (0.100 * 50000)
	actual := mr.Apply(income, child1, child2)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 100000.0
	expected = 0.0
	actual = mr.Apply(income)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 500000.0
	expected = 0.0
	actual = mr.Apply(income, child1)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}
	// actualReduction = sr.Reduce(100000, -1.0)
	// expectedReduction = 0.070*(65976-30450) + 0.032*(100000-65976)
	// if actualReduction != expectedReduction {
	// 	t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expectedReduction, actualReduction)
	// }
	//
	// actualReduction = sr.Reduce(100000, 100.0)
	// expectedReduction = 0.230*(65976-30450) + 0.095*(100000-65976)
	// if actualReduction != expectedReduction {
	// 	t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expectedReduction, actualReduction)
	// }
}

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

func TestCCBMaxReducer_MaxAnnualAmount(t *testing.T) {

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

//
// func TestCCBMaxReducer_Clone(t *testing.T) {
//
// 	formula := &CCBMaxReducer{
// 		BenefitClasses: []AgeGroupBenefits{
// 			{
// 				AgesMonths:      calc.AgeRange{0, 10},
// 				AmountsPerMonth: calc.Bracket{50, 100},
// 			},
// 			{
// 				AgesMonths:      calc.AgeRange{11, 20},
// 				AmountsPerMonth: calc.Bracket{25, 50},
// 			},
// 		},
// 	}
//
// 	originalResults := formula.Apply(50000, children)(calc.Person{AgeMonths: 5})
//
// 	clone := formula.Clone()
// 	formula.BenefitClasses = nil
// 	formula.Reducers = nil
//
// 	clone.Apply(income, children)
// 	if actual != expected {
// 		t.Errorf(
// 			"expected a 5 month old child to be entitled to %.2f, got %.2f",
// 			expected, actual,
// 		)
// 	}
