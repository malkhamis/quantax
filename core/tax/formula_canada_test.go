package tax

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/core/finance"
)

func TestCanadianFormula_Apply(t *testing.T) {

	formulaCanada2018 := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			-0.150: finance.Bracket{0, 11809},
			0.150:  finance.Bracket{0, 46605},
			0.205:  finance.Bracket{46606, 93208},
			0.260:  finance.Bracket{93209, 144489},
			0.290:  finance.Bracket{144490, 205842},
			0.330:  finance.Bracket{205843, math.Inf(1)},
		},
	}
	err := formulaCanada2018.Validate()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		netIncome   float64
		expectedTax float64
	}{
		{
			netIncome:   300000.0,
			expectedTax: 76969,
		},
		{
			netIncome:   0.0,
			expectedTax: 0.0,
		},
		{
			netIncome:   9000.0,
			expectedTax: 0.0,
		},
		{
			netIncome:   12000.0,
			expectedTax: 28.65,
		},
	}

	for i, c := range cases {
		i, c := i, c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			actualTax := formulaCanada2018.Apply(c.netIncome)
			if !areEqual(actualTax, c.expectedTax, 1e-9) {
				t.Errorf(
					"difference between actual and expected total "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTax, actualTax,
				)
			}

		})

	}
}

func TestCanadianFormula_Clone(t *testing.T) {

	original := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{0.1: finance.Bracket{0, 10}},
	}

	clone := original.Clone()
	originalResults := original.Apply(5)
	original.WeightedBrackets[0.1] = finance.Bracket{100, 1000}
	cloneResults := clone.Apply(5)
	if originalResults != cloneResults {
		t.Errorf("expected clone results to be equal to results of original formula prior to mutation")
	}

}

func TestCanadianFormula_Clone_Nil(t *testing.T) {

	var mr *CanadianFormula
	clone := mr.Clone()
	if clone != nil {
		t.Fatal("cloning a nil formula should return nil")
	}
}

func TestCanadianFormula_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianFormula{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
