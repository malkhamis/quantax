// Package calc defines interfaces for various tax-realted calculators
package calc

// IncomeTaxCalculator is used to calculate income tax
type IncomeTaxCalculator interface {
	// CalcTaxTotal returns the sum of all due taxes (e.g. fed and provincial)
	CalcTaxTotal() float64
	// Validate checks if the calculator is valid for use
	Validate() error
}

// IncomeTaxCalculatorCA is used to calculate Canadian income tax
type IncomeTaxCalculatorCA interface {
	CalcTaxFederal() float64
	CalcTaxProvincial() float64
	IncomeTaxCalculator
}
