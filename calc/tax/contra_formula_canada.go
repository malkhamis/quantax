package tax

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

var (
	// compile-time check for interface implementation
	_ ContraFormula = (*CanadianContraFormula)(nil)
)

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
func (cf *CanadianContraFormula) Apply(finances *finance.IndividualFinances, netIncome float64) []Credits {

	// creditSources := make(map[CreditSource]Credits)
	//
	// for source, srcIncome := range finances.Income {
	//
	// 	creditor := cf.CreditsFromIncome[source]
	// 	credits := creditor.TaxCredits(srcIncome, netIncome)
	// 	if credits.Amount == 0 {
	// 		continue
	// 	}
	// }

	return nil // Not implemented
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

	if cf.ApplicationOrder != nil {
		clone.ApplicationOrder = make([]CreditSource, len(cf.ApplicationOrder))
		for i, s := range cf.ApplicationOrder {
			clone.ApplicationOrder[i] = s
		}
	}

	return clone
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {

	appOrderCreditSrcSet, dups := creditSources(cf.ApplicationOrder).makeSetAndGetDuplicates()
	if len(dups) > 0 {
		return errors.Wrapf(ErrDupCreditSource, "%v", dups)
	}

	for incomeSrc, creditor := range cf.CreditsFromIncome {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "income source %q", incomeSrc)
		}

		_, exist := appOrderCreditSrcSet[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"income source %q -> credit source %q: must be in application order list",
				incomeSrc, creditor.Source(),
			)
		}

	}

	for deducSrc, creditor := range cf.CreditsFromDeduction {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "deduction source %q", deducSrc)
		}

		_, exist := appOrderCreditSrcSet[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"deduction source %q -> credit source %q: must be in application order list",
				deducSrc, creditor.Source(),
			)
		}

	}

	return nil
}
