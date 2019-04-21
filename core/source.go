package core

// Source represents a financial source
type FinancialSource int

const (
	SrcUnknown FinancialSource = iota

	IncomeSourcesBegin
	IncSrcEarned        // employment and labour income
	IncSrcInterest      // e.g. saving account interest
	IncSrcCapitalGainCA // Canadian-sourced capital gain on sold assets
	IncSrcRRSP          // withdrawal from RRSP
	IncSrcUCCB          // universal child care benefits
	IncSrcRDSP          // registered disability saving plan
	IncSrcTFSA          // tax-free saving account
	IncomeSourcesEnd

	DeductionSourcesBegin
	DeducSrcRRSP    // contribution to RRSP
	DeducSrcMedical // medical expenses
	DeductionSourcesEnd

	MiscSourcesBegin
	MiscSrcMedical // medical expenses
	MiscSourcesEnd
)

// AmountBySource maps amounts to their financial source types
type AmountBySource map[FinancialSource]float64

// Sum returns the sum of all amounts
func (s AmountBySource) Sum() float64 {
	var total float64
	for _, amount := range s {
		total += amount
	}
	return total
}

// Clone returns a copy of this instance
func (s AmountBySource) Clone() AmountBySource {

	var clone AmountBySource

	if s != nil {
		clone = make(AmountBySource)
	}

	for source, amount := range s {
		clone[source] = amount
	}

	return clone
}
