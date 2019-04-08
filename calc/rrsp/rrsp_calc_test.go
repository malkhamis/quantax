package rrsp

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

func TestNewCalculator(t *testing.T) {

	formula := &testFormula{}
	_, err := NewCalculator(formula, nil)
	if errors.Cause(err) != ErrNoTaxCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoTaxCalc, err)
	}

	taxCalc := &testTaxCalculator{}
	_, err = NewCalculator(formula, taxCalc)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = NewCalculator(nil, nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	simulatedErr := errors.New("test error")
	formula = &testFormula{onValidate: simulatedErr}
	_, err = NewCalculator(formula, nil)
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}

func TestCalculator_TaxPaid(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onCalc: []float64{1000.0, 1015.0},
	}

	formula := &testFormula{}

	c, err := NewCalculator(formula, taxCalc)
	if err != nil {
		t.Fatal(err)
	}

	actual, expected := c.TaxPaid(150.0), 15.0
	if actual != expected {
		t.Errorf(
			"unexpected tax paid\nwant: %.2f\n got: %.2f", expected, actual,
		)
	}
}

func TestCalculator_TaxRefund(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onCalc: []float64{1000.0, 985.0},
	}

	formula := &testFormula{}

	c, err := NewCalculator(formula, taxCalc)
	if err != nil {
		t.Fatal(err)
	}

	finances := finance.NewEmptyIndividualFinances(0)
	finances.RRSPContributionRoom = 1000.0
	c.SetFinances(finances)

	expected := 15.0
	actual, err := c.TaxRefund(150.0)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Errorf(
			"unexpected tax refund\nwant: %.2f\n got: %.2f", expected, actual,
		)
	}
}

func TestCalculator_ContributionEarned(t *testing.T) {

	dummyIncSrc := finance.IncomeSource(1111)
	formula := &testFormula{
		onAllowedIncomeSources: []finance.IncomeSource{dummyIncSrc},
	}

	c, err := NewCalculator(formula, &testTaxCalculator{})
	if err != nil {
		t.Fatal(err)
	}

	finances := finance.NewEmptyIndividualFinances(0)
	finances.RRSPContributionRoom = 1000.0
	finances.AddIncome(dummyIncSrc, 1500.0)
	c.SetFinances(finances)
	formula.onContributionEarned = finances.TotalIncome(dummyIncSrc)

	actual := c.ContributionEarned()
	expected := formula.onContributionEarned
	if actual != expected {
		t.Errorf(
			"unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual,
		)
	}
}

func TestCalculator_RefundError(t *testing.T) {

	c := &Calculator{
		finances: &finance.IndividualFinances{
			RRSPContributionRoom: 1000.0,
		},
	}

	_, err := c.TaxRefund(2000.0)
	if err != ErrNoRRSPRoom {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoRRSPRoom, err)
	}
}

func TestCalculator_SetFinances_Nil(t *testing.T) {

	c := &Calculator{}
	c.SetFinances(nil)
	if c.finances == nil {
		t.Fatal("expected a non nil finances")
	}

	if c.finances.EOY != 0 {
		t.Fatalf("expected EOY to be initialized to zero, got: %d", c.finances.EOY)
	}
}
