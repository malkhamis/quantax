package finance

import (
	"github.com/malkhamis/quantax/core"
)

// compile-time check for inteface implementation
var _ core.HouseholdFinances = (*HouseholdFinances)(nil)

// HouseholdFinances represents financial data for a couple, family etc
type HouseholdFinances struct {
	spouseA *IndividualFinances
	spouseB *IndividualFinances
}

// NewHouseholdFinances returns a new household finance instance. Future change
// to the given individual finances are reflected in the returned instance
func NewHouseholdFinances(spouseA, spouseB *IndividualFinances) *HouseholdFinances {

	return &HouseholdFinances{
		spouseA: spouseA,
		spouseB: spouseB,
	}
}

// SpouseA returns a reference to the individual finances of the first spouse.
// If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) SpouseA() core.Financer {
	if hf == nil || hf.spouseA == nil {
		return nil
	}
	return hf.spouseA
}

// SpouseA returns a reference to the individual finances of the second spouse.
// If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) SpouseB() core.Financer {
	if hf == nil || hf.spouseB == nil {
		return nil
	}
	return hf.spouseB
}

// MutableSpouseA returns a reference to the individual finances of the first
// spouse for mutations. If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) MutableSpouseA() core.FinanceMutator {
	if hf == nil || hf.spouseA == nil {
		return nil
	}
	return hf.spouseA
}

// MutableSpouseB returns a reference to the individual finances of the second
// spouse for mutations. If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) MutableSpouseB() core.FinanceMutator {
	if hf == nil || hf.spouseB == nil {
		return nil
	}
	return hf.spouseB
}

// Clone returns a copy of this instance. If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) Clone() core.HouseholdFinanceMutator {

	if hf == nil {
		return nil
	}

	return &HouseholdFinances{
		spouseA: hf.spouseA.clone(),
		spouseB: hf.spouseB.clone(),
	}
}

// clone returns a copy of this instance. If 'hf' is nil, it returns nil
func (hf *HouseholdFinances) clone() *HouseholdFinances {

	if hf == nil {
		return nil
	}

	clone := &HouseholdFinances{
		spouseA: hf.spouseA.clone(),
		spouseB: hf.spouseB.clone(),
	}

	return clone
}
