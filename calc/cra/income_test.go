package cra

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/facts"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestTaxCalculator_Calc(t *testing.T) {

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

			taxParamsProv, err := history.GetProvincial(c.year, c.prov)
			if err != nil {
				t.Fatal(err)
			}

			taxParamsFed, err := history.GetFederal(c.year)
			if err != nil {
				t.Fatal(err)
			}

			taxParams := facts.Facts{
				Year:      c.year,
				FactsProv: taxParamsProv,
				FactsFed:  taxParamsFed,
			}
			finNums := FinancialNumbers{
				Income:         c.income,
				DeductionsFed:  c.deductionsFed,
				DeductionsProv: c.deductionsProv,
				CreditsFed:     c.creditsFed,
				CreditsProv:    c.creditsProv,
			}

			calculator, err := NewTaxCalculator(finNums, taxParams)
			if err != nil {
				t.Fatal(err)
			}

			actualTaxProv := calculator.CalcTaxProvincial()
			if !areEqual(actualTaxProv, c.expectedTaxProv, c.errMargin) {
				t.Errorf(
					"difference between actual and expected provincial "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTaxProv, actualTaxProv,
				)
			}

			actualTaxFed := calculator.CalcTaxFederal()
			if !areEqual(actualTaxFed, c.expectedTaxFed, c.errMargin) {
				t.Errorf(
					"difference between actual and expected federal "+
						"tax exceeds error margin\nwant: %04f\n got: %04f",
					c.expectedTaxFed, actualTaxFed,
				)
			}

			actualTaxTot := calculator.CalcTaxTotal()
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

func TestNewTaxCalculator_Error(t *testing.T) {

	invalidTaxParams := facts.Facts{
		FactsFed: facts.FactsFed{
			facts.BracketRates{
				0.10: facts.Bracket{-100, 200},
			},
		},
	}

	_, err := NewTaxCalculator(FinancialNumbers{}, invalidTaxParams)
	cause := errors.Cause(err)
	if cause != facts.ErrValNeg {
		t.Errorf("unexpected error\nwant: %v\n got: %v", facts.ErrValNeg, err)
	}

}
