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

func TestHouseholdFinances_Version_Change(t *testing.T) {

	spouseA, spouseB := NewIndividualFinances(), NewIndividualFinances()
	hf := NewHouseholdFinances(spouseA, spouseB)

	spouseA.AddAmount(core.IncSrcEarned, 0)
	spouseA.AddAmount(core.DeducSrcRRSP, 0)
	spouseA.AddAmount(core.MiscSrcOthers, 0)
	spouseA.SetAmount(core.IncSrcEarned, 0)
	spouseA.SetAmount(core.DeducSrcRRSP, 0)
	spouseA.SetAmount(core.MiscSrcOthers, 0)

	spouseB.AddAmount(core.IncSrcEarned, 0)
	spouseB.AddAmount(core.DeducSrcRRSP, 0)
	spouseB.AddAmount(core.MiscSrcOthers, 0)
	spouseB.SetAmount(core.IncSrcEarned, 0)
	spouseB.SetAmount(core.DeducSrcRRSP, 0)
	spouseB.SetAmount(core.MiscSrcOthers, 0)

	actual := hf.Version()
	var expected uint64 = 15
	if actual != expected {
		t.Errorf("unexpected version\nwant: %d\n got: %d", expected, actual)
	}
}

func TestHouseholdFinances_Version(t *testing.T) {

	var hf *HouseholdFinances
	if hf.Version() != 0 {
		t.Errorf(
			"expected the version number of nil finances to be zero: got %d",
			hf.Version(),
		)
	}

	hf = NewHouseholdFinances(nil, nil)
	actual := hf.Version()
	var expected uint64 = 1
	if actual != expected {
		t.Errorf("unexpected version\nwant: %d\n got: %d", expected, actual)
	}

	hf = NewHouseholdFinances(NewIndividualFinances(), nil)
	actual = hf.Version()
	expected = 2
	if actual != expected {
		t.Errorf("unexpected version\nwant: %d\n got: %d", expected, actual)
	}

	hf = NewHouseholdFinances(NewIndividualFinances(), NewIndividualFinances())
	actual = hf.Version()
	expected = 3
	if actual != expected {
		t.Errorf("unexpected version\nwant: %d\n got: %d", expected, actual)
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
	// if clone values is unexpectedly set to nil, it will panic and recover
	_, _ = clone.SpouseA(), clone.SpouseB()

}

func TestHouseholdFinances_Clone(t *testing.T) {

	var original *HouseholdFinances
	clone := original.Clone()
	if clone != nil {
		t.Errorf("cloning nil finances should return nil")
	}

	original = NewHouseholdFinances(NewIndividualFinances(), NewIndividualFinances())
	clone = original.Clone()
	if clone.Version() != original.Version() {
		t.Errorf(
			"expected clone's version '%d' to match original's '%d'",
			clone.Version(), original.Version(),
		)
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
