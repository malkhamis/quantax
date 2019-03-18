package cra

import (
	"github.com/pkg/errors"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/facts"
)

// compile-time check for interface implementation
var (
	_ calc.IncomeTaxCalculator   = TaxCalculator{}
	_ calc.IncomeTaxCalculatorCA = TaxCalculator{}
)

type FinancialNumbers struct {
	Income         float64
	DeductionsProv float64
	DeductionsFed  float64
	CreditsProv    float64
	CreditsFed     float64
}

// TaxCalculator represents a calculator that simplifies the calculation
// of Candian income tax.
type TaxCalculator struct {
	FinancialNumbers
	taxParams facts.Facts
}

func NewTaxCalculator(finNums FinancialNumbers, taxParams facts.Facts) (TaxCalculator, error) {

	c := TaxCalculator{
		FinancialNumbers: finNums,
		taxParams:        taxParams.Clone(),
	}

	return c, c.Validate()
}

func (c TaxCalculator) Validate() error {

	err := c.taxParams.Validate()
	return errors.Wrap(err, "invalid tax parameters")
}

// CalcTaxTotal computes the sum of federal and provincial tax on the income
func (c TaxCalculator) CalcTaxTotal() float64 {
	total := c.CalcTaxFederal() + c.CalcTaxProvincial()
	return total
}

// CalcTaxFederal computes the federal tax on the income
func (c TaxCalculator) CalcTaxFederal() float64 {

	totalIncome := c.Income - c.DeductionsFed
	payableTax := calcTax(totalIncome, c.taxParams.FactsFed.Rates)
	payableTax -= c.CreditsFed
	return payableTax
}

// CalcTaxProvincial computes the provincial tax on the income
func (c TaxCalculator) CalcTaxProvincial() float64 {

	totalIncome := c.Income - c.DeductionsProv
	payableTax := calcTax(totalIncome, c.taxParams.FactsProv.Rates)
	payableTax -= c.CreditsProv
	return payableTax
}

// calcTax calculate the tax given a total income and tax bracket rates
func calcTax(income float64, rates facts.BracketRates) float64 {

	var total float64

	for rate, bracket := range rates {

		if income < bracket.Lower() {
			continue
		}

		if income >= bracket.Upper() {
			total += rate * bracket.Amount()
			continue
		}

		total += rate * (income - bracket.Lower())
	}

	return total
}
