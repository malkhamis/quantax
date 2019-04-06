package finance

// IncomeSource represents the source type of income
type IncomeSource int

const (
	IncSrcUnknown IncomeSource = iota
	// Recognized sources of income
	IncSrcEarned        // employment and labour income
	IncSrcInterest      // e.g. saving account interest
	IncSrcCapitalGainCA // Canadian-sourced capital gain on sold assets
	IncSrcRRSP          // withdrawal from RRSP
	IncSrcUCCB          // universal child care benefits
	IncSrcRDSP          // registered disability saving plan
	IncSrcTFSA          // tax-free saving account
)

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

// DeductionSource represents the source type of deduction
type DeductionSource int

const (
	DeducSrcUnknown DeductionSource = iota
	// Recognized sources of deduction
	DeducSrcRRSP    // contribution to RRSP
	DeducSrcMedical // medical expenses
)

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

// IncomeSourceSet is a convenience type used for income source lookup
type IncomeSourceSet map[IncomeSource]struct{}

// NewIncomeSourceSet returns a new set from the given income sources
func NewIncomeSourceSet(sources ...IncomeSource) IncomeSourceSet {

	set := make(IncomeSourceSet)
	for _, s := range sources {
		set[s] = struct{}{}
	}
	return set
}

// Has returns true if the given source exists in this set
func (set IncomeSourceSet) Has(source IncomeSource) bool {
	_, ok := set[source]
	return ok
}

// DeductionSourceSet is a convenience type used for deduction source lookup
type DeductionSourceSet map[DeductionSource]struct{}

// NewDeductionSourceSet returns a new set from the given deduciton sources
func NewDeductionSourceSet(sources ...DeductionSource) DeductionSourceSet {

	set := make(DeductionSourceSet)
	for _, s := range sources {
		set[s] = struct{}{}
	}
	return set
}

// Has returns true if the given source exists in this set
func (set DeductionSourceSet) Has(source DeductionSource) bool {
	_, ok := set[source]
	return ok
}
