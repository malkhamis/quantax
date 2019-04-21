package history

import (
	"github.com/malkhamis/quantax/core/benefits"
	"github.com/malkhamis/quantax/core/income"
	"github.com/malkhamis/quantax/core/rrsp"
	"github.com/malkhamis/quantax/core/tax"
)

// TaxParams represents the tax parameters associated with a tax jurisdiction
// for a specific tax year
type TaxParams struct {
	Formula       tax.Formula
	ContraFormula tax.ContraFormula
	IncomeRecipe  *income.Recipe
}

// Clone returns a copy of these parameters
func (p TaxParams) Clone() TaxParams {
	return TaxParams{
		Formula:       p.Formula.Clone(),
		ContraFormula: p.ContraFormula.Clone(),
		IncomeRecipe:  p.IncomeRecipe.Clone(),
	}
}

// RRSPParams represents the RRSP parameters associated with a jurisdiction
// for a specific tax year
type RRSPParams struct {
	Formula rrsp.Formula
}

// Clone returns a copy of these parameters
func (p RRSPParams) Clone() RRSPParams {
	return RRSPParams{
		Formula: p.Formula.Clone(),
	}
}

// CBParams represents the child benefit parameters associated with a
// jurisdiction for a specific tax year
type CBParams struct {
	Formula      benefits.ChildBenefitFormula
	IncomeRecipe *income.Recipe
}

func (p CBParams) Clone() CBParams {
	return CBParams{
		Formula:      p.Formula.Clone(),
		IncomeRecipe: p.IncomeRecipe.Clone(),
	}
}

type (
	yearlyTaxParams  = map[uint]TaxParams
	yearlyCBParams   = map[uint]CBParams
	yearlyRRSPParams = map[uint]RRSPParams
)

const monthsInYear = 12
