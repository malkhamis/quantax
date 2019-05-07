package tax

import "github.com/malkhamis/quantax/core"

// CreditorConst is a Creditor that returns a constant amount
type CreditorConst struct {
	// the amount to return when TaxCredit is called
	Amount float64
	// the credit rule to return when Rule is called
	CreditRule core.CreditRule
	// the target financial source to return when FinancialSource is called
	TargetFinancialSource core.FinancialSource
	// a short description of the tax credit
	CreditDescription string
}

// TaxCredit returns the Amount set in this creditor
func (cc CreditorConst) TaxCredit(_ *TaxPayer) float64 {
	return cc.Amount
}

// Rule returns the credit rule set in this creditor
func (cc CreditorConst) Rule() core.CreditRule {
	return cc.CreditRule
}

// FinancialSource returns the financial source set in this creditor
func (cc CreditorConst) FinancialSource() core.FinancialSource {
	return cc.TargetFinancialSource
}

// Description returns a description of the tax credit
func (cc CreditorConst) Description() string {
	return cc.CreditDescription
}

// Clone returns a deep copy of this creditor
func (cc CreditorConst) Clone() Creditor {
	return cc.clone()
}

// clone returns a copy of this creditor
func (cc CreditorConst) clone() CreditorConst {
	return cc
}
