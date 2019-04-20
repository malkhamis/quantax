package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestCreditSources_makeSetAndGetDuplicates(t *testing.T) {

	cases := []struct {
		name        string
		sources     creditRuleGroup
		expectedDup []string
		expectedSet map[string]struct{}
	}{
		{
			name: "one-duplicate",
			sources: creditRuleGroup{
				{Source: "1"},
				{Source: "2"},
				{Source: "3"},
				{Source: "3"},
			},
			expectedDup: []string{"3"},
			expectedSet: map[string]struct{}{
				"1": struct{}{}, "2": struct{}{}, "3": struct{}{},
			},
		},
		{
			name: "no-duplicates",
			sources: creditRuleGroup{
				{Source: "1"},
				{Source: "2"},
				{Source: "3"},
			},
			expectedDup: []string{},
			expectedSet: map[string]struct{}{
				"1": struct{}{}, "2": struct{}{}, "3": struct{}{},
			},
		},
		{
			name: "all-duplicates",
			sources: creditRuleGroup{
				{Source: "1"},
				{Source: "1"},
				{Source: "1"},
			},
			expectedDup: []string{"1", "1"},
			expectedSet: map[string]struct{}{"1": struct{}{}},
		},
	}

	for i, c := range cases {

		actualSet, actualDup := c.sources.makeSrcSetAndGetDuplicates()
		diff := deep.Equal(actualSet, c.expectedSet)
		if diff != nil {
			t.Errorf(
				"case %d-%s: actual does not match expected\n%s",
				i, c.name, strings.Join(diff, "\n"),
			)
		}

		diff = deep.Equal(actualDup, c.expectedDup)
		if diff != nil {
			t.Errorf(
				"case %d-%s: actual does not match expected\n%s",
				i, c.name, strings.Join(diff, "\n"),
			)
		}

	}

}

func Test_taxCredit_Amount(t *testing.T) {

	tc := &taxCredit{amount: 10}
	if tc.Amount() != tc.amount {
		t.Fatalf(
			"unexpected tax credit amount\nwant: %.2f\n got: %.2f",
			tc.amount, tc.Amount(),
		)
	}
}

func Test_taxCredit_Source(t *testing.T) {

	tc := &taxCredit{rule: CreditRule{Source: t.Name()}}
	if tc.Source() != tc.rule.Source {
		t.Fatalf(
			"unexpected tax credit source\nwant: %q\n got: %q",
			tc.rule.Source, tc.Source(),
		)
	}
}

func Test_taxCredit_clone(t *testing.T) {

	c := new(Calculator)
	original := &taxCredit{
		amount: 123,
		owner:  c,
		rule:   CreditRule{Source: t.Name(), Type: CrRuleTypeCashable},
	}

	clone := original.clone()
	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	original.owner = nil
	if clone.owner == nil {
		t.Fatal("expected changes to original to not affect clone")
	}
}

func Test_taxCredit_clone_nil(t *testing.T) {

	var nilTaxCredit *taxCredit
	clone := nilTaxCredit.clone()
	if clone != nil {
		t.Fatal("cloning a nil tax credit should return nil")
	}
}

func Test_taxCreditGroup_clone(t *testing.T) {

	c := new(Calculator)
	original := []*taxCredit{
		&taxCredit{
			amount: 123,
			owner:  c,
			rule:   CreditRule{Source: t.Name(), Type: CrRuleTypeCashable},
		},
	}

	clone := taxCreditGroup(original).clone()
	diff := deep.Equal(original, clone)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}
