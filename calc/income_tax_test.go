package calc

import (
	"fmt"
	"testing"

	"github.com/malkhamis/tax/facts"
	"github.com/malkhamis/tax/facts/history"
	"github.com/pkg/errors"
)

func TestIncomeTaxCalculator_Calc(t *testing.T) {

	cases := []struct {
		year            uint
		prov            history.Province
		income          float64
		deductionsFed   float64
		deductionsProv  float64
		creditsFed      float64
		creditsProv     float64
		expectedTaxProv float64
		expectedTaxFed  float64
		expectedTaxTot  float64
		errMargin       float64
	}{
		{
			year:            2018,
			prov:            history.BC,
			income:          300000,
			expectedTaxProv: 39156,
			expectedTaxFed:  76969,
			expectedTaxTot:  116125,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			prov:            history.BC,
			income:          0,
			expectedTaxProv: 0,
			expectedTaxFed:  0,
			expectedTaxTot:  0,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			prov:            history.BC,
			income:          9000,
			expectedTaxProv: 0,
			expectedTaxFed:  0,
			expectedTaxTot:  0,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			prov:            history.BC,
			income:          12000,
			expectedTaxProv: 80.35,
			expectedTaxFed:  28.65,
			expectedTaxTot:  109,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			prov:            history.BC,
			income:          85000,
			expectedTaxProv: 5128,
			expectedTaxFed:  13090,
			expectedTaxTot:  18218,
			errMargin:       1e-9,
		},
	}

	for i, c := range cases {
		i, c := i, c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			f, err := history.Get(c.year, c.prov)
			if err != nil {
				t.Fatal(err)
			}

			calculator := &IncomeTaxCalculator{
				Income:         c.income,
				DeductionsFed:  c.deductionsFed,
				DeductionsProv: c.deductionsProv,
				CreditsFed:     c.creditsFed,
				CreditsProv:    c.creditsProv,
				Facts:          f,
			}

			actualTaxProv, err := calculator.CalcProv()
			if err != nil {
				t.Error(err)
			}
			if !areEqual(actualTaxProv, c.expectedTaxProv, c.errMargin) {
				t.Errorf(
					"difference between actual and expected provincial "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTaxProv, actualTaxProv,
				)
			}

			actualTaxFed, err := calculator.CalcFed()
			if err != nil {
				t.Error(err)
			}
			if !areEqual(actualTaxFed, c.expectedTaxFed, c.errMargin) {
				t.Errorf(
					"difference between actual and expected federal "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTaxFed, actualTaxFed,
				)
			}

			actualTaxTot, err := calculator.CalcTotal()
			if err != nil {
				t.Error(err)
			}
			if !areEqual(actualTaxTot, c.expectedTaxTot, c.errMargin) {
				t.Errorf(
					"difference between actual and expected total "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTaxTot, actualTaxTot,
				)
			}

		})

	}

}

func TestIncomeTaxCalculator_Calc_InvalidRatesFed(t *testing.T) {

	invalidFacts := facts.Facts{
		FactsFed: facts.FactsFed{
			facts.BracketRates{
				0.10: facts.Bracket{-100, 200},
			},
		},
	}

	calculator := IncomeTaxCalculator{
		Facts: invalidFacts,
	}

	_, err := calculator.CalcTotal()
	cause := errors.Cause(err)
	if cause != facts.ErrValNeg {
		t.Errorf("unexpected error\nwant: %v\n got: %v", facts.ErrValNeg, err)
	}

}

func TestIncomeTaxCalculator_Calc_InvalidRatesProv(t *testing.T) {

	invalidFacts := facts.Facts{
		FactsProv: facts.FactsProv{
			facts.BracketRates{
				0.10: facts.Bracket{-100, 200},
			},
		},
	}

	calculator := IncomeTaxCalculator{
		Facts: invalidFacts,
	}

	_, err := calculator.CalcTotal()
	cause := errors.Cause(err)
	if cause != facts.ErrValNeg {
		t.Errorf("unexpected error\nwant: %v\n got: %v", facts.ErrValNeg, err)
	}

}
