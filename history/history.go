// Package history provides historical tax params for various jurisdictions
package history

import (
	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/rrsp"
	"github.com/malkhamis/quantax/calc/tax"
)

type (
	yearlyTaxFormulas  = map[uint]tax.Formula
	yearlyCBFormulas   = map[uint]benefits.ChildBenefitFormula
	yearlyRRSPFormulas = map[uint]rrsp.Formula
)

const MonthsInYear = 12

var (
	taxFormulasAll = map[Jurisdiction]yearlyTaxFormulas{
		BC:     taxFormulasBC,
		Canada: taxFormulasCanada,
	}

	cbFormulasAll = map[Jurisdiction]yearlyCBFormulas{
		BC:     cbFormulasBC,
		Canada: cbFormulasCanada,
	}

	rrspFormulasAll = map[Jurisdiction]yearlyRRSPFormulas{
		Canada: rrspFormulasCanada,
	}
)

// GetTaxFormula returns a copy of the tax formula for the given year and region
func GetTaxFormula(year uint, region Jurisdiction) (tax.Formula, error) {

	jurisdictionFormulas, ok := taxFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionFormulas[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula.Clone(), nil
}

// GetChildBenefitFormula returns a copy of the child benefit formula for the
// given year and region
func GetChildBenefitFormula(year uint, region Jurisdiction) (benefits.ChildBenefitFormula, error) {

	jurisdictionFormulas, ok := cbFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionFormulas[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula.Clone(), nil
}

// GetRRSOFormula returns a copy of the RRSP formula for the given year/region
func GetRRSPFormula(year uint, region Jurisdiction) (rrsp.Formula, error) {

	jurisdictionFormulas, ok := rrspFormulasAll[region]
	if !ok {
		return nil, ErrJurisdictionNotExist
	}

	formula, ok := jurisdictionFormulas[year]
	if !ok {
		return nil, ErrFormulaNotExist
	}

	return formula.Clone(), nil
}
