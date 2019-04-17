package tax

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ calc.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula          Formula
	contraFormula    ContraFormula
	incomeCalculator calc.IncomeCalculator
}

// NewCalculator returns a new tax calculator for the given tax formula and the
// income calculator
func NewCalculator(cfg CalcConfig) (*Calculator, error) {

	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := &Calculator{
		formula:          cfg.TaxFormula.Clone(),
		contraFormula:    cfg.ContraTaxFormula.Clone(),
		incomeCalculator: cfg.IncomeCalc,
	}

	return c, nil
}

// Calc computes the tax on the net income for the given finances
func (c *Calculator) Calc(finances *finance.IndividualFinances) float64 {

	if finances == nil {
		return 0.0
	}

	netIncome := c.incomeCalculator.NetIncome(finances)
	totalTax := c.formula.Apply(netIncome)
	credits := c.contraFormula.Apply(finances, netIncome)
	netPayableTax, _ := c.netPayableTax(totalTax, credits)

	return netPayableTax
}

// TODO: use if statements and panic on unidentified control
// netPayableTax returns the payable tax and any unsable/remaining credits
func (c *Calculator) netPayableTax(taxAmount float64, crGroup []TaxCredit) (float64, []finance.TaxCredit) {

	var remainingCredits []finance.TaxCredit

	for _, cr := range crGroup {

		switch {

		case taxAmount >= cr.Amount, cr.Control == ControlTypeCashable:
			taxAmount -= cr.Amount
			newCrBalance := finance.TaxCredit{Amount: 0, Source: cr.Source}
			remainingCredits = append(remainingCredits, newCrBalance)

		case taxAmount <= 0.0 && cr.Control == ControlTypeCanCarryForward:
			newCrBalance := finance.TaxCredit{Amount: cr.Amount, Source: cr.Source}
			remainingCredits = append(remainingCredits, newCrBalance)

		case taxAmount <= 0.0 && cr.Control == ControlTypeNotCarryForward:
			newCrBalance := finance.TaxCredit{Amount: 0, Source: cr.Source}
			remainingCredits = append(remainingCredits, newCrBalance)

		default:
			diff := cr.Amount - taxAmount
			taxAmount = 0.0
			newCrBalance := finance.TaxCredit{Amount: diff, Source: cr.Source}
			if cr.Control == ControlTypeNotCarryForward {
				newCrBalance.Amount = 0.0
			}
			remainingCredits = append(remainingCredits, newCrBalance)
		}
	}

	return taxAmount, remainingCredits
}
