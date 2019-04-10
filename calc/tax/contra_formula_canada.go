package tax

import "github.com/malkhamis/quantax/calc/finance"

// compile-time check for interface implementation
var _ ContraFormula = (*CanadianContraFormula)(nil)

// CanadianContraFormula is used to calculate Canadian tax credits
type CanadianContraFormula struct {
	// CreditsFromIncome returns income sources that provide tax credits
	CreditsFromIncome map[finance.IncomeSource]Creditor
	// CreditsFromDeduction returns income sources that provide tax credits
	CreditsFromDeduction map[finance.DeductionSource]Creditor
	// ApplicationOrder stores the order in which tax credits are used
	ApplicationOrder []CreditSource
}

// Apply applies the contra-formula on the income and the set finances
func (cf *CanadianContraFormula) Apply(finances *finance.IndividualFinances, netIncome float64) map[CreditSource]Credits {
	return nil // Not Implemented
}

// OrderOfUse returns the order in which credit sources are used to reduce tax
func (cf *CanadianContraFormula) OrderOfUse() []CreditSource {
	return cf.ApplicationOrder
}

// Clone returns a copy of this contra-formula
func (cf *CanadianContraFormula) Clone() ContraFormula {

	if cf == nil {
		return nil
	}

	clone := new(CanadianContraFormula)

	if cf.CreditsFromDeduction != nil {
		clone.CreditsFromDeduction = make(map[finance.DeductionSource]Creditor)
		for source, creditor := range cf.CreditsFromDeduction {
			clone.CreditsFromDeduction[source] = creditor.Clone()
		}
	}

	if cf.CreditsFromIncome != nil {
		clone.CreditsFromIncome = make(map[finance.IncomeSource]Creditor)
		for source, creditor := range cf.CreditsFromIncome {
			clone.CreditsFromIncome[source] = creditor.Clone()
		}
	}

	return clone
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {
	return nil
}
