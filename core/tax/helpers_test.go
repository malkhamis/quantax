package tax

import (
	"math"

	"github.com/malkhamis/quantax/core"
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

func (tic testIncomeCalculator) TotalIncome() float64 {
	return tic.onTotalIncome
}
func (tic testIncomeCalculator) TotalDeductions() float64 {
	return tic.onTotalDeductions
}
func (tic testIncomeCalculator) NetIncome() float64 {
	return tic.onNetIncome
}
func (tic testIncomeCalculator) SetFinances(core.Financer) {

}

type testTaxFormula struct {
	onApply    float64
	onValidate error
	onTaxInfo  core.TaxInfo
}

func (tcb *testTaxFormula) Apply(_ float64) float64 {
	return tcb.onApply
}
func (tcb *testTaxFormula) Validate() error {
	return tcb.onValidate
}
func (tcb *testTaxFormula) Clone() Formula {
	return tcb
}
func (tcb *testTaxFormula) TaxInfo() core.TaxInfo {
	return tcb.onTaxInfo
}

type testTaxContraFormula struct {
	onApply         []*TaxCredit
	onFilterAndSort []core.TaxCredit
	onTaxInfo       core.TaxInfo
	onValidate      error
}

func (tcf *testTaxContraFormula) Apply(_ *TaxPayer) []*TaxCredit {
	return tcf.onApply
}
func (tcf *testTaxContraFormula) FilterAndSort(_ []core.TaxCredit) []core.TaxCredit {
	return tcf.onFilterAndSort
}
func (tcf *testTaxContraFormula) Validate() error {
	return tcf.onValidate
}
func (tcf *testTaxContraFormula) TaxInfo() core.TaxInfo {
	return tcf.onTaxInfo
}
func (tcf *testTaxContraFormula) Clone() ContraFormula {
	return tcf
}

type testCreditor struct {
	onTaxCredit       float64
	onCrSourceName    string
	onFinancialSource core.FinancialSource
}

func (tc testCreditor) TaxCredit(_ *TaxPayer) float64 {
	return tc.onTaxCredit
}
func (tc testCreditor) CrSourceName() string {
	return tc.onCrSourceName
}
func (tc testCreditor) FinancialSource() core.FinancialSource {
	return tc.onFinancialSource
}
func (tc testCreditor) Description() string {
	return "test"
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
