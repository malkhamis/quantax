package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestBCECTBMaxReducer_Apply(t *testing.T) {

	bracket := calc.WeightedBracketFormula{
		0.0132: calc.Bracket{100000, math.Inf(1)},
	}

	mr := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: calc.Bracket{0, 55},
			},
		},
	}

	err := mr.Validate()
	if err != nil {
		t.Fatal(err)
	}

	child1, child2 := calc.Person{AgeMonths: 0}, calc.Person{AgeMonths: 6}
	income := 110000.0
	max := 2.0 * 12.0 * 55
	expected := max - (2 * 0.0132 * 10000)

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

func TestBCECTBMaxReducer_Validate_InvalidAgeRanges(t *testing.T) {

	formula := BCECTBMaxReducer{
		BeneficiaryClasses: []AgeGroupBenefits{
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

func TestBCECTBMaxReducer_Validate_NilFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		BeneficiaryClasses: nil,
		ReducerFormula:     nil,
	}

	err := formula.Validate()
	if errors.Cause(err) != calc.ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}

}

func TestBCECTBMaxReducer_Validate_InvalidFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		ReducerFormula: calc.WeightedBracketFormula{
			0.0132: calc.Bracket{100000, 1},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}
}

func TestBCECTBMaxReducer_Clone(t *testing.T) {

	bracket := calc.WeightedBracketFormula{
		0.0132: calc.Bracket{100000, math.Inf(1)},
	}

	originalFormula := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: calc.Bracket{0, 55},
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
	originalFormula.BeneficiaryClasses = nil
	originalFormula.ReducerFormula = nil

	actualResults := clone.Apply(income, child1, child2)
	if actualResults != originalResults {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", originalResults, actualResults)
	}

}

func TestBCECTBMaxReducer_IncomeCalcMethod(t *testing.T) {

	incomeType := (&BCECTBMaxReducer{}).IncomeCalcMethod()
	if incomeType != AFNI {
		t.Errorf("unexpected income type\nwant: %s\n got: %s", AFNI, incomeType)
	}
}
