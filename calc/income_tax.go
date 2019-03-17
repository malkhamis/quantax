package calc

import (
	"github.com/pkg/errors"

	"github.com/malkhamis/tax/facts"
)

// IncomeTaxCalculator represents a calculator that simplifies the calculation
// of Candian income tax
type IncomeTaxCalculator struct {
	Income         float64
	DeductionsProv float64
	DeductionsFed  float64
	CreditsProv    float64
	CreditsFed     float64
	facts.Facts
}

// CalcTotal computes the sum of the federal and provincial tax on the
// income set in this calculator
func (c *IncomeTaxCalculator) CalcTotal() (float64, error) {

	taxFed, err := c.CalcFed()
	if err != nil {
		return taxFed, errors.Wrap(err, "error calculating federal tax")
	}

	taxProv, err := c.CalcProv()
	if err != nil {
		return taxProv, errors.Wrap(err, "error calculating provincial tax")
	}

	total := taxFed + taxProv
	return total, nil
}

// CalcFed computes the federal tax on the income set in this calculator
func (c *IncomeTaxCalculator) CalcFed() (float64, error) {

	totalIncome := c.Income - c.DeductionsFed
	payableTax, err := calcTax(totalIncome, c.FactsFed.Rates)
	payableTax -= c.CreditsFed
	return payableTax, err
}

// CalcProv computes the provincial tax on the income set in this calculator
func (c *IncomeTaxCalculator) CalcProv() (float64, error) {

	totalIncome := c.Income - c.DeductionsProv
	payableTax, err := calcTax(totalIncome, c.FactsProv.Rates)
	payableTax -= c.CreditsProv
	return payableTax, err
}

// calcTax calculate the tax given a total income and tax bracket rates
func calcTax(income float64, rates facts.BracketRates) (float64, error) {

	err := rates.Validate()
	if err != nil {
		return 0, errors.Wrap(err, "invalid bracket rates")
	}

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

	return total, nil
}
