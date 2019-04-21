package tax

import (
	"sort"

	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

var (
	// compile-time check for interface implementation
	_ ContraFormula = (*CanadianContraFormula)(nil)
)

// CanadianContraFormula is used to calculate Canadian tax credits
type CanadianContraFormula struct {
	// CreditsFromIncome stores creditors associated with income sources
	CreditsFromIncome map[core.FinancialSource]Creditor
	// CreditsFromDeduction stores creditors associated with deduction sources
	CreditsFromDeduction map[core.FinancialSource]Creditor
	// CreditsFromMiscAmounts stores creditors associated with misc sources
	CreditsFromMiscAmounts map[core.FinancialSource]Creditor
	// PersistentCredits are source-named credits that are available by default in
	// this contra-formula. They must be accounted for in ApplicationOrder.
	PersistentCredits map[string]float64
	// ApplicationOrder stores the order in which tax credits are used
	ApplicationOrder []CreditRule
}

// Apply applies the contra-formula on the given finances and net income and
// extracts tax credits from them
func (cf *CanadianContraFormula) Apply(finances *core.IndividualFinances, netIncome float64) []*taxCredit {

	if finances == nil {
		return nil
	}

	allCredits := cf.extractCredits(finances, netIncome)
	convertedCredits := cf.convertCredits(allCredits)
	cf.orderCredits(convertedCredits)

	return convertedCredits
}

// extractCredits extracts tax credits from the given finances and net income.
// The returns credits also include persistent credits
func (cf *CanadianContraFormula) extractCredits(finances *core.IndividualFinances, netIncome float64) []*creditBySource {

	pCredits := cf.persistentCredits()
	incSrcCredits := cf.creditsFromIncSrcs(finances, netIncome)
	deducSrcCredits := cf.creditsFromDeducSrcs(finances, netIncome)
	miscSrcCredits := cf.creditsFromMiscSrcs(finances, netIncome)

	allCredits := append(pCredits, incSrcCredits...)
	allCredits = append(allCredits, deducSrcCredits...)
	allCredits = append(allCredits, miscSrcCredits...)

	return allCredits
}

// persistentCredits converts cf.PersistentCredits to []*creditBySource
func (cf *CanadianContraFormula) persistentCredits() []*creditBySource {

	pCredits := make([]*creditBySource, 0, len(cf.PersistentCredits))
	for src, amount := range cf.PersistentCredits {
		pCredits = append(pCredits, &creditBySource{src, amount})
	}
	return pCredits
}

// creditsFromIncSrcs returns a list of credits extracted from the income
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromIncSrcs(finances *core.IndividualFinances, netIncome float64) []*creditBySource {

	credits := make([]*creditBySource, 0, len(finances.Income))

	for src, income := range finances.Income {

		creditor := cf.CreditsFromIncome[src]
		if creditor == nil {
			continue
		}

		amount := creditor.TaxCredit(income, netIncome)
		if amount == 0 {
			continue
		}

		cr := &creditBySource{creditor.Source(), amount}
		credits = append(credits, cr)
	}

	return credits
}

// creditsFromDeducSrcs returns a list of credits extracted from the deduciton
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromDeducSrcs(finances *core.IndividualFinances, netIncome float64) []*creditBySource {

	credits := make([]*creditBySource, 0, len(finances.Deductions))

	for src, deduc := range finances.Deductions {

		creditor := cf.CreditsFromDeduction[src]
		if creditor == nil {
			continue
		}

		amount := creditor.TaxCredit(deduc, netIncome)
		if amount == 0 {
			continue
		}

		cr := &creditBySource{creditor.Source(), amount}
		credits = append(credits, cr)
	}

	return credits
}

// creditsFromMiscSrcs returns a list of credits extracted from the misc
// sources in this contra-formula. It assumes that finances is never nil and
// that the contra formula was validated
func (cf *CanadianContraFormula) creditsFromMiscSrcs(finances *core.IndividualFinances, netIncome float64) []*creditBySource {

	credits := make([]*creditBySource, 0, len(finances.MiscAmounts))

	for src, misc := range finances.MiscAmounts {

		creditor := cf.CreditsFromMiscAmounts[src]
		if creditor == nil {
			continue
		}

		amount := creditor.TaxCredit(misc, netIncome)
		if amount == 0 {
			continue
		}

		cr := &creditBySource{creditor.Source(), amount}
		credits = append(credits, cr)
	}

	return credits
}

// orderCredits sort the credit group according to the application
// order of this contra formula. It assumes that cf wa validated before use.
func (cf *CanadianContraFormula) orderCredits(credits []*taxCredit) {

	if len(credits) == 0 {
		return
	}

	priority := make(map[string]int)
	for i, rule := range cf.ApplicationOrder {
		priority[rule.Source] = i
	}

	sort.SliceStable(credits, func(i int, j int) bool {
		iSrc := credits[i].Source()
		jSrc := credits[j].Source()
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
		clone.CreditsFromDeduction = make(map[core.FinancialSource]Creditor)
		for source, creditor := range cf.CreditsFromDeduction {
			clone.CreditsFromDeduction[source] = creditor.Clone()
		}
	}

	if cf.CreditsFromIncome != nil {
		clone.CreditsFromIncome = make(map[core.FinancialSource]Creditor)
		for source, creditor := range cf.CreditsFromIncome {
			clone.CreditsFromIncome[source] = creditor.Clone()
		}
	}

	if cf.CreditsFromMiscAmounts != nil {
		clone.CreditsFromMiscAmounts = make(map[core.FinancialSource]Creditor)
		for source, creditor := range cf.CreditsFromMiscAmounts {
			clone.CreditsFromMiscAmounts[source] = creditor.Clone()
		}
	}

	if cf.PersistentCredits != nil {
		clone.PersistentCredits = make(map[string]float64)
		for src, amount := range cf.PersistentCredits {
			clone.PersistentCredits[src] = amount
		}
	}

	if cf.ApplicationOrder != nil {
		clone.ApplicationOrder = make([]CreditRule, len(cf.ApplicationOrder))
		copy(clone.ApplicationOrder, cf.ApplicationOrder)
	}

	return clone
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {

	appOrderCreditSrcSet, dups := creditRuleGroup(
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
func (cf *CanadianContraFormula) checkIncSrcCreditorsInSet(set map[string]struct{}) error {

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
func (cf *CanadianContraFormula) checkDeducSrcCreditorsInSet(set map[string]struct{}) error {

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
func (cf *CanadianContraFormula) checkMiscSrcCreditorsInSet(set map[string]struct{}) error {

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
func (cf *CanadianContraFormula) checkPersistentCrSrcsInSet(set map[string]struct{}) error {

	for src := range cf.PersistentCredits {

		_, exist := set[src]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"persistent credit source: %q",
				src,
			)
		}

	}

	return nil
}

// convertCredits convert amount-source pairs to tax credtis to include rules.
// Sources that are not present in the application order are omitted
func (cf *CanadianContraFormula) convertCredits(credits []*creditBySource) []*taxCredit {

	knownRules := make(map[string]CreditRule)
	for _, rule := range cf.ApplicationOrder {
		knownRules[rule.Source] = rule
	}

	convertedCredits := make([]*taxCredit, 0, len(credits))
	for _, cr := range credits {

		rule, ok := knownRules[cr.source]
		if !ok {
			continue
		}

		newCr := &taxCredit{
			amount: cr.amount,
			rule:   rule,
		}
		convertedCredits = append(convertedCredits, newCr)
	}

	return convertedCredits
}
