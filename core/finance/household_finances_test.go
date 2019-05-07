package finance

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestNewHouseholdFinances(t *testing.T) {

	f := NewHouseholdFinances(nil, nil)
	if f == nil {
		t.Fatal("expected a non-nil instance")
	}
}

func TestHouseholdFinances_SpouseA(t *testing.T) {

	var f *HouseholdFinances
	if f.SpouseA() != nil {
		t.Error("expected nil financial data when calling method on nil instance")
	}

	f = NewHouseholdFinances(nil, nil)
	actual := f.SpouseA()
	if actual != nil {
		t.Fatalf("expected a nil instance")
	}

	newSpouseA := NewIndividualFinances()
	f = NewHouseholdFinances(newSpouseA, nil)
	actual = f.SpouseA()
	if actual != newSpouseA {
		t.Errorf("actual spouseA reference does not match expected")
	}

}

func TestHouseholdFinances_MutableSpouseA(t *testing.T) {

	newSpouseA := NewIndividualFinances()
	f := NewHouseholdFinances(newSpouseA, nil)
	actual := f.MutableSpouseA()
	if actual != newSpouseA {
		t.Errorf("actual spouseA reference does not match expected")
	}

	actual.SetAmount(core.IncSrcEarned, 1000)
	spouseA := f.SpouseA()
	actualAmount := spouseA.TotalAmount(core.IncSrcEarned)
	if actualAmount != 1000.0 {
		t.Errorf(
			"mutating spouse finances should be reflected in household finances"+
				"\nwant: %.2f\n got: %.2f", 1000.0, actualAmount,
		)
	}

}

func TestHouseholdFinances_SpouseB(t *testing.T) {

	var f *HouseholdFinances
	if f.SpouseB() != nil {
		t.Error("expected nil financial data when calling method on nil instance")
	}

	f = NewHouseholdFinances(nil, nil)
	actual := f.SpouseB()
	if actual != nil {
		t.Fatalf("expected a nil instance")
	}

	newSpouseB := NewIndividualFinances()
	f = NewHouseholdFinances(nil, newSpouseB)
	actual = f.SpouseB()
	if actual != newSpouseB {
		t.Errorf("actual spouseB reference does not match expected")
	}

}

func TestHouseholdFinances_MutableSpouseB(t *testing.T) {

	newSpouseB := NewIndividualFinances()
	f := NewHouseholdFinances(nil, newSpouseB)
	actual := f.MutableSpouseB()
	if actual != newSpouseB {
		t.Errorf("actual spouseB reference does not match expected")
	}

	actual.SetAmount(core.IncSrcEarned, 1000)
	spouseB := f.SpouseB()
	actualAmount := spouseB.TotalAmount(core.IncSrcEarned)
	if actualAmount != 1000.0 {
		t.Errorf(
			"mutating spouse finances should be reflected in household finances"+
				"\nwant: %.2f\n got: %.2f", 1000.0, actualAmount,
		)
	}

}

func TestHouseholdFinances_MutableSpouseA_nils(t *testing.T) {

	var hf *HouseholdFinances
	spouseA := hf.MutableSpouseA()
	if spouseA != nil {
		t.Error("expected mutable finances to be nil")
	}

	hf = NewHouseholdFinances(nil, NewIndividualFinances())
	spouseA = hf.MutableSpouseA()
	if spouseA != nil {
		t.Error("expected mutable finances to be nil")
	}
}

func TestHouseholdFinances_MutableSpouseB_nils(t *testing.T) {

	var hf *HouseholdFinances
	spouseB := hf.MutableSpouseB()
	if spouseB != nil {
		t.Error("expected mutable finances to be nil")
	}

	hf = NewHouseholdFinances(NewIndividualFinances(), nil)
	spouseB = hf.MutableSpouseB()
	if spouseB != nil {
		t.Error("expected mutable finances to be nil")
	}
}

func TestHouseholdFinances_clone(t *testing.T) {

	var original *HouseholdFinances
	clone := original.clone()
	if clone != nil {
		t.Errorf("cloning nil finances should return nil")
	}

	spouseA, spouseB := NewIndividualFinances(), NewIndividualFinances()
	original = NewHouseholdFinances(spouseA, spouseB)
	clone = original.clone()
	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Fatal("clone does not match original\n", strings.Join(diff, "\n"))
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatal("expected change to original to not affect clone")
		}
	}()
	original.spouseA = nil
	original.spouseB = nil
	// if clone values is unexpectedly set to nil, it will panic and fail the test
	_, _ = clone.SpouseA(), clone.SpouseB()

}

func TestHouseholdFinances_Clone(t *testing.T) {

	var original *HouseholdFinances
	clone := original.Clone()
	if clone != nil {
		t.Errorf("cloning nil finances should return nil")
	}

	original = NewHouseholdFinances(nil, nil)
	clone = original.Clone()
	if original == clone {
		t.Fatal("expected clone to return a new instance")
	}
}

func TestHouseholdFinances_NumFieldsUnchanged(t *testing.T) {

	dummy := HouseholdFinances{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 2 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
