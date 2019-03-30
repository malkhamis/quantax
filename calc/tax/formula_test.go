package tax

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestFormula_Clone(t *testing.T) {

	original := CanadianFormula{0.1: finance.Bracket{0, 10}}
	clone := original.Clone()
	originalResults := original.Apply(5)
	original[0.1] = finance.Bracket{100, 1000}
	cloneResults := clone.Apply(5)
	if originalResults != cloneResults {
		t.Errorf("expected clone results to be equal to results of original formula prior to mutation")
	}
}
