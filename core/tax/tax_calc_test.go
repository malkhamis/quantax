package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"

	"github.com/pkg/errors"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestCalculator_SetFinances_Nil(t *testing.T) {

	c := &Calculator{incomeCalculator: testIncomeCalculator{}}
	c.SetFinances(nil)

	if c.finances == nil {
		t.Fatal("expected empty finances to be set")
	}
}

func TestCalculator_SetCredits(t *testing.T) {

	mainCalc := new(Calculator)
	anotherCalc := new(Calculator)

	credits := []core.TaxCredit{
		testTaxCredit{},
		&taxCredit{owner: anotherCalc},
		&taxCredit{amount: 123, owner: mainCalc},
	}

	mainCalc.SetCredits(credits)
	actual := mainCalc.credits
	expected := []*taxCredit{&taxCredit{amount: 123, owner: mainCalc}}
	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_ownCredit(t *testing.T) {

	c := new(Calculator)
	credits := []*taxCredit{&taxCredit{amount: 123}, &taxCredit{amount: 456}}
	c.ownCredits(credits)

	expected := []*taxCredit{
		&taxCredit{owner: c, amount: 123},
		&taxCredit{owner: c, amount: 456},
	}

	diff := deep.Equal(credits, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_TaxPayable(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testTaxFormula{onApply: incCalc.TotalIncome() / 2.0}
	cformula := &testTaxContraFormula{
		onApply: []*taxCredit{
			&taxCredit{
				amount: 50,
				rule:   CreditRule{Source: "tuition", Type: CrRuleTypeCashable},
			},
		},
	}

	cfg := CalcConfig{
		TaxFormula:       formula,
		ContraTaxFormula: cformula,
		IncomeCalc:       incCalc,
	}

	c, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	c.SetFinances(core.NewEmptyIndividualFinances())
	actualTax, actualCr := c.TaxPayable()

	expectedTax := formula.onApply - 50.0
	expectedCr := []core.TaxCredit{
		&taxCredit{
			owner:  c,
			amount: 0.0,
			rule:   CreditRule{Source: "tuition", Type: CrRuleTypeCashable},
		},
	}

	if actualTax != expectedTax {
		t.Fatalf("unexpected tax\nwant: %.2f\n got: %.2f", expectedTax, actualTax)
	}

	diff := deep.Equal(actualCr, expectedCr)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_netPayableTax(t *testing.T) {

	crGroup := []*taxCredit{
		&taxCredit{amount: 5000, rule: CreditRule{Source: "1", Type: CrRuleTypeCashable}},
		&taxCredit{amount: 4000, rule: CreditRule{Source: "2", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 2000, rule: CreditRule{Source: "3", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 1000, rule: CreditRule{Source: "4", Type: CrRuleTypeCashable}},
		&taxCredit{amount: 500, rule: CreditRule{Source: "5", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 500, rule: CreditRule{Source: "6", Type: CrRuleTypeCanCarryForward}},
	}

	actualNetTax, actualRemainingCrs := (&Calculator{}).netPayableTax(10000, crGroup)
	expectedNetTax := -1000.0
	expectedRemainingCrs := []*taxCredit{
		&taxCredit{amount: 0, rule: CreditRule{Source: "1", Type: CrRuleTypeCashable}},
		&taxCredit{amount: 0, rule: CreditRule{Source: "2", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 0, rule: CreditRule{Source: "3", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 0, rule: CreditRule{Source: "4", Type: CrRuleTypeCashable}},
		&taxCredit{amount: 0, rule: CreditRule{Source: "5", Type: CrRuleTypeNotCarryForward}},
		&taxCredit{amount: 500, rule: CreditRule{Source: "6", Type: CrRuleTypeCanCarryForward}},
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
