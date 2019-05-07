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
	// OrderedCreditors are the objects that extract tax credits. The order of
	// creditors is used for sorting tax credits according to their priority of
	// use, where the first creditor has the highest priority before the next one
	OrderedCreditors []Creditor
	// TaxYear is the tax year this contra formula is associated with
	TaxYear uint
	// TaxRegion is the tax region this contra formula is associated with
	TaxRegion core.Region
}

// Apply applies the contra-formula for the given tax payer and returns the tax
// credits that are sorted according to priority of use, where the first tax
// credit has a higher priority of use than the next one
func (cf *CanadianContraFormula) Apply(taxPayer *TaxPayer) []*TaxCredit {

	if taxPayer == nil || taxPayer.Finances == nil {
		return nil
	}

	credits := make([]*TaxCredit, 0, len(cf.OrderedCreditors))
	for _, creditor := range cf.OrderedCreditors {

		amount := creditor.TaxCredit(taxPayer)
		if amount == 0 {
			continue
		}

		cr := &TaxCredit{
			AmountInitial:   amount,
			AmountRemaining: amount,
			AmountUsed:      0,
			Ref:             taxPayer.Finances,
			CrRule:          creditor.Rule(),
			FinancialSource: creditor.FinancialSource(),
			TaxYear:         cf.TaxYear,
			TaxRegion:       cf.TaxRegion,
			Desc:            creditor.Description(),
		}

		credits = append(credits, cr)
	}

	return credits
}

// FilterAndSort removes tax credits whose credit source names aren't recognized
// by this contra formula and sort them according to their priority of use. It
// assumes that cf was validated before use.
func (cf *CanadianContraFormula) FilterAndSort(credits *[]core.TaxCredit) {

	if credits == nil {
		return
	}

	// FIXME: we can avoid this by having NewCanadaianContraFormula
	// to initialize an internal variable containing the same info
	// which can be used again and again
	priority := make(map[core.CreditRule]int)
	for i, creditor := range cf.OrderedCreditors {
		priority[creditor.Rule()] = i
	}

	filtered := make([]core.TaxCredit, 0, len(*credits))
	for _, cr := range *credits {

		if cr == nil {
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

	*credits = filtered
}

// Validate checks if the formula is valid for use. It ensures that there are no
// duplicate in creditors' source names
func (cf *CanadianContraFormula) Validate() error {

	seen := make(map[string]bool)
	for _, creditor := range cf.OrderedCreditors {

		if creditor == nil {
			return ErrNoCreditor
		}

		src := creditor.Rule().CrSource
		if seen[src] {
			return errors.Wrap(ErrDupCreditSource, src)
		}

		seen[src] = true
	}

	return nil
}

// Clone returns a copy of this contra-formula
func (cf *CanadianContraFormula) Clone() ContraFormula {

	if cf == nil {
		return nil
	}

	clone := &CanadianContraFormula{
		TaxYear:   cf.TaxYear,
		TaxRegion: cf.TaxRegion,
	}

	if cf.OrderedCreditors != nil {
		clone.OrderedCreditors = make([]Creditor, len(cf.OrderedCreditors))
		for i, creditor := range cf.OrderedCreditors {
			clone.OrderedCreditors[i] = creditor.Clone()
		}
	}

	return clone
}

// Year returns the tax year for this formula
func (cf *CanadianContraFormula) Year() uint {
	return cf.TaxYear
}

// TaxYear returns the tax region for this formula
func (cf *CanadianContraFormula) Region() core.Region {
	return cf.TaxRegion
}
