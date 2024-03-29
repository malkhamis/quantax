package factory

import (
	"fmt"
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/history"

	"github.com/pkg/errors"
)

func TestRRSPFactory_Uninitialized(t *testing.T) {

	_, err := (&RRSPFactory{}).NewCalculator()
	if err != ErrFactoryNotInit {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFactoryNotInit, err)
	}

}

func TestRRSPFactory_Errors(t *testing.T) {

	cases := []struct {
		name   string
		config RRSPFactoryConfig
		err    error
	}{
		{
			name: "invalid-year",
			config: RRSPFactoryConfig{
				Year:       1000,
				RRSPRegion: core.RegionCA,
				TaxRegions: []core.Region{core.RegionCA, core.RegionBC},
			},
			err: history.ErrParamsNotExist,
		},
		{
			name: "invalid-rrsp-region",
			config: RRSPFactoryConfig{
				Year:       2018,
				RRSPRegion: "1000",
				TaxRegions: []core.Region{core.RegionBC},
			},
			err: history.ErrRegionNotExist,
		},
		{
			name: "invalid-tax-region",
			config: RRSPFactoryConfig{
				Year:       2018,
				RRSPRegion: core.RegionCA,
				TaxRegions: []core.Region{"1000"},
			},
			err: history.ErrRegionNotExist,
		},
		{
			name: "valid",
			config: RRSPFactoryConfig{
				Year:       2018,
				RRSPRegion: core.RegionCA,
				TaxRegions: []core.Region{core.RegionCA, core.RegionBC},
			},
			err: nil,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			f := NewRRSPFactory(c.config)
			_, err := f.NewCalculator()
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}
