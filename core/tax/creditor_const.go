package tax

import "github.com/malkhamis/quantax/core"

// TODO docstring
type CreditorConst struct {
	Amount                float64
	CreditRule            core.CreditRule
	TargetFinancialSource core.FinancialSource
	CreditDescription     string
}

// TODO: docstring
func (cc CreditorConst) TaxCredit(_ *TaxPayer) float64 {
	return cc.Amount
}

// TODO: docstring
func (cc CreditorConst) Rule() core.CreditRule {
	return cc.CreditRule
}

// TODO: docstring
func (cc CreditorConst) FinancialSource() core.FinancialSource {
	return cc.TargetFinancialSource
}

// TODO: docstring
func (cc CreditorConst) Description() string {
	return cc.CreditDescription
}

// TODO: docstring
func (cc CreditorConst) Clone() Creditor {
	return cc.clone()
}

// TODO: docstring
func (cc CreditorConst) clone() CreditorConst {
	return cc
}
