package tax

import "github.com/malkhamis/quantax/calc/finance"

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given amount. The net
	// income may or may not be used by the underlying implementation
	TaxCredit(amount, netIncome float64) finance.TaxCredit
	// Source returns the credit source for this creditor
	Source() finance.CreditSource
	// Clone returns a copy of this creditor
	Clone() Creditor
}

// ControlType is an enum type for credit source control
type ControlType int

const (
	// unknown/uninitialized
	_ ControlType = iota
	// ControlTypeCashable indicates cashable credits
	// that may result in negative payable tax amount
	// when used to reduce payable tax
	ControlTypeCashable
	// ControlTypeCanCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero
	// and the remaining balance may be carried forward
	// to the future
	ControlTypeCanCarryForward
	// ControlTypeNotCarryForward indicates non-cashable
	// credits which may only reduce payable tax to zero.
	// Unused balance cannot be carried forward to the future
	ControlTypeNotCarryForward
)

// CreditSourceControl is used to constrain/guide a contra-formula user when
// adding/using the credit amount for a given source
type CreditSourceControl struct {
	// the controlled credit source
	Source  finance.CreditSource
	Control ControlType
}

// TaxCredit represents an amount owed to tax payer alongside the means of
// of recieving the amount owed (i.e. control)
type TaxCredit struct {
	Amount float64
	CreditSourceControl
}

// ConstCreditor returns a constant amount of tax credits
type ConstCreditor struct {
	Const finance.TaxCredit
}

// TaxCredits returns constant credits disregarding the given amount to extract
// credits from and the given net income
func (cc ConstCreditor) TaxCredit(_, _ float64) finance.TaxCredit {
	return cc.Const
}

// Source returns the credit source this creditor
func (cc ConstCreditor) Source() finance.CreditSource {
	return cc.Const.Source
}

// Clone returns a copy of this creditor
func (cc ConstCreditor) Clone() Creditor {
	return cc
}

type creditSourceControlGroup []CreditSourceControl

// makeSet convert 'cs' into a set of unique items. It also returns duplicates
func (crcg creditSourceControlGroup) makeSrcSetAndGetDuplicates() (map[finance.CreditSource]struct{}, []finance.CreditSource) {

	srcSet := make(map[finance.CreditSource]struct{})
	srcDup := make([]finance.CreditSource, 0, len(crcg))

	for _, crc := range crcg {

		if _, ok := srcSet[crc.Source]; ok {
			srcDup = append(srcDup, crc.Source)
			continue
		}

		srcSet[crc.Source] = struct{}{}
	}

	return srcSet, srcDup
}
