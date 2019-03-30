package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
	"github.com/pkg/errors"
)

func TestNewTaxCalculatorFactory_NewCalculator_SingleFormula(t *testing.T) {

	f, err := NewTaxCalcFactory(2018, BC)
	if err != nil {
		t.Fatal(err)
	}

	c, err := f.NewCalculator()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := c.(*tax.Calculator)
	if !ok {
		t.Fatalf("unexpected type\nwant: %T\n got: %T", (&tax.Calculator{}), c)
	}

}

func TestTaxCalcFactory_NewCalculator_Errors(t *testing.T) {

	_, err := (&TaxCalcFactory{}).NewCalculator()
	if err != ErrFactoryNotInit {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFactoryNotInit, err)
	}

}

func TestNewChildBenefitCalculatorFactory_NewCalculator_SingleFormula(t *testing.T) {

	f, err := NewChildBenefitCalcFactory(2017, BC)
	if err != nil {
		t.Fatal(err)
	}

	c, err := f.NewCalculator()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := c.(*benefits.Calculator)
	if !ok {
		t.Fatalf("unexpected type\nwant: %T\n got: %T", (&benefits.Calculator{}), c)
	}

}

func TestChildBenefitCalcFactory_NewCalculator_Errors(t *testing.T) {

	_, err := (&ChildBenefitCalcFactory{}).NewCalculator()
	if err != ErrFactoryNotInit {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFactoryNotInit, err)
	}

}

func TestNewTaxCalculatorFactory_Errors(t *testing.T) {

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
			err:     history.ErrFormulaNotExist,
		},
		{
			name:    "invalid-region",
			year:    1000,
			regions: []Region{Region(1000)},
			err:     ErrRegionNotExist,
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

			_, err := NewTaxCalcFactory(c.year, c.regions[0], c.regions[1:]...)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}

func TestNewChildBenefitCalculatorFactory_Errors(t *testing.T) {

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
			err:     history.ErrFormulaNotExist,
		},
		{
			name:    "invalid-region",
			year:    2018,
			regions: []Region{Region(1000)},
			err:     ErrRegionNotExist,
		},
		{
			name:    "valid",
			year:    2017,
			regions: []Region{BC},
			err:     nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			_, err := NewChildBenefitCalcFactory(c.year, c.regions[0], c.regions...)
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}
