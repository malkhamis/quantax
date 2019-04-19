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
