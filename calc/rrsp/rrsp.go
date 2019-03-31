// Package rrsp provides implementations for the RRSPCalculator interface
// defined in package calc
package rrsp

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula  = errors.New("not formula given/set")
	ErrNoRRSPRoom = errors.New("no enough RRSP contribution room")
	ErrNoTaxCalc  = errors.New("no tax calculator given")
)

// compile-time check for interface implementation
var _ calc.RRSPCalculator = (*Calculator)(nil)

// Calculator is a type used to calculate tax paid and refunded when making RRSP
// withdrawal or contribution.
type Calculator struct {
	taxCalculator calc.TaxCalculator
	formula       Formula
	finances      finance.IndividualFinances
}

// NewCalculator returns a new RRSP calculator from the given formula and tax
// tax calculator
func NewCalculator(taxCalc calc.TaxCalculator, formula Formula) (*Calculator, error) {

	if formula == nil {
		return nil, ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	if taxCalc == nil {
		return nil, ErrNoTaxCalc
	}
	c := &Calculator{
		formula:       formula.Clone(),
		taxCalculator: taxCalc,
	}
	return c, nil
}

// TaxPaid calculates the extra tax payable given the finances set in this
// calculator for the given withdrawal amount
func (c *Calculator) TaxPaid(withdrawal float64) float64 {

	oldFinances := c.finances
	taxOnOldFinances := c.taxCalculator.Calc(oldFinances)

	newFinances := c.finances
	newFinances.Income += withdrawal
	taxOnNewFinances := c.taxCalculator.Calc(newFinances)

	diff := taxOnNewFinances - taxOnOldFinances
	return diff
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator
func (c *Calculator) TaxRefund(contribution float64) (float64, error) {

	if contribution > c.finances.RRSPRoom {
		return 0.0, ErrNoRRSPRoom
	}

	oldFinances := c.finances
	taxOnOldFinances := c.taxCalculator.Calc(oldFinances)

	newFinances := c.finances
	newFinances.Income -= contribution
	taxOnNewFinances := c.taxCalculator.Calc(newFinances)

	diff := taxOnOldFinances - taxOnNewFinances
	return diff, nil
}

// ContributionEarned calculates the newly acquired contribution room
func (c *Calculator) ContributionEarned() float64 {

	income := c.formula.IncomeCalcMethod().CalcForIndividials(c.finances)
	return c.formula.Contribution(income)
}

// SetFinances makes subsequent calculations based on the given finances
func (c *Calculator) SetFinances(newFinances finance.IndividualFinances) {
	c.finances = newFinances
}
