// Package incometax provides implementations for calc.TaxCalculator
package incometax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var (
	_ calc.TaxCalculator = (*Calculator)(nil)
)

// Calculator is used to calculate payable tax
type Calculator struct {
	taxRates calc.BracketRates
	calc.FinancialNumbers
}

// NewCalculator returns a new calculator for the given financial numbers
// and tax brackets.
func NewCalculator(finNums calc.FinancialNumbers, rates calc.BracketRates) (*Calculator, error) {

	c := &Calculator{
		FinancialNumbers: finNums,
		taxRates:         rates.Clone(),
	}

	return c, c.validate()
}

// validate ensures that this calculator is valid for use.
func (c *Calculator) validate() error {

	err := c.taxRates.Validate()
	return errors.Wrap(err, "invalid tax parameters")
}

// Calc computes the tax on the taxable amount set in this calculator
func (c *Calculator) Calc() float64 {

	var (
		payableTax     float64
		taxableAmntTot = c.TaxableAmount - c.Deductions
	)

	for rate, bracket := range c.taxRates {

		if taxableAmntTot < bracket.Lower() {
			continue
		}

		if taxableAmntTot >= bracket.Upper() {
			payableTax += rate * bracket.Amount()
			continue
		}

		payableTax += rate * (taxableAmntTot - bracket.Lower())
	}

	return payableTax - c.Credits
}

// Update sets the financial numbers which the tax will be calculated for
func (c *Calculator) Update(newFinNums calc.FinancialNumbers) {
	c.FinancialNumbers = newFinNums
}
