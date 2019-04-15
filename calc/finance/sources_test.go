package finance

import "testing"

func TestSets_Has(t *testing.T) {

	ds := DeductionSourceSet{DeducSrcRRSP: struct{}{}}
	if !ds.Has(DeducSrcRRSP) {
		t.Error("expected a call to Has() for existing source to return true")
	}

	is := IncomeSourceSet{IncSrcRDSP: struct{}{}}
	if !is.Has(IncSrcRDSP) {
		t.Error("expected a call to Has() for existing source to return true")
	}

	ms := MiscSourceSet{MiscSrcMedical: struct{}{}}
	if !ms.Has(MiscSrcMedical) {
		t.Error("expected a call to Has() for existing source to return true")
	}

}
