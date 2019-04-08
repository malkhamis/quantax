package finance

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestHouseholdFinances_Full(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances(2018)
	spouse1.AddIncome(IncSrcEarned, 6)
	spouse1.AddIncome(IncSrcUCCB, 4)
	spouse1.AddDeduction(DeducSrcRRSP, 20)
	spouse1.Cash = 24

	spouse2 := NewEmptyIndividualFinances(2018)
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

func TestHouseholdFinances_AddIncome(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances(2018)
	spouse1.AddIncome(IncSrcEarned, 6)

	actual := spouse1.Income[IncSrcEarned]
	if actual != 6 {
		t.Fatalf("expected AdddIncome to set the map key, got: %.2f", actual)
	}

	spouse1.AddIncome(IncSrcEarned, 6)
	actual = spouse1.Income[IncSrcEarned]
	if actual != 12 {
		t.Fatalf("expected AdddIncome to accumulate amount, got: %.2f", actual)
	}

}

func TestHouseholdFinances_AddDeduction(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances(2018)

	spouse1.AddDeduction(DeducSrcRRSP, 6)
	actual := spouse1.Deductions[DeducSrcRRSP]
	if actual != 6 {
		t.Fatalf("expected AdddDeduction to set the map key, got: %.2f", actual)
	}

	spouse1.AddDeduction(DeducSrcRRSP, 6)
	actual = spouse1.Deductions[DeducSrcRRSP]
	if actual != 12 {
		t.Fatalf("expected AdddDeduction to set the map key, got: %.2f", actual)
	}

}

func TestHouseholdFinances_Sources(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances(2018)
	spouse1.AddIncome(IncSrcEarned, 6)
	spouse1.AddIncome(IncSrcUCCB, 4)
	spouse1.AddDeduction(DeducSrcRRSP, 20)

	actualIncSrcs := spouse1.IncomeSources()
	expectedIncSrcs := NewIncomeSourceSet(IncSrcEarned, IncSrcUCCB)
	diff := deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualDeducSrcs := spouse1.DeductionSources()
	expectedDeducSrcs := NewDeductionSourceSet(DeducSrcRRSP)
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	spouse2 := NewEmptyIndividualFinances(2018)
	spouse2.AddIncome(IncSrcEarned, 12)
	spouse2.AddIncome(IncSrcInterest, 3)
	spouse2.AddDeduction(DeducSrcRRSP, 25)
	spouse2.AddDeduction(DeducSrcMedical, 25)

	finances := NewHouseholdFinances(spouse1, spouse2)

	actualIncSrcs = finances.IncomeSources()
	expectedIncSrcs = NewIncomeSourceSet(IncSrcEarned, IncSrcUCCB, IncSrcInterest)
	diff = deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualDeducSrcs = finances.DeductionSources()
	expectedDeducSrcs = NewDeductionSourceSet(DeducSrcRRSP, DeducSrcMedical)
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestIndividualFinances_Nil_Sources(t *testing.T) {

	var nilFinances *IndividualFinances

	expectedIncSrcs := make(IncomeSourceSet)
	actualIncSrcs := nilFinances.IncomeSources()
	diff := deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedDeducSrcs := make(DeductionSourceSet)
	actualDeducSrcs := nilFinances.DeductionSources()
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestHouseholdFinances_Clone(t *testing.T) {

	f1 := NewEmptyIndividualFinances(2018)
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
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
