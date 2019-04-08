package income

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"

	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
)

// compile-time check for interface implementation
var _ = (*calc.IncomeCalculator)(nil)

// Calculator is used to calculate net income as per the underlying formula
type Calculator struct {
	incomeAdjusters map[finance.IncomeSource]Adjuster
	deducAdjusters  map[finance.DeductionSource]Adjuster
	formula         Formula
}

// NewCalculator returns a new income calculator for the given formula
func NewCalculator(formula Formula) (*Calculator, error) {

	if formula == nil {
		return nil, ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	c := new(Calculator)
	c.initialize(formula)
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
// adjustments as per the underlying formula without subtracting deductions.
// If the given finances is nil, it returns 0.0
func (c *Calculator) TotalIncome(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	var totalIncome float64
	for source := range finances.IncomeSources() {

		adjuster, isAdjustable := c.incomeAdjusters[source]
		if isAdjustable {
			totalIncome += adjuster.Adjusted(finances)
			continue
		}

		totalIncome += finances.TotalIncome(source)
	}

	return totalIncome
}

// TotalDeductions returns the total income of given finances, applying any
// needed adjustments as per the underlying formula without adding income.
// If the given finances is nil, it returns 0.0
func (c *Calculator) TotalDeductions(finances finance.IncomeDeductor) float64 {

	if finances == nil {
		return 0.0
	}

	var totalDeductions float64
	for source := range finances.DeductionSources() {

		adjuster, isAdjustable := c.deducAdjusters[source]
		if isAdjustable {
			totalDeductions += adjuster.Adjusted(finances)
			continue
		}

		totalDeductions += finances.TotalDeductions(source)
	}

	return totalDeductions
}

// initialize is used to initialize this calculator from the given formula
func (c *Calculator) initialize(formula Formula) {

	incomeAdjusters := make(map[finance.IncomeSource]Adjuster)
	for source, adjuster := range formula.IncomeAdjusters() {
		incomeAdjusters[source] = adjuster.Clone()
	}

	deducAdjusters := make(map[finance.DeductionSource]Adjuster)
	for source, adjuster := range formula.DeductionAdjusters() {
		deducAdjusters[source] = adjuster.Clone()
	}

	c.incomeAdjusters = incomeAdjusters
	c.deducAdjusters = deducAdjusters
	c.formula = formula.Clone()
}
