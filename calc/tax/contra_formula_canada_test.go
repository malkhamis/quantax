package tax

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

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
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				CreditsFromDeduction: map[finance.DeductionSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				ApplicationOrder: []CreditSource{1000, 2000},
			},
			err: nil,
		},
		//
		{
			name: "duplicates-in-app-order",
			formula: &CanadianContraFormula{
				ApplicationOrder: []CreditSource{1000, 1000},
			},
			err: ErrDupCreditSource,
		},
		//
		{
			name: "unknown-income-creditor",
			formula: &CanadianContraFormula{
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2222},
				},
				ApplicationOrder: []CreditSource{1000, 2000},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "unknown-deduc-creditor",
			formula: &CanadianContraFormula{
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				CreditsFromDeduction: map[finance.DeductionSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2222},
				},
				ApplicationOrder: []CreditSource{1000, 2000},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "unknown-misc-creditor",
			formula: &CanadianContraFormula{
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				CreditsFromDeduction: map[finance.DeductionSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2000},
				},
				CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
					123: testCreditor{onSource: 1000},
					456: testCreditor{onSource: 2222},
				},
				ApplicationOrder: []CreditSource{1000, 2000},
			},
			err: ErrUnknownCreditSource,
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

func TestCanadianContraFormula_creditsFromIncSrcs(t *testing.T) {

	cf := CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			123: testCreditor{
				onSource:     1000,
				onTaxCredits: Credits{Amount: 115, Source: 1000, IsRefundable: true}},
			456: testCreditor{
				onSource:     2000,
				onTaxCredits: Credits{Amount: 95, Source: 2000, IsRefundable: false},
			},
			111: testCreditor{
				onSource:     3000,
				onTaxCredits: Credits{Amount: 0, Source: 3000, IsRefundable: false},
			},
		},
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddIncome(123, 15000) // has creditor
	finances.AddIncome(111, 20000) // zero credits
	finances.AddIncome(999, 8000)  // no creditor

	actual := cf.creditsFromIncSrcs(finances, 0)
	expected := []Credits{
		{Amount: 115, Source: 1000, IsRefundable: true},
	}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_creditsFromDeducSrcs(t *testing.T) {

	cf := CanadianContraFormula{
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			123: testCreditor{
				onSource:     1000,
				onTaxCredits: Credits{Amount: 115, Source: 1000, IsRefundable: true}},
			456: testCreditor{
				onSource:     2000,
				onTaxCredits: Credits{Amount: 95, Source: 2000, IsRefundable: false},
			},
			111: testCreditor{
				onSource:     3000,
				onTaxCredits: Credits{Amount: 0, Source: 3000, IsRefundable: false},
			},
		},
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddDeduction(123, 15000) // has creditor
	finances.AddDeduction(111, 20000) // zero credits
	finances.AddDeduction(999, 8000)  // no creditor

	actual := cf.creditsFromDeducSrcs(finances, 0)
	expected := []Credits{
		{Amount: 115, Source: 1000, IsRefundable: true},
	}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}
func TestCanadianContraFormula_checkMiscSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: 123}
	cf := &CanadianContraFormula{
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkMiscSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[CreditSource]struct{}{123: struct{}{}}

	cf.CreditsFromMiscAmounts[9] = nil
	err = cf.checkMiscSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromMiscAmounts, 9)
	err = cf.checkMiscSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_checkDeducSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: 123}
	cf := &CanadianContraFormula{
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkDeducSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[CreditSource]struct{}{123: struct{}{}}

	cf.CreditsFromDeduction[9] = nil
	err = cf.checkDeducSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromDeduction, 9)
	err = cf.checkDeducSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_checkIncSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: 123}
	cf := &CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkIncSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[CreditSource]struct{}{123: struct{}{}}

	cf.CreditsFromIncome[9] = nil
	err = cf.checkIncSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromIncome, 9)
	err = cf.checkIncSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_orderCreditGroupInPlace(t *testing.T) {

	cf := &CanadianContraFormula{
		ApplicationOrder: []CreditSource{1000, 2000, 3000, 5000, 4000},
	}
	err := cf.Validate()
	if err != nil {
		t.Fatal(err)
	}

	crGrp := []Credits{
		Credits{Source: 3000},
		Credits{Source: 2000},
		Credits{Source: 5000},
		Credits{Source: 1000},
		Credits{Source: 4000},
	}

	cf.orderCreditGroupInPlace(crGrp)

	expectedOrder := []Credits{
		Credits{Source: 1000},
		Credits{Source: 2000},
		Credits{Source: 3000},
		Credits{Source: 5000},
		Credits{Source: 4000},
	}

	diff := deep.Equal(crGrp, expectedOrder)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

	var emptyGrp []Credits
	cf.orderCreditGroupInPlace(emptyGrp)
	diff = deep.Equal(emptyGrp, []Credits(nil))
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_Clone(t *testing.T) {

	var original *CanadianContraFormula
	if original.Clone() != nil {
		t.Error("cloning nil contra formula should return nil")
	}

	original = &CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			123: testCreditor{onSource: 1000},
			456: testCreditor{onSource: 2000},
		},
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			123: testCreditor{onSource: 1000},
			456: testCreditor{onSource: 2000},
		},
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			123: testCreditor{onSource: 1000},
			456: testCreditor{onSource: 2000},
		},
		ApplicationOrder: []CreditSource{1000, 2000},
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

	original.ApplicationOrder[0], original.ApplicationOrder[1] = 1, 2
	original.CreditsFromIncome[123] = testCreditor{onSource: 1111}
	original.CreditsFromDeduction[123] = testCreditor{onSource: 1111}
	original.CreditsFromMiscAmounts[123] = testCreditor{onSource: 1111}
	err = original.Validate()
	if err == nil {
		t.Fatal("invalid formula did not return error when validated")
	}
	err = clone.Validate()
	if err != nil {
		t.Fatal("changes to original contra formula should not affect clone")
	}
}

func TestCanadianContraFormula_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianContraFormula{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 4 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
