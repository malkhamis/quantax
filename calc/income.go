package calc

import "github.com/malkhamis/quantax/calc/finance"

// IncomeCalculator is used to calculate income used for benefits and tax
type IncomeCalculator interface {
	// TotalIncome returns the total adjusted income
	TotalIncome(finance.IncomeDeductor) float64
	// TotalDeductions returns total adjusted deductions
	TotalDeductions(finance.IncomeDeductor) float64
	// NetIncome returns total income less total deduction
	NetIncome(finance.IncomeDeductor) float64
}
