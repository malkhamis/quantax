package calc

import (
	"math"
	"testing"
)

func TestWeightedBracketFormula_Clone(t *testing.T) {

	originalBracket := Bracket{100, 200}
	original := WeightedBracketFormula{0.10: originalBracket}
	err := original.Validate()
	if err != nil {
		t.Fatal(err)
	}

	clone := original.Clone()

	if len(clone) != len(original) {
		t.Fatalf(
			"expected clone len to be %d, got: %d",
			len(original), len(clone),
		)
	}

	actualBracket, ok := clone[0.10]
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

}

func TestWeightedBracketFormula_Validate(t *testing.T) {

	rateMap := WeightedBracketFormula{0.10: Bracket{100, 200}}

	err := rateMap.Validate()
	if err != nil {
		t.Fatalf("expected facts object to validate with no errors, got: %v", err)
	}

	invalid := WeightedBracketFormula{
		math.Inf(-1): Bracket{10, 20},
		math.Inf(1):  Bracket{10, 20},
	}

	err = invalid.Validate()
	if err == nil {
		t.Fatal("expected an error validating an bracket rates with infinity rates")
	}

	invalid = WeightedBracketFormula{0.15: Bracket{20, 10}}
	err = invalid.Validate()
	if err == nil {
		t.Fatal("expected an error validating an bracket rates with invalid brackets")
	}
}
