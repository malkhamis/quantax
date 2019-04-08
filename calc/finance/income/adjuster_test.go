package income

import (
	"testing"
)

func TestWeightedAdjuster_Adjusted(t *testing.T) {

	wa := WeightedAdjuster(0.5)
	expected := 50.0
	actual := wa.Adjusted(100.0)
	if actual != expected {
		t.Fatalf("unexpected result\nwant: %.2f\n got: %.2f", actual, expected)
	}

}

func TestCanadianCapitalGainAdjuster_Clone(t *testing.T) {

	wa := WeightedAdjuster(0.50)
	clone := wa.Clone()
	if clone != wa {
		t.Fatalf("unexpected result\nwant: %v\n got: %v", wa, clone)
	}

	wa = -1.0
	if clone == wa {
		t.Fatal("expected changes to original to not affect clone")
	}
}
