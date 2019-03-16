package tax

import (
	"fmt"
	"testing"
)

func TestIncomeTaxCalculator_Calc(t *testing.T) {

	cases := []struct {
		year            uint
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
			income:          300000,
			expectedTaxProv: 39156,
			expectedTaxFed:  76969,
			expectedTaxTot:  116125,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			income:          0,
			expectedTaxProv: 0,
			expectedTaxFed:  0,
			expectedTaxTot:  0,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			income:          9000,
			expectedTaxProv: 0,
			expectedTaxFed:  0,
			expectedTaxTot:  0,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			income:          12000,
			expectedTaxProv: 80.35,
			expectedTaxFed:  28.65,
			expectedTaxTot:  109,
			errMargin:       1e-9,
		},
		{
			year:            2018,
			income:          85000,
			expectedTaxProv: 5128,
			expectedTaxFed:  13090,
			expectedTaxTot:  18218,
			errMargin:       1e-9,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			calculator := &IncomeTaxCalculator{
				Income:         c.income,
				DeductionsFed:  c.deductionsFed,
				DeductionsProv: c.deductionsProv,
				CreditsFed:     c.creditsFed,
				CreditsProv:    c.creditsProv,
				Facts:          facts[c.year],
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
