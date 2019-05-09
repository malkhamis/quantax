package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestWeightedCreditor_TaxAmount(t *testing.T) {

	creditor := WeightedCreditor{Weight: 0.15}
	creditor.TargetFinancialSource = core.MiscSrcTuition

	finances := &testFinancer{onTotalAmount: 1000}
	taxPayer := &TaxPayer{
		Finances: finances,
	}

	actual, expected := creditor.TaxCredit(taxPayer), 0.15*1000
	if actual != expected {
		t.Errorf("expected result\nwant: %.2f\n got: %.2f", expected, actual)
	}

	if finances.onTotalAmountCapturedArg[0] != core.MiscSrcTuition {
		t.Errorf(
			"expected passed arg to be %v, got: %v",
			core.MiscSrcTuition, finances.onTotalAmountCapturedArg,
		)
	}
}

func TestWeightedCreditor_TaxAmount_NilTaxPayer(t *testing.T) {

	creditor := WeightedCreditor{}
	actual, expected := creditor.TaxCredit(nil), 0.0
	if actual != expected {
		t.Errorf("expected result\nwant: %.2f\n got: %.2f", expected, actual)
	}

}

func TestWeightedCreditor_clone(t *testing.T) {

	original := WeightedCreditor{
		Weight: 0.50,
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
