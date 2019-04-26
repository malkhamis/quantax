package tax

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

// compile-time check for interface implementation
var _ core.TaxCalculator = (*Calculator)(nil)

// Calculator is used to calculate payable tax for individuals
type Calculator struct {
	formula          Formula
	contraFormula    ContraFormula
	incomeCalculator core.IncomeCalculator
	finances         core.HouseholdFinances
	credits          []*taxCredit
	dependents       []*human.Person
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
		finances:         core.NewHouseholdFinancesNop(),
	}

	return c, nil
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil or both spouses' finances are nil, a noop
// instance is set
func (c *Calculator) SetFinances(f core.HouseholdFinances) {
	if f == nil {
		f = core.NewHouseholdFinancesNop()
	}
	if f.SpouseA() == nil && f.SpouseB() == nil {
		f = core.NewHouseholdFinancesNop()
	}
	c.finances = f
}

// SetCredits stores relevent credits from the given credits in this calculator.
// Subsequent calls to other calculator functions may or may not be influenced
// by these credits.
func (c *Calculator) SetCredits(credits []core.TaxCredit) {

	c.credits = make([]*taxCredit, 0, len(credits))

	for _, cr := range credits {

		if cr.Amount() == 0 {
			continue
		}

		typed, ok := cr.(*taxCredit)
		if ok && typed.owner != c {
			continue
		}

		c.credits = append(c.credits, typed)
	}

}

// SetDependents sets the dependents which the calculator might use for tax-
// related calculations
func (c *Calculator) SetDependents(dependents ...*human.Person) {

	c.dependents = make([]*human.Person, 0, len(dependents))
	for _, d := range dependents {
		if d == nil {
			continue
		}
		c.dependents = append(c.dependents, d)
	}
}

// TaxPayable computes the tax on the net income for the previously set finances
// and any relevent credits.
func (c *Calculator) TaxPayable() (spouseA, spouseB float64, unusedCredits []core.TaxCredit) {

	netIncomeA, netIncomeB := c.netIncome()
	totalTaxA, totalTaxB := c.totalTax(netIncomeA, netIncomeB)
	newCredits := c.totalCombinedCredits(netIncomeA, netIncomeB)

	netPayableTaxA, unusedCrA := c.netPayableTax(totalTaxA, newCredits)
	netPayableTaxB, unusedCrB := c.netPayableTax(totalTaxB, newCredits)
	unusedCr := taxCreditGroup(append(unusedCrA, unusedCrB...)).typecast()

	return netPayableTaxA, netPayableTaxB, unusedCr
}

// netIncome returns the net income for both spouses in the set finances
func (c *Calculator) netIncome() (spouseA, spouseB float64) {

	c.incomeCalculator.SetFinances(c.finances.SpouseA())
	spouseA = c.incomeCalculator.NetIncome()

	c.incomeCalculator.SetFinances(c.finances.SpouseB())
	spouseB = c.incomeCalculator.NetIncome()

	return spouseA, spouseB
}

// totalTax returns the total tax amount for both spouses in the set finances
func (c *Calculator) totalTax(netIncomeA, netIncomeB float64) (totalTaxA, totalTaxB float64) {
	totalTaxA = c.formula.Apply(netIncomeA)
	totalTaxB = c.formula.Apply(netIncomeB)
	return totalTaxA, totalTaxB
}

// totalCombinedCredits returns tax credits for the set household finances. Each
// individual tax credit instance is associated with its own financer
func (c *Calculator) totalCombinedCredits(netIncomeA, netIncomeB float64) []*taxCredit {

	taxPayerA, taxPayerB := c.makeTaxPayers(netIncomeA, netIncomeB)
	creditsA := c.taxCredits(taxPayerA)
	creditsB := c.taxCredits(taxPayerB)

	allCredits := append(c.credits, creditsA...)
	allCredits = append(allCredits, creditsB...)
	return allCredits
}

// taxCredits returns the tax credits for the given tax payer. individual tax
// credits are associated with this calculator and tax payer's finanaces. If
// a tax credit is nil or its amount is zero, it is skipped. If the given tax
// payer is nil, it returns nil
func (c *Calculator) taxCredits(taxPayer *TaxPayer) []*taxCredit {

	if taxPayer == nil {
		return nil
	}

	allCredits := c.contraFormula.Apply(taxPayer)
	usableCredits := make([]*taxCredit, 0, len(allCredits))
	for _, cr := range allCredits {

		if cr == nil {
			continue
		}
		if cr.amount == 0 {
			continue
		}

		cr.owner = c
		cr.ref = taxPayer.Finances

		usableCredits = append(usableCredits, cr)
	}

	return usableCredits
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

// makeTaxPayers returns dual tax payers from the given net income amounts and
// the finances stored in this calculator
func (c *Calculator) makeTaxPayers(netIncomeA, netIncomeB float64) (A, B *TaxPayer) {

	var (
		taxPayerA, taxPayerB *TaxPayer
		financesA            = c.finances.SpouseA()
		financesB            = c.finances.SpouseB()
	)

	if financesA != nil {
		taxPayerA = &TaxPayer{
			Finances:        financesA,
			NetIncome:       netIncomeA,
			Dependents:      c.dependents,
			HasSpouse:       financesB != nil,
			SpouseNetIncome: netIncomeB,
		}
	}

	if financesB != nil {
		taxPayerB = &TaxPayer{
			Finances:        financesB,
			NetIncome:       netIncomeB,
			Dependents:      c.dependents,
			HasSpouse:       financesA != nil,
			SpouseNetIncome: netIncomeA,
		}
	}

	return taxPayerA, taxPayerB
}
