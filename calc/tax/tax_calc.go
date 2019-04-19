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
	finances         *finance.IndividualFinances
	credits          []*taxCredit
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
		finances:         finance.NewEmptyIndividualFinances(0),
	}

	return c, nil
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil, a non-nil, empty finances is set
func (c *Calculator) SetFinances(f *finance.IndividualFinances) {

	if f == nil {
		c.finances = finance.NewEmptyIndividualFinances(0)
	} else {
		c.finances = f
	}
}

// SetCredits stores the given credits in this calculator. Subsequent calls to
// other calculator functions may or may not be be influenced by these credits.
func (c *Calculator) SetCredits(credits []calc.TaxCredit) {

	c.credits = make([]*taxCredit, 0, len(credits))
	for _, cr := range credits {
		typed, ok := cr.(*taxCredit)
		if ok {
			c.credits = append(c.credits, typed)
		}
	}
}

// TaxPayable computes the tax on the net income for the previously set finances
func (c *Calculator) TaxPayable() (float64, []calc.TaxCredit) {

	netIncome := c.incomeCalculator.NetIncome(c.finances)
	totalTax := c.formula.Apply(netIncome)
	newCredits := c.contraFormula.Apply(c.finances, netIncome)
	allCredits := append(c.credits, newCredits...) // TODO merge similar based on source and owner?
	netPayableTax, remainingCredits := c.netPayableTax(totalTax, allCredits)

	return netPayableTax, taxCreditGroup(remainingCredits).typecast()
}

// TODO: use if statements and panic on unidentified control
// netPayableTax returns the payable tax and any unsable/remaining credits
func (c *Calculator) netPayableTax(taxAmount float64, credits []*taxCredit) (float64, []*taxCredit) {

	var remainingCredits []*taxCredit

	for _, cr := range credits {

		switch {

		case taxAmount >= cr.amount, cr.rule.Type == CrRuleTypeCashable:
			taxAmount -= cr.amount
			newCr := cr.clone()
			newCr.amount = 0
			remainingCredits = append(remainingCredits, newCr)

		case taxAmount <= 0.0 && cr.rule.Type == CrRuleTypeCanCarryForward:
			remainingCredits = append(remainingCredits, cr.clone())

		case taxAmount <= 0.0 && cr.rule.Type == CrRuleTypeNotCarryForward:
			newCr := cr.clone()
			newCr.amount = 0
			remainingCredits = append(remainingCredits, newCr)

		default:
			diff := cr.amount - taxAmount
			taxAmount = 0.0
			newCr := cr.clone()
			newCr.amount = diff
			if cr.rule.Type == CrRuleTypeNotCarryForward {
				newCr.amount = 0.0
			}
			remainingCredits = append(remainingCredits, newCr)
		}
	}

	return taxAmount, remainingCredits
}
