package tax

import (
	"testing"

	"github.com/pkg/errors"
)

func TestCalcConfig_validate(t *testing.T) {

	simulatedErr := errors.New("simulated error")
	cases := []struct {
		name string
		cfg  CalcConfig
		err  error
	}{
		//
		{
			name: "valid",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       testTaxFormula{},
				ContraTaxFormula: &testTaxContraFormula{},
			},
			err: nil,
		},
		//
		{
			name: "invalid-tax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       testTaxFormula{onValidate: simulatedErr},
				ContraTaxFormula: &testTaxContraFormula{},
			},
			err: simulatedErr,
		},
		//
		{
			name: "invalid-contratax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       testTaxFormula{},
				ContraTaxFormula: &testTaxContraFormula{onValidate: simulatedErr},
			},
			err: simulatedErr,
		},
		//
		{
			name: "nil-tax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       nil,
				ContraTaxFormula: &testTaxContraFormula{},
			},
			err: ErrNoFormula,
		},
		//
		{
			name: "nil-contratax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       testTaxFormula{},
				ContraTaxFormula: nil,
			},
			err: ErrNoContraFormula,
		},
		//
		{
			name: "invalid-contratax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       nil,
				ContraTaxFormula: &testTaxContraFormula{},
			},
			err: ErrNoFormula,
		},
		{
			name: "nil-income-calculator",
			cfg: CalcConfig{
				IncomeCalc:       nil,
				TaxFormula:       testTaxFormula{},
				ContraTaxFormula: &testTaxContraFormula{},
			},
			err: ErrNoIncCalc,
		},
	}

	for i, c := range cases {
		err := c.cfg.validate()
		if errors.Cause(err) != c.err {
			t.Errorf(
				"case-%d-%s: unexpected error\nwant: %v\n got: %v",
				i, c.name, c.err, err,
			)
		}
	}

}
