// Package history provides historical tax params for various jurisdictions
package history

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
)

type (
	yearlyTaxFormulas = map[uint]calc.TaxFormula
	yearlyCBFormulas  = map[uint]benefits.ChildBenefitFormula
)

const MonthsInYear = 12

var taxFormulasAll = map[Jurisdiction]yearlyTaxFormulas{
	BC:     taxFormulasBC,
	Canada: taxFormulasCanada,
}

var cbFormulasAll = map[Jurisdiction]yearlyCBFormulas{
	BC:     cbFormulasBC,
	Canada: cbFormulasCanada,
}

// GetTaxFormula returns the tax formula for the given year and region
func GetTaxFormula(year uint, region Jurisdiction) (calc.TaxFormula, error) {

	jurisdictionFormulas, ok := taxFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionFormulas[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula, nil
}

// GetChildBenefitFormula returns the child benefit formula for the given year
// and region
func GetChildBenefitFormula(year uint, region Jurisdiction) (benefits.ChildBenefitFormula, error) {

	jurisdictionFormulas, ok := cbFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionFormulas[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula, nil
}
