package rrsp

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ calc.RRSPCalculator = (*Calculator)(nil)

// Calculator is a type used to calculate tax paid and refunded when making RRSP
// withdrawal or contribution. It also computes added contributution room for
// a given income
type Calculator struct {
	taxCalculator calc.TaxCalculator
	formula       Formula
	finances      *finance.IndividualFinances
}

// NewCalculator returns a new RRSP calculator from the given options with
// an empty finances instance
func NewCalculator(formula Formula, taxCalc calc.TaxCalculator) (*Calculator, error) {

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
		finances:      finance.NewEmptyIndividualFinances(0),
	}
	return c, nil
}

// TaxPaid calculates the extra tax payable given the finances set in this
// calculator for the given withdrawal amount
func (c *Calculator) TaxPaid(withdrawal float64) float64 {

	taxBefore := c.taxCalculator.Calc(c.finances)

	incomeSrc := c.formula.TargetSourceForWithdrawl()
	c.finances.AddIncome(incomeSrc, withdrawal)
	defer c.finances.AddIncome(incomeSrc, -withdrawal)

	taxAfter := c.taxCalculator.Calc(c.finances)

	diff := taxAfter - taxBefore
	return diff
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator
func (c *Calculator) TaxRefund(contribution float64) (float64, error) {

	if contribution > c.finances.RRSPContributionRoom {
		return 0.0, ErrNoRRSPRoom
	}

	taxBefore := c.taxCalculator.Calc(c.finances)

	deducSrc := c.formula.TargetSourceForContribution()
	c.finances.AddDeduction(deducSrc, contribution)
	defer c.finances.AddDeduction(deducSrc, -contribution)

	taxAfter := c.taxCalculator.Calc(c.finances)

	diff := taxBefore - taxAfter
	return diff, nil
}

// ContributionEarned calculates the newly acquired contribution room
func (c *Calculator) ContributionEarned() float64 {

	var netIncome float64
	incSrcs := c.formula.AllowedIncomeSources()

	if len(incSrcs) > 0 {
		netIncome = c.finances.TotalIncome(incSrcs...)
	}

	return c.formula.ContributionEarned(netIncome)
}

// SetFinances makes subsequent calculations based on the given finances.
// if newFinances is nil, an empty finances instance is set
func (c *Calculator) SetFinances(newFinances *finance.IndividualFinances) {
	c.finances = newFinances
	if c.finances == nil {
		c.finances = finance.NewEmptyIndividualFinances(0)
	}
}
