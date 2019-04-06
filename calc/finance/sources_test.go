package finance

import "testing"

func TestSets_Has(t *testing.T) {

	ds := NewDeductionSourceSet(DeducSrcRRSP)
	if !ds.Has(DeducSrcRRSP) {
		t.Error("expected a call to Has() for existing source to return true")
	}

	is := NewIncomeSourceSet(IncSrcRDSP)
	if !is.Has(IncSrcRDSP) {
		t.Error("expected a call to Has() for existing source to return true")
	}

}
