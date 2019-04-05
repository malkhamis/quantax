package tax

import "github.com/malkhamis/quantax/calc/finance"

// compile-time check for interface implementatino
var (
	_ Adjuster = CanadianCapitalGainAdjuster{}
)

// Adjuster is a type that adjusts any amount given individual finances. The
// type is mainly used by tax calculators in conjuction with tax forumlas to
//adjust specific income/deduction sources for taxation
type Adjuster interface {
	// Adjusted returns an adjustment amount from the given finances
	Adjusted(*finance.IndividualFinances) Adjustment
	// Clone returns a copy of this adjuster
	Clone() Adjuster
}

// TaxCredits represents tax credits that can be applied to payable tax
type TaxCredits struct {
	// amounts that is used to reduce payable tax.
	// Remaining amount are paid back
	Refundable float64
	// amounts that reduces taxes only. Remaining
	// amount is lost and is never paid back
	NonRefundable float64
}

// Adjustment stores the result of computing an adjustment for a source of
// income/deduction as well as any associated tax credits
type Adjustment struct {
	AdjustedAmount float64
	Credits        TaxCredits
}

// CanadianCapitalGainAdjuster returns the adjusted income for Canadian-sourced
// capital gain income
type CanadianCapitalGainAdjuster struct {
	// the proportion of the taxable capital gain amount
	TaxableProportion float64
}

// Adjusted returns the adjusted amount of capital gain income by setting the
// adjusted amount as follows:
//  AdjustedAmount = (TaxableProportion) x (Capital Gain Income)
func (cg CanadianCapitalGainAdjuster) Adjusted(finances *finance.IndividualFinances) Adjustment {

	income := finances.TotalIncome(finance.IncSrcCapitalGainCA)
	taxableAmount := cg.TaxableProportion * income
	return Adjustment{AdjustedAmount: taxableAmount}
}

// Clone returns a copy of this instance
func (cg CanadianCapitalGainAdjuster) Clone() Adjuster {
	return cg
}
