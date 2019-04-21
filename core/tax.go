package core

import "github.com/malkhamis/quantax/core/finance"

// TaxCalculator is used to calculate payable tax on earnings
type TaxCalculator interface {
	// TaxPayable returns the payable amount of tax for the set finances.
	// The tax credit represent any amount owed to the tax payer without
	// implications for how they might be used.
	TaxPayable() (float64, []TaxCredit)
	// SetFinances stores the given financial data in the underlying tax
	// calculator. Subsequent calls to other functions are based on the
	// the given finances. Changes to the given finances after calling
	// this function should affect future calculations
	SetFinances(*finance.IndividualFinances)
	// SetCredits stores the given credits in the underlying tax calculator.
	// Subsequent calls to other functions will be influenced by the given tax
	// credits. Treatment of given credits is implementation-specific. Ideally,
	// These given credits are originated by the same tax calculator.
	SetCredits([]TaxCredit)
}

// TaxCredit represents an amount that is owed to the tax payer
type TaxCredit interface {
	// Amount is the amount owed to tax payer
	Amount() float64
	// Source describes the source of this tax credit
	Source() string
}
