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
	// Creditors stores creditors associated with financial sources
	Creditors map[core.FinancialSource]Creditor
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

// Clone returns a copy of this contra-formula
func (cf *CanadianContraFormula) Clone() ContraFormula {

	if cf == nil {
		return nil
	}

	clone := new(CanadianContraFormula)

	if cf.Creditors != nil {
		clone.Creditors = make(map[core.FinancialSource]Creditor)
		for source, creditor := range cf.Creditors {
			clone.Creditors[source] = creditor.Clone()
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

	err := cf.checkFinSrcCreditorsInSet(appOrderCreditSrcSet)
	if errors.Cause(err) == ErrUnknownCreditSource {
		err = errors.Wrap(
			err,
			"financial source creditors must return credit sources which are known "+
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

	incSrcs := finances.IncomeSources()
	credits := make([]*creditBySource, 0, len(incSrcs))

	for src := range incSrcs {

		creditor := cf.Creditors[src]
		if creditor == nil {
			continue
		}

		income := finances.TotalIncome(src)
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

	deducSrcs := finances.DeductionSources()
	credits := make([]*creditBySource, 0, len(deducSrcs))

	for src := range deducSrcs {

		creditor := cf.Creditors[src]
		if creditor == nil {
			continue
		}

		deduc := finances.TotalDeductions(src)
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

	miscSrcs := finances.MiscSources()
	credits := make([]*creditBySource, 0, len(miscSrcs))

	for src := range miscSrcs {

		creditor := cf.Creditors[src]
		if creditor == nil {
			continue
		}

		miscAmnt := finances.MiscAmount(src)
		amount := creditor.TaxCredit(miscAmnt, netIncome)
		if amount == 0 {
			continue
		}

		cr := &creditBySource{creditor.Source(), amount}
		credits = append(credits, cr)
	}

	return credits
}

// checkFinSrcCreditorsInSet returns ErrUnknownCreditSource if a single creditor
// in income-source creditos is not in the given set. If a creditor is nil, it
// returns ErrNoCreditor
func (cf *CanadianContraFormula) checkFinSrcCreditorsInSet(set map[string]struct{}) error {

	for financialSrc, creditor := range cf.Creditors {

		if creditor == nil {
			return errors.Wrapf(ErrNoCreditor, "financial source %q", financialSrc)
		}

		_, exist := set[creditor.Source()]
		if !exist {
			return errors.Wrapf(
				ErrUnknownCreditSource,
				"financial source %q -> credit source %q",
				financialSrc, creditor.Source(),
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
