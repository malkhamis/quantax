package core

// compile-time check for interface implementation
var (
	_ Financer          = (*financerNop)(nil)
	_ HouseholdFinances = (*householdFinancesNop)(nil)
)

// NewFinancerNop returns a no-op Financer instance. Packages implementing the
// interfaces defined by this package may use instances created by this function
// if the user passes nil finances
func NewFinancerNop() Financer {
	return &financerNop{}
}

// NewFinancerNop returns a no-op HouseholdFinances instance. Packages
// implementing the interfaces defined by this package may use instances created
// by this function if the user passes nil finances
func NewHouseholdFinancesNop() HouseholdFinances {
	return &householdFinancesNop{
		spouseA: &financerNop{},
		spouseB: &financerNop{},
	}
}

type financerNop struct {
	_ byte // to guarantee a new instance
}

func (nop *financerNop) TotalAmount(_ ...FinancialSource) float64 {
	return 0
}
func (nop *financerNop) IncomeSources() []FinancialSource {
	return nil
}
func (nop *financerNop) DeductionSources() []FinancialSource {
	return nil
}
func (nop *financerNop) MiscSources() []FinancialSource {
	return nil
}
func (nop *financerNop) AllSources() []FinancialSource {
	return nil
}
func (nop *financerNop) Clone() FinanceMutator {
	return &financerNop{}
}
func (nop *financerNop) Version() uint64 {
	return 0
}
func (nop *financerNop) SetAmount(_ FinancialSource, _ float64) {}
func (nop *financerNop) AddAmount(_ FinancialSource, _ float64) {}
func (nop *financerNop) RemoveAmounts(_ ...FinancialSource)     {}

type householdFinancesNop struct {
	spouseA *financerNop
	spouseB *financerNop
}

func (nop *householdFinancesNop) SpouseA() Financer {
	return nop.spouseA
}
func (nop *householdFinancesNop) SpouseB() Financer {
	return nop.spouseB
}
func (nop *householdFinancesNop) Clone() HouseholdFinances {
	return NewHouseholdFinancesNop()
}
func (nop *householdFinancesNop) Version() uint64 {
	return 0
}
