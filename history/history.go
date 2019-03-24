// Package history provides historical tax params for various jurisdictions
package history

import (
	"github.com/malkhamis/quantax/calc"
)

type (
	yearlyTaxFormulas = map[uint]calc.TaxFormula
	yearlyCBFormulas  = map[uint]calc.ChildBenefitFormula
)

const MonthsInYear = 12

var taxFormulasAll = map[Jurisdiction]yearlyTaxFormulas{
	BC:     taxFormulasBC,
	Canada: taxFormulasCanada,
}

var cbFormulasAll = map[Jurisdiction]yearlyCBFormulas{
	Canada: cbFormulasCanada,
}

// GetFormula returns the tax formula for the given region in a specific year
func GetTaxFormula(year uint, region Jurisdiction) (calc.TaxFormula, error) {

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
