package rrsp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core/finance"
)

func TestMaxCapper_Contribution(t *testing.T) {

	formula := &MaxCapper{
		Cap:  1000.0,
		Rate: 0.10,
	}

	actual := formula.ContributionEarned(500)
	expected := 50.0
	if actual != expected {
		t.Errorf("unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual)
	}

	actual = formula.ContributionEarned(1000)
	expected = 100.0
	if actual != expected {
		t.Errorf("unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual)
	}

	actual = formula.ContributionEarned(10000)
	expected = 1000.0
	if actual != expected {
		t.Errorf("unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual)
	}

	actual = formula.ContributionEarned(100000)
	expected = 1000.0
	if actual != expected {
		t.Errorf("unexpected contribution\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

func TestMaxCapper_Validate(t *testing.T) {
	err := (&MaxCapper{}).Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMaxCapper_Clone(t *testing.T) {

	original := MaxCapper{
		Rate:                           0.10,
		Cap:                            1000,
		IncomeSources:                  []finance.IncomeSource{finance.IncSrcEarned},
		DeductionSourceForContribution: finance.DeductionSource(2000),
		IncomeSourceForWithdrawal:      finance.IncomeSource(1000),
	}

	income := 5000.0
	originaResults := original.ContributionEarned(income)

	clone := original.Clone()
	original.Rate = 0.25

	cloneResults := clone.ContributionEarned(income)
	if cloneResults != originaResults {
		t.Fatalf("expected changes to original formula to not affect clone formula")
	}
}

func TestMaxCapper_AllowedIncomeSources(t *testing.T) {

	expected := []finance.IncomeSource{1, 2, 3, 6}
	f := &MaxCapper{
		IncomeSources: expected,
	}
	actual := f.AllowedIncomeSources()

	diff := deep.Equal(expected, actual)
	if diff != nil {
		t.Fatalf("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestMaxCapper_TargetSourceForWithdrawl(t *testing.T) {

	expected := finance.IncomeSource(1234)
	f := &MaxCapper{
		IncomeSourceForWithdrawal: expected,
	}
	actual := f.TargetSourceForWithdrawl()

	diff := deep.Equal(expected, actual)
	if diff != nil {
		t.Fatalf("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestMaxCapper_TargetSourceForContribution(t *testing.T) {

	expected := finance.DeductionSource(1234)
	f := &MaxCapper{
		DeductionSourceForContribution: expected,
	}
	actual := f.TargetSourceForContribution()

	diff := deep.Equal(expected, actual)
	if diff != nil {
		t.Fatalf("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestMaxCapper_NumFieldsUnchanged(t *testing.T) {

	dummy := MaxCapper{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 5 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
