package calc

// Payments represent a recurring payment schedule where the index of each
// payment in the slice represents a time period, i.e. month, year etc.
type Payments []float64

// Total returns the total payments in this schedule
func (p Payments) Total() float64 {

	var total float64
	for _, payment := range p {
		total += payment
	}
	return total
}

// Normalize returns a copy of these payments such that the total amount is the
// the same as 'p', but all payments in the returned copy are equal
func (p Payments) Normalize() Payments {

	if p == nil {
		return nil
	}

	avg := p.Total() / float64(len(p))
	clone := make(Payments, len(p))
	for i := range clone {
		clone[i] = avg
	}

	return clone
}

// Clone returns a copy of these payments
func (p Payments) Clone() Payments {

	if p == nil {
		return nil
	}

	clone := make(Payments, len(p))
	for i, payment := range p {
		clone[i] = payment
	}

	return clone
}

// TODO: deductions should be a map[Type]float64 to allow for exclusion of
// certain deductions under different criteria
// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	Income     float64 `json:"income"`
	Deductions float64 `json:"deductions"`
}

// NetIncome calculate the total income of +/- adjustments. The income is
// calculated as the sum of taxable amounts less the sum of deductions,
// plus/minus the given adjustments (if any)
// TODO: this might not be the responsibility of this type?
func (f IndividualFinances) NetIncome(adjustments ...float64) float64 {

	total := f.Income - f.Deductions
	for _, adj := range adjustments {
		total += adj
	}
	return total
}

// FamilyFinances is used by types implementing the ChildBenefitCalculator
// interface to recieve input needed to calculate benefits
type FamilyFinances [2]IndividualFinances

// NetIncome calculate the total income of the family +/- adjustments. The
// income is calculated as the sum of taxable amounts less the sum of
// deductions. Adjustments are added/subtracted from the total income.
func (f FamilyFinances) NetIncome(adjustments ...float64) float64 {

	total := f[0].NetIncome() + f[1].NetIncome()
	for _, adj := range adjustments {
		total += adj
	}
	return total
}

func (f FamilyFinances) NetAdjIncome() float64 {
	panic("not implemented")
}

// Split returns the individual finances that jointly represent this object
func (f FamilyFinances) Split() (IndividualFinances, IndividualFinances) {
	return f[0], f[1]
}
