package tax

import "github.com/malkhamis/quantax/calc/finance"

// CanadianContraFormula is used to calculate Canadian tax credits
type CanadianContraFormula struct {
	// CreditsFromIncome returns income sources that provide tax credits
	CreditsFromIncome map[finance.IncomeSource]Creditor
	// CreditsFromDeduction returns income sources that provide tax credits
	CreditsFromDeduction map[finance.DeductionSource]Creditor
	// ApplicationOrder stores the order in which tax credits are used
	ApplicationOrder []CreditSource
}
