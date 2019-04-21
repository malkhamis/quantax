package finance

// IncomeSource represents the source type of income
type IncomeSource int

// DeductionBySource represents deduction amounts according to their sources
type IncomeBySource map[IncomeSource]float64

// Sum returns the sum of all income sources
func (s IncomeBySource) Sum() float64 {
	var total float64
	for _, income := range s {
		total += income
	}
	return total
}

// Clone returns a copy of this instance
func (s IncomeBySource) Clone() IncomeBySource {

	var clone IncomeBySource

	if s != nil {
		clone = make(IncomeBySource)
	}

	for source, income := range s {
		clone[source] = income
	}

	return clone
}

// IncomeSourceSet is a convenience type used for income source lookup
type IncomeSourceSet map[IncomeSource]struct{}

// Has returns true if the given source exists in this set
func (set IncomeSourceSet) Has(source IncomeSource) bool {
	_, ok := set[source]
	return ok
}
