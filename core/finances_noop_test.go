package core

import "testing"

func TestNewFinancerNop(t *testing.T) {

	f1 := NewFinancerNop()
	f2 := NewFinancerNop()
	if f1 == f2 {
		t.Fatal("multiple calls should return different instances")
	}
}

func TestNewHouseholdFinancesNop(t *testing.T) {

	f1 := NewHouseholdFinancesNop()
	f2 := NewHouseholdFinancesNop()
	if f1 == f2 {
		t.Fatal("multiple calls should return different instances")
	}
}

func Test_financerNop(t *testing.T) {

	f := NewFinancerNop()

	// Nop calls
	f.AddAmount(1, 10)
	f.RemoveAmounts(0)
	f.SetAmount(1, 10)

	if actual := f.TotalAmount(1); actual != 0 {
		t.Errorf("expected zero total amount, got: %.2f", actual)
	}

	if actual := f.IncomeSources(); actual != nil {
		t.Errorf("expected no income sources, got: %v", actual)
	}

	if actual := f.DeductionSources(); actual != nil {
		t.Errorf("expected no deduction sources, got: %v", actual)
	}

	if actual := f.MiscSources(); actual != nil {
		t.Errorf("expected no misc sources, got: %v", actual)
	}

	if actual := f.AllSources(); actual != nil {
		t.Errorf("expected no sources, got: %v", actual)
	}

	clone := f.Clone()
	if clone == f {
		t.Error("expected a new instance when calling clone")
	}
}

func Test_householdFinancesNop(t *testing.T) {

	f := NewHouseholdFinancesNop()

	if ref := f.SpouseA(); ref == nil {
		t.Error("expected method call to return a valid reference")
	}

	if ref := f.SpouseB(); ref == nil {
		t.Error("expected method call to return a valid reference")
	}

	if f.SpouseA() == f.SpouseB() {
		t.Error("expected each spouse to be a unique instance")
	}

	if ref := f.MutableSpouseA(); ref != f.SpouseA() {
		t.Error("expected mutable instance to equal the read-only one")
	}

	if ref := f.MutableSpouseB(); ref != f.SpouseB() {
		t.Error("expected mutable instance to equal the read-only one")
	}

	clone := f.Clone()
	if clone == f {
		t.Fatal("expected clone to be a new instance")
	}

	if clone.SpouseA() == f.SpouseA() {
		t.Error("expected clone's spouseA to be different from original")
	}

	if clone.SpouseB() == f.SpouseB() {
		t.Error("expected clone's spouseA to be different from original")
	}
}
