package core

// Financer is a type that holds the financial data for an individual
type Financer interface {
	// TotalAmount returns the sum of of the given sources only. If no sources
	// given, it returns zero
	TotalAmount(sources ...FinancialSource) float64
	// IncomeSources returns a set of all income sources in this instance. The
	// returned map is never nil
	IncomeSources() []FinancialSource
	// DeductionSources returns a set of all deduciton sources in this instance
	// The returned map is never nil
	DeductionSources() []FinancialSource
	// MiscSources returns a set of all miscellaneous sources in this instance
	// The returned map is never nil
	MiscSources() []FinancialSource
	// AllSources returns all financial sources in this instance
	AllSources() []FinancialSource
	// Clone returns a deep copy of this instance
	Clone() FinanceMutator
}

// FinanceMutator exposes financial data for mutations
type FinanceMutator interface {
	// SetAmount force-set the given amount for the given source. If the source
	// is within the range of recognized income sources, it is set as an income
	// source.If it is within the range of recognized deduction sources, it is
	// added as a deduction. Other sources, including unknown ones are added as
	// misc amounts
	SetAmount(source FinancialSource, amount float64)
	// AddAmount adds the given amount to the finances. If the source is within
	// the range of recognized income sources, it is added as an income source.If
	// it is within the range of recognized deduction sources, it is added as a
	// deduction. Other sources, including unknown ones are added as misc amounts
	AddAmount(source FinancialSource, amount float64)
	// RemoveAmounts removes the given stored sources. If sources is empty, the
	// call is noop. This operation ensures subsequent calls to IncomeSources,
	// DeductionSources, and MiscSources returns a list that does not contain
	// the given source(s)
	RemoveAmounts(sources ...FinancialSource)
	Financer
}

// HouseholdFinances holds the financial data for a household of two spouses
type HouseholdFinances interface {
	// SpouseA returns the financia data of the first spouse
	SpouseA() Financer
	// SpouseB returns the financia data of the second spouse
	SpouseB() Financer
	// Clone returns a deep copy of the instance
	Clone() HouseholdFinanceMutator
}

type HouseholdFinanceMutator interface {
	MutableSpouseA() FinanceMutator
	MutableSpouseB() FinanceMutator
	HouseholdFinances
}
