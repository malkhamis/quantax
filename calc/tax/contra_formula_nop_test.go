package tax

import "testing"

func TestNopContraFormula(t *testing.T) {

	nop := NopContraFormula{}
	cr := nop.Apply(nil, 0)
	if cr != nil {
		t.Errorf(
			"expected nop contra-formula type to return nil credits, got: %v", cr,
		)
	}

	err := nop.Validate()
	if err != nil {
		t.Errorf("expected no error validating noop contra formula, got: %v", err)
	}

	nop.Clone()
}
