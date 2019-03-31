package rrsp

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestMaxCapper_Contribution(t *testing.T) {

	cases := []struct {
		maxCapper *MaxCapper
		income    float64
		expected  float64
	}{
		{
			maxCapper: &MaxCapper{Rate: 0.15, Cap: 2000.0, IncomeType: finance.EARNED},
			income:    0.0,
			expected:  0.0,
		},
		{
			maxCapper: &MaxCapper{Rate: 0.15, Cap: 2000.0, IncomeType: finance.EARNED},
			income:    1000.0,
			expected:  150.0,
		},
		{
			maxCapper: &MaxCapper{Rate: 0.15, Cap: 2000.0, IncomeType: finance.EARNED},
			income:    13333.34,
			expected:  2000.0,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			actual := c.maxCapper.Contribution(c.income)
			if actual != c.expected {
				t.Errorf("unexpected contribution\nwant: %.2f\n got: %.2f", c.expected, actual)
			}

		})
	}
}

func TestMaxCapper_IncomeCalcMethod(t *testing.T) {
	mc := &MaxCapper{IncomeType: finance.EARNED}
	it := mc.IncomeCalcMethod()
	if it != finance.EARNED {
		t.Fatalf("unexpected income type\n want: %s\n got: %s", finance.EARNED, it)
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
		Rate:       0.10,
		Cap:        1000,
		IncomeType: finance.EARNED,
	}

	income := 5000.0
	originaResults := original.Contribution(income)

	clone := original.Clone()
	original.Rate = 0.25

	cloneResults := clone.Contribution(income)
	if cloneResults != originaResults {
		t.Fatalf("expected changes to original formula to not affect clone formula")
	}
}
