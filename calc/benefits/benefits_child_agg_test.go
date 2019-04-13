package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/calc/human"
	"github.com/pkg/errors"
)

func TestNewChildBenefitCalculatorAgg(t *testing.T) {

	c0, c1, c2 := &ChildBenfitCalculator{}, &ChildBenfitCalculator{}, &ChildBenfitCalculator{}
	_, err := NewChildBenefitAggregator(c0, c1, c2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = NewChildBenefitAggregator(nil, nil, nil)
	if errors.Cause(err) != ErrNoCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCalc, err)
	}

}

func TestCalculatorAgg_Calc(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testCBFormula{onApply: incCalc.TotalIncome(nil) / 2.0}
	formula.onClone = formula

	c0, err := NewChildBenefitCalculator(formula, incCalc)
	if err != nil {
		t.Fatal(err)
	}

	c1, err := NewChildBenefitCalculator(formula, incCalc)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := NewChildBenefitCalculator(formula, incCalc)
	if err != nil {
		t.Fatal(err)
	}

	aggregator, err := NewChildBenefitAggregator(c0, c1, c2)
	if err != nil {
		t.Fatal(err)
	}

	actual := aggregator.Calc(nil)
	expected := (3000.0 / 2.0) * float64(len(aggregator.calculators))
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestAggregator_SetBeneficiaries(t *testing.T) {

	c0, c1, c2 := &ChildBenfitCalculator{}, &ChildBenfitCalculator{}, &ChildBenfitCalculator{}
	children := []human.Person{{AgeMonths: 1}, {AgeMonths: 2}}

	aggregator, err := NewChildBenefitAggregator(c0, c1, c2)
	if err != nil {
		t.Fatal(err)
	}

	aggregator.SetBeneficiaries(children...)

	for i, c := range aggregator.calculators {
		typed, ok := c.(*ChildBenfitCalculator)
		if !ok {
			t.Fatal("failed to typecast")
		}
		if len(typed.children) != len(children) {
			t.Fatalf(
				"calculator %d: expected %d children, got: %d",
				i, len(children), len(typed.children),
			)
		}
	}

	for _, c := range aggregator.calculators {
		typed, ok := c.(*ChildBenfitCalculator)
		if !ok {
			t.Fatal("failed to typecast")
		}
		for i, actual := range typed.children {
			if actual != children[i] {
				t.Errorf(
					"actual does not match expected\nwant: %v\n got: %v",
					children[i], actual,
				)
			}
		}
	}

}