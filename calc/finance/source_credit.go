package finance

// CreditSource represents the name of a tax credit source
type CreditSource int

// TaxCredit represent an amount that reduces payable tax
type TaxCredit struct {
	// the name of the source of the tax credits
	Source CreditSource
	// the amount owed to tax payer
	Amount float64
}

// TaxCreditGroup is convenience type used to manipulate tax credit group
type TaxCreditGroup []TaxCredit

// MergeSimilars aggregates the amount of credits with similar credit sources
func (crs TaxCreditGroup) MergeSimilars() []TaxCredit {

	set := make(map[CreditSource]TaxCredit)
	for _, cr := range crs {

		aggCr, exist := set[cr.Source]
		if !exist {
			set[cr.Source] = cr
			continue
		}

		set[cr.Source] = TaxCredit{
			Source: cr.Source,
			Amount: aggCr.Amount + cr.Amount,
		}
	}

	mergedCrs := make([]TaxCredit, 0, len(set))
	for _, cr := range set {
		mergedCrs = append(mergedCrs, cr)
	}
	return mergedCrs
}
