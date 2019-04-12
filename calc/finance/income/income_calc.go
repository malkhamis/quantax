package income

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
)

// compile-time check for interface implementation
var _ = (*calc.IncomeCalculator)(nil)

// Calculator is used to calculate net income as per the underlying recipe
type Calculator struct {
	incomeAdjusters map[finance.IncomeSource]Adjuster
	deducAdjusters  map[finance.DeductionSource]Adjuster
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
func (c *Calculator) NetIncome(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	netIncome := c.TotalIncome(finances) - c.TotalDeductions(finances)
	return netIncome
}

// TotalIncome returns the total income of given finances, applying any needed
// adjustments as per the underlying recipe without subtracting deductions.
// If the given finances is nil, it returns 0.0
func (c *Calculator) TotalIncome(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	var totalIncome float64
	for source := range finances.IncomeSources() {

		incomeFromSrc := finances.TotalIncome(source)

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
func (c *Calculator) TotalDeductions(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	var totalDeductions float64
	for source := range finances.DeductionSources() {

		deducFromSrc := finances.TotalDeductions(source)

		adjuster, isAdjustable := c.deducAdjusters[source]
		if isAdjustable {
			totalDeductions += adjuster.Adjusted(deducFromSrc)
			continue
		}

		totalDeductions += deducFromSrc
	}

	return totalDeductions
}

// initialize is used to initialize this calculator from the given recipe
func (c *Calculator) initialize(recipe *Recipe) {

	incomeAdjusters := make(map[finance.IncomeSource]Adjuster)
	for source, adjuster := range recipe.IncomeAdjusters {
		incomeAdjusters[source] = adjuster.Clone()
	}

	deducAdjusters := make(map[finance.DeductionSource]Adjuster)
	for source, adjuster := range recipe.DeductionAdjusters {
		deducAdjusters[source] = adjuster.Clone()
	}

	c.incomeAdjusters = incomeAdjusters
	c.deducAdjusters = deducAdjusters
}
