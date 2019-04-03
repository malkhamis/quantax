package tax

import (
	"fmt"
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

func TestCalculator_Calc(t *testing.T) {

	formulaCanada2018 := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			-0.150: finance.Bracket{0, 11809},
			0.150:  finance.Bracket{0, 46605},
			0.205:  finance.Bracket{46606, 93208},
			0.260:  finance.Bracket{93209, 144489},
			0.290:  finance.Bracket{144490, 205842},
			0.330:  finance.Bracket{205843, math.Inf(1)},
		},
		ExcludedIncome:     []finance.IncomeSource{finance.IncSrcTFSA},
		ExcludedDeductions: []finance.DeductionSource{finance.DeducSrcMedical},
	}

	cases := []struct {
		finances    finance.IndividualFinances
		formula     Formula
		expectedTax float64
		errMargin   float64
	}{
		{
			finances: finance.IndividualFinances{
				Income:     finance.IncomeBySource{finance.IncSrcEarned: 400000},
				Deductions: finance.DeductionBySource{finance.DeducSrcRRSP: 100000},
			},
			formula:     formulaCanada2018,
			expectedTax: 76969,
			errMargin:   1e-9,
		},
		{
			finances: finance.IndividualFinances{
				Income: finance.IncomeBySource{
					finance.IncSrcEarned: 0,
					finance.IncSrcTFSA:   50000,
				},
				Deductions: finance.DeductionBySource{
					finance.DeducSrcRRSP:    100000,
					finance.DeducSrcMedical: 1000,
				},
			},
			formula:     formulaCanada2018,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			finances: finance.IndividualFinances{
				Income: finance.IncomeBySource{finance.IncSrcEarned: 9000},
			},
			formula:     formulaCanada2018,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			finances: finance.IndividualFinances{
				Income: finance.IncomeBySource{
					finance.IncSrcEarned: 12000,
					finance.IncSrcTFSA:   13000,
				},
			},
			formula:     formulaCanada2018,
			expectedTax: 28.65,
			errMargin:   1e-9,
		},
		{
			finances: finance.IndividualFinances{
				Income:     finance.IncomeBySource{finance.IncSrcEarned: 85000},
				Deductions: finance.DeductionBySource{finance.DeducSrcMedical: 1000},
			},
			formula:     formulaCanada2018,
			expectedTax: 13090,
			errMargin:   1e-9,
		},
		{
			finances: finance.IndividualFinances{
				Income:     finance.IncomeBySource{finance.IncSrcEarned: 90000},
				Deductions: finance.DeductionBySource{finance.DeducSrcRRSP: 5000},
			},
			formula:     formulaCanada2018,
			expectedTax: 13090,
			errMargin:   1e-9,
		},
	}

	for i, c := range cases {
		i, c := i, c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			calculator, err := NewCalculator(c.formula)
			if err != nil {
				t.Fatal(err)
			}

			actualTax := calculator.Calc(&c.finances)
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

func TestNewCalculator_InvalidFormula(t *testing.T) {

	invalidTaxParams := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			0.10: finance.Bracket{300, 200},
		},
	}
	_, err := NewCalculator(invalidTaxParams)
	cause := errors.Cause(err)
	if cause != finance.ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %v", finance.ErrValNeg, err)
	}

}

func TestCalculator_NilFormula(t *testing.T) {

	_, err := NewCalculator(nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}
}
