package rrsp

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
)

var _ core.TaxCalculator = (*testTaxCalculator)(nil)

type testTaxCalculator struct {
	_currentIndex               int // do not set
	onTaxPayableSpouseA         []float64
	onTaxPayableSpouseB         []float64
	financesPassedOnSetFinances core.HouseholdFinances
	creditsPassedOnSetFinances  []core.TaxCredit
	depsPassedOnSetDependents   []*human.Person
}

func (ttc *testTaxCalculator) TaxPayable() (float64, float64, []core.TaxCredit) {
	spouseA := ttc.onTaxPayableSpouseA[ttc._currentIndex]
	spouseB := ttc.onTaxPayableSpouseB[ttc._currentIndex]
	ttc._currentIndex++
	return spouseA, spouseB, nil
}
func (ttc *testTaxCalculator) SetFinances(f core.HouseholdFinances, cr []core.TaxCredit) {
	ttc.financesPassedOnSetFinances = f
	ttc.creditsPassedOnSetFinances = cr
}
func (ttc *testTaxCalculator) SetDependents(deps []*human.Person) {
	ttc.depsPassedOnSetDependents = deps
}
func (ttc *testTaxCalculator) Regions() []core.Region {
	return nil
}
func (ttc *testTaxCalculator) Year() uint {
	return 0
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
