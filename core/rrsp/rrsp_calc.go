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
// calculator for the given withdrawal amount. The returned tax credit are
// after the withdrawal
func (c *Calculator) TaxPaid(withdrawal float64) (float64, []core.TaxCredit) {

	incomeSrc := c.formula.TargetSourceForWithdrawl()
	diff, credits := c.taxDiffForTargetSpouse(incomeSrc, withdrawal)
	return -diff, credits
}

// TaxRefund calculates the refundable tax proportion of deposit/contribution
// given the finances set in this calculator. The returned tax credit are after
// the contribution amount
func (c *Calculator) TaxRefund(contribution float64) (float64, []core.TaxCredit) {

	deducSrc := c.formula.TargetSourceForContribution()
	diff, credits := c.taxDiffForTargetSpouse(deducSrc, contribution)
	return diff, credits
}

// ContributionEarned calculates the newly acquired contribution room
func (c *Calculator) ContributionEarned() float64 {

	targetSpouse := c.targetSpouse()
	if targetSpouse == nil {
		return 0
	}

	incSrcs := c.formula.AllowedIncomeSources()
	totalIncome := targetSpouse.TotalAmount(incSrcs...)
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

// taxDiffForTargetSpouse returns the tax difference after adding the given source amount to
// the target spouse set in this calculator. The returned amount is calculated
// as follows:
//   diff = (tax before adding the amount) - (tax after adding the amount)
//
// The returned credits are for after adding the given source amount when
// computing the tax
func (c *Calculator) taxDiffForTargetSpouse(src core.FinancialSource, amount float64) (float64, []core.TaxCredit) {

	householdFinancesClone, targetSpouseClone := c.cloneFinancesAndGetTargetRef()
	if targetSpouseClone == nil {
		return 0, nil
	}

	c.setupTaxCalculator(householdFinancesClone)

	taxBeforeA, taxBeforeB, _ := c.taxCalculator.TaxPayable()
	targetSpouseClone.AddAmount(src, amount)
	taxAfterA, taxAfterB, credits := c.taxCalculator.TaxPayable()

	var diff float64
	if c.isTargetSpouseB {
		diff = taxBeforeB - taxAfterB
	} else {
		diff = taxBeforeA - taxAfterA
	}

	return diff, credits
}

// targetSpouse returns a read-only reference to the target spouse. If the set
// target points to a nil spouse, it returns nil
func (c *Calculator) targetSpouse() core.Financer {
	if c.isTargetSpouseB {
		return c.householdFinances.SpouseB()
	}
	return c.householdFinances.SpouseA()
}

// cloneFinancesAndGetTargetRef clones the household finances set in calculator.
// In addition, it returns a mutable reference to the target spouse from the
// clone household finances
func (c *Calculator) cloneFinancesAndGetTargetRef() (core.HouseholdFinanceMutator, core.FinanceMutator) {

	var (
		finances = c.householdFinances.Clone()
		target   core.FinanceMutator
	)

	if c.isTargetSpouseB {
		target = finances.MutableSpouseB()
	} else {
		target = finances.MutableSpouseA()
	}

	return finances, target
}

// setupTaxCalculator set the given finances as well as dependents and tax
// credits stored in the RRSP calculator into the tax calculator
func (c *Calculator) setupTaxCalculator(finances core.HouseholdFinances) {
	c.taxCalculator.SetDependents(c.dependents)
	c.taxCalculator.SetFinances(finances, c.taxCredits)
}
