package core

// Source represents a financial source
type FinancialSource int

const (
	SrcNone FinancialSource = iota

	IncomeSourcesBegin
	IncSrcEarned                 // employment and labour income
	IncSrcInterest               // e.g. saving account interest
	IncSrcCapitalGainCA          // Canadian-sourced capital gain on sold assets
	IncSrcEligibleDividendsCA    // Canadian eligible dividends
	IncSrcNonEligibleDividendsCA // Canadian non-eligible dividends
	IncSrcForeignDividends       // non-Canadian sourced dividends
	IncSrcRRSP                   // withdrawal from RRSP
	IncSrcUCCB                   // universal child care benefits
	IncSrcRDSP                   // registered disability saving plan
	IncSrcTFSA                   // tax-free saving account
	IncomeSourcesEnd

	DeductionSourcesBegin
	DeducSrcChildCareExpense // child-care expenses (TODO: adjuster)
	DeducSrcRRSP             // contribution to RRSP
	DeducSrcOthers           // other deduction
	DeductionSourcesEnd

	MiscSourcesBegin
	MiscSrcMedical // medical expenses (TODO: creditor)
	MiscSrcTuition // tuition expenses
	MiscSrcOthers  // other amounts (unaccounted for)
	MiscSourcesEnd
)

// IsIncomeSource returns true if this source is an identified income source
func (s FinancialSource) IsIncomeSource() bool {
	return s > IncomeSourcesBegin && s < IncomeSourcesEnd
}

// IsDeductionSource returns true if this source is an identified deduction
// source
func (s FinancialSource) IsDeductionSource() bool {
	return s > DeductionSourcesBegin && s < DeductionSourcesEnd
}

// IsMiscSource returns true if this source is an identified misc source
func (s FinancialSource) IsMiscSource() bool {
	return s > MiscSourcesBegin && s < MiscSourcesEnd
}

// IsUnknownSource returns true if this source is an unidentified source
func (s FinancialSource) IsUnknownSource() bool {
	return !s.IsIncomeSource() && !s.IsDeductionSource() && !s.IsMiscSource()
}
