package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

func TestNewChildBenefitCalculator(t *testing.T) {

	formula := testCBFormula{}
	_, err := NewChildBenefitCalculator(CalcConfigCB{formula, nil})
	if errors.Cause(err) != ErrNoIncCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoIncCalc, err)
	}

	incCalc := testIncomeCalculator{}
	_, err = NewChildBenefitCalculator(CalcConfigCB{formula, incCalc})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = NewChildBenefitCalculator(CalcConfigCB{nil, nil})
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	simulatedErr := errors.New("test error")
	formula = testCBFormula{onValidate: simulatedErr}
	_, err = NewChildBenefitCalculator(CalcConfigCB{formula, nil})
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}

func TestCalculator_Calc(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testCBFormula{onApply: incCalc.TotalIncome(nil) / 2.0}
	formula.onClone = formula

	calculator, err := NewChildBenefitCalculator(CalcConfigCB{formula, incCalc})
	if err != nil {
		t.Fatal(err)
	}

	actual := calculator.Calc(nil)
	expected := 3000.0 / 2.0
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestCalculator_SetBeneficiaries(t *testing.T) {

	c := &ChildBenfitCalculator{}
	children := []human.Person{{AgeMonths: 1}, {AgeMonths: 2}}

	c.SetBeneficiaries(children...)
	if len(c.children) != len(children) {
		t.Fatalf("expected %d children, got: %d", len(children), len(c.children))
	}

	for i, actual := range c.children {
		if actual != children[i] {
			t.Errorf(
				"actual does not match expected\nwant: %v\n got: %v",
				children[i], actual,
			)
		}
	}

	c.SetBeneficiaries()
	if len(c.children) != 0 {
		t.Errorf("expected no children, got: %d", len(c.children))
	}
}