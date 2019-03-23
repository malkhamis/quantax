// Package history provides historical tax params for various jurisdictions
package history

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type TaxFormula = calc.TaxFormula
type yearlyTaxFormulas = map[uint]TaxFormula

type WeightedBracketFormula = calc.WeightedBracketFormula
type Bracket = calc.Bracket

var taxFormulasAll = map[Jurisdiction]yearlyTaxFormulas{
	BC:     taxFormulasBC,
	Canada: taxFormulasCanada,
}

func init() {

	err := validateAllFormulas()
	if err != nil {
		panic(err)
	}

}

// GetFormula returns the tax formula for the given region in a specific year
func GetFormula(year uint, region Jurisdiction) (TaxFormula, error) {

	jurisdictionRates, ok := taxFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionRates[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula, nil
}

func validateAllFormulas() error {

	for jursdiction, formulasAllYears := range taxFormulasAll {
		for year, formula := range formulasAllYears {

			err := formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}
