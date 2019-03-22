package calc

import (
	"math"
	"testing"
)

func TestRateAdjBracketFormula_Clone(t *testing.T) {

	originalBracket := Bracket{100, 200}
	originalRateMap := WeightedBrackets{0.10: originalBracket}
	original, err := NewRateAdjBracketFormula(originalRateMap, 2)
	if err != nil {
		t.Fatal(err)
	}

	clone := original.Clone()

	if len(clone.RateMap) != len(original.RateMap) {
		t.Fatalf(
			"expected clone len to be %d, got: %d",
			len(original.RateMap), len(clone.RateMap),
		)
	}

	actualBracket, ok := clone.RateMap[0.10]
	if !ok {
		t.Fatal("expected bracket to exist in cloned object")
	}
	if actualBracket[0] != originalBracket[0] {
		t.Errorf(
			"unexpected bracket value\nwant: %.2f\n got: %.2f",
			originalBracket[0], actualBracket[0],
		)
	}
	if actualBracket[1] != originalBracket[1] {
		t.Errorf(
			"unexpected bracket value\nwant: %.2f\n got: %.2f",
			originalBracket[1], actualBracket[1],
		)
	}

	if clone.Param != original.Param {
		t.Fatalf(
			"actual formula param (%.2f) is not equal to expected (%.2f)",
			original.Param, clone.Param,
		)
	}
}

func TestRateAdjBracketFormula_Validate(t *testing.T) {

	rateMap := WeightedBrackets{0.10: Bracket{100, 200}}

	err := rateMap.Validate()
	if err != nil {
		t.Fatalf("expected facts object to validate with no errors, got: %v", err)
	}

	invalid := WeightedBrackets{
		math.Inf(-1): Bracket{10, 20},
		math.Inf(1):  Bracket{10, 20},
	}

	err = invalid.Validate()
	if err == nil {
		t.Fatal("expected an error validating an bracket rates with infinity rates")
	}

	invalid = WeightedBrackets{0.15: Bracket{20, 10}}
	err = invalid.Validate()
	if err == nil {
		t.Fatal("expected an error validating an bracket rates with invalid brackets")
	}
}
