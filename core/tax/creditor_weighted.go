package tax

// WeightedCreditor is a Creditor that returns a weighted amount for the target
// financial source of the tax payer's finances
type WeightedCreditor struct {
	// the weight to apply on the target financial source
	Weight float64
	CreditDescriptor
}

// TaxCredit returns the tax credit amount multiplied by the weight set in this
// creditor. If tax payer is nil, it returns zero. This creditor does not refer
// to the tax payer's spouse
func (wc WeightedCreditor) TaxCredit(tp *TaxPayer) float64 {

	if tp == nil {
		return 0.0
	}

	amount := tp.Finances.TotalAmount(wc.TargetFinancialSource)
	return wc.Weight * amount
}

// Clone returns a deep copy of this creditor
func (wc WeightedCreditor) Clone() Creditor {
	return wc.clone()
}

// clone returns a copy of this creditor
func (wc WeightedCreditor) clone() WeightedCreditor {
	return wc
}
