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

type IndividualFinances = calc.IndividualFinances
type TaxFormula = calc.TaxFormula

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula TaxFormula
	IndividualFinances
}

// NewCalculator returns a new calculator for the given financial numbers
// and tax brackets.
func NewCalculator(finNums IndividualFinances, formula TaxFormula) (*Calculator, error) {

	if formula == nil {
		return nil, calc.ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	c := &Calculator{
		formula:            formula,
		IndividualFinances: finNums,
	}

	return c, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (c *Calculator) Calc(taxCredits ...float64) float64 {

	netIncome := c.Income - c.Deductions
	payableTax := c.formula.Apply(netIncome)

	for _, credit := range taxCredits {
		payableTax -= credit
	}

	return payableTax
}

// Update sets the financial numbers which the tax will be calculated for
func (c *Calculator) UpdateFinances(newFinNums IndividualFinances) {
	c.IndividualFinances = newFinNums
}

// UpdateFormula sets this calculator up with the given formula. If the new
// formula is nil, the formula is not changed
func (c *Calculator) UpdateFormula(newFormula TaxFormula) {
	if newFormula != nil {
		c.formula = newFormula
	}
}
