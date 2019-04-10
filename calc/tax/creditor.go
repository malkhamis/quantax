package tax

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given amount. The net
	// income may or may not be used by the underlying implementation
	TaxCredits(amount, netIncome float64) Credits
	// Clone returns a copy of this creditor
	Clone() Creditor
}
