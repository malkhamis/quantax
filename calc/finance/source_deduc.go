package finance

// DeductionSource represents the source type of deduction
type DeductionSource int

// DeductionBySource represents deduction amounts according to their sources
type DeductionBySource map[DeductionSource]float64

// Sum returns the sum of all deduction sources
func (s DeductionBySource) Sum() float64 {
	var total float64
	for _, deduction := range s {
		total += deduction
	}
	return total
}

// Clone returns a copy of this instance
func (s DeductionBySource) Clone() DeductionBySource {

	var clone DeductionBySource

	if s != nil {
		clone = make(DeductionBySource)
	}

	for source, deduction := range s {
		clone[source] = deduction
	}

	return clone
}

// DeductionSourceSet is a convenience type used for deduction source lookup
type DeductionSourceSet map[DeductionSource]struct{}

// Has returns true if the given source exists in this set
func (set DeductionSourceSet) Has(source DeductionSource) bool {
	_, ok := set[source]
	return ok
}
