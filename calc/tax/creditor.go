package tax

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given amount. The net
	// income may or may not be used by the underlying implementation
	TaxCredits(amount, netIncome float64) Credits
	// Source returns the credit source this creditor
	Source() CreditSource
	// Clone returns a copy of this creditor
	Clone() Creditor
}

// CreditSource represents the name of a tax credit source
type CreditSource string

// credits represent an amount that reduces payable tax
type Credits struct {
	Source       CreditSource // the name of the source of the credits
	IsRefundable bool         // if true, the amount is paid back if not used
	Amount       float64      // the amount owed to tax payer
}
