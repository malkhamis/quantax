package incometax

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestCalculator_Calc(t *testing.T) {

	cases := []struct {
		year        uint
		region      history.Jurisdiction
		taxableAmnt float64
		deductions  float64
		credits     float64
		expectedTax float64
		errMargin   float64
	}{
		{
			year:        2018,
			region:      history.Canada,
			taxableAmnt: 300000,
			expectedTax: 76969,
			errMargin:   1e-9,
		},
		{
			year:        2018,
			region:      history.Canada,
			taxableAmnt: 0,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			year:        2018,
			region:      history.Canada,
			taxableAmnt: 9000,
			expectedTax: 0,
			errMargin:   1e-9,
		},
		{
			year:        2018,
			region:      history.Canada,
			taxableAmnt: 12000,
			expectedTax: 28.65,
			errMargin:   1e-9,
		},
		{
			year:        2018,
			region:      history.Canada,
			taxableAmnt: 85000,
			expectedTax: 13090,
			errMargin:   1e-9,
		},
	}

	for i, c := range cases {
		i, c := i, c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			taxRates, err := history.Get(c.year, c.region)
			if err != nil {
				t.Fatal(err)
			}

			finNums := calc.Finances{
				TaxableAmount: c.taxableAmnt,
				Deductions:    c.deductions,
				Credits:       c.credits,
			}

			calculator, err := NewCalculator(finNums, taxRates)
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

	invalidTaxParams := calc.BracketRates{0.10: calc.Bracket{-100, 200}}
	_, err := NewCalculator(calc.Finances{}, invalidTaxParams)
	cause := errors.Cause(err)
	if cause != calc.ErrValNeg {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrValNeg, err)
	}

}

func TestCalculator_Setters(t *testing.T) {

	c := Calculator{
		Finances: calc.Finances{
			TaxableAmount: 0,
			Credits:       0,
			Deductions:    0,
		},
	}

	newFinNums := calc.Finances{
		TaxableAmount: 10,
		Credits:       20,
		Deductions:    30,
	}

	c.Update(newFinNums)

	if c.TaxableAmount != 10 {
		t.Error("expected Update() to mutate the calculator")
	}
	if c.Credits != 20 {
		t.Error("expected Update() to mutate the calculator")
	}
	if c.Deductions != 30 {
		t.Error("expected Update() to mutate the calculator")
	}
}
