package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

func TestNewTaxFactory_NewCalculator_SingleFormula(t *testing.T) {

	f := NewTaxFactory(2018, BC)
	c, err := f.NewCalculator()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := c.(*tax.Calculator)
	if !ok {
		t.Fatalf("unexpected type\nwant: %T\n got: %T", (&tax.Calculator{}), c)
	}

}

func TestTaxFactory_Uninitialized(t *testing.T) {

	_, err := (&TaxFactory{}).NewCalculator()
	if err != ErrFactoryNotInit {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFactoryNotInit, err)
	}

}

func TestTaxFactory_Errors(t *testing.T) {

	cases := []struct {
		name    string
		year    uint
		regions []Region
		err     error
	}{
		{
			name:    "invalid-year",
			year:    1000,
			regions: []Region{BC},
			err:     history.ErrParamsNotExist,
		},
		{
			name:    "invalid-region",
			year:    1000,
			regions: []Region{Region(1000)},
			err:     ErrRegionNotExist,
		},
		{
			name:    "no-regions",
			year:    1000,
			regions: nil,
			err:     tax.ErrNoFormula,
		},
		{
			name:    "valid",
			year:    2018,
			regions: []Region{BC},
			err:     nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			f := NewTaxFactory(c.year, c.regions...)
			_, err := f.NewCalculator()
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}
