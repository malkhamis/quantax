package tax

import (
	"testing"

	"github.com/malkhamis/quantax/core"
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
				TaxFormula:       &testTaxFormula{},
				ContraTaxFormula: &testContraTaxFormula{},
			},
			err: nil,
		},
		//
		{
			name: "invalid-tax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       &testTaxFormula{onValidate: simulatedErr},
				ContraTaxFormula: &testContraTaxFormula{},
			},
			err: simulatedErr,
		},
		//
		{
			name: "invalid-contratax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       &testTaxFormula{},
				ContraTaxFormula: &testContraTaxFormula{onValidate: simulatedErr},
			},
			err: simulatedErr,
		},
		//
		{
			name: "nil-tax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       nil,
				ContraTaxFormula: &testContraTaxFormula{},
			},
			err: ErrNoFormula,
		},
		//
		{
			name: "nil-contratax-formula",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       &testTaxFormula{},
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
				ContraTaxFormula: &testContraTaxFormula{},
			},
			err: ErrNoFormula,
		},
		{
			name: "year-mismatch",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       &testTaxFormula{onYear: 2019},
				ContraTaxFormula: &testContraTaxFormula{onYear: 2020},
			},
			err: ErrInvalidTaxArg,
		},
		{
			name: "region-mismatch",
			cfg: CalcConfig{
				IncomeCalc:       testIncomeCalculator{},
				TaxFormula:       &testTaxFormula{onRegion: core.Region("Somewhere")},
				ContraTaxFormula: &testContraTaxFormula{onRegion: core.Region("NoMan")},
			},
			err: ErrInvalidTaxArg,
		},
		{
			name: "nil-income-calculator",
			cfg: CalcConfig{
				IncomeCalc:       nil,
				TaxFormula:       &testTaxFormula{},
				ContraTaxFormula: &testContraTaxFormula{},
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
