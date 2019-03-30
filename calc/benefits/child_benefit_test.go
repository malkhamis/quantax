package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestNewChildBenefitCalculator_Full(t *testing.T) {

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
	calculator, err := NewChildBenefitCalculator(formulaBC)
	if err != nil {
		t.Fatal(err)
	}

	calculator.SetBeneficiaries(children...)
	actual := calculator.Calc(finances)
	expected := (55*12 + 55*2) - (2 * 0.0132 * 10000)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.SetBeneficiaries()
	if len(calculator.children) != 0 {
		t.Errorf(
			"expected length of children slice to be zero, got: %d",
			len(calculator.children),
		)
	}

	actual = calculator.Calc(finances)
	expected = 0.0
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.SetBeneficiaries(calc.Person{AgeMonths: 0})
	actual = calculator.Calc(calc.FamilyFinances{{}, {}})
	expected = 55 * 12
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestNewChildBenefitCalculator_Errors(t *testing.T) {

	bracket := calc.WeightedBracketFormula{
		0.0132: calc.Bracket{math.Inf(1), 100000},
	}
	formulaBC := &BCECTBMaxReducer{ReducerFormula: bracket}

	_, err := NewChildBenefitCalculator(formulaBC)
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrBoundsReversed, err)
	}

	_, err = NewChildBenefitCalculator(nil)
	if errors.Cause(err) != calc.ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}

}