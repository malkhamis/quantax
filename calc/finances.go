package calc

// TODO: deductions should be a map[Type]float64 to allow for exclusion of
// certain deductions under different criteria
// Finances is used by types implementing the TaxCalculator interface
// to recieve input needed to calculate payable taxes
type Finances struct {
	TaxableAmount float64 `json:"taxable-amount"` // amount to be taxed
	Deductions    float64 `json:"deductions"`     // subtracted from taxableAmount
	Credits       float64 `json:"credits"`        // subtracted from total tax
}

// NetIncome calculate the total income of +/- adjustments. The income is
// calculated as the sum of taxable amounts less the sum of deductions,
// plus/minus the given adjustments (if any)
func (f Finances) NetIncome(adjustments ...float64) float64 {

	total := f.TaxableAmount - f.Deductions
	for _, adj := range adjustments {
		total += adj
	}
	return total
}

// FamilyFinances is used by types implementing the ChildBenefitCalculator
// interface to recieve input needed to calculate benefits
type FamilyFinances [2]Finances

// NetIncome calculate the total income of the family +/- adjustments. The
// income is calculated as the sum of taxable amounts less the sum of
// deductions. Adjustments are added/subtracted from the total income
func (f FamilyFinances) NetIncome(adjustments ...float64) float64 {

	total := f[0].NetIncome() + f[1].NetIncome()
	for _, adj := range adjustments {
		total += adj
	}
	return total
}
