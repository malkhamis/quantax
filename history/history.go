// Package history provides historical tax params for various jurisdictions
package history

var (
	taxParamsAll = map[Jurisdiction]yearlyTaxParams{
		BC:     taxParamsBC,
		Canada: taxParamsCanada,
	}

	cbParamsAll = map[Jurisdiction]yearlyCBParams{
		BC:     cbParamsBC,
		Canada: cbParamsCanada,
	}

	rrspParamsAll = map[Jurisdiction]yearlyRRSPParams{
		Canada: rrspParamsCanada,
	}
)

// GetTaxParams returns a copy of the tax params for the given year and region
func GetTaxParams(year uint, region Jurisdiction) (TaxParams, error) {

	jurisdictionParams, ok := taxParamsAll[region]
	if !ok {
		return TaxParams{}, ErrJurisdictionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return TaxParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}

// GetChildBenefitParams returns a copy of the child benefit parameters for
// the given year and region
func GetChildBenefitParams(year uint, region Jurisdiction) (CBParams, error) {

	jurisdictionParams, ok := cbParamsAll[region]
	if !ok {
		return CBParams{}, ErrJurisdictionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return CBParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}

// GetRRSOParams returns a copy of the RRSP parameters for the given year/region
func GetRRSPParams(year uint, region Jurisdiction) (RRSPParams, error) {

	jurisdictionParams, ok := rrspParamsAll[region]
	if !ok {
		return RRSPParams{}, ErrJurisdictionNotExist
	}

	params, ok := jurisdictionParams[year]
	if !ok {
		return RRSPParams{}, ErrParamsNotExist
	}

	return params.Clone(), nil
}
