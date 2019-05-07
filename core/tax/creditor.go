package tax

import (
	"github.com/malkhamis/quantax/core"
)

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given financial source
	TaxCredit(taxPayer *TaxPayer) float64
	// Rule is the rule of using the returned tax credit amount
	Rule() core.CreditRule
	// FinancialSource returns the target financial source of this creditor
	FinancialSource() core.FinancialSource
	// Description returns a short description of the tax credit
	Description() string
	// Clone returns a copy of this creditor
	Clone() Creditor
}
