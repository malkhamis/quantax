package rrsp

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.RRSPCalculator = (*Calculator)(nil)

// Calculator is a type used to calculate tax paid and refunded when making RRSP
// withdrawal or contribution. It also computes added contributution room for
// a given income
type Calculator struct {
	taxCalculator core.TaxCalculator
	formula       Formula
	finances      *finance.IndividualFinances
}

// NewCalculator returns a new RRSP calculator from the given options with
// an empty finances instance
func NewCalculator(formula Formula, taxCalc core.TaxCalculator) (*Calculator, error) {

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

	taxBefore, _ := c.taxCalculator.TaxPayable()

	incomeSrc := c.formula.TargetSourceForWithdrawl()
	c.finances.AddIncome(incomeSrc, withdrawal)
	defer c.finances.AddIncome(incomeSrc, -withdrawal)

	taxAfter, _ := c.taxCalculator.TaxPayable()

	diff := taxAfter - taxBefore
	return diff
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator
func (c *Calculator) TaxRefund(contribution float64) (float64, error) {

	if contribution > c.finances.RRSPContributionRoom {
		return 0.0, ErrNoRRSPRoom
	}

	taxBefore, _ := c.taxCalculator.TaxPayable()

	deducSrc := c.formula.TargetSourceForContribution()
	c.finances.AddDeduction(deducSrc, contribution)
	defer c.finances.AddDeduction(deducSrc, -contribution)

	taxAfter, _ := c.taxCalculator.TaxPayable()

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
// if new finances is nil, an empty finances instance is set. Change to the
// given finances will affect the results of future calls on this calculator
func (c *Calculator) SetFinances(f *finance.IndividualFinances) {

	if f == nil {
		f = finance.NewEmptyIndividualFinances(0)
	}

	c.finances = f
	c.taxCalculator.SetFinances(f)
}
