package rrsp

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.RRSPCalculator = (*Calculator)(nil)

// Calculator is a type used to calculate tax paid and refunded when making RRSP
// withdrawal or contribution. It also computes added contributution room for
// a given income
type Calculator struct {
	formula           Formula
	householdFinances core.HouseholdFinances
	isTargetSpouseB   bool // default to SpouseA
	taxCredits        []core.TaxCredit
	dependents        []*human.Person
	taxCalculator     core.TaxCalculator
}

// NewCalculator returns a new RRSP calculator from the given options with
// an empty finances instance
func NewCalculator(cfg CalcConfig) (*Calculator, error) {

	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := &Calculator{
		formula:           cfg.Formula.Clone(),
		taxCalculator:     cfg.TaxCalc,
		householdFinances: core.NewHouseholdFinancesNop(),
	}

	return c, nil
}

// TaxPaid calculates the extra tax payable given the finances set in this
// calculator for the given withdrawal amount
func (c *Calculator) TaxPaid(withdrawal float64) (float64, []core.TaxCredit) {

	var (
		finances = c.householdFinances.Clone()
		target   core.FinanceMutator
	)

	if c.isTargetSpouseB {
		target = finances.MutableSpouseB()
	} else {
		target = finances.MutableSpouseA()
	}
	if target == nil {
		return 0, nil
	}

	incomeSrc := c.formula.TargetSourceForWithdrawl()
	c.taxCalculator.SetDependents(c.dependents)
	c.taxCalculator.SetFinances(finances, c.taxCredits)

	taxBeforeA, taxBeforeB, _ := c.taxCalculator.TaxPayable()
	target.AddAmount(incomeSrc, withdrawal)
	taxAfterA, taxAfterB, credits := c.taxCalculator.TaxPayable()

	var diff float64
	if c.isTargetSpouseB {
		diff = taxAfterB - taxBeforeB
	} else {
		diff = taxAfterA - taxBeforeA
	}

	return diff, credits
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator
func (c *Calculator) TaxRefund(contribution float64) (float64, []core.TaxCredit) {

	var (
		finances = c.householdFinances.Clone()
		target   core.FinanceMutator
	)

	if c.isTargetSpouseB {
		target = finances.MutableSpouseB()
	} else {
		target = finances.MutableSpouseA()
	}
	if target == nil {
		return 0, nil
	}

	deducSrc := c.formula.TargetSourceForContribution()
	c.taxCalculator.SetDependents(c.dependents)
	c.taxCalculator.SetFinances(finances, c.taxCredits)

	taxBeforeA, taxBeforeB, _ := c.taxCalculator.TaxPayable()
	target.AddAmount(deducSrc, contribution)
	taxAfterA, taxAfterB, credits := c.taxCalculator.TaxPayable()

	var diff float64
	if c.isTargetSpouseB {
		diff = taxBeforeB - taxAfterB
	} else {
		diff = taxBeforeA - taxAfterA
	}

	return diff, credits
}

// ContributionEarned calculates the newly acquired contribution room
func (c *Calculator) ContributionEarned() float64 {

	var target core.Financer
	if c.isTargetSpouseB {
		target = c.householdFinances.SpouseA()
	} else {
		target = c.householdFinances.SpouseB()
	}
	if target == nil {
		return 0
	}

	incSrcs := c.formula.AllowedIncomeSources()
	totalIncome := target.TotalAmount(incSrcs...)
	return c.formula.ContributionEarned(totalIncome)
}

// SetDependents sets the dependents which the calculator might use for tax-
// related calculations
func (c *Calculator) SetDependents(dependents []*human.Person) {
	c.dependents = dependents
}

// SetFinances makes subsequent calculations based on the given finances.
// if new finances is nil, an empty finances instance is set. Change to the
// given finances will affect the results of future calls on this calculator
func (c *Calculator) SetFinances(f core.HouseholdFinances, credits []core.TaxCredit) {

	if f == nil {
		f = core.NewHouseholdFinancesNop()
	}
	c.householdFinances = f
	c.taxCredits = credits
}

// SetTargetSpouseA makes subsequent calculations based on SpouseA of the
// previously set finances. This is the default target of the calculator
func (c *Calculator) SetTargetSpouseA() {
	c.isTargetSpouseB = false
}

// SetTargetSpouseB makes subsequent calculations based on SpouseA of the
// previously set finances
func (c *Calculator) SetTargetSpouseB() {
	c.isTargetSpouseB = true
}
