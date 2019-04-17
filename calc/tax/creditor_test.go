package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/calc/finance"
)

func TestCreditSources_makeSetAndGetDuplicates(t *testing.T) {

	cases := []struct {
		name        string
		sources     creditSourceControlGroup
		expectedDup []finance.CreditSource
		expectedSet map[finance.CreditSource]struct{}
	}{
		{
			name: "one-duplicate",
			sources: creditSourceControlGroup{
				{Source: 1},
				{Source: 2},
				{Source: 3},
				{Source: 3},
			},
			expectedDup: []finance.CreditSource{3},
			expectedSet: map[finance.CreditSource]struct{}{
				1: struct{}{}, 2: struct{}{}, 3: struct{}{},
			},
		},
		{
			name: "no-duplicates",
			sources: creditSourceControlGroup{
				{Source: 1},
				{Source: 2},
				{Source: 3},
			},
			expectedDup: []finance.CreditSource{},
			expectedSet: map[finance.CreditSource]struct{}{
				1: struct{}{}, 2: struct{}{}, 3: struct{}{},
			},
		},
		{
			name: "all-duplicates",
			sources: creditSourceControlGroup{
				{Source: 1},
				{Source: 1},
				{Source: 1},
			},
			expectedDup: []finance.CreditSource{1, 1},
			expectedSet: map[finance.CreditSource]struct{}{1: struct{}{}},
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

func TestConstCreditor(t *testing.T) {

	cc := ConstCreditor{
		Const: finance.TaxCredit{Amount: 1000.0, Source: 2},
	}

	actualSrc := cc.Source()
	expectedSrc := finance.CreditSource(2)
	if actualSrc != expectedSrc {
		t.Errorf("unexpected source\nwant: %v\n got: %v", expectedSrc, actualSrc)
	}

	actualCr := cc.TaxCredit(0, 0)
	expectedCr := finance.TaxCredit{Amount: 1000.0, Source: 2}
	if actualCr != expectedCr {
		t.Errorf("unexpected source\nwant: %v\n got: %v", expectedCr, actualCr)
	}

	clone := cc.Clone()
	cc.Const = finance.TaxCredit{}
	if clone.TaxCredit(0, 0) == cc.Const {
		t.Error("expected change to original to not affect clone")
	}
}
