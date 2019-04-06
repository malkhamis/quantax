package income

import (
	"reflect"
	"testing"
)

func TestAdjuster_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianCapitalGainAdjuster{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}

func TestCanadianCapitalGainAdjuster_Adjusted(t *testing.T) {

	finances := testIncomeDeductor{totalIncome: 500}
	ccgAdj := CanadianCapitalGainAdjuster{Proportion: 0.50}
	expected := 250.0
	actual := ccgAdj.Adjusted(finances)
	if actual != expected {
		t.Fatalf("unexpected result\nwant: %.2f\n got: %.2f", actual, expected)
	}

}

func TestCanadianCapitalGainAdjuster_Adjusted_NilFinances(t *testing.T) {

	ccgAdj := CanadianCapitalGainAdjuster{Proportion: 0.50}
	expected := 0.0
	actual := ccgAdj.Adjusted(nil)
	if actual != expected {
		t.Fatalf("unexpected result\nwant: %.2f\n got: %.2f", actual, expected)
	}

}

func TestCanadianCapitalGainAdjuster_Clone(t *testing.T) {

	ccgAdj := CanadianCapitalGainAdjuster{Proportion: 0.50}
	clone := ccgAdj.Clone()
	if clone != ccgAdj {
		t.Fatalf("unexpected result\nwant: %v\n got: %v", ccgAdj, clone)
	}

	ccgAdj.Proportion = -1.0
	if clone == ccgAdj {
		t.Fatal("expected changes to original to not affect clone")
	}
}
