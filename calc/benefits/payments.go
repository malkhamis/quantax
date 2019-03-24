package benefits

// payments represent a recurring payment schedule where the index of each
// payment in the slice represents a time period, i.e. month, year etc.
type payments []float64

// Total returns the total payments in this schedule
func (p payments) Total() float64 {

	var total float64
	for _, payment := range p {
		total += payment
	}
	return total
}

// Clone returns a copy of these payments
func (p payments) Clone() payments {

	if p == nil {
		return nil
	}

	clone := make(payments, len(p))
	for i, payment := range p {
		clone[i] = payment
	}

	return clone
}
