// Package calc defines interfaces for various tax-related calculators
package calc

// Finances is used by types implementing the TaxCalculator interface
// to recieve input needed to calculate payable taxes
type Finances struct {
	TaxableAmount float64 `json:"taxable-amount"` // amount to be taxed
	Deductions    float64 `json:"deductions"`     // subtracted from taxableAmount
	Credits       float64 `json:"credits"`        // subtracted from total tax
}

// TaxCalculator is used to calculate payable tax. Constructors for types
// implementing this interface should typically accept 'Finances' as argument
type TaxCalculator interface {
	// Calc returns the payable amount of tax set/initialized in this calculator
	Calc() float64
	// Update sets the financial numbers which the tax will be calculated for in
	// subsequent calls to Calc(). Users may call this method to set the financial
	// numbers to anything other than what the calculator was initialized with
	Update(Finances)
}
