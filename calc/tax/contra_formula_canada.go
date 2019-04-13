package tax

import (
	"sort"

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

	if finances == nil {
		return nil
	}

	incSrcCredits := cf.creditsFromIncSrcs(finances, netIncome)
	deducSrcCredits := cf.creditsFromDeducSrcs(finances, netIncome)
	miscSrcCredits := cf.creditsFromMiscSrcs(finances, netIncome)

	allCredits := append(
		incSrcCredits,
		append(
			deducSrcCredits,
			miscSrcCredits...)...,
	)

	cf.orderCreditGroupInPlace(allCredits)
	return allCredits
}

// creditsFromIncSrcs returns a list of credits extracted from the income
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromIncSrcs(finances *finance.IndividualFinances, netIncome float64) []Credits {

	var creditGroup []Credits

	for source, srcIncome := range finances.Income {

		creditor := cf.CreditsFromIncome[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredits(srcIncome, netIncome)
		if credits.Amount == 0 {
			continue
		}

		creditGroup = append(creditGroup, credits)
	}

	return creditGroup
}

// creditsFromDeducSrcs returns a list of credits extracted from the deduciton
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromDeducSrcs(finances *finance.IndividualFinances, netIncome float64) []Credits {

	var creditGroup []Credits

	for source, srcDeduc := range finances.Deductions {

		creditor := cf.CreditsFromDeduction[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredits(srcDeduc, netIncome)
		if credits.Amount == 0 {
			continue
		}

		creditGroup = append(creditGroup, credits)
	}

	return creditGroup
}

// creditsFromMiscSrcs returns a list of credits extracted from the misc
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromMiscSrcs(finances *finance.IndividualFinances, netIncome float64) []Credits {

	var creditGroup []Credits

	for source, srcMisc := range finances.MiscAmounts {

		creditor := cf.CreditsFromMiscAmounts[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredits(srcMisc, netIncome)
		if credits.Amount == 0 {
			continue
		}

		creditGroup = append(creditGroup, credits)
	}

	return creditGroup
}

// orderCreditGroupInPlace sort the credit group according to the application
// order of this contra formula. It assumes that cf wa validated before use.
func (cf *CanadianContraFormula) orderCreditGroupInPlace(creditGroup []Credits) {

	if len(creditGroup) == 0 {
		return
	}

	priority := make(map[CreditSource]int)
	for i, src := range cf.ApplicationOrder {
		priority[src] = i
	}

	sort.SliceStable(creditGroup, func(i int, j int) bool {
		iSrc := creditGroup[i].Source
		jSrc := creditGroup[j].Source
		return priority[iSrc] < priority[jSrc]
	})

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

// checkIncSrcCreditorsInSet returns ErrUnknownCreditSource if a single creditor
// in income-source creditos is not in the given set. If a creditor is nil, it
// returns ErrNoCreditor
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

// checkDeducSrcCreditorsInSet returns ErrUnknownCreditSource if a single
// creditor in deduction-source creditos is not in the given set. If a creditor
// is nil, it returns ErrNoCreditor
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

// checkMiscSrcCreditorsInSet returns ErrUnknownCreditSource if a single
// creditor in misc-source creditos is not in the given set. If a creditor is
// nil, it returns ErrNoCreditor
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
