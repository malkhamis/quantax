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

// MiscSource represents the source type of any miscellaneous amount that is
// anything other than income or deduction
type MiscSource int

const (
	MiscSrcUnknown MiscSource = iota
	// recognized sources of miscellaneous amounts
	MiscSrcMedical // medical expenses
)

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

// IncomeSourceSet is a convenience type used for income source lookup
type IncomeSourceSet map[IncomeSource]struct{}

// Has returns true if the given source exists in this set
func (set IncomeSourceSet) Has(source IncomeSource) bool {
	_, ok := set[source]
	return ok
}

// DeductionSourceSet is a convenience type used for deduction source lookup
type DeductionSourceSet map[DeductionSource]struct{}

// Has returns true if the given source exists in this set
func (set DeductionSourceSet) Has(source DeductionSource) bool {
	_, ok := set[source]
	return ok
}

// MiscSourceSet is a convenience type used for misc source lookup
type MiscSourceSet map[MiscSource]struct{}

// Has returns true if the given source exists in this set
func (set MiscSourceSet) Has(source MiscSource) bool {
	_, ok := set[source]
	return ok
}
