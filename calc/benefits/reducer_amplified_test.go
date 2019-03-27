package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestAmplifiedReducer(t *testing.T) {

	formula := AmplifiedReducer{
		0.0132: calc.Bracket{100000, math.Inf(1)},
	}

	actualReduction := formula.Reduce(105000.0, 3.5)
	expectedReduction := (105000.0 - 100000.0) * 0.0132 * 3.5

	if actualReduction != expectedReduction {
		t.Errorf(
			"actual does not match expected\n want: %.2f\n got: %.2f",
			expectedReduction, actualReduction,
		)
	}
}

func TestAmplifiedReducer_Validate(t *testing.T) {

	err := AmplifiedReducer{0.0132: calc.Bracket{100000, 0}}.Validate()
	if errors.Cause(err) != calc.ErrBoundsReversed {
		t.Errorf(
			"unexpected error\n want: %v\n got: %v", calc.ErrBoundsReversed, err)
	}

}

func TestAmplifiedReducer_Clone(t *testing.T) {

	originalFormula := AmplifiedReducer{
		0.1000: calc.Bracket{0, 100000},
		0.2000: calc.Bracket{100000, math.Inf(1)},
	}

	originalReduction := originalFormula.Reduce(200000, 2.0)

	clone := originalFormula.Clone()
	cloneReduction := clone.Reduce(200000, 2.0)
	if cloneReduction != originalReduction {
		t.Error(
			"expected clone formula to return the same results as original"+
				"\nwant: %.2f\n got: %.2f",
			originalReduction, cloneReduction,
		)
	}

	originalFormula[0.1000] = calc.Bracket{300000, 400000}
	originalFormula[0.2000] = calc.Bracket{400000, 500000}

	cloneReduction = clone.Reduce(200000, 2.0)
	if cloneReduction != originalReduction {
		t.Error("expected modification to original to not affect the clone formula")
	}

}
