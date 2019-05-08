package factory

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/finance"
)

// FinanceFactory is a type used to conveniently create finances
type FinanceFactory struct{}

// NewFinanceFactory returns a new factory for finances
func NewFinanceFactory() *FinanceFactory {
	return &FinanceFactory{}
}

// NewFinances returns a new finance mutator initialized with the given amounts
func (f *FinanceFactory) NewFinances(amounts map[core.FinancialSource]float64) core.FinanceMutator {

	finances := finance.NewIndividualFinances()
	for src, amount := range amounts {
		finances.SetAmount(src, amount)
	}
	return finances
}

// NewHouseholFinancesForCouple returns a household finances instance. References
// for the finances of both spouses are initialized with the given amounts and
// are never nil
func (f *FinanceFactory) NewHouseholdFinancesForCouple(amountsSpouseA, amountsSpouseB map[core.FinancialSource]float64) core.HouseholdFinanceMutator {

	financesA := finance.NewIndividualFinances()
	for src, amount := range amountsSpouseA {
		financesA.SetAmount(src, amount)
	}

	financesB := finance.NewIndividualFinances()
	for src, amount := range amountsSpouseB {
		financesB.SetAmount(src, amount)
	}

	finances := finance.NewHouseholdFinances(financesA, financesB)
	return finances
}

// NewHouseholdFinancesForSingle returns a household finances instance. The
// reference for spouses A is initialized with the given amounts and is never
// nil. The reference for Spouse B is always nil
func (f *FinanceFactory) NewHouseholdFinancesForSingle(amounts map[core.FinancialSource]float64) core.HouseholdFinanceMutator {

	financesA := finance.NewIndividualFinances()
	for src, amount := range amounts {
		financesA.SetAmount(src, amount)
	}

	finances := finance.NewHouseholdFinances(financesA, nil)
	return finances
}
