package finance

// IncomeSource represents the source type of income
type IncomeSource int

// Recognized sources of income
const (
	IncSrcUnknown IncomeSource = iota

	IncSrcEarned              // e.g. employment and labour income
	IncSrcRRSP                // registered retirement saving plan
	IncSrcTFSA                // tax-free saving account
	IncSrcUCCB                // universal child care benefits
	IncSrcRDSP                // registered disability saving plan
	IncSrcInterest            // e.g. saving account interest
	IncSrcDividendsRegular    // regular Canadian corp dividens
	IncSrcDividendsNonRegular // non-regular Canadian corp dividens
	IncSrcCapitalGain         // e.g. capital gain on sale of assets
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

// DeductionSource represents the source type of deduction
type DeductionSource int

// Recognized sources of deducitons
const (
	DeducSrcUnknown DeductionSource = iota
	// Recognized sources of deduction
	DeducSrcRRSP // registered retirement saving plan
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
