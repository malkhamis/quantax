package tax

// ConstCreditor is a Creditor that returns a constant amount
type ConstCreditor struct {
	// the amount to return when TaxCredit is called
	Amount float64
	CreditDescriptor
}

// TaxCredit returns the Amount set in this creditor
func (cc ConstCreditor) TaxCredit(_ *TaxPayer) float64 {
	return cc.Amount
}

// Clone returns a deep copy of this creditor
func (cc ConstCreditor) Clone() Creditor {
	return cc.clone()
}

// clone returns a copy of this creditor
func (cc ConstCreditor) clone() ConstCreditor {
	return cc
}
