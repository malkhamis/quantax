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
	// PersistentCredits are credits that are available by default in this
	// contra-formula. They still need to be accounted for in ApplicationOrder
	PersistentCredits []finance.TaxCredit
	// ApplicationOrder stores the order in which tax credits are used
	ApplicationOrder []CreditSourceControl
}

// Apply applies the contra-formula on the income and the set finances
func (cf *CanadianContraFormula) Apply(finances *finance.IndividualFinances, netIncome float64) []TaxCredit {

	if finances == nil {
		return nil
	}

	incSrcCredits := cf.creditsFromIncSrcs(finances, netIncome)
	deducSrcCredits := cf.creditsFromDeducSrcs(finances, netIncome)
	miscSrcCredits := cf.creditsFromMiscSrcs(finances, netIncome)

	allCredits := append(cf.PersistentCredits, incSrcCredits...)
	allCredits = append(allCredits, deducSrcCredits...)
	allCredits = append(allCredits, miscSrcCredits...)

	allCredits = finance.TaxCreditGroup(allCredits).MergeSimilars()
	controlledCredits := cf.makeControlledTaxCreditsFrom(allCredits)
	cf.orderCreditGroupInPlace(controlledCredits)

	return controlledCredits
}

// creditsFromIncSrcs returns a list of credits extracted from the income
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromIncSrcs(finances *finance.IndividualFinances, netIncome float64) []finance.TaxCredit {

	var creditGroup []finance.TaxCredit

	for source, srcIncome := range finances.Income {

		creditor := cf.CreditsFromIncome[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredit(srcIncome, netIncome)
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
func (cf *CanadianContraFormula) creditsFromDeducSrcs(finances *finance.IndividualFinances, netIncome float64) []finance.TaxCredit {

	var creditGroup []finance.TaxCredit

	for source, srcDeduc := range finances.Deductions {

		creditor := cf.CreditsFromDeduction[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredit(srcDeduc, netIncome)
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
func (cf *CanadianContraFormula) creditsFromMiscSrcs(finances *finance.IndividualFinances, netIncome float64) []finance.TaxCredit {

	var creditGroup []finance.TaxCredit

	for source, srcMisc := range finances.MiscAmounts {

		creditor := cf.CreditsFromMiscAmounts[source]
		if creditor == nil {
			continue
		}

		credits := creditor.TaxCredit(srcMisc, netIncome)
		if credits.Amount == 0 {
			continue
		}

		creditGroup = append(creditGroup, credits)
	}

	return creditGroup
}

// orderCreditGroupInPlace sort the credit group according to the application
// order of this contra formula. It assumes that cf wa validated before use.
func (cf *CanadianContraFormula) orderCreditGroupInPlace(creditGroup []TaxCredit) {

	if len(creditGroup) == 0 {
		return
	}

	priority := make(map[finance.CreditSource]int)
	for i, crSrcCtrl := range cf.ApplicationOrder {
		priority[crSrcCtrl.Source] = i
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

	if cf.PersistentCredits != nil {
		clone.PersistentCredits = make([]finance.TaxCredit, len(cf.PersistentCredits))
		copy(clone.PersistentCredits, cf.PersistentCredits)
	}

	if cf.ApplicationOrder != nil {
		clone.ApplicationOrder = make([]CreditSourceControl, len(cf.ApplicationOrder))
		copy(clone.ApplicationOrder, cf.ApplicationOrder)
	}

	return clone
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {

	appOrderCreditSrcSet, dups := creditSourceControlGroup(
		cf.ApplicationOrder,
	).makeSrcSetAndGetDuplicates()

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

	err = cf.checkPersistentCrSrcsInSet(appOrderCreditSrcSet)
	if errors.Cause(err) == ErrUnknownCreditSource {
		err = errors.Wrap(
			err,
			"persistent credit sources must be known in the application order list",
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
func (cf *CanadianContraFormula) checkIncSrcCreditorsInSet(set map[finance.CreditSource]struct{}) error {

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
func (cf *CanadianContraFormula) checkDeducSrcCreditorsInSet(set map[finance.CreditSource]struct{}) error {

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
func (cf *CanadianContraFormula) checkMiscSrcCreditorsInSet(set map[finance.CreditSource]struct{}) error {

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

// checkPersistentCrSrcsInSet returns ErrUnknownCreditSource if a source of
// credits in cf.PersistentCredits is not in the given set
func (cf *CanadianContraFormula) checkPersistentCrSrcsInSet(set map[finance.CreditSource]struct{}) error {

	for _, cr := range cf.PersistentCredits {

		_, exist := set[cr.Source]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"persistent credit source: %q",
				cr.Source,
			)
		}

	}

	return nil
}

// makeControlledTaxCreditsFrom convert tax credits to controlled tax credtis.
// Sources that are not present in the application order aren't included
func (cf *CanadianContraFormula) makeControlledTaxCreditsFrom(credits []finance.TaxCredit) []TaxCredit {

	knownCrSrcCtrl := make(map[finance.CreditSource]CreditSourceControl)
	for _, crSrcCtrl := range cf.ApplicationOrder {
		knownCrSrcCtrl[crSrcCtrl.Source] = crSrcCtrl
	}

	var controlledCredits []TaxCredit
	for _, cr := range credits {

		crSrcCtrl, ok := knownCrSrcCtrl[cr.Source]

		if !ok {
			continue
		}

		ctrlCr := TaxCredit{Amount: cr.Amount, CreditSourceControl: crSrcCtrl}
		controlledCredits = append(controlledCredits, ctrlCr)
	}

	return controlledCredits
}
