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
	Creditors []Creditor
	// ApplicationOrder stores the order in which tax credits are used
	ApplicationOrder []core.CreditRule
	// TaxInfo TODO
	RelatedTaxInfo core.TaxInfo
}

// Apply applies the contra-formula on the given on the given tax payer and
// returns tax credits
func (cf *CanadianContraFormula) Apply(taxPayer *TaxPayer) []*TaxCredit {

	if taxPayer == nil || taxPayer.Finances == nil {
		return nil
	}

	// FIXME: we can avoid this by having NewCanadaianContraFormula
	// to initialize an internal variable containing the same info
	// which can be used again and again
	ruleByCrSrc := make(map[string]core.CreditRule)
	for _, rule := range cf.ApplicationOrder {
		ruleByCrSrc[rule.CrSource] = rule
	}

	credits := make([]*TaxCredit, 0, len(cf.Creditors))
	for _, creditor := range cf.Creditors {

		amount := creditor.TaxCredit(taxPayer)
		if amount == 0 {
			continue
		}

		cr := &TaxCredit{
			AmountInitial:   amount,
			AmountRemaining: amount,
			AmountUsed:      0,
			Ref:             taxPayer.Finances,
			CrRule:          ruleByCrSrc[creditor.CrSourceName()],
			FinancialSource: creditor.FinancialSource(),
			RelatedTaxInfo:  cf.RelatedTaxInfo,
			Desc:            creditor.Description(),
		}

		credits = append(credits, cr)
	}

	return credits
}

// FilterAndSort removes tax credits whose credit source names aren't recognized
// by this contra formula and sort them according to the application order. It
// assumes that cf wa validated before use.
func (cf *CanadianContraFormula) FilterAndSort(credits []core.TaxCredit) []core.TaxCredit {

	// FIXME: we can avoid this by having NewCanadaianContraFormula
	// to initialize an internal variable containing the same info
	// which can be used again and again
	priority := make(map[core.CreditRule]int)
	for i, rule := range cf.ApplicationOrder {
		priority[rule] = i
	}

	filtered := make([]core.TaxCredit, 0, len(credits))
	for _, cr := range credits {

		if cr.TaxInfo().TaxRegion != cf.RelatedTaxInfo.TaxRegion {
			continue
		}

		_, ok := priority[cr.Rule()]
		if !ok {
			continue
		}

		filtered = append(filtered, cr)
	}

	sort.SliceStable(filtered, func(i int, j int) bool {
		iSrc := filtered[i].Rule()
		jSrc := filtered[j].Rule()
		return priority[iSrc] < priority[jSrc]
	})

	return filtered
}

// Validate checks if the formula is valid for use
func (cf *CanadianContraFormula) Validate() error {

	appOrderCreditSrcSet, dups := creditRuleGroup(
		cf.ApplicationOrder,
	).makeSrcSetAndGetDuplicates()

	if len(dups) > 0 {
		return errors.Wrapf(ErrDupCreditSource, "%v", dups)
	}

	err := cf.checkCreditorCrSrcNamesInSet(appOrderCreditSrcSet)
	return err
}

// Clone returns a copy of this contra-formula
func (cf *CanadianContraFormula) Clone() ContraFormula {

	if cf == nil {
		return nil
	}

	clone := &CanadianContraFormula{RelatedTaxInfo: cf.RelatedTaxInfo}

	if cf.Creditors != nil {
		clone.Creditors = make([]Creditor, len(cf.Creditors))
		for i, creditor := range cf.Creditors {
			clone.Creditors[i] = creditor.Clone()
		}
	}

	if cf.ApplicationOrder != nil {
		clone.ApplicationOrder = make([]core.CreditRule, len(cf.ApplicationOrder))
		copy(clone.ApplicationOrder, cf.ApplicationOrder)
	}

	return clone
}

// TaxInfo TODO
func (cf *CanadianContraFormula) TaxInfo() core.TaxInfo {
	return cf.RelatedTaxInfo
}

// checkCreditorCrSrcsInSet returns ErrUnknownCreditSource if a single creditor
// returns a credit source that is not in the given set. If a creditor is nil,
// it returns ErrNoCreditor
func (cf *CanadianContraFormula) checkCreditorCrSrcNamesInSet(crSrcNames map[string]struct{}) error {

	for _, creditor := range cf.Creditors {

		if creditor == nil {
			return ErrNoCreditor
		}

		_, exist := crSrcNames[creditor.CrSourceName()]
		if !exist {
			return errors.Wrapf(ErrUnknownCreditSource, "%q", creditor.CrSourceName())
		}

	}

	return nil
}
