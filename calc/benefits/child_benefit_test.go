package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"

	"github.com/pkg/errors"
)

func TestNewChildBenefitCalculator_Full(t *testing.T) {

	bracket := finance.WeightedBrackets{
		0.0132: finance.Bracket{100000, math.Inf(1)},
	}

	formulaBC := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: finance.Bracket{0, 55},
			},
		},
	}

	finances := finance.HouseholdFinances{
		{
			Income: finance.IncomeBySource{
				finance.IncSrcEarned: 120000,
			},
			Deductions: finance.DeductionBySource{
				finance.DeducSrcRRSP: 10000,
			},
		},
		{
			Income: finance.IncomeBySource{
				finance.IncSrcEarned: 20000,
			},
			Deductions: finance.DeductionBySource{
				finance.DeducSrcRRSP: 20000,
			},
		},
	}

	children := []human.Person{{AgeMonths: 0}, {AgeMonths: 6*12 - 2}}
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

	calculator.SetBeneficiaries(human.Person{AgeMonths: 0})
	actual = calculator.Calc(finance.HouseholdFinances{})
	expected = 55 * 12
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestNewChildBenefitCalculator_Errors(t *testing.T) {

	bracket := finance.WeightedBrackets{
		0.0132: finance.Bracket{math.Inf(1), 100000},
	}
	formulaBC := &BCECTBMaxReducer{ReducerFormula: bracket}

	_, err := NewChildBenefitCalculator(formulaBC)
	if errors.Cause(err) != finance.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", finance.ErrBoundsReversed, err)
	}

	_, err = NewChildBenefitCalculator(nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

}
