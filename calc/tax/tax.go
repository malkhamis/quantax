// Package tax provides implementations for calc.TaxCalculator
package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var (
	_ calc.TaxCalculator = (*Calculator)(nil)
)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula  Formula
	finances calc.IndividualFinances
}

// NewCalculator returns a new tax calculator for the given financial numbers
// and tax formula.
func NewCalculator(finances calc.IndividualFinances, formula Formula) (*Calculator, error) {

	if formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	c := &Calculator{formula: formula}
	c.UpdateFinances(finances)
	return c, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (c *Calculator) Calc() float64 {

	netIncome := c.finances.Income - c.finances.Deductions
	payableTax := c.formula.Apply(netIncome)
	return payableTax
}

// Update sets the financial numbers which the tax will be calculated for
func (c *Calculator) UpdateFinances(newFinances calc.IndividualFinances) {
	c.finances = newFinances
}
