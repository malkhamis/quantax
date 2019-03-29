package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestNewCBCalculator_Full(t *testing.T) {

	bracket := calc.WeightedBracketFormula{
		0.0132: calc.Bracket{100000, math.Inf(1)},
	}

	formulaBC := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      calc.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: calc.Bracket{0, 55},
			},
		},
	}

	finances := calc.FamilyFinances{
		{Income: 120000.0, Deductions: 10000},
		{Income: 20000, Deductions: 20000},
	}

	children := []calc.Person{{AgeMonths: 0}, {AgeMonths: 6*12 - 2}}
	calculator, err := NewCBCalculator(formulaBC, finances, children...)
	if err != nil {
		t.Fatal(err)
	}

	actual := calculator.Calc()
	expected := (55*12 + 55*2) - (2 * 0.0132 * 10000)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.UpdateBeneficiaries()
	if len(calculator.children) != 0 {
		t.Errorf(
			"expected length of children slice to be zero, got: %d",
			len(calculator.children),
		)
	}

	actual = calculator.Calc()
	expected = 0.0
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.UpdateBeneficiaries(calc.Person{AgeMonths: 0})
	calculator.UpdateFinances(calc.FamilyFinances{{}, {}})
	actual = calculator.Calc()
	expected = 55 * 12
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestNewCBCalculator_Errors(t *testing.T) {

	bracket := calc.WeightedBracketFormula{
		0.0132: calc.Bracket{math.Inf(1), 100000},
	}
	formulaBC := &BCECTBMaxReducer{ReducerFormula: bracket}
	finances := calc.FamilyFinances{{}, {}}

	_, err := NewCBCalculator(formulaBC, finances)
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}

	_, err = NewCBCalculator(nil, finances)
	if errors.Cause(err) != calc.ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}

}
