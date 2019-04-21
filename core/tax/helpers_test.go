package tax

import (
	"math"

	"github.com/malkhamis/quantax/core/finance"
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
}

func (tcb testTaxFormula) Apply(_ float64) float64 {
	return tcb.onApply
}
func (tcb testTaxFormula) Validate() error {
	return tcb.onValidate
}
func (tcb testTaxFormula) Clone() Formula {
	return tcb
}

type testTaxContraFormula struct {
	onApply    []*taxCredit
	onValidate error
}

func (tcf *testTaxContraFormula) Apply(_ *finance.IndividualFinances, _ float64) []*taxCredit {
	return tcf.onApply
}
func (tcf *testTaxContraFormula) Validate() error {
	return tcf.onValidate
}
func (tcf *testTaxContraFormula) Clone() ContraFormula {
	return tcf
}

type testCreditor struct {
	onTaxCredits float64
	onSource     string
}

func (tc testCreditor) TaxCredit(_, _ float64) float64 {
	return tc.onTaxCredits
}
func (tc testCreditor) Source() string {
	return tc.onSource
}
func (tc testCreditor) Clone() Creditor {
	return tc
}

type testTaxCredit struct {
	onAmount float64
	onSource string
}

func (ttc testTaxCredit) Amount() float64 {
	return ttc.onAmount
}

func (ttc testTaxCredit) Source() string {
	return ttc.onSource
}