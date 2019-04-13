package tax

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

func TestCalculator_Calc(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testTaxFormula{onApply: incCalc.TotalIncome(nil) / 2.0}

	cfg := CalcConfig{
		TaxFormula:       formula,
		ContraTaxFormula: testTaxContraFormula{},
		IncomeCalc:       incCalc,
	}

	c, err := NewCalculator(cfg)
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

func TestCalculator_netPayableTax(t *testing.T) {

	crGroup := []Credits{
		{Amount: 5000, IsRefundable: true},
		{Amount: 4000, IsRefundable: false},
		{Amount: 2000, IsRefundable: false},
		{Amount: 1000, IsRefundable: true},
		{Amount: 500, IsRefundable: false},
	}

	actualNetTax, actualLostCr := (&Calculator{}).netPayableTax(10000, crGroup)
	expectedNetTax, expectedLostCr := -1000.0, 1500.0

	if actualNetTax != expectedNetTax {
		t.Errorf(
			"actual net tax does not match expected\nwant: %.2f\ngot: %.2f",
			expectedNetTax, actualNetTax,
		)
	}

	if actualLostCr != expectedLostCr {
		t.Fatalf(
			"actual lost credits does not match expected\nwant: %.2f\ngot: %.2f",
			expectedLostCr, actualLostCr,
		)
	}
}

func TestNewCalculator_Error(t *testing.T) {

	cfg := CalcConfig{
		TaxFormula:       testTaxFormula{},
		ContraTaxFormula: testTaxContraFormula{},
		IncomeCalc:       nil,
	}
	_, err := NewCalculator(cfg)
	if errors.Cause(err) != ErrNoIncCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoIncCalc, err)
	}

}
