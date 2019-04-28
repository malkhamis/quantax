package tax

import "github.com/malkhamis/quantax/core"

// compile-time check for interface implementation
var _ core.TaxCredit = (*TaxCredit)(nil)

// creditRuleGroup is a type used to encapsulate slice-specific logic
type creditRuleGroup []core.CreditRule

// makeSrcSetAndGetDuplicates convert 'crg' into a set of unique items and
// returns duplicates
func (crg creditRuleGroup) makeSrcSetAndGetDuplicates() (map[string]struct{}, []string) {

	srcSet := make(map[string]struct{})
	srcDup := make([]string, 0, len(crg))

	for _, crRule := range crg {

		if _, ok := srcSet[crRule.CrSource]; ok {
			srcDup = append(srcDup, crRule.CrSource)
			continue
		}

		srcSet[crRule.CrSource] = struct{}{}
	}

	return srcSet, srcDup
}

// creditBySource is used to pass around credit amount alongside its source
type creditBySource struct {
	source core.FinancialSource
	amount float64
}

// TaxCredit represents an amount owed to tax payer
type TaxCredit struct {
	// the initial amount
	AmountInitial float64
	// the remaining usable amount
	AmountRemaining float64
	// the amount used by tax calculator
	AmountUsed float64
	// the associated financial source
	FinancialSource core.FinancialSource
	// the rule of crediting the amount
	CrRule core.CreditRule
	// the financer this tax credit belongs to
	Ref core.Financer
	// description/reason for the tax credit
	Desc string
	// TODO
	RelatedTaxInfo core.TaxInfo
}

// Amounts returns the initial, used, and remaining credit amounts
func (tc *TaxCredit) Amounts() (initial, used, remaining float64) {
	if tc == nil {
		return 0, 0, 0
	}
	return tc.AmountInitial, tc.AmountUsed, tc.AmountRemaining
}

// SetAmounts sets the initial, used, and remaining credit amounts
func (tc *TaxCredit) SetAmounts(initial, used, remaining float64) {
	if tc == nil {
		return
	}
	tc.AmountInitial = initial
	tc.AmountUsed = used
	tc.AmountRemaining = remaining
}

// RemainingAmount is the remaining credit amount owed to tax payer
func (tc *TaxCredit) RemainingAmount() float64 {
	if tc == nil {
		return 0
	}
	return tc.AmountRemaining
}

// UsedAmount is the amount previously used to pay the tax payer
func (tc *TaxCredit) UsedAmount() float64 {
	if tc == nil {
		return 0
	}
	return tc.AmountUsed
}

// Source is the financial source of this tax credit. If the credit is
// not associated with a specific source, it should return 'SrcNone'
func (tc *TaxCredit) Source() core.FinancialSource {
	if tc == nil {
		return core.SrcNone
	}
	return tc.FinancialSource
}

// ReferenceFinancer returns the owner of this tax credit
func (tc *TaxCredit) ReferenceFinancer() core.Financer {
	if tc == nil {
		return nil
	}
	return tc.Ref
}

// Rule returns the rule of using the tax credit
func (tc *TaxCredit) Rule() core.CreditRule {
	if tc == nil {
		return core.CreditRule{}
	}
	return tc.CrRule
}

// Description is a short description for the reason of the tax credit
func (tc *TaxCredit) Description() string {
	if tc == nil {
		return ""
	}
	return tc.Desc
}

// TaxInfo returns the tax information this credit is associated with
func (tc *TaxCredit) TaxInfo() core.TaxInfo {
	if tc == nil {
		return core.TaxInfo{}
	}
	return tc.RelatedTaxInfo
}

// ShallowCopy returns a copy of this credit without cloning the references
func (tc *TaxCredit) ShallowCopy() core.TaxCredit {
	return tc.shallowCopy()
}

// shallowCopy returns a copy of this tax credit instance
func (cr *TaxCredit) shallowCopy() *TaxCredit {
	if cr == nil {
		return nil
	}
	return &(*cr)
}

// taxCreditGroup is a type used to encapsulate slice-specific logic
type taxCreditGroup []*TaxCredit

// typecast convert []*TaxCredit to []core.TaxCredit
func (tcg taxCreditGroup) typecast() []core.TaxCredit {

	typed := make([]core.TaxCredit, len(tcg))
	for i, cr := range tcg {
		typed[i] = cr
	}
	return typed
}

// clone returns a copy of this tax credit group
func (tcg taxCreditGroup) clone() []*TaxCredit {

	if tcg == nil {
		return nil
	}

	c := make([]*TaxCredit, len(tcg))
	for i, cr := range tcg {
		c[i] = cr.shallowCopy()
	}
	return c
}
