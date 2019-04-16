package finance

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

const (
	DeducSrcUnknown DeductionSource = iota
	// Recognized sources of deduction
	DeducSrcRRSP    // contribution to RRSP
	DeducSrcMedical // medical expenses
)

const (
	MiscSrcUnknown MiscSource = iota
	// recognized sources of miscellaneous amounts
	MiscSrcMedical // medical expenses
)
