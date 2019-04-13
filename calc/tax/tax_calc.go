package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula          Formula
	contraFormula    ContraFormula
	incomeCalculator calc.IncomeCalculator
}

// NewCalculator returns a new tax calculator for the given tax formula and the
// income calculator
func NewCalculator(cfg CalcConfig) (*Calculator, error) {

	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := &Calculator{
		formula:          cfg.TaxFormula.Clone(),
		contraFormula:    cfg.ContraTaxFormula.Clone(),
		incomeCalculator: cfg.IncomeCalc,
	}

	return c, nil
}

// Calc computes the tax on the net income for the given finances
func (c *Calculator) Calc(finances *finance.IndividualFinances) float64 {

	if finances == nil {
		return 0.0
	}

	netIncome := c.incomeCalculator.NetIncome(finances)
	payableTax := c.formula.Apply(netIncome)
	// TODO: apply credits
	return payableTax
}
