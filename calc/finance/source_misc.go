package finance

// MiscSource represents the source type of any miscellaneous amount that is
// anything other than income or deduction
type MiscSource int

// MiscAmountsBySource represents misc amounts according to their sources
type MiscAmountsBySource map[MiscSource]float64

// Sum returns the sum of all miscellaneous sources
func (s MiscAmountsBySource) Sum() float64 {
	var total float64
	for _, amount := range s {
		total += amount
	}
	return total
}

// Clone returns a copy of this instance
func (s MiscAmountsBySource) Clone() MiscAmountsBySource {

	var clone MiscAmountsBySource

	if s != nil {
		clone = make(MiscAmountsBySource)
	}

	for source, amount := range s {
		clone[source] = amount
	}

	return clone
}

// MiscSourceSet is a convenience type used for misc source lookup
type MiscSourceSet map[MiscSource]struct{}

// Has returns true if the given source exists in this set
func (set MiscSourceSet) Has(source MiscSource) bool {
	_, ok := set[source]
	return ok
}
