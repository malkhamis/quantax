package rrsp

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
)

var (
	_ core.TaxCalculator           = (*testTaxCalculator)(nil)
	_ core.HouseholdFinanceMutator = (*testHouseholdFinances)(nil)
)

type testTaxCalculator struct {
	_currentIndex               int // do not set
	onTaxPayableSpouseA         []float64
	onTaxPayableSpouseB         []float64
	onTaxPayableCredits         [][]core.TaxCredit
	financesPassedOnSetFinances []core.HouseholdFinances
	creditsPassedOnSetFinances  [][]core.TaxCredit
	depsPassedOnSetDependents   [][]*human.Person
}

func (ttc *testTaxCalculator) TaxPayable() (float64, float64, []core.TaxCredit) {
	spouseA := ttc.onTaxPayableSpouseA[ttc._currentIndex]
	spouseB := ttc.onTaxPayableSpouseB[ttc._currentIndex]
	credits := ttc.onTaxPayableCredits[ttc._currentIndex]
	ttc._currentIndex++
	return spouseA, spouseB, credits
}
func (ttc *testTaxCalculator) SetFinances(f core.HouseholdFinances, cr []core.TaxCredit) {
	ttc.financesPassedOnSetFinances = append(ttc.financesPassedOnSetFinances, f)
	ttc.creditsPassedOnSetFinances = append(ttc.creditsPassedOnSetFinances, cr)
}
func (ttc *testTaxCalculator) SetDependents(deps []*human.Person) {
	ttc.depsPassedOnSetDependents = append(ttc.depsPassedOnSetDependents, deps)
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
