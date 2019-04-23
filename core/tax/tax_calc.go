package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula          Formula
	contraFormula    ContraFormula
	incomeCalculator core.IncomeCalculator
	finances         *core.IndividualFinances
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
		finances:         core.NewEmptyIndividualFinances(),
	}

	return c, nil
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil, a non-nil, empty finances is set
func (c *Calculator) SetFinances(f *core.IndividualFinances) {

	if f == nil {
		f = core.NewEmptyIndividualFinances()
	}
	c.incomeCalculator.SetFinances(f)
	c.finances = f
}

// SetCredits stores relevent credits from the given credits in this calculator.
// Subsequent calls to other calculator functions may or may not be influenced
// by these credits.
func (c *Calculator) SetCredits(credits []core.TaxCredit) {

	c.credits = make([]*taxCredit, 0, len(credits))

	for _, cr := range credits {
		typed, ok := cr.(*taxCredit)
		if ok && typed.owner == c {
			c.credits = append(c.credits, typed)
		}
	}

}

// TaxPayable computes the tax on the net income for the previously set finances
// and any relevent credits.
func (c *Calculator) TaxPayable() (float64, []core.TaxCredit) {

	netIncome := c.incomeCalculator.NetIncome()
	totalTax := c.formula.Apply(netIncome)

	newCredits := c.contraFormula.Apply(c.finances, netIncome)
	c.ownCredits(newCredits)
	allCredits := append(c.credits, newCredits...)

	netPayableTax, remainingCredits := c.netPayableTax(totalTax, allCredits)
	return netPayableTax, taxCreditGroup(remainingCredits).typecast()
}

// ownCredits sets the owner of the given credits to this calculator
func (c *Calculator) ownCredits(credits []*taxCredit) {
	for _, cr := range credits {
		cr.owner = c
	}
}

// netPayableTax returns the payable tax and any unsable/remaining credits. It
// assumes that the given credits are owned by this calculator
func (c *Calculator) netPayableTax(taxAmount float64, credits []*taxCredit) (float64, []*taxCredit) {

	newCredits := taxCreditGroup(credits).clone()

	for _, cr := range newCredits {

		if taxAmount <= 0.0 && cr.rule.Type == CrRuleTypeCanCarryForward {
			continue
		}

		if taxAmount <= 0.0 && cr.rule.Type == CrRuleTypeNotCarryForward {
			cr.amount = 0
			continue
		}

		if taxAmount >= cr.amount || cr.rule.Type == CrRuleTypeCashable {
			taxAmount -= cr.amount
			cr.amount = 0
			continue
		}

		// reached at most once
		diff := cr.amount - taxAmount
		taxAmount = 0.0
		cr.amount = diff
		if cr.rule.Type == CrRuleTypeNotCarryForward {
			cr.amount = 0.0
		}
	}

	return taxAmount, newCredits
}
