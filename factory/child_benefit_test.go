package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/benefits"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

func TestNewChildBenefitFactory_NewCalculator_SingleFormula(t *testing.T) {

	f := NewChildBenefitFactory(2018, core.RegionBC)
	c, err := f.NewCalculator()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := c.(*benefits.ChildBenfitCalculator)
	if !ok {
		t.Fatalf("unexpected type\nwant: %T\n got: %T", (&benefits.ChildBenfitCalculator{}), c)
	}

}

func TestChildBenefitFactory_NewCalculator_Errors(t *testing.T) {

	_, err := (&ChildBenefitFactory{}).NewCalculator()
	if err != ErrFactoryNotInit {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFactoryNotInit, err)
	}
}

func TestNewChildBenefitCalculatorFactory_Errors(t *testing.T) {

	cases := []struct {
		name    string
		year    uint
		regions []core.Region
		err     error
	}{
		{
			name:    "invalid-year",
			year:    1000,
			regions: []core.Region{core.RegionBC},
			err:     history.ErrParamsNotExist,
		},
		{
			name:    "invalid-region",
			year:    2018,
			regions: []core.Region{core.Region("1000")},
			err:     history.ErrRegionNotExist,
		},
		{
			name:    "no-regions",
			year:    2018,
			regions: nil,
			err:     benefits.ErrNoFormula,
		},
		{
			name:    "valid",
			year:    2018,
			regions: []core.Region{core.RegionBC},
			err:     nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			f := NewChildBenefitFactory(c.year, c.regions...)
			_, err := f.NewCalculator()
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}
