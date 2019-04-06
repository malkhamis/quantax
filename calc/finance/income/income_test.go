package income

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

func TestCalculator_NetIncome_Adjusted(t *testing.T) {

	f := testFormula{
		incomeAdjusters: map[finance.IncomeSource]Adjuster{
			finance.IncomeSource(1000): testAdjuster{adjusted: 250.0},
		},
		deducAdjusters: map[finance.DeductionSource]Adjuster{
			finance.DeductionSource(2000): testAdjuster{adjusted: 100.0},
		},
	}

	c, err := NewCalculator(f)
	if err != nil {
		t.Fatal(err)
	}

	finances := testIncomeDeductor{
		deducSrcs:       finance.NewDeductionSourceSet(finance.DeductionSource(2000)),
		incomeSrcs:      finance.NewIncomeSourceSet(finance.IncomeSource(1000)),
		totalDeductions: 15000.0,
		totalIncome:     30000.0,
	}

	actual, expected := c.NetIncome(finances), (250.0 - 100.0)
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestCalculator_NetIncome_Unadjusted(t *testing.T) {

	f := testFormula{}

	c, err := NewCalculator(f)
	if err != nil {
		t.Fatal(err)
	}

	finances := testIncomeDeductor{
		deducSrcs:       finance.NewDeductionSourceSet(finance.DeductionSource(2000)),
		incomeSrcs:      finance.NewIncomeSourceSet(finance.IncomeSource(1000)),
		totalDeductions: 10000.0,
		totalIncome:     30000.0,
	}

	actual, expected := c.NetIncome(finances), (30000.0 - 10000.0)
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestCalculator_NilFinances(t *testing.T) {

	f := testFormula{}

	c, err := NewCalculator(f)
	if err != nil {
		t.Fatal(err)
	}

	actual, expected := c.NetIncome(nil), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}

	actual, expected = c.TotalIncome(nil), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}

	actual, expected = c.TotalDeductions(nil), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestNewCalculator_Errors(t *testing.T) {

	_, err := NewCalculator(nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	expectedErr := errors.New("injected error")
	f := testFormula{err: expectedErr}
	_, err = NewCalculator(f)
	if errors.Cause(err) != expectedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", expectedErr, err)
	}

	f.err = nil
	c, err := NewCalculator(f)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}
	if c == nil {
		t.Errorf("expected non-nil calculator if no error was returned")
	}
}
