package tax

import (
	"math"

	"github.com/malkhamis/quantax/calc/finance"
)

// areEqual returns true if the difference between floor(actual) and
// floor(expected) is within the given +/- error margin of expcted. Negative
// error margins are converted to a positive number
func areEqual(actual, expected, errMargin float64) bool {

	actual, expected = math.Floor(actual), math.Floor(expected)
	allowedDiff := math.Abs(errMargin * expected)
	actualDiff := math.Abs(actual - expected)

	return actualDiff <= allowedDiff
}

type testIncomeCalculator struct {
	onNetIncome       float64
	onTotalDeductions float64
	onTotalIncome     float64
}

func (tic testIncomeCalculator) TotalIncome(_ finance.IncomeDeductor) float64 {
	return tic.onTotalIncome
}
func (tic testIncomeCalculator) TotalDeductions(_ finance.IncomeDeductor) float64 {
	return tic.onTotalDeductions
}
func (tic testIncomeCalculator) NetIncome(_ finance.IncomeDeductor) float64 {
	return tic.onNetIncome
}

type testTaxFormula struct {
	onApply    float64
	onValidate error
	onClone    Formula
}

func (tcb testTaxFormula) Apply(_ float64) float64 {
	return tcb.onApply
}
func (tcb testTaxFormula) Validate() error {
	return tcb.onValidate
}
func (tcb testTaxFormula) Clone() Formula {
	return tcb.onClone
}
