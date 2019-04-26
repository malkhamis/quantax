package tax

import (
	"github.com/malkhamis/quantax/core"
)

// compile-time check for interface implementation
var _ core.TaxCredit = (*taxCredit)(nil)

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given amount. The net
	// income may or may not be used by the underlying implementation
	TaxCredit(fromAmount, netIncome float64) float64
	// Source returns the name of the credit source for this creditor
	Source() string
	// Clone returns a copy of this creditor
	Clone() Creditor
}

// CreditRuleType is an enum type for credit application rules
type CreditRuleType int

const (
	// unknown/uninitialized
	_ CreditRuleType = iota
	// CrRuleTypeCashable indicates cashable credits
	// that may result in negative payable tax amount
	// when used to reduce payable tax
	CrRuleTypeCashable
	// CrRuleTypeCanCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero
	// and the remaining balance may be carried forward
	// to the future
	CrRuleTypeCanCarryForward
	// CrRuleTypeNotCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero.
	// Unused balance cannot be carried forward to the future
	CrRuleTypeNotCarryForward
)

// CreditRule is used to constrain/guide a contra-formula user when adding or
// using the credit amount for a given source
type CreditRule struct {
	// the credit source which the rule is applied on
	Source string
	// the way, or rule, of using the credit source
	Type CreditRuleType
}

// creditRuleGroup is a type used to encapsulate slice-specific logic
type creditRuleGroup []CreditRule

// makeSrcSetAndGetDuplicates convert 'crg' into a set of unique items and
// returns duplicates
func (crg creditRuleGroup) makeSrcSetAndGetDuplicates() (map[string]struct{}, []string) {

	srcSet := make(map[string]struct{})
	srcDup := make([]string, 0, len(crg))

	for _, crRule := range crg {

		if _, ok := srcSet[crRule.Source]; ok {
			srcDup = append(srcDup, crRule.Source)
			continue
		}

		srcSet[crRule.Source] = struct{}{}
	}

	return srcSet, srcDup
}

// creditBySource is used to pass around credit amount alongside its source
type creditBySource struct {
	source string
	amount float64
}

// taxCredit represents an amount owed to tax payer
type taxCredit struct {
	// the amount owed to the tax payer
	amount float64
	// the rule of crediting the amount
	rule CreditRule
	// the originator of this tax credit
	owner core.TaxCalculator
	// the financer this tax credit belonngs to
	ref core.Financer
}

// Amount returns the credit amount
func (cr *taxCredit) Amount() float64 {
	return cr.amount
}

// Source returns the source name of this tax credit
func (cr *taxCredit) Source() string {
	return cr.rule.Source
}

// Reference returns the financer instance which this tax credit belongs to
func (cr *taxCredit) Reference() core.Financer {
	return cr.ref
}

// clone returns a copy of this tax credit instance
func (cr *taxCredit) clone() *taxCredit {
	if cr == nil {
		return nil
	}
	return &taxCredit{
		amount: cr.amount,
		rule:   cr.rule,
		owner:  cr.owner,
		ref:    cr.ref,
	}
}

// taxCreditGroup is a type used to encapsulate slice-specific logic
type taxCreditGroup []*taxCredit

// typecast convert []*taxCredit to []core.TaxCredit
func (tcg taxCreditGroup) typecast() []core.TaxCredit {

	typed := make([]core.TaxCredit, len(tcg))
	for i, cr := range tcg {
		typed[i] = cr
	}
	return typed
}

// clone returns a copy of this tax credit group
func (tcg taxCreditGroup) clone() []*taxCredit {

	if tcg == nil {
		return nil
	}

	c := make([]*taxCredit, len(tcg))
	for i, cr := range tcg {
		c[i] = cr.clone()
	}
	return c
}
