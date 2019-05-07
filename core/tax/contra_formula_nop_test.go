package tax

import (
	"testing"

	"github.com/malkhamis/quantax/core"
)

func TestNopContraFormula(t *testing.T) {

	nop := NopContraFormula{}
	cr := nop.Apply(nil)
	if cr != nil {
		t.Errorf(
			"expected nop contra-formula type to return nil credits, got: %v", cr,
		)
	}

	err := nop.Validate()
	if err != nil {
		t.Errorf("expected no error validating noop contra formula, got: %v", err)
	}

	if nop.Year() != 0 {
		t.Errorf("expected nop formula to return 0 year, got: %d", nop.Year())
	}

	if nop.Region() != core.Region("") {
		t.Errorf("expected nop formula to return empty-string region, got: %q", nop.Region())
	}

	nop.Clone()
}
