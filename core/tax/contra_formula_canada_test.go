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
		OrderedCreditors: []Creditor{
			&testCreditor{
				onTaxCredit:       0,
				onFinancialSource: 2,
				onRule:            core.CreditRule{CrSource: "zero-creditor", Type: 1},
			},
			&testCreditor{
				onTaxCredit:       1000,
				onFinancialSource: 1,
				onRule:            core.CreditRule{CrSource: t.Name(), Type: 123},
			},
		},

		TaxYear:   2019,
		TaxRegion: "UnitTest",
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
			TaxYear:         2019,
			TaxRegion:       "UnitTest",
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
		OrderedCreditors: []Creditor{
			&testCreditor{onRule: core.CreditRule{CrSource: "1000", Type: 1}},
			&testCreditor{onRule: core.CreditRule{CrSource: "2000", Type: 2}},
			&testCreditor{onRule: core.CreditRule{CrSource: "3000", Type: 3}},
			&testCreditor{onRule: core.CreditRule{CrSource: "5000", Type: 5}},
			&testCreditor{onRule: core.CreditRule{CrSource: "4000", Type: 4}},
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

	credits := []core.TaxCredit{cr3000, nil, cr2000, cr5000, cr1000, cr2000, noRuleCr, cr4000}
	cf.FilterAndSort(&credits)
	expected := []core.TaxCredit{cr1000, cr2000, cr2000, cr3000, cr5000, cr4000}
	diff := deep.Equal(credits, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

	cf.FilterAndSort(nil) // should not panic
}

func TestCanadianContraFormula_Validate(t *testing.T) {

	cases := []struct {
		name  string
		cform *CanadianContraFormula
		err   error
	}{
		//
		{
			name: "valid",
			cform: &CanadianContraFormula{
				OrderedCreditors: []Creditor{
					&testCreditor{onRule: core.CreditRule{CrSource: "1000"}},
					&testCreditor{onRule: core.CreditRule{CrSource: "2000"}},
				},
			},
			err: nil,
		},
		//
		{
			name: "invalid-duplicate-creditors",
			cform: &CanadianContraFormula{
				OrderedCreditors: []Creditor{
					&testCreditor{onRule: core.CreditRule{CrSource: "1000"}},
					&testCreditor{onRule: core.CreditRule{CrSource: "2000", Type: 1}},
					&testCreditor{onRule: core.CreditRule{CrSource: "2000", Type: 2}},
				},
			},
			err: ErrDupCreditSource,
		},
		//
		{
			name: "nil-creditor",
			cform: &CanadianContraFormula{
				OrderedCreditors: []Creditor{nil},
			},
			err: ErrNoCreditor,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			err := c.cform.Validate()
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

	originalRule := core.CreditRule{CrSource: t.Name(), Type: 123}
	originalCreditor := &testCreditor{
		onTaxCredit:       1000,
		onFinancialSource: 1,
		onRule:            originalRule,
	}

	original = &CanadianContraFormula{
		OrderedCreditors: []Creditor{originalCreditor},
		TaxYear:          2019,
		TaxRegion:        "UnitTest",
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

	original.OrderedCreditors = append(original.OrderedCreditors, originalCreditor)
	err = original.Validate()
	if err == nil {
		t.Fatal("invalid formula did not return error when validated")
	}
	err = clone.Validate()
	if err != nil {
		t.Fatal("changes to original contra formula should not affect clone")
	}

	expected := &CanadianContraFormula{
		OrderedCreditors: []Creditor{originalCreditor},
		TaxYear:          2019,
		TaxRegion:        "UnitTest",
	}

	diff := deep.Equal(clone, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_TaxYear(t *testing.T) {

	cf := &CanadianContraFormula{TaxYear: 2019}

	if cf.Year() != 2019 {
		t.Errorf("actual year does not match expected\nwant: %d\n got: %d",
			2019, cf.Year(),
		)
	}

}

func TestCanadianContraFormula_TaxRegion(t *testing.T) {

	cf := &CanadianContraFormula{TaxRegion: core.Region(t.Name())}

	if cf.Region() != core.Region(t.Name()) {
		t.Errorf("actual region does not match expected\nwant: %q\n got: %q",
			t.Name(), string(cf.Region()),
		)
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
