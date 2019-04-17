package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
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

	crGroup := []TaxCredit{
		{
			Amount: 5000,
			CreditSourceControl: CreditSourceControl{
				Source:  1,
				Control: ControlTypeCashable,
			},
		},
		{
			Amount: 4000,
			CreditSourceControl: CreditSourceControl{
				Source:  2,
				Control: ControlTypeNotCarryForward,
			},
		},
		{
			Amount: 2000, CreditSourceControl: CreditSourceControl{
				Source:  3,
				Control: ControlTypeNotCarryForward,
			},
		},
		{
			Amount: 1000,
			CreditSourceControl: CreditSourceControl{
				Source:  4,
				Control: ControlTypeCashable,
			},
		},
		{
			Amount: 500,
			CreditSourceControl: CreditSourceControl{
				Source:  5,
				Control: ControlTypeNotCarryForward,
			},
		},
		{
			Amount: 500,
			CreditSourceControl: CreditSourceControl{
				Source:  6,
				Control: ControlTypeCanCarryForward,
			},
		},
	}

	actualNetTax, actualRemainingCrs := (&Calculator{}).netPayableTax(10000, crGroup)
	expectedNetTax := -1000.0
	expectedRemainingCrs := []finance.TaxCredit{
		{Source: 1, Amount: 0.0},
		{Source: 2, Amount: 0.0},
		{Source: 3, Amount: 0.0},
		{Source: 4, Amount: 0.0},
		{Source: 5, Amount: 0.0},
		{Source: 6, Amount: 500},
	}

	if actualNetTax != expectedNetTax {
		t.Errorf(
			"actual net tax does not match expected\nwant: %.2f\ngot: %.2f",
			expectedNetTax, actualNetTax,
		)
	}

	diff := deep.Equal(actualRemainingCrs, expectedRemainingCrs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
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
