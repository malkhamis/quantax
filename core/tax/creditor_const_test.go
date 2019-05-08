package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestConstCreditor_TaxCredit(t *testing.T) {

	cc := &ConstCreditor{Amount: 1000}
	if actual := cc.TaxCredit(nil); actual != 1000.0 {
		t.Errorf("unexpected tax credit amount\nwant: %.2f\n got: %.2f", 1000.0, actual)
	}
}

func TestConstCreditor_Rule(t *testing.T) {

	rule := core.CreditRule{
		CrSource: t.Name(),
		Type:     1,
	}
	cc := &ConstCreditor{CreditRule: rule}
	if actual := cc.Rule(); actual != rule {
		t.Errorf("unexpected tax rule\nwant: %v\n got: %v", rule, actual)
	}
}

func TestConstCreditor_FinancialSource(t *testing.T) {

	cc := &ConstCreditor{TargetFinancialSource: 1}
	if actual := cc.FinancialSource(); actual != 1 {
		t.Errorf("unexpected financial source\nwant: %d\n got: %d", 1, actual)
	}
}

func TestConstCreditor_Description(t *testing.T) {

	cc := &ConstCreditor{CreditDescription: t.Name()}
	if actual := cc.Description(); actual != t.Name() {
		t.Errorf("unexpected description\nwant: %q\n got: %q", t.Name(), actual)
	}
}

func TestConstCreditor_clone(t *testing.T) {

	original := ConstCreditor{
		Amount:                1000,
		CreditDescription:     t.Name(),
		TargetFinancialSource: 2,
		CreditRule:            core.CreditRule{CrSource: "test", Type: 3},
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
