// Package tax provides implementations for calc.TaxCalculator
package tax

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
var _ calc.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula Formula
}

// NewCalculator returns a new tax calculator for the given financial numbers
// and tax formula
func NewCalculator(formula Formula) (*Calculator, error) {

	if formula == nil {
		return nil, ErrNoFormula
	}

	err := formula.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula")
	}

	return &Calculator{formula: formula.Clone()}, nil
}

// Calc computes the tax on the taxable amount set in this calculator
func (c *Calculator) Calc(finances *finance.IndividualFinances) float64 {

	totalIncome, creditsInc := c.calcTotalIncome(finances)
	totalDeductions, creditsDeduc := c.calcTotalDeductions(finances)

	adjustedNetIncome := totalIncome - totalDeductions
	payableTax := c.formula.Apply(adjustedNetIncome)

	totalCredits := TaxCredits{
		Refundable:    creditsInc.Refundable + creditsDeduc.Refundable,
		NonRefundable: creditsInc.NonRefundable + creditsDeduc.NonRefundable,
	}
	payableTax -= totalCredits.NonRefundable
	if payableTax < 0 {
		return -totalCredits.Refundable
	}

	payableTax -= totalCredits.Refundable
	return payableTax
}

// calcTotalIncome returns the total, unadjusted income for the given finances
// according to the underlying formula. It also returns the adjustment amount
// that should be added to the total income and the final tax amount
func (c *Calculator) calcTotalIncome(finances *finance.IndividualFinances) (float64, TaxCredits) {

	adjustable := c.formula.AdjustableIncomeSources()
	if adjustable == nil {
		// to avoid panics and nil-checks later
		adjustable = make(map[finance.IncomeSource]Adjuster)
	}

	excluded := make(map[finance.IncomeSource]struct{})
	for _, source := range c.formula.ExcludedIncomeSources() {
		// to avoid slice searches
		excluded[source] = struct{}{}
	}

	var (
		totalIncome  float64
		totalCredits TaxCredits
	)

	for source, income := range finances.Income {

		_, isExcluded := excluded[source]
		if isExcluded {
			continue
		}

		adjuster, isAdjustable := adjustable[source]
		if isAdjustable {
			adj := adjuster.Adjusted(finances)
			totalIncome += adj.AdjustedAmount
			totalCredits.Refundable += adj.Credits.Refundable
			totalCredits.NonRefundable += adj.Credits.NonRefundable
			continue
		}

		totalIncome += income
	}

	return totalIncome, totalCredits
}

// calcTotalDeductions returns the total, unadjusted deduction for the given
// finances according to the underlying formula. It also returns the adjustment
// amount that should beadded to the total income and the final tax amount
func (c *Calculator) calcTotalDeductions(finances *finance.IndividualFinances) (float64, TaxCredits) {

	adjustable := c.formula.AdjustableDeductionSources()
	if adjustable == nil {
		// to avoid panics and nil-checks later
		adjustable = make(map[finance.DeductionSource]Adjuster)
	}

	excluded := make(map[finance.DeductionSource]struct{})
	for _, source := range c.formula.ExcludedDeductionSources() {
		// to avoid slice searches
		excluded[source] = struct{}{}
	}

	var (
		totalDeductions float64
		totalCredits    TaxCredits
	)

	for source, deduction := range finances.Deductions {

		_, isExcluded := excluded[source]
		if isExcluded {
			continue
		}

		adjuster, isAdjustable := adjustable[source]
		if isAdjustable {
			adj := adjuster.Adjusted(finances)
			totalDeductions += adj.AdjustedAmount
			totalCredits.Refundable += adj.Credits.Refundable
			totalCredits.NonRefundable += adj.Credits.NonRefundable
			continue
		}

		totalDeductions += deduction
	}

	return totalDeductions, totalCredits
}
