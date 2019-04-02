package finance

// IncomeSource represents the source type of income
type IncomeSource int

const (
	IncSrcUnknown IncomeSource = iota
	// Recognized sources of income
	IncSrcEarned   // employment and labour income
	IncSrcInterest // e.g. saving account interest
	IncSrcRRSP     // withdrawal from RRSP
	IncSrcUCCB     // universal child care benefits
	IncSrcRDSP     // registered disability saving plan
	IncSrcTFSA     // tax-free saving account
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