package tax

import (
	"math"

	"github.com/malkhamis/quantax/core"
)

var (
	_ Formula                = (*testTaxFormula)(nil)
	_ ContraFormula          = (*testContraTaxFormula)(nil)
	_ Creditor               = (*testCreditor)(nil)
	_ core.IncomeCalculator  = (*testIncomeCalculator)(nil)
	_ core.HouseholdFinances = (*testHouseholdFinances)(nil)
	_ core.TaxCredit         = (*testTaxCredit)(nil)
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

func (tic *testIncomeCalculator) TotalIncome() float64 {
	return tic.onTotalIncome
}
func (tic *testIncomeCalculator) TotalDeductions() float64 {
	return tic.onTotalDeductions
}
func (tic *testIncomeCalculator) NetIncome() float64 {
	return tic.onNetIncome
}
func (tic *testIncomeCalculator) SetFinances(core.Financer) {

}

type testTaxFormula struct {
	onApply    float64
	onValidate error
	onYear     uint
	onRegion   core.Region
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
func (tcb *testTaxFormula) Year() uint {
	return tcb.onYear
}
func (tcb *testTaxFormula) Region() core.Region {
	return tcb.onRegion
}

type testContraTaxFormula struct {
	onApply         []*TaxCredit
	onFilterAndSort func()
	onYear          uint
	onRegion        core.Region
	onValidate      error
}

func (tcf *testContraTaxFormula) Apply(_ *TaxPayer) []*TaxCredit {

	clone := make([]*TaxCredit, len(tcf.onApply))
	for i, cr := range tcf.onApply {
		cloneCr := *cr
		clone[i] = &cloneCr
	}
	return clone
}
func (tcf *testContraTaxFormula) FilterAndSort(_ *[]core.TaxCredit) {
	if tcf.onFilterAndSort != nil {
		tcf.onFilterAndSort()
	}
}
func (tcf *testContraTaxFormula) Validate() error {
	return tcf.onValidate
}
func (tcf *testContraTaxFormula) Region() core.Region {
	return tcf.onRegion
}
func (tcf *testContraTaxFormula) Year() uint {
	return tcf.onYear
}
func (tcf *testContraTaxFormula) Clone() ContraFormula {
	return tcf
}

type testCreditor struct {
	onTaxCredit       float64
	onRule            core.CreditRule
	onFinancialSource core.FinancialSource
}

func (tc *testCreditor) TaxCredit(_ *TaxPayer) float64 {
	return tc.onTaxCredit
}
func (tc *testCreditor) Rule() core.CreditRule {
	return tc.onRule
}
func (tc *testCreditor) FinancialSource() core.FinancialSource {
	return tc.onFinancialSource
}
func (tc *testCreditor) Description() string {
	return "test"
}
func (tc *testCreditor) Clone() Creditor {
	return tc
}

type testTaxCredit struct {
	onAmounts           [3]float64
	onRule              core.CreditRule
	onReferenceFinancer core.Financer
	onSource            core.FinancialSource
	onYear              uint
	onRegion            core.Region
}

func (ttc *testTaxCredit) SetAmounts(_, _, _ float64) {}
func (ttc *testTaxCredit) Amounts() (initial, used, remaining float64) {
	return ttc.onAmounts[0], ttc.onAmounts[1], ttc.onAmounts[2]
}
func (ttc *testTaxCredit) Rule() core.CreditRule {
	return ttc.onRule
}
func (ttc *testTaxCredit) ReferenceFinancer() core.Financer {
	return ttc.onReferenceFinancer
}
func (ttc *testTaxCredit) Region() core.Region {
	return ttc.onRegion
}
func (ttc *testTaxCredit) Year() uint {
	return ttc.onYear
}
func (ttc *testTaxCredit) Description() string {
	return "test"
}
func (ttc *testTaxCredit) Source() core.FinancialSource {
	return ttc.onSource
}
func (ttc *testTaxCredit) ShallowCopy() core.TaxCredit {
	clone := *ttc
	return &clone
}

type testHouseholdFinances struct {
	onSpouseA core.FinanceMutator
	onSpouseB core.FinanceMutator
	onVersion uint64
}

func (thf *testHouseholdFinances) SpouseA() core.Financer {
	return thf.onSpouseA
}
func (thf *testHouseholdFinances) SpouseB() core.Financer {
	return thf.onSpouseB
}
func (thf *testHouseholdFinances) MutableSpouseA() core.FinanceMutator {
	return thf.onSpouseA
}
func (thf *testHouseholdFinances) MutableSpouseB() core.FinanceMutator {
	return thf.onSpouseB
}
func (thf *testHouseholdFinances) Clone() core.HouseholdFinanceMutator {
	return thf
}
func (thf *testHouseholdFinances) Version() uint64 {
	return thf.onVersion
}
