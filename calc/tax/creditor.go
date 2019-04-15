package tax

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given amount. The net
	// income may or may not be used by the underlying implementation
	TaxCredits(amount, netIncome float64) Credits
	// Source returns the credit source for this creditor
	Source() CreditSource
	// Clone returns a copy of this creditor
	Clone() Creditor
}

// CreditSource represents the name of a tax credit source
type CreditSource int

// credits represent an amount that reduces payable tax
type Credits struct {
	Source CreditSource // the name of the source of the credits
	// TDO: IsRefundable -> CanCarryForward
	IsRefundable bool    // if true, the amount is paid back if not used
	Amount       float64 // the amount owed to tax payer
}

// ConstCreditor returns a constant amount of tax credits
type ConstCreditor struct {
	Const Credits
}

// TaxCredits returns constant credits disregarding the given amount to extract
// credits from and the given net income
func (cc ConstCreditor) TaxCredits(_, _ float64) Credits {
	return cc.Const
}

// Source returns the credit source this creditor
func (cc ConstCreditor) Source() CreditSource {
	return cc.Const.Source
}

// Clone returns a copy of this creditor
func (cc ConstCreditor) Clone() Creditor {
	return cc
}

type creditSources []CreditSource

// makeSet convert 'cs' into a set of unique items. It also returns duplicates
func (cs creditSources) makeSetAndGetDuplicates() (map[CreditSource]struct{}, creditSources) {

	srcSet := make(map[CreditSource]struct{})
	srcDup := make(creditSources, 0, len(cs))

	for _, creditSrc := range cs {

		if _, ok := srcSet[creditSrc]; ok {
			srcDup = append(srcDup, creditSrc)
			continue
		}

		srcSet[creditSrc] = struct{}{}
	}

	return srcSet, srcDup
}
