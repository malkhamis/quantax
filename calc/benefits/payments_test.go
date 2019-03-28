package benefits

import "testing"

func TestPayments_Total(t *testing.T) {

	p := payments{10, 10, 20, 30}
	actual := p.Total()
	expected := float64(10 + 10 + 20 + 30)
	if actual != expected {
		t.Errorf(
			"actual total payments (%.2f) does not match expected (%.2f)",
			actual, expected,
		)
	}
}

func TestPayments_Clone(t *testing.T) {

	p := payments{10, 10, 20, 30}
	clone := p.Clone()
	clone[0] = 123.321
	if p[0] != 10 {
		t.Errorf("expected changes to clone to not affect original payments")
	}

	p = nil
	clone = p.Clone()
	if clone != nil {
		t.Error("expected cloned nil paymeny to also be nil")
	}
}
