package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestCanadianSpouseCreditor_TaxAmount(t *testing.T) {

	creditor := CanadianSpouseCreditor{Weight: 0.15, BaseAmount: 500}
	taxPayer := &TaxPayer{
		SpouseFinances:  &testFinancer{},
		NetIncome:       10000,
		SpouseNetIncome: 150,
	}

	actual, expected := creditor.TaxCredit(taxPayer), (500-150)*0.15
	if actual != expected {
		t.Errorf("unexpected result\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestCanadianSpouseCreditor_TaxAmount_HigherSpouseNetIncome(t *testing.T) {

	creditor := CanadianSpouseCreditor{Weight: 0.15, BaseAmount: 500}
	taxPayer := &TaxPayer{
		SpouseFinances:  &testFinancer{},
		NetIncome:       150,
		SpouseNetIncome: 10000,
	}

	actual, expected := creditor.TaxCredit(taxPayer), 0.0
	if actual != expected {
		t.Errorf("unexpected result\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestCanadianSpouseCreditor_TaxAmount_SpouseNetIncomeExceedingBase(t *testing.T) {

	creditor := CanadianSpouseCreditor{Weight: 0.15, BaseAmount: 500}
	taxPayer := &TaxPayer{
		SpouseFinances:  &testFinancer{},
		NetIncome:       10000,
		SpouseNetIncome: 600,
	}

	actual, expected := creditor.TaxCredit(taxPayer), 0.0
	if actual != expected {
		t.Errorf("unexpected result\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestCanadianSpouseCreditor_TaxAmount_Nils(t *testing.T) {

	creditor := CanadianSpouseCreditor{Weight: 0.15, BaseAmount: 500}

	actual, expected := creditor.TaxCredit(nil), 0.0
	if actual != expected {
		t.Errorf("expected amount to be %.2f for nil taxpayer, got: %.2f", expected, actual)
	}

	actual, expected = creditor.TaxCredit(&TaxPayer{}), 0.0
	if actual != expected {
		t.Errorf("expected amount to be %.2f for nil nil spouse, got: %.2f", expected, actual)
	}

}

func TestCanadianSpouseCreditor_clone(t *testing.T) {

	original := CanadianSpouseCreditor{
		BaseAmount: 1000,
		Weight:     0.50,
		CreditDescriptor: CreditDescriptor{
			CreditDescription:     t.Name(),
			TargetFinancialSource: 2,
			CreditRule:            core.CreditRule{CrSource: "test", Type: 3},
		},
	}

	cloneInternal := original.clone()
	diff := deep.Equal(original, cloneInternal)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	cloneExported := original.Clone()
	diff = deep.Equal(original, cloneExported)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}
