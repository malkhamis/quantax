package core

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestHouseholdFinances_Full(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances()
	spouse1.AddAmount(IncSrcEarned, 6)
	spouse1.AddAmount(IncSrcUCCB, 4)
	spouse1.AddAmount(DeducSrcRRSP, 20)
	spouse1.AddAmount(MiscSrcMedical, 11)
	spouse1.SetCash(24)

	spouse2 := NewEmptyIndividualFinances()
	spouse2.AddAmount(IncSrcEarned, 12)
	spouse2.AddAmount(IncSrcInterest, 3)
	spouse2.AddAmount(DeducSrcRRSP, 25)
	spouse2.AddAmount(DeducSrcMedical, 25)
	spouse2.AddAmount(SrcUnknown, 12)
	spouse2.SetCash(66)

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

	actualMiscAmount := finances.MiscAmount()
	expectedMiscAmount := 11.0 + 12.0
	if actualMiscAmount != expectedMiscAmount {
		t.Errorf(
			"unexpected deduction total\nwant: %.2f\n got: %.2f",
			expectedMiscAmount, actualMiscAmount,
		)
	}

	actualMiscAmount = finances.MiscAmount(MiscSrcMedical)
	expectedMiscAmount = 11.0
	if actualMiscAmount != expectedMiscAmount {
		t.Errorf(
			"unexpected deduction total\nwant: %.2f\n got: %.2f",
			expectedMiscAmount, actualMiscAmount,
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

func TestHouseholdFinances_AddAmount_Income(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances()
	spouse1.AddAmount(IncSrcEarned, 6)

	actual := spouse1.income[IncSrcEarned]
	if actual != 6 {
		t.Fatalf("expected AddIncome to set the map key, got: %.2f", actual)
	}

	spouse1.AddAmount(IncSrcEarned, 6)
	actual = spouse1.income[IncSrcEarned]
	if actual != 12 {
		t.Fatalf("expected AddIncome to accumulate amount, got: %.2f", actual)
	}

}

func TestHouseholdFinances_AddAmount_Deduction(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances()

	spouse1.AddAmount(DeducSrcRRSP, 6)
	actual := spouse1.deductions[DeducSrcRRSP]
	if actual != 6 {
		t.Fatalf("expected AdddDeduction to set the map key, got: %.2f", actual)
	}

	spouse1.AddAmount(DeducSrcRRSP, 6)
	actual = spouse1.deductions[DeducSrcRRSP]
	if actual != 12 {
		t.Fatalf("expected AdddDeduction to set the map key, got: %.2f", actual)
	}

}

func TestHouseholdFinances_AddAmount_Misc(t *testing.T) {

	finances := NewEmptyIndividualFinances()
	finances.AddAmount(MiscSrcMedical, 6)

	actual := finances.miscAmounts[MiscSrcMedical]
	if actual != 6 {
		t.Fatalf("expected AddMiscAmount to set the map key, got: %.2f", actual)
	}

	finances.AddAmount(MiscSrcMedical, 6)
	actual = finances.miscAmounts[MiscSrcMedical]
	if actual != 12 {
		t.Fatalf("expected AddMiscAmount to accumulate amount, got: %.2f", actual)
	}

}

func TestHouseholdFinances_Sources(t *testing.T) {

	spouse1 := NewEmptyIndividualFinances()
	spouse1.AddAmount(IncSrcEarned, 6)
	spouse1.AddAmount(IncSrcUCCB, 4)
	spouse1.AddAmount(DeducSrcRRSP, 20)
	spouse1.AddAmount(MiscSrcMedical, 11)

	actualIncSrcs := spouse1.IncomeSources()
	expectedIncSrcs := map[FinancialSource]struct{}{
		IncSrcEarned: struct{}{},
		IncSrcUCCB:   struct{}{},
	}
	diff := deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualDeducSrcs := spouse1.DeductionSources()
	expectedDeducSrcs := map[FinancialSource]struct{}{DeducSrcRRSP: struct{}{}}
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualMiscSrcs := spouse1.MiscSources()
	expectedMiscSrcs := map[FinancialSource]struct{}{MiscSrcMedical: struct{}{}}
	diff = deep.Equal(actualMiscSrcs, expectedMiscSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	spouse2 := NewEmptyIndividualFinances()
	spouse2.AddAmount(IncSrcEarned, 12)
	spouse2.AddAmount(IncSrcInterest, 3)
	spouse2.AddAmount(DeducSrcRRSP, 25)
	spouse2.AddAmount(DeducSrcMedical, 25)
	spouse2.AddAmount(SrcUnknown, 12)

	finances := NewHouseholdFinances(spouse1, spouse2)

	actualIncSrcs = finances.IncomeSources()
	expectedIncSrcs = map[FinancialSource]struct{}{
		IncSrcEarned:   struct{}{},
		IncSrcUCCB:     struct{}{},
		IncSrcInterest: struct{}{},
	}
	diff = deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualDeducSrcs = finances.DeductionSources()
	expectedDeducSrcs = map[FinancialSource]struct{}{
		DeducSrcRRSP:    struct{}{},
		DeducSrcMedical: struct{}{},
	}
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actualMiscSrcs = finances.MiscSources()
	expectedMiscSrcs = map[FinancialSource]struct{}{
		MiscSrcMedical: struct{}{},
		SrcUnknown:     struct{}{},
	}
	diff = deep.Equal(actualMiscSrcs, expectedMiscSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestIndividualFinances_RemoveAmounts(t *testing.T) {

	f1 := NewEmptyIndividualFinances()

	f1.AddAmount(IncSrcTFSA, 1000)
	_, ok := f1.income[IncSrcTFSA]
	if !ok {
		t.Errorf("expected key to exist")
	}
	f1.RemoveAmounts(IncSrcTFSA)
	_, ok = f1.income[IncSrcTFSA]
	if ok {
		t.Errorf("expected key to not exist following a call to remove")
	}

	f1.AddAmount(DeducSrcRRSP, 1000)
	_, ok = f1.deductions[DeducSrcRRSP]
	if !ok {
		t.Errorf("expected key to exist")
	}
	f1.RemoveAmounts(DeducSrcRRSP)
	_, ok = f1.deductions[DeducSrcRRSP]
	if ok {
		t.Errorf("expected key to not exist following a call to remove")
	}

	f1.AddAmount(MiscSrcMedical, 1000)
	_, ok = f1.miscAmounts[MiscSrcMedical]
	if !ok {
		t.Errorf("expected key to exist")
	}
	f1.RemoveAmounts(MiscSrcMedical)
	_, ok = f1.miscAmounts[MiscSrcMedical]
	if ok {
		t.Errorf("expected key to not exist following a call to remove")
	}

}

func TestIndividualFinances_Nil_Sources(t *testing.T) {

	var nilFinances *IndividualFinances

	expectedIncSrcs := make(map[FinancialSource]struct{})
	actualIncSrcs := nilFinances.IncomeSources()
	diff := deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedDeducSrcs := make(map[FinancialSource]struct{})
	actualDeducSrcs := nilFinances.DeductionSources()
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedMiscSrcs := make(map[FinancialSource]struct{})
	actualMiscSrcs := nilFinances.MiscSources()
	diff = deep.Equal(actualMiscSrcs, expectedMiscSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestHouseholdFinances_Nil_Sources(t *testing.T) {

	nilFinances := HouseholdFinances{nil, nil}

	expectedIncSrcs := make(map[FinancialSource]struct{})
	actualIncSrcs := nilFinances.IncomeSources()
	diff := deep.Equal(actualIncSrcs, expectedIncSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedDeducSrcs := make(map[FinancialSource]struct{})
	actualDeducSrcs := nilFinances.DeductionSources()
	diff = deep.Equal(actualDeducSrcs, expectedDeducSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedMiscSrcs := make(map[FinancialSource]struct{})
	actualMiscSrcs := nilFinances.MiscSources()
	diff = deep.Equal(actualMiscSrcs, expectedMiscSrcs)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestHouseholdFinances_Clone(t *testing.T) {

	f1 := NewEmptyIndividualFinances()
	f1.AddAmount(IncSrcEarned, 123)
	f1.AddAmount(DeducSrcRRSP, 456)
	f1.AddAmount(MiscSrcMedical, 789)
	original := HouseholdFinances{f1, nil}

	clone := original.Clone()
	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Fatal("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	f1.AddAmount(IncSrcEarned, 500)
	f1.AddAmount(DeducSrcRRSP, 500)
	f1.AddAmount(MiscSrcMedical, 500)
	f1.AddAmount(SrcUnknown, 12)

	if original.TotalIncome() == clone.TotalIncome() {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}

	if original.TotalDeductions() == clone.TotalDeductions() {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}

	if original.MiscAmount() == clone.MiscAmount() {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}

	originalMiscSrcs := original.MiscAmount(MiscSrcMedical, SrcUnknown)
	cloneMiscSrcs := clone.MiscAmount(MiscSrcMedical, SrcUnknown)
	if originalMiscSrcs == cloneMiscSrcs {
		t.Fatal("expected changes made to original to not be reflected in clone")
	}
}

func TestHouseholdFinances_Version(t *testing.T) {

	f1 := NewEmptyIndividualFinances()
	f1.AddAmount(IncSrcEarned, 100)

	f2 := NewEmptyIndividualFinances()
	f2.AddAmount(DeducSrcRRSP, 100)

	hf := NewHouseholdFinances(f1, f2)
	// TODO: This will fail once version is implemented
	if hf.Version() != 0 {
		t.Fatalf("unexpected version\nwant: %d\n got: %d", 0, hf.Version())
	}
}

func TestIndividualFinances_NumFieldsUnchanged(t *testing.T) {

	dummy := IndividualFinances{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 7 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
