package finance

import (
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

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

func TestTaxCreditGroup_MergeSimilars(t *testing.T) {

	crs := TaxCreditGroup{
		{1, 10},
		{1, 20},
		{2, 30},
		{1, 40},
		{2, 50},
		{3, 60},
	}

	expected := []TaxCredit{{1, 70}, {2, 80}, {3, 60}}
	actual := crs.MergeSimilars()

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Source < expected[j].Source
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].Source < actual[j].Source
	})

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestTaxCreditGroup_MergeSimilars_EdgeCases(t *testing.T) {

	crs := TaxCreditGroup{}

	expected := []TaxCredit{}
	actual := crs.MergeSimilars()

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	actual = (TaxCreditGroup)(nil).MergeSimilars()

	diff = deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}
