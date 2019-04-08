package rrsp

import "github.com/malkhamis/quantax/calc/finance"

type testTaxCalculator struct {
	currentIndex int
	onCalc       []float64
}

func (ttc *testTaxCalculator) Calc(_ *finance.IndividualFinances) float64 {
	currentVal := ttc.onCalc[ttc.currentIndex]
	ttc.currentIndex++
	return currentVal
}

type testFormula struct {
	onValidate                    error
	onContributionEarned          float64
	onTargetSourceForWithdrawl    finance.IncomeSource
	onTargetSourceForContribution finance.DeductionSource
	onAllowedIncomeSources        []finance.IncomeSource
}

func (f *testFormula) ContributionEarned(income float64) float64 {
	return f.onContributionEarned
}
func (f *testFormula) TargetSourceForWithdrawl() finance.IncomeSource {
	return f.onTargetSourceForWithdrawl
}
func (f *testFormula) TargetSourceForContribution() finance.DeductionSource {
	return f.onTargetSourceForContribution
}
func (f *testFormula) AllowedIncomeSources() []finance.IncomeSource {
	return f.onAllowedIncomeSources
}
func (f *testFormula) Validate() error {
	return f.onValidate
}
func (f *testFormula) Clone() Formula {
	return f
}
