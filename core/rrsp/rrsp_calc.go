package rrsp

import (
	"github.com/malkhamis/quantax/core"
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
	finances      *core.IndividualFinances
}

// NewCalculator returns a new RRSP calculator from the given options with
// an empty finances instance
func NewCalculator(cfg CalcConfig) (*Calculator, error) {

	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := &Calculator{
		formula:       cfg.Formula.Clone(),
		taxCalculator: cfg.TaxCalc,
		finances:      core.NewEmptyIndividualFinances(),
	}
	return c, nil
}

// TaxPaid calculates the extra tax payable given the finances set in this
// calculator for the given withdrawal amount
func (c *Calculator) TaxPaid(withdrawal float64) float64 {

	original := c.finances
	clone := original.Clone()
	c.SetFinances(clone)
	defer c.SetFinances(original)

	incomeSrc := c.formula.TargetSourceForWithdrawl()

	taxBefore, _ := c.taxCalculator.TaxPayable()
	c.finances.AddAmount(incomeSrc, withdrawal)
	taxAfter, _ := c.taxCalculator.TaxPayable()

	diff := taxAfter - taxBefore
	return diff
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator
func (c *Calculator) TaxRefund(contribution float64) (float64, error) {

	if contribution > c.finances.RRSPAmounts().ContributionRoom {
		return 0.0, ErrNoRRSPRoom
	}

	original := c.finances
	clone := original.Clone()
	c.SetFinances(clone)
	defer c.SetFinances(original)

	deducSrc := c.formula.TargetSourceForContribution()

	taxBefore, _ := c.taxCalculator.TaxPayable()
	c.finances.AddAmount(deducSrc, contribution)
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
func (c *Calculator) SetFinances(f *core.IndividualFinances) {

	if f == nil {
		f = core.NewEmptyIndividualFinances()
	}

	c.finances = f
	c.taxCalculator.SetFinances(f)
}
