package calc

import "github.com/malkhamis/quantax/calc/finance"

// TaxCalculator is used to calculate payable tax on individual earnings
type TaxCalculator interface {
	// Calc returns the payable amount of tax for the given finances
	Calc(*finance.IndividualFinances) float64
}
