package rrsp

import (
	"testing"

	"github.com/malkhamis/quantax/core"

	"github.com/pkg/errors"
)

func TestNewCalculator_Error(t *testing.T) {

	_, err := NewCalculator(CalcConfig{nil, nil})
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}
}

func TestCalcConfig_validate(t *testing.T) {

	formula := &testFormula{}
	err := CalcConfig{formula, nil}.validate()
	if errors.Cause(err) != ErrNoTaxCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoTaxCalc, err)
	}

	taxCalc := &testTaxCalculator{}
	err = CalcConfig{formula, taxCalc}.validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = CalcConfig{nil, nil}.validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	simulatedErr := errors.New("test error")
	formula = &testFormula{onValidate: simulatedErr}
	err = CalcConfig{formula, nil}.validate()
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}

func TestCalculator_TaxPaid(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onTaxPayable: []float64{1000.0, 1015.0},
	}

	formula := &testFormula{}

	c, err := NewCalculator(CalcConfig{formula, taxCalc})
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
		onTaxPayable: []float64{1000.0, 985.0},
	}

	formula := &testFormula{}

	c, err := NewCalculator(CalcConfig{formula, taxCalc})
	if err != nil {
		t.Fatal(err)
	}

	finances := core.NewEmptyIndividualFinances()
	finances.SetRRSPAmounts(core.RRSPAmounts{ContributionRoom: 1000.0})
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

	formula := &testFormula{
		onAllowedIncomeSources: []core.FinancialSource{core.IncSrcEarned},
	}

	c, err := NewCalculator(CalcConfig{formula, &testTaxCalculator{}})
	if err != nil {
		t.Fatal(err)
	}

	finances := core.NewEmptyIndividualFinances()
	finances.SetRRSPAmounts(core.RRSPAmounts{ContributionRoom: 1000})
	finances.AddAmount(core.IncSrcEarned, 1500.0)
	c.SetFinances(finances)
	formula.onContributionEarned = finances.TotalIncome(core.IncSrcEarned)

	actual := c.ContributionEarned()
	expected := formula.onContributionEarned
	if actual != expected {
		t.Errorf(
			"unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual,
		)
	}
}

func TestCalculator_RefundError(t *testing.T) {

	f := core.NewEmptyIndividualFinances()
	f.SetRRSPAmounts(core.RRSPAmounts{ContributionRoom: 1000})
	c := &Calculator{finances: f}

	_, err := c.TaxRefund(2000.0)
	if err != ErrNoRRSPRoom {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoRRSPRoom, err)
	}
}

func TestCalculator_SetFinances_Nil(t *testing.T) {

	c := &Calculator{
		taxCalculator: &testTaxCalculator{},
	}
	c.SetFinances(nil)
	if c.finances == nil {
		t.Fatal("expected a non nil finances")
	}

}
