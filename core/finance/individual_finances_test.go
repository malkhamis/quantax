package finance

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
)

func TestNewIndividualFinances(t *testing.T) {

	f := NewIndividualFinances()
	if f.income == nil {
		t.Error("expected map to not be nil")
	}
	if f.deductions == nil {
		t.Error("expected map to not be nil")
	}
	if f.miscAmounts == nil {
		t.Error("expected map to not be nil")
	}
}

func TestIndividualFinances_TotalAmount(t *testing.T) {

	var nilFin *IndividualFinances
	actual := nilFin.TotalAmount(core.IncomeSourcesEnd)
	expected := 0.0
	if actual != expected {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", expected, actual)
	}

	f := NewIndividualFinances()
	f.income[core.IncSrcEarned] = 10
	f.deductions[core.DeducSrcRRSP] = 20
	f.miscAmounts[core.MiscSrcOthers] = 30

	actual = f.TotalAmount(
		core.IncSrcEarned,
		core.DeducSrcRRSP,
		core.MiscSrcOthers,
		1111111,
	)
	expected = 10.0 + 20 + 30
	if actual != expected {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

func TestIndividualFinances_AddAmount(t *testing.T) {

	f := NewIndividualFinances()

	f.AddAmount(core.IncSrcEarned, 10)
	f.AddAmount(core.IncSrcEarned, 10)

	f.AddAmount(core.DeducSrcRRSP, 20)
	f.AddAmount(core.DeducSrcRRSP, 20)

	f.AddAmount(core.MiscSrcOthers, 30)
	f.AddAmount(core.MiscSrcOthers, 30)

	actual := f.income[core.IncSrcEarned]
	if actual != 20.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 20.0, actual)
	}

	actual = f.deductions[core.DeducSrcRRSP]
	if actual != 40.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 40.0, actual)
	}

	actual = f.miscAmounts[core.MiscSrcOthers]
	if actual != 60.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 60.0, actual)
	}
}

func TestIndividualFinances_SetAmount(t *testing.T) {

	f := NewIndividualFinances()

	f.SetAmount(core.IncSrcEarned, 10)
	f.SetAmount(core.IncSrcEarned, 20)

	f.SetAmount(core.DeducSrcRRSP, 20)
	f.SetAmount(core.DeducSrcRRSP, 40)

	f.SetAmount(core.MiscSrcOthers, 30)
	f.SetAmount(core.MiscSrcOthers, 60)

	actual := f.income[core.IncSrcEarned]
	if actual != 20.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 20.0, actual)
	}

	actual = f.deductions[core.DeducSrcRRSP]
	if actual != 40.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 40.0, actual)
	}

	actual = f.miscAmounts[core.MiscSrcOthers]
	if actual != 60.0 {
		t.Errorf("unexpected total amount\nwant: %.2f\n got: %.2f", 60.0, actual)
	}
}

func TestIndividualFinances_RemoveAmount(t *testing.T) {

	f := NewIndividualFinances()
	f.income[core.IncSrcEarned] = 10
	f.deductions[core.DeducSrcRRSP] = 20
	f.miscAmounts[core.MiscSrcOthers] = 30

	f.RemoveAmounts(core.IncSrcEarned, core.DeducSrcRRSP, core.MiscSrcOthers, 1111111)

	if _, ok := f.income[core.IncSrcEarned]; ok {
		t.Errorf("expected removed source to not exist in map")
	}

	if _, ok := f.deductions[core.DeducSrcRRSP]; ok {
		t.Errorf("expected removed source to not exist in map")
	}

	if _, ok := f.miscAmounts[core.MiscSrcOthers]; ok {
		t.Errorf("expected removed source to not exist in map")
	}

}

func TestIndividualFinances_Sources(t *testing.T) {

	var f *IndividualFinances
	if f.IncomeSources() != nil {
		t.Error("expected nil sources if finances is nile")
	}
	if f.DeductionSources() != nil {
		t.Error("expected nil sources if finances is nile")
	}
	if f.MiscSources() != nil {
		t.Error("expected nil sources if finances is nile")
	}

	f = NewIndividualFinances()

	f.income[core.IncSrcEarned] = 10
	f.income[core.IncSrcInterest] = 15

	f.deductions[core.DeducSrcRRSP] = 20
	f.deductions[core.DeducSrcOthers] = 25

	f.miscAmounts[core.MiscSrcOthers] = 30
	f.miscAmounts[core.MiscSrcTuition] = 35

	actual := f.IncomeSources()
	expected := []core.FinancialSource{core.IncSrcEarned, core.IncSrcInterest}
	sort.SliceStable(actual, func(i int, j int) bool {
		return int(actual[i]) < int(actual[j])
	})
	sort.SliceStable(expected, func(i int, j int) bool {
		return int(expected[i]) < int(expected[j])
	})

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actual = f.DeductionSources()
	expected = []core.FinancialSource{core.DeducSrcRRSP, core.DeducSrcOthers}
	sort.SliceStable(actual, func(i int, j int) bool {
		return int(actual[i]) < int(actual[j])
	})
	sort.SliceStable(expected, func(i int, j int) bool {
		return int(expected[i]) < int(expected[j])
	})

	diff = deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actual = f.MiscSources()
	expected = []core.FinancialSource{core.MiscSrcOthers, core.MiscSrcTuition}
	sort.SliceStable(actual, func(i int, j int) bool {
		return int(actual[i]) < int(actual[j])
	})
	sort.SliceStable(expected, func(i int, j int) bool {
		return int(expected[i]) < int(expected[j])
	})

	diff = deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actual = f.AllSources()
	expected = []core.FinancialSource{
		core.IncSrcEarned, core.IncSrcInterest,
		core.DeducSrcRRSP, core.DeducSrcOthers,
		core.MiscSrcOthers, core.MiscSrcTuition,
	}
	sort.SliceStable(actual, func(i int, j int) bool {
		return int(actual[i]) < int(actual[j])
	})
	sort.SliceStable(expected, func(i int, j int) bool {
		return int(expected[i]) < int(expected[j])
	})

	diff = deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestIndividualFinances_clone(t *testing.T) {

	var original *IndividualFinances
	clone := original.clone()
	if clone != nil {
		t.Errorf("cloning nil finances should return nil")
	}

	original = NewIndividualFinances()
	original.SetAmount(core.IncSrcEarned, 10)
	original.SetAmount(core.DeducSrcRRSP, 20)
	original.SetAmount(core.MiscSrcOthers, 30)

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
	original.income = nil
	original.deductions = nil
	original.miscAmounts = nil
	// if clone values are unexpectedly set to nil, it will panic and recover
	_ = clone.income[0]
	_ = clone.deductions[0]
	_ = clone.miscAmounts[0]

}

func TestIndividualFinances_Clone(t *testing.T) {

	var original *IndividualFinances
	clone := original.Clone()
	if clone != nil {
		t.Errorf("cloning nil finances should return nil")
	}

	original = NewIndividualFinances()
	clone = original.Clone()
	if original == clone {
		t.Fatal("expected clone to return a new instance")
	}
}

func TestIndividualFinances_NumFieldsUnchanged(t *testing.T) {

	dummy := IndividualFinances{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 3 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
