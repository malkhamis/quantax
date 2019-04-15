package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestCreditSources_makeSetAndGetDuplicates(t *testing.T) {

	cases := []struct {
		name        string
		sources     creditSources
		expectedDup creditSources
		expectedSet map[CreditSource]struct{}
	}{
		{
			name:        "one-duplicate",
			sources:     creditSources{1, 2, 3, 3},
			expectedDup: creditSources{3},
			expectedSet: map[CreditSource]struct{}{
				1: struct{}{}, 2: struct{}{}, 3: struct{}{},
			},
		},
		{
			name:        "no-duplicates",
			sources:     creditSources{1, 2, 3},
			expectedDup: creditSources{},
			expectedSet: map[CreditSource]struct{}{
				1: struct{}{}, 2: struct{}{}, 3: struct{}{},
			},
		},
		{
			name:        "all-duplicates",
			sources:     creditSources{1, 1, 1},
			expectedDup: creditSources{1, 1},
			expectedSet: map[CreditSource]struct{}{1: struct{}{}},
		},
	}

	for i, c := range cases {

		actualSet, actualDup := c.sources.makeSetAndGetDuplicates()
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
		Const: Credits{Amount: 1000.0, Source: 2, IsRefundable: true},
	}

	actualSrc := cc.Source()
	expectedSrc := CreditSource(2)
	if actualSrc != expectedSrc {
		t.Errorf("unexpected source\nwant: %v\n got: %v", expectedSrc, actualSrc)
	}

	actualCr := cc.TaxCredits(0, 0)
	expectedCr := Credits{Amount: 1000.0, Source: 2, IsRefundable: true}
	if actualCr != expectedCr {
		t.Errorf("unexpected source\nwant: %v\n got: %v", expectedCr, actualCr)
	}

	clone := cc.Clone()
	cc.Const = Credits{}
	if clone.TaxCredits(0, 0) == cc.Const {
		t.Error("expected change to original to not affect clone")
	}
}
