package rrsp

import (
	"github.com/malkhamis/quantax/core"
)

type testTaxCalculator struct {
	currentIndex int
	onTaxPayable []float64
}

func (ttc *testTaxCalculator) TaxPayable() (float64, []core.TaxCredit) {
	currentVal := ttc.onTaxPayable[ttc.currentIndex]
	ttc.currentIndex++
	return currentVal, nil
}
func (ttc *testTaxCalculator) SetCredits(_ []core.TaxCredit) {
}
func (ttc *testTaxCalculator) SetFinances(_ *core.IndividualFinances) {
}

type testFormula struct {
	onValidate                    error
	onContributionEarned          float64
	onTargetSourceForWithdrawl    core.FinancialSource
	onTargetSourceForContribution core.FinancialSource
	onAllowedIncomeSources        []core.FinancialSource
}

func (f *testFormula) ContributionEarned(income float64) float64 {
	return f.onContributionEarned
}
func (f *testFormula) TargetSourceForWithdrawl() core.FinancialSource {
	return f.onTargetSourceForWithdrawl
}
func (f *testFormula) TargetSourceForContribution() core.FinancialSource {
	return f.onTargetSourceForContribution
}
func (f *testFormula) AllowedIncomeSources() []core.FinancialSource {
	return f.onAllowedIncomeSources
}
func (f *testFormula) Validate() error {
	return f.onValidate
}
func (f *testFormula) Clone() Formula {
	return f
}
