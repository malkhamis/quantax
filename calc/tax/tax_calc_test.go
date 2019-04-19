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
		ContraTaxFormula: &testTaxContraFormula{},
		IncomeCalc:       incCalc,
	}

	c, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	expected := formula.onApply
	c.SetFinances(finance.NewEmptyIndividualFinances(2018))
	actual, _ := c.TaxPayable()
	t.Fatal("^^ missing check for credits")
	if actual != expected {
		t.Fatalf("unexpected tax\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

func TestCalculator_netPayableTax(t *testing.T) {

	crGroup := []*taxCredit{
		&taxCredit{
			amount: 5000,
			rule: CreditRule{
				Source: "1",
				Type:   CrRuleTypeCashable,
			},
		},
		&taxCredit{
			amount: 4000,
			rule: CreditRule{
				Source: "2",
				Type:   CrRuleTypeNotCarryForward,
			},
		},
		&taxCredit{
			amount: 2000,
			rule: CreditRule{
				Source: "3",
				Type:   CrRuleTypeNotCarryForward,
			},
		},
		&taxCredit{
			amount: 1000,
			rule: CreditRule{
				Source: "4",
				Type:   CrRuleTypeCashable,
			},
		},
		&taxCredit{
			amount: 500,
			rule: CreditRule{
				Source: "5",
				Type:   CrRuleTypeNotCarryForward,
			},
		},
		&taxCredit{
			amount: 500,
			rule: CreditRule{
				Source: "6",
				Type:   CrRuleTypeCanCarryForward,
			},
		},
	}

	actualNetTax, actualRemainingCrs := (&Calculator{}).netPayableTax(10000, crGroup)
	expectedNetTax := -1000.0
	expectedRemainingCrs := []*taxCredit{
		&taxCredit{rule: CreditRule{Source: "1"}, owner: nil, amount: 0.0},
		&taxCredit{rule: CreditRule{Source: "2"}, owner: nil, amount: 0.0},
		&taxCredit{rule: CreditRule{Source: "3"}, owner: nil, amount: 0.0},
		&taxCredit{rule: CreditRule{Source: "4"}, owner: nil, amount: 0.0},
		&taxCredit{rule: CreditRule{Source: "5"}, owner: nil, amount: 0.0},
		&taxCredit{rule: CreditRule{Source: "6"}, owner: nil, amount: 500},
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
		ContraTaxFormula: &testTaxContraFormula{},
		IncomeCalc:       nil,
	}
	_, err := NewCalculator(cfg)
	if errors.Cause(err) != ErrNoIncCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoIncCalc, err)
	}

}
