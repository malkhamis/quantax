// Package history provides historical tax params for various jurisdictions
package history

import "github.com/malkhamis/quantax/core"

var (
	taxParamsAll = map[core.Region]yearlyTaxParams{
		core.RegionBC: taxParamsBC,
		core.RegionCA: taxParamsCanada,
	}

	cbParamsAll = map[core.Region]yearlyCBParams{
		core.RegionBC: cbParamsBC,
		core.RegionCA: cbParamsCanada,
	}

	rrspParamsAll = map[core.Region]yearlyRRSPParams{
		core.RegionCA: rrspParamsCanada,
	}
)

// GetTaxParams returns a copy of the tax params for the given year and region
func GetTaxParams(year uint, region core.Region) (TaxParams, error) {

	jurisdictionParams, ok := taxParamsAll[region]
	if !ok {
		return TaxParams{}, ErrRegionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return TaxParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}

// GetChildBenefitParams returns a copy of the child benefit parameters for
// the given year and region
func GetChildBenefitParams(year uint, region core.Region) (CBParams, error) {

	jurisdictionParams, ok := cbParamsAll[region]
	if !ok {
		return CBParams{}, ErrRegionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return CBParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}

// GetRRSOParams returns a copy of the RRSP parameters for the given year/region
func GetRRSPParams(year uint, region core.Region) (RRSPParams, error) {

	jurisdictionParams, ok := rrspParamsAll[region]
	if !ok {
		return RRSPParams{}, ErrRegionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return RRSPParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}
