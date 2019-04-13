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
	// CreditsFromIncome stores creditors associated with income sources
	CreditsFromIncome map[finance.IncomeSource]Creditor
	// CreditsFromDeduction stores creditors associated with deduction sources
	CreditsFromDeduction map[finance.DeductionSource]Creditor
	// CreditsFromMiscAmounts stores creditors associated with misc sources
	CreditsFromMiscAmounts map[finance.MiscSource]Creditor
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

	if cf.CreditsFromMiscAmounts != nil {
		clone.CreditsFromMiscAmounts = make(map[finance.MiscSource]Creditor)
		for source, creditor := range cf.CreditsFromMiscAmounts {
			clone.CreditsFromMiscAmounts[source] = creditor.Clone()
		}
	}

	if cf.ApplicationOrder != nil {
		clone.ApplicationOrder = make([]CreditSource, len(cf.ApplicationOrder))
		copy(clone.ApplicationOrder, cf.ApplicationOrder)
	}

	return clone
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {

	appOrderCreditSrcSet, dups := creditSources(
		cf.ApplicationOrder,
	).makeSetAndGetDuplicates()

	if len(dups) > 0 {
		return errors.Wrapf(ErrDupCreditSource, "%v", dups)
	}

	err := cf.checkIncSrcCreditorsInSet(appOrderCreditSrcSet)
	if errors.Cause(err) == ErrUnknownCreditSource {
		err = errors.Wrap(
			err,
			"income-source creditors must return credit sources which are known "+
				"in the application order list",
		)
	}
	if err != nil {
		return err
	}

	err = cf.checkDeducSrcCreditorsInSet(appOrderCreditSrcSet)
	if errors.Cause(err) == ErrUnknownCreditSource {
		err = errors.Wrap(
			err,
			"deduction-source creditors must return credit sources which are known "+
				"in the application order list",
		)
	}
	if err != nil {
		return err
	}

	err = cf.checkMiscSrcCreditorsInSet(appOrderCreditSrcSet)
	if errors.Cause(err) == ErrUnknownCreditSource {
		err = errors.Wrap(
			err,
			"misc-source creditors must return credit sources which are known "+
				"in the application order list",
		)
	}
	if err != nil {
		return err
	}

	return nil
}

func (cf *CanadianContraFormula) checkIncSrcCreditorsInSet(set map[CreditSource]struct{}) error {

	for incomeSrc, creditor := range cf.CreditsFromIncome {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "income source %q", incomeSrc)
		}

		_, exist := set[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"income source %q -> credit source %q",
				incomeSrc, creditor.Source(),
			)
		}

	}

	return nil
}

func (cf *CanadianContraFormula) checkDeducSrcCreditorsInSet(set map[CreditSource]struct{}) error {

	for deducSrc, creditor := range cf.CreditsFromDeduction {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "deduction source %q", deducSrc)
		}

		_, exist := set[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"deduction source %q -> credit source %q",
				deducSrc, creditor.Source(),
			)
		}

	}

	return nil
}

func (cf *CanadianContraFormula) checkMiscSrcCreditorsInSet(set map[CreditSource]struct{}) error {

	for miscSrc, creditor := range cf.CreditsFromMiscAmounts {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "misc source %q", miscSrc)
		}

		_, exist := set[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"misc source %q -> credit source %q",
				miscSrc, creditor.Source(),
			)
		}

	}

	return nil
}
