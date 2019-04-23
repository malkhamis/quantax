package income

import (
	"github.com/malkhamis/quantax/core"
)

// compile-time check for interface implementation
var _ core.IncomeCalculator = (*Calculator)(nil)

// Calculator is used to calculate net income as per the underlying recipe
type Calculator struct {
	incomeAdjusters map[core.FinancialSource]Adjuster
	deducAdjusters  map[core.FinancialSource]Adjuster
	finances        core.Financer
}

// NewCalculator returns a new income calculator for the given recipe
func NewCalculator(recipe *Recipe) (*Calculator, error) {

	if recipe == nil {
		return nil, ErrNoRecipe
	}

	c := new(Calculator)
	c.initialize(recipe)
	return c, nil
}

// NetIncome returns the net income of the given finances as follows:
//  NetIncome = (Total Adjusted Income) - (Total Adjusted Deductions)
// If the given finances is nil, it returns 0.0
func (c *Calculator) NetIncome() float64 {

	if c.finances == nil {
		return 0.0
	}

	netIncome := c.TotalIncome() - c.TotalDeductions()
	return netIncome
}

// TotalIncome returns the total income of given finances, applying any needed
// adjustments as per the underlying recipe without subtracting deductions.
// If the given finances is nil, it returns 0.0
func (c *Calculator) TotalIncome() float64 {

	if c.finances == nil {
		return 0.0
	}

	var totalIncome float64
	for source := range c.finances.IncomeSources() {

		incomeFromSrc := c.finances.TotalIncome(source)

		adjuster, isAdjustable := c.incomeAdjusters[source]
		if isAdjustable {
			totalIncome += adjuster.Adjusted(incomeFromSrc)
			continue
		}

		totalIncome += incomeFromSrc
	}

	return totalIncome
}

// TotalDeductions returns the total deductions of given finances, applying any
// needed adjustments as per the underlying recipe without adding income.
// If the given finances is nil, it returns 0.0
func (c *Calculator) TotalDeductions() float64 {

	if c.finances == nil {
		return 0.0
	}

	var totalDeductions float64
	for source := range c.finances.DeductionSources() {

		deducFromSrc := c.finances.TotalDeductions(source)

		adjuster, isAdjustable := c.deducAdjusters[source]
		if isAdjustable {
			totalDeductions += adjuster.Adjusted(deducFromSrc)
			continue
		}

		totalDeductions += deducFromSrc
	}

	return totalDeductions
}

// SetFinances sets the given finances to this calculator, making subsequent
// calls based on them
func (c *Calculator) SetFinances(finances core.Financer) {
	c.finances = finances
}

// initialize is used to initialize this calculator from the given recipe
func (c *Calculator) initialize(recipe *Recipe) {

	incomeAdjusters := make(map[core.FinancialSource]Adjuster)
	for source, adjuster := range recipe.IncomeAdjusters {
		incomeAdjusters[source] = adjuster.Clone()
	}

	deducAdjusters := make(map[core.FinancialSource]Adjuster)
	for source, adjuster := range recipe.DeductionAdjusters {
		deducAdjusters[source] = adjuster.Clone()
	}

	c.incomeAdjusters = incomeAdjusters
	c.deducAdjusters = deducAdjusters
}
