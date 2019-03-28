package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestNewTaxCalculatorFactory_Errors(t *testing.T) {

	cases := []struct {
		name string
		opts Options
		err  error
	}{
		{
			name: "invalid-year",
			opts: Options{Year: 1000, Region: BC},
			err:  history.ErrFormulaNotExist,
		},
		{
			name: "invalid-region",
			opts: Options{Year: 2018, Region: Region(1000)},
			err:  ErrRegionNotExist,
		},
		{
			name: "valid",
			opts: Options{Year: 2018, Region: BC},
			err:  nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			_, err := NewTaxCalcFactory(c.opts)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}

func TestNewChildBenefitCalculatorFactory_Errors(t *testing.T) {

	cases := []struct {
		name string
		opts Options
		err  error
	}{
		{
			name: "invalid-year",
			opts: Options{Year: 1000, Region: BC},
			err:  history.ErrFormulaNotExist,
		},
		{
			name: "invalid-region",
			opts: Options{Year: 2018, Region: Region(1000)},
			err:  ErrRegionNotExist,
		},
		{
			name: "valid",
			opts: Options{Year: 2017, Region: BC},
			err:  nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			_, err := NewChildBenefitCalcFactory(c.opts)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}
