// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula Formula
}

// NewCalculator returns a new tax calculator for the given financial numbers
// and tax formula
func NewCalculator(formula Formula) (*Calculator, error) {

	if formula == nil {
		return nil, ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	return &Calculator{formula: formula.Clone()}, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (c *Calculator) Calc(finances *finance.IndividualFinances) float64 {

	excludedIncSrcs := c.formula.ExcludedIncomeSources()
	excludedDeducSrcs := c.formula.ExcludedDeductionSources()

	adjustedTotalIncome := finances.TotalIncome()
	if len(excludedIncSrcs) > 0 {
		adjustedTotalIncome -= finances.TotalIncome(excludedIncSrcs...)
	}

	adjustedTotalDeductions := finances.TotalDeductions()
	if len(excludedDeducSrcs) > 0 {
		adjustedTotalDeductions -= finances.TotalDeductions(excludedDeducSrcs...)
	}

	adjustedNetIncome := adjustedTotalIncome - adjustedTotalDeductions
	payableTax := c.formula.Apply(adjustedNetIncome)
	return payableTax
}
