package factory

import (
	"testing"

	"github.com/malkhamis/quantax/core"
)

func TestNewFinanceFactory(t *testing.T) {
	f := NewFinanceFactory()
	if f == nil {
		t.Error("expected a non nil finance factory")
	}
}

func TestFinanceFactory_NewFinances(t *testing.T) {

	amounts := map[core.FinancialSource]float64{
		core.IncSrcEarned: 10000.0,
	}
	finances := (&FinanceFactory{}).NewFinances(amounts)

	actual := finances.TotalAmount(core.IncSrcEarned)
	expected := 10000.0
	if actual != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expected, actual,
		)
	}
}

func TestFinanceFactory_NewHouseholFinancesForCouple(t *testing.T) {

	amountsA := map[core.FinancialSource]float64{
		core.IncSrcEarned: 10000.0,
	}
	amountsB := map[core.FinancialSource]float64{
		core.IncSrcEarned: 20000.0,
	}

	finances := (&FinanceFactory{}).NewHouseholFinancesForCouple(amountsA, amountsB)

	actualA := finances.SpouseA().TotalAmount(core.IncSrcEarned)
	expectedA := 10000.0
	if actualA != expectedA {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expectedA, actualA,
		)
	}

	actualB := finances.SpouseB().TotalAmount(core.IncSrcEarned)
	expectedB := 20000.0
	if actualB != expectedB {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expectedB, actualB,
		)
	}
}

func TestFinanceFactory_NewHouseholFinancesForSingle(t *testing.T) {

	amounts := map[core.FinancialSource]float64{
		core.IncSrcEarned: 10000.0,
	}

	finances := (&FinanceFactory{}).NewHouseholdFinancesForSingle(amounts)

	actual := finances.SpouseA().TotalAmount(core.IncSrcEarned)
	expected := 10000.0
	if actual != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expected, actual,
		)
	}

	if finances.SpouseB() != nil {
		t.Errorf("expected spouse B finances to be nil")
	}
}
