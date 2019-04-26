package income

import (
	"testing"

	"github.com/malkhamis/quantax/core"

	"github.com/pkg/errors"
)

func TestCalculator_NetIncome_Adjusted(t *testing.T) {

	r := &Recipe{
		IncomeAdjusters: map[core.FinancialSource]Adjuster{
			core.FinancialSource(1000): testAdjuster{adjusted: 250.0},
		},
		DeductionAdjusters: map[core.FinancialSource]Adjuster{
			core.FinancialSource(2000): testAdjuster{adjusted: 100.0},
		},
	}

	c, err := NewCalculator(r)
	if err != nil {
		t.Fatal(err)
	}

	finances := &testFinancer{
		onIncomeSources:    []core.FinancialSource{core.FinancialSource(1000)},
		onDeductionSources: []core.FinancialSource{core.FinancialSource(2000)},
		onTotalAmount:      []float64{2500, 1000},
	}

	c.SetFinances(finances)
	actual, expected := c.NetIncome(), (250.0 - 100.0)
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestCalculator_NetIncome_Unadjusted(t *testing.T) {

	r := new(Recipe)

	c, err := NewCalculator(r)
	if err != nil {
		t.Fatal(err)
	}

	finances := &testFinancer{
		onDeductionSources: []core.FinancialSource{core.FinancialSource(123)},
		onIncomeSources:    []core.FinancialSource{core.FinancialSource(123)},
		onTotalAmount:      []float64{30000, 10000},
	}

	c.SetFinances(finances)
	actual, expected := c.NetIncome(), (30000.0 - 10000.0)
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestCalculator_NilFinances(t *testing.T) {

	r := new(Recipe)

	c, err := NewCalculator(r)
	if err != nil {
		t.Fatal(err)
	}

	c.SetFinances(nil)
	actual, expected := c.NetIncome(), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}

	actual, expected = c.TotalIncome(), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}

	actual, expected = c.TotalDeductions(), 0.0
	if actual != expected {
		t.Fatalf("unexpected net income\nwant: %.2f\ngot: %.2f", expected, actual)
	}
}

func TestNewCalculator_Errors(t *testing.T) {

	_, err := NewCalculator(nil)
	if errors.Cause(err) != ErrNoRecipe {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoRecipe, err)
	}

	c, err := NewCalculator(new(Recipe))
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}
	if c == nil {
		t.Errorf("expected non-nil calculator if no error was returned")
	}
}
