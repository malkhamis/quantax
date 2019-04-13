package tax

import (
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

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
