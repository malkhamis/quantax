package tax

type CanadianSpouseCreditor struct {
	BaseAmount float64
	Weight     float64
	CreditDescriptor
}

// TaxCredit returns the tax credit amount for the given taxpayer's spouse. It
// assumes that the spouse is a Canadian resident for tax porpuses. If tax payer
// is nil it returns zero. If the taxpayer's spouse is nil, it assume there is
// no spouse, and hence, it returns zero
func (csc CanadianSpouseCreditor) TaxCredit(tp *TaxPayer) float64 {

	if tp == nil {
		return 0.0
	}

	if tp.SpouseFinances == nil {
		return 0.0
	}

	if tp.NetIncome < tp.SpouseNetIncome {
		return 0.0 // spouse should claim this credit
	}

	amount := csc.BaseAmount - tp.SpouseNetIncome
	if amount <= 0.0 {
		return 0.0
	}

	return csc.Weight * amount
}

// Clone returns a deep copy of this creditor
func (csc CanadianSpouseCreditor) Clone() Creditor {
	return csc.clone()
}

// clone returns a copy of this creditor
func (csc CanadianSpouseCreditor) clone() CanadianSpouseCreditor {
	return csc
}
