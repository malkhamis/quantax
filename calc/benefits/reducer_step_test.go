package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
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

//
// func TestStepReducer_Clone(t *testing.T) {
//
// 	aboveMaxStep := calc.WeightedBracketFormula{
// 		0.000: calc.Bracket{0, 30450},
// 		0.230: calc.Bracket{30450, 65976},
// 		0.095: calc.Bracket{65976, math.Inf(1)},
// 	}
//
// 	step0 := calc.WeightedBracketFormula{}
// 	step1 := calc.WeightedBracketFormula{
// 		0.000: calc.Bracket{0, 30450},
// 		0.070: calc.Bracket{30450, 65976},
// 		0.032: calc.Bracket{65976, math.Inf(1)},
// 	}
//
// 	original, err := NewStepReducer(aboveMaxStep, step0, step1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	originalReduction := original.Reduce(100000, 1.0)
// 	clone := original.Clone()
// 	cloneReduction := clone.Reduce(100000, 1.0)
//
// 	if cloneReduction != originalReduction {
// 		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", originalReduction, cloneReduction)
// 	}
//
// 	original.StepFormulas = nil
// 	original.AboveMaxStepFormula = calc.WeightedBracketFormula{}
//
// 	cloneReduction = clone.Reduce(100000, 1.0)
// 	if cloneReduction != originalReduction {
// 		t.Errorf("expected changes to original formula to not affect clone formula")
// 	}
//
// }
//
// func TestStepReducer_Validate_NilFormulas(t *testing.T) {
//
// 	_, err := NewStepReducer(nil, nil)
// 	if errors.Cause(err) != calc.ErrNoFormula {
// 		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
// 	}
//
// 	_, err = NewStepReducer(
// 		calc.WeightedBracketFormula{0.000: calc.Bracket{0, 30450}}, nil,
// 	)
// 	if errors.Cause(err) != calc.ErrNoFormula {
// 		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
// 	}
// }
//
// func TestStepReducer_Validate_InvalidFormulas(t *testing.T) {
//
// 	_, err := NewStepReducer(
// 		calc.WeightedBracketFormula{0.000: calc.Bracket{100, 0}}, nil,
// 	)
// 	if errors.Cause(err) != calc.ErrBoundsReversed {
// 		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
// 	}
//
// 	_, err = NewStepReducer(
// 		calc.WeightedBracketFormula{0.000: calc.Bracket{0, 1000}},
// 		calc.WeightedBracketFormula{0.000: calc.Bracket{100, 0}},
// 	)
// 	if errors.Cause(err) != calc.ErrBoundsReversed {
// 		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
// 	}
//
// }
