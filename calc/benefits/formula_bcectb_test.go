package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
	"github.com/pkg/errors"
)

func TestBCECTBMaxReducer_Apply(t *testing.T) {

	bracket := finance.WeightedBrackets{
		0.0132: finance.Bracket{100000, math.Inf(1)},
	}

	mr := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: finance.Bracket{0, 55},
			},
		},
		IncomeType: finance.AFNI,
	}

	err := mr.Validate()
	if err != nil {
		t.Fatal(err)
	}

	child1, child2 := human.Person{AgeMonths: 0}, human.Person{AgeMonths: 6}
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
				AgesMonths:      human.AgeRange{10, 0},
				AmountsPerMonth: finance.Bracket{0, 55},
			},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != human.ErrInvalidAgeRange {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", human.ErrInvalidAgeRange, err)
	}
}

func TestBCECTBMaxReducer_Validate_NilFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		BeneficiaryClasses: nil,
		ReducerFormula:     nil,
		IncomeType:         finance.AFNI,
	}

	err := formula.Validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

}

func TestBCECTBMaxReducer_Validate_InvalidFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		ReducerFormula: finance.WeightedBrackets{
			0.0132: finance.Bracket{100000, 1},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != finance.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", finance.ErrBoundsReversed, err)
	}
}

func TestBCECTBMaxReducer_Clone(t *testing.T) {

	bracket := finance.WeightedBrackets{
		0.0132: finance.Bracket{100000, math.Inf(1)},
	}

	originalFormula := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: finance.Bracket{0, 55},
			},
		},
		IncomeType: finance.AFNI,
	}

	err := originalFormula.Validate()
	if err != nil {
		t.Fatal(err)
	}

	income := 100000.0
	child1, child2 := human.Person{AgeMonths: 0}, human.Person{AgeMonths: 6}
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

	incomeType := (&BCECTBMaxReducer{IncomeType: finance.AFNI}).IncomeCalcMethod()
	if incomeType != finance.AFNI {
		t.Errorf("unexpected income type\nwant: %s\n got: %s", finance.AFNI, incomeType)
	}
}
