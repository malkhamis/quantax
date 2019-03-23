package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
)

// IncomeTaxParams is used to pass tax parameters to relevant factory functions
type IncomeTaxParams struct {
	Year   uint
	Region history.Jurisdiction
}

// NewIncomeTaxCalculator returns a tax calculator for the given parameters
func NewIncomeTaxCalculator(finNums calc.IndividualFinances, params IncomeTaxParams) (calc.TaxCalculator, error) {

	rates, err := history.GetFormula(params.Year, params.Region)
	if err != nil {
		return nil, err
	}

	return tax.NewCalculator(finNums, rates)
}
