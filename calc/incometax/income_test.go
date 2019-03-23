package incometax

import (
	"fmt"
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

func TestCalculator_Calc(t *testing.T) {

	formulaCanada2018 := calc.WeightedBracketFormula{
		-0.150: calc.Bracket{0, 11809},
		0.150:  calc.Bracket{0, 46605},
		0.205:  calc.Bracket{46606, 93208},
		0.260:  calc.Bracket{93209, 144489},
		0.290:  calc.Bracket{144490, 205842},
		0.330:  calc.Bracket{205843, math.Inf(1)},
	}

	cases := []struct {
		finances    IndividualFinances
		formula     TaxFormula
		expectedTax float64
		errMargin   float64
	}{
		{
			finances:    IndividualFinances{Income: 300000},
			formula:     formulaCanada2018,
			expectedTax: 76969,
			errMargin:   1e-9,
		},
		{
			finances:    IndividualFinances{},
			formula:     formulaCanada2018,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			finances:    IndividualFinances{Income: 9000},
			formula:     formulaCanada2018,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			finances:    IndividualFinances{Income: 12000},
			formula:     formulaCanada2018,
			expectedTax: 28.65,
			errMargin:   1e-9,
		},
		{
			finances:    IndividualFinances{Income: 85000},
			formula:     formulaCanada2018,
			expectedTax: 13090,
			errMargin:   1e-9,
		},
	}

	for i, c := range cases {
		i, c := i, c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			calculator, err := NewCalculator(c.finances, c.formula)
			if err != nil {
				t.Fatal(err)
			}

			actualTax := calculator.Calc()
			if !areEqual(actualTax, c.expectedTax, c.errMargin) {
				t.Errorf(
					"difference between actual and expected total "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTax, actualTax,
				)
			}

		})

	}
}

func TestNewCalculator_Error(t *testing.T) {

	invalidTaxParams := calc.WeightedBracketFormula{0.10: calc.Bracket{300, 200}}
	_, err := NewCalculator(IndividualFinances{}, invalidTaxParams)
	cause := errors.Cause(err)
	if cause != calc.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrValNeg, err)
	}

}

func TestCalculator_Setters(t *testing.T) {

	c := Calculator{
		IndividualFinances: IndividualFinances{
			Income:     0,
			Deductions: 0,
		},
	}

	newFinNums := IndividualFinances{
		Income:     10,
		Deductions: 30,
	}

	c.UpdateFinances(newFinNums)

	if c.Income != 10 {
		t.Error("expected Update() to mutate the calculator")
	}
	if c.Deductions != 30 {
		t.Error("expected Update() to mutate the calculator")
	}
}
