package tax

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

func TestCalculator_Calc(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testTaxFormula{onApply: incCalc.TotalIncome(nil) / 2.0}
	formula.onClone = formula

	c, err := NewCalculator(formula, incCalc)
	if err != nil {
		t.Fatal(err)
	}

	expected := formula.onApply
	actual := c.Calc(finance.NewEmptyIndividualFinances(2018))
	if actual != expected {
		t.Fatalf("unexpected tax\nwant: %.2f\n got: %.2f", expected, actual)
	}

	expected = 0.0
	actual = c.Calc(nil)
	if actual != expected {
		t.Fatalf("unexpected tax\nwant: %.2f\n got: %.2f", expected, actual)
	}

}
func TestNewCalculator_InvalidFormula(t *testing.T) {

	simulatedErr := errors.New("test error")
	invalidFormula := testTaxFormula{onValidate: simulatedErr}

	_, err := NewCalculator(invalidFormula, nil)
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}

func TestCalculator_NilFormula(t *testing.T) {

	_, err := NewCalculator(nil, nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}
}

func TestNewCalculator_NilIncomeCalculator(t *testing.T) {

	_, err := NewCalculator(testTaxFormula{}, nil)
	if errors.Cause(err) != ErrNoIncCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoIncCalc, err)
	}

}
