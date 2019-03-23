package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestNewIncomeTaxCalculator_Errors(t *testing.T) {

	cases := []struct {
		name   string
		params IncomeTaxParams
		err    error
	}{
		{
			name:   "invalid-year",
			params: IncomeTaxParams{Year: 1000, Region: history.BC},
			err:    history.ErrFormulaNotExist,
		},
		{
			name:   "invalid-region",
			params: IncomeTaxParams{Year: 2018, Region: history.Jurisdiction("o'lala")},
			err:    history.ErrJurisdictionNotExist,
		},
		{
			name:   "valid",
			params: IncomeTaxParams{Year: 2018, Region: history.BC},
			err:    nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {
			dummy := calc.IndividualFinances{}
			_, err := NewIncomeTaxCalculator(dummy, c.params)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}
		})
	}
}
