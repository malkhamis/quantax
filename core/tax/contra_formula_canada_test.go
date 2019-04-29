package tax

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

func TestCanadianContraFormula_Apply_Nil_finances(t *testing.T) {

	cf := new(CanadianContraFormula)

	actual := cf.Apply(nil)
	if actual != nil {
		t.Errorf("expected nil []credits if finances is nil")
	}

	actual = cf.Apply(&TaxPayer{})
	if actual != nil {
		t.Errorf("expected nil []credits if finances is nil")
	}

}

func TestCanadianContraFormula_Apply(t *testing.T) {

	cf := &CanadianContraFormula{
		Creditors: []Creditor{
			&testCreditor{
				onTaxCredit:       1000,
				onFinancialSource: 1,
				onCrSourceName:    t.Name(),
			},
			&testCreditor{
				onTaxCredit:       0,
				onFinancialSource: 2,
				onCrSourceName:    "zero-creditor",
			},
		},
		ApplicationOrder: []core.CreditRule{
			core.CreditRule{CrSource: "zero-creditor", Type: 456},
			core.CreditRule{CrSource: t.Name(), Type: 123},
		},
		RelatedTaxInfo: core.TaxInfo{2019, "UnitTest"},
	}

	err := cf.Validate()
	if err != nil {
		t.Fatal(err)
	}

	f := core.NewFinancerNop()
	expected := []*TaxCredit{
		&TaxCredit{
			AmountInitial:   1000,
			AmountRemaining: 1000,
			AmountUsed:      0,
			CrRule:          core.CreditRule{t.Name(), 123},
			Desc:            "test",
			FinancialSource: 1,
			Ref:             f,
			RelatedTaxInfo:  core.TaxInfo{2019, "UnitTest"},
		},
	}

	actual := cf.Apply(&TaxPayer{Finances: f})
	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_FilterAndSort(t *testing.T) {

	cf := &CanadianContraFormula{
		ApplicationOrder: []core.CreditRule{
			{CrSource: "1000", Type: 1},
			{CrSource: "2000", Type: 2},
			{CrSource: "3000", Type: 3},
			{CrSource: "5000", Type: 5},
			{CrSource: "4000", Type: 4},
		},
	}

	err := cf.Validate()
	if err != nil {
		t.Fatal(err)
	}

	cr1000 := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "1000", Type: 1},
	}
	cr2000 := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "2000", Type: 2},
	}
	cr3000 := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "3000", Type: 3},
	}
	cr4000 := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "4000", Type: 4},
	}
	cr5000 := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "5000", Type: 5},
	}
	noRuleCr := &testTaxCredit{
		onRule: core.CreditRule{CrSource: "5000", Type: 0},
	}

	credits := []core.TaxCredit{cr3000, nil, cr2000, cr5000, cr1000, noRuleCr, cr4000}
	actual := cf.FilterAndSort(credits)
	expected := []core.TaxCredit{cr1000, cr2000, cr3000, cr5000, cr4000}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

	actual = cf.FilterAndSort(nil)
	diff = deep.Equal(actual, []core.TaxCredit{})
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_Validate(t *testing.T) {

	cases := []struct {
		name    string
		formula *CanadianContraFormula
		err     error
	}{
		//
		{
			name: "valid",
			formula: &CanadianContraFormula{
				Creditors: []Creditor{
					&testCreditor{onCrSourceName: "1000"},
					&testCreditor{onCrSourceName: "2000"},
				},
				ApplicationOrder: []core.CreditRule{{CrSource: "1000"}, {CrSource: "2000"}},
			},
			err: nil,
		},
		//
		{
			name: "valid-duplicate-creditors",
			formula: &CanadianContraFormula{
				Creditors: []Creditor{
					&testCreditor{onCrSourceName: "1000"},
					&testCreditor{onCrSourceName: "2000"},
					&testCreditor{onCrSourceName: "2000"}, // duplicates here is ok
				},
				ApplicationOrder: []core.CreditRule{{CrSource: "1000"}, {CrSource: "2000"}},
			},
			err: nil,
		},
		//
		{
			name: "duplicates-in-app-order",
			formula: &CanadianContraFormula{
				ApplicationOrder: []core.CreditRule{{CrSource: "1000"}, {CrSource: "1000"}},
			},
			err: ErrDupCreditSource,
		},
		//
		{
			name: "unknown-creditor",
			formula: &CanadianContraFormula{
				Creditors: []Creditor{
					&testCreditor{onCrSourceName: "1000"},
					&testCreditor{onCrSourceName: "2222"},
				},
				ApplicationOrder: []core.CreditRule{{CrSource: "1000"}, {CrSource: "2000"}},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "nil-creditor",
			formula: &CanadianContraFormula{
				Creditors: []Creditor{nil},
			},
			err: ErrNoCreditor,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			err := c.formula.Validate()
			if errors.Cause(err) != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}

func TestCanadianContraFormula_Clone(t *testing.T) {

	var original *CanadianContraFormula
	if original.Clone() != nil {
		t.Error("cloning nil contra formula should return nil")
	}

	originalCreditor := &testCreditor{
		onTaxCredit:       1000,
		onFinancialSource: 1,
		onCrSourceName:    t.Name(),
	}
	originalRule := core.CreditRule{CrSource: t.Name(), Type: 123}
	originalTaxInfo := core.TaxInfo{2019, "UnitTest"}

	original = &CanadianContraFormula{
		Creditors:        []Creditor{originalCreditor},
		ApplicationOrder: []core.CreditRule{originalRule},
		RelatedTaxInfo:   originalTaxInfo,
	}

	err := original.Validate()
	if err != nil {
		t.Fatal(err)
	}

	clone := original.Clone()
	err = clone.Validate()
	if err != nil {
		t.Fatal(err)
	}

	original.ApplicationOrder[0] = core.CreditRule{}
	original.Creditors[0] = &testCreditor{onCrSourceName: "1111"}
	err = original.Validate()
	if err == nil {
		t.Fatal("invalid formula did not return error when validated")
	}
	err = clone.Validate()
	if err != nil {
		t.Fatal("changes to original contra formula should not affect clone")
	}

	expected := &CanadianContraFormula{
		Creditors:        []Creditor{originalCreditor},
		ApplicationOrder: []core.CreditRule{originalRule},
		RelatedTaxInfo:   originalTaxInfo,
	}

	diff := deep.Equal(clone, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_TaxInfo(t *testing.T) {

	cf := &CanadianContraFormula{
		RelatedTaxInfo: core.TaxInfo{TaxYear: 2019, TaxRegion: core.Region(t.Name())},
	}

	expected := core.TaxInfo{TaxYear: 2019, TaxRegion: core.Region(t.Name())}
	actual := cf.TaxInfo()
	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_checkFinSrcCreditorsInSet(t *testing.T) {

	cf := &CanadianContraFormula{
		Creditors: []Creditor{
			&testCreditor{onCrSourceName: t.Name()},
		},
	}

	err := cf.checkCreditorCrSrcNamesInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	err = cf.checkCreditorCrSrcNamesInSet(
		map[string]struct{}{
			t.Name(): struct{}{},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	cf.Creditors[0] = nil
	err = cf.checkCreditorCrSrcNamesInSet(nil)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

}

func TestCanadianContraFormula_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianContraFormula{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 3 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
