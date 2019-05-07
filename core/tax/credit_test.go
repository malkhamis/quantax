package tax

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestTaxCredit_Amounts(t *testing.T) {

	var cr *TaxCredit

	initial, used, remaining := cr.Amounts()
	if initial != 0.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			0.0, initial,
		)
	}
	if used != 0.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			0.0, used,
		)
	}
	if remaining != 0.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			0.0, remaining,
		)
	}

	cr = &TaxCredit{
		AmountInitial:   30,
		AmountUsed:      10,
		AmountRemaining: 20,
	}

	initial, used, remaining = cr.Amounts()
	if initial != 30.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			30.0, initial,
		)
	}
	if used != 10.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			10.0, used,
		)
	}
	if remaining != 20.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			20.0, remaining,
		)
	}

}

func TestTaxCredit_SetAmounts(t *testing.T) {

	var cr *TaxCredit
	cr.SetAmounts(0, 0, 0) // should not panic

	cr = new(TaxCredit)

	cr.SetAmounts(30, 10, 20)
	if cr.AmountInitial != 30.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			30.0, cr.AmountInitial,
		)
	}
	if cr.AmountUsed != 10.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			10.0, cr.AmountUsed,
		)
	}
	if cr.AmountRemaining != 20.0 {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\ngot: %.2f",
			20.0, cr.AmountRemaining,
		)
	}

}

func TestTaxCredit_Source(t *testing.T) {

	var cr *TaxCredit
	src := cr.Source() // should not panic
	if src != core.SrcNone {
		t.Errorf(
			"actual does not match expected\nwant: %d\ngot: %d",
			core.SrcNone, src,
		)
	}

	cr = &TaxCredit{FinancialSource: 2}
	src = cr.Source()
	if src != 2 {
		t.Errorf(
			"actual does not match expected\nwant: %d\ngot: %d",
			2, src,
		)
	}
}

func TestTaxCredit_ReferenceFinancer(t *testing.T) {

	var cr *TaxCredit

	ref := cr.ReferenceFinancer() // should not panic
	if ref != nil {
		t.Errorf(
			"actual does not match expected\nwant: %v\ngot: %v",
			nil, ref,
		)
	}

	financer := core.NewFinancerNop()
	cr = &TaxCredit{Ref: financer}
	ref = cr.ReferenceFinancer()
	if ref != financer {
		t.Error("actual reference financer is different")
	}
}

func TestTaxCredit_Rule(t *testing.T) {

	var cr *TaxCredit
	rule := cr.Rule() // should not panic

	if (rule != core.CreditRule{}) {
		t.Errorf(
			"actual does not match expected\nwant: %v\ngot: %v",
			core.CreditRule{}, rule,
		)
	}

	expected := core.CreditRule{CrSource: t.Name(), Type: 1}
	cr = &TaxCredit{CrRule: expected}
	rule = cr.Rule()
	if rule != expected {
		t.Errorf(
			"actual does not match expected\nwant: %v\ngot: %v",
			expected, rule,
		)
	}
}

func TestTaxCredit_Description(t *testing.T) {

	var cr *TaxCredit
	desc := cr.Description() // should not panic
	if desc != "" {
		t.Errorf("expected no description, got: %q", desc)
	}

	cr = &TaxCredit{Desc: t.Name()}
	desc = cr.Description()
	if desc != t.Name() {
		t.Errorf(
			"actual does not match expected\nwant: %q\ngot: %q",
			t.Name(), desc,
		)
	}
}

func TestTaxCredit_Year(t *testing.T) {

	var cr *TaxCredit
	year := cr.Year() // should not panic
	if year != 0 {
		t.Errorf("expected year to be zero, got: %d", year)
	}

	cr = &TaxCredit{TaxYear: 2020}
	year = cr.Year()
	if year != 2020 {
		t.Errorf(
			"actual does not match expected\nwant: %d\ngot: %d",
			2020, year,
		)
	}
}

func TestTaxCredit_Region(t *testing.T) {

	var cr *TaxCredit
	region := cr.Region() // should not panic
	if region != core.Region("") {
		t.Errorf("expected region to be empty, got: %q", region)
	}

	cr = &TaxCredit{TaxRegion: core.Region(t.Name())}
	region = cr.Region()
	if region != core.Region(t.Name()) {
		t.Errorf(
			"actual does not match expected\nwant: %q\ngot: %q",
			t.Name(), region,
		)
	}
}

func TestTaxCredit_shallowCopy(t *testing.T) {

	var original *TaxCredit
	clone := original.shallowCopy() // should not panic
	if clone != nil {
		t.Errorf("expected nil TaxCredit clone to also be nil, got: %v", clone)
	}

	original = &TaxCredit{
		AmountInitial:   30,
		AmountRemaining: 10,
		AmountUsed:      20,
		CrRule:          core.CreditRule{CrSource: "something"},
		Desc:            t.Name(),
		FinancialSource: 1,
		Ref:             core.NewFinancerNop(),
		TaxRegion:       core.Region("somewhere"),
		TaxYear:         2020,
	}

	clone = original.shallowCopy()
	if clone == original {
		t.Fatal("expcted a shallow copy to have a different pointer")
	}

	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestTaxCredit_ShallowCopy(t *testing.T) {

	var original *TaxCredit
	clone := original.ShallowCopy() // should not panic
	if clone != nil {
		t.Errorf("expected nil TaxCredit clone to also be nil, got: %v", clone)
	}

	original = &TaxCredit{}

	clone = original.ShallowCopy()
	_, ok := clone.(*TaxCredit)
	if !ok {
		t.Fatalf(
			"expected the underlying type to be %T, got: %T",
			&TaxCredit{}, clone,
		)
	}
}

func TestTaxCredit_NumFieldsUnchanged(t *testing.T) {

	dummy := TaxCredit{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 9 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
