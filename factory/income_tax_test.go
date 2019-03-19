package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestNewIncomeTaxCalculatorAgg_Errors(t *testing.T) {

	cases := []struct {
		name   string
		first  IncomeTaxParams
		second IncomeTaxParams
		extras []IncomeTaxParams
		err    error
	}{
		{
			name:  "invalid-first",
			first: IncomeTaxParams{Year: 1000, Region: history.BC},
			err:   history.ErrNoRates,
		},
		{
			name:   "invalid-second",
			first:  IncomeTaxParams{Year: 2018, Region: history.BC},
			second: IncomeTaxParams{Year: 1000, Region: history.BC},
			err:    history.ErrNoRates,
		},
		{
			name:   "invalid-extra",
			first:  IncomeTaxParams{Year: 2018, Region: history.BC},
			second: IncomeTaxParams{Year: 2018, Region: history.Canada},
			extras: []IncomeTaxParams{{Year: 1000, Region: history.BC}},
			err:    history.ErrNoRates,
		},
		{
			name:   "valid",
			first:  IncomeTaxParams{Year: 2018, Region: history.BC},
			second: IncomeTaxParams{Year: 2018, Region: history.Canada},
			extras: []IncomeTaxParams{{Year: 2018, Region: history.Canada}},
			err:    nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {
			dummy := calc.FinancialNumbers{}
			_, err := NewIncomeTaxCalculatorAgg(dummy, c.first, c.second, c.extras...)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}
		})
	}
}
