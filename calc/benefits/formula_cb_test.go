package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestCCBMaxReducer_Apply(t *testing.T) {

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

	err := mr.Validate()
	if err != nil {
		t.Fatal(err)
	}

	child1, child2 := calc.Person{AgeMonths: 0}, calc.Person{AgeMonths: 6}
	max := (12.0 * 500) + (6.0*500 + 6.0*250)

	income := 100000.0
	expected := max - (0.050 * 40000) - (0.100 * 50000)
	actual := mr.Apply(income, child1, child2)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 0.0
	expected = 0.0
	actual = mr.Apply(income) // no children
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 500000.0
	expected = 0.0
	actual = mr.Apply(income, child1)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
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

	formula = CCBMaxReducer{
		BenefitClasses: nil,
		Reducers:       []calc.WeightedBracketFormula{nil},
	}

	err = formula.Validate()
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

func TestBCECTBReducer_Clone(t *testing.T) {

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

	originalFormula := &CCBMaxReducer{
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

	err := originalFormula.Validate()
	if err != nil {
		t.Fatal(err)
	}

	income := 100000.0
	child1, child2 := calc.Person{AgeMonths: 0}, calc.Person{AgeMonths: 6}
	originalResults := originalFormula.Apply(income, child1, child2)

	clone := originalFormula.Clone()
	originalFormula.BenefitClasses = nil
	originalFormula.Reducers = nil

	actualResults := clone.Apply(income, child1, child2)
	if actualResults != originalResults {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", originalResults, actualResults)
	}

}

func TestCCBMaxReducer_IncomeCalcMethod(t *testing.T) {

	incomeType := (&CCBMaxReducer{}).IncomeCalcMethod()
	if incomeType != AFNI {
		t.Errorf("unexpected income type\nwant: %s\n got: %s", AFNI, incomeType)
	}
}
