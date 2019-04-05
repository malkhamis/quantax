package tax

import "github.com/malkhamis/quantax/calc/finance"

// Adjuster is a type that adjusts any amount given individual finances. The
// type is mainly used to adjust specific income/deduction sources for taxation
type Adjuster interface {
	// Adjusted returns an adjusted amount from given finances
	Adjusted(*finance.IndividualFinances) Adjustment
	// Clone returns a copy of this adjuster
	Clone() Adjuster
}

// Adjustment is used to represent adjustment amounts that are applied to tax
// calculations
type Adjustment struct {
	AdjustedAmount          float64
	TaxCreditsRefundable    float64
	TaxCreditsNonRefundable float64
}

// CanadianCapitalGainAdjuster returns the required income adjustment to
// accurately compute the tax on Canadian-sourced capital gain income
type CanadianCapitalGainAdjuster struct {
	// the proportion of the taxable capital gain amount
	TaxableProportion float64
}

// Adjusted returns the capital gain amount that should be added to the net
// income before calculating the tax
func (cg CanadianCapitalGainAdjuster) Adjusted(finances *finance.IndividualFinances) Adjustment {

	adjIncSrc := finances.TotalIncome(finance.IncSrcCapitalGainCA) * cg.TaxableProportion
	return Adjustment{AdjustedAmount: adjIncSrc}
}

// Clone returns a copy of this instance
func (cg CanadianCapitalGainAdjuster) Clone() Adjuster {
	return cg
}
