package finance

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestHouseholdFinances_Full(t *testing.T) {

	spouse1 := NewEmptyIndividialFinances(2018)
	spouse1.AddIncome(IncSrcEarned, 6)
	spouse1.AddIncome(IncSrcUCCB, 4)
	spouse1.AddDeduction(DeducSrcRRSP, 20)
	spouse1.Cash = 24

	spouse2 := NewEmptyIndividialFinances(2018)
	spouse2.AddIncome(IncSrcEarned, 12)
	spouse2.AddIncome(IncSrcInterest, 3)
	spouse2.AddDeduction(DeducSrcRRSP, 25)
	spouse2.AddDeduction(DeducSrcMedical, 25)
	spouse2.Cash = 66

	finances := NewHouseholdFinances(spouse1, spouse2)

	actualIncome := finances.TotalIncome()
	expectedIncome := 10.0 + 15.0
	if actualIncome != expectedIncome {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedIncome, actualIncome,
		)
	}

	actualIncome = finances.TotalIncome(IncSrcInterest, IncSrcUCCB)
	expectedIncome = 4.0 + 3.0
	if actualIncome != expectedIncome {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedIncome, actualIncome,
		)
	}

	actualDeductions := finances.TotalDeductions()
	expectedDeductions := 20.0 + 50.0
	if actualDeductions != expectedDeductions {
		t.Errorf(
			"unexpected deduction total\nwant: %.2f\n got: %.2f",
			expectedDeductions, actualDeductions,
		)
	}

	actualDeductions = finances.TotalDeductions(DeducSrcMedical)
	expectedDeductions = 25.0
	if actualDeductions != expectedDeductions {
		t.Errorf(
			"unexpected deduction total\nwant: %.2f\n got: %.2f",
			expectedDeductions, actualDeductions,
		)
	}

	actualCash := finances.Cash()
	expectedCash := 90.0
	if actualCash != expectedCash {
		t.Errorf(
			"unexpected cash total\nwant: %.2f\n got: %.2f",
			expectedCash, actualCash,
		)
	}
}

func TestHouseholdFinances_Clone(t *testing.T) {

	f1 := NewEmptyIndividialFinances(2018)
	f1.AddIncome(IncSrcEarned, 123)
	f1.AddDeduction(DeducSrcRRSP, 456)
	original := HouseholdFinances{f1, nil}

	clone := original.Clone()
	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Fatal("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	f1.AddIncome(IncSrcEarned, 500)
	diff = deep.Equal(original, clone)
	if diff == nil {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}

	f1.AddDeduction(DeducSrcMedical, 500)
	diff = deep.Equal(original, clone)
	if diff == nil {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}
}

func TestIndividualFinances_NumFieldsUnchanged(t *testing.T) {

	dummy := IndividualFinances{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 6 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type. Next, update this test with the new " +
				"number of fields",
		)
	}
}
