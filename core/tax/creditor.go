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

// CreditDescriptor is a convenience type used by the types that implement the
// Creditor interface
type CreditDescriptor struct {
	// the credit rule to return when Rule is called
	CreditRule core.CreditRule
	// the target financial source to return when FinancialSource is called
	TargetFinancialSource core.FinancialSource
	// a short description of the tax credit
	CreditDescription string
}

// Rule returns the credit rule set in this credit descriptor
func (cd CreditDescriptor) Rule() core.CreditRule {
	return cd.CreditRule
}

// FinancialSource returns the financial source set in this credit descriptor
func (cd CreditDescriptor) FinancialSource() core.FinancialSource {
	return cd.TargetFinancialSource
}

// Description returns a description set in this credit descriptor
func (cd CreditDescriptor) Description() string {
	return cd.CreditDescription
}
