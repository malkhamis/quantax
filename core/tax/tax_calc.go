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
	crSpouseA        []core.TaxCredit
	crSpouseB        []core.TaxCredit
	dependents       []*human.Person
	taxYear          uint
	taxRegion        core.Region
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
		taxYear:          cfg.TaxFormula.Year(),
		taxRegion:        cfg.TaxFormula.Region(),
	}

	return c, nil
}

// Year returns the tax year for which this calculator is configured
func (c *Calculator) Year() uint {
	return c.taxYear
}

// Regions returns the tax region for which this calculator is configured
func (c *Calculator) Regions() []core.Region {
	return []core.Region{c.taxRegion}
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil or both spouses' finances are nil, a noop
// instance is set. The given finances are never modified in the calculator
func (c *Calculator) SetFinances(f core.HouseholdFinances, credits []core.TaxCredit) {

	if f == nil {
		f = core.NewHouseholdFinancesNop()
	}

	c.finances = f
	c.panicIfEqNonNilSpouses()
	c.setCredits(credits)
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
func (c *Calculator) TaxPayable() (spouseA, spouseB float64, combinedCredits []core.TaxCredit) {

	c.panicIfEqNonNilSpouses()

	netIncomeA, netIncomeB := c.netIncome()
	totalTaxA, totalTaxB := c.totalTax(netIncomeA, netIncomeB)
	taxCrA, taxCrB := c.totalCredits(netIncomeA, netIncomeB)

	netPayableTaxA := c.netPayableTax(totalTaxA, taxCrA)
	netPayableTaxB := c.netPayableTax(totalTaxB, taxCrB)
	finalCr := append(taxCrA, taxCrB...)

	return netPayableTaxA, netPayableTaxB, finalCr
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

func (c *Calculator) totalCredits(netIncomeA, netIncomeB float64) (totalCrA, totalCrB []core.TaxCredit) {

	taxPayerA, taxPayerB := c.makeTaxPayers(netIncomeA, netIncomeB)

	newCrA := taxCreditGroup(
		c.contraFormula.Apply(taxPayerA),
	).typecast()

	newCrB := taxCreditGroup(
		c.contraFormula.Apply(taxPayerB),
	).typecast()

	crSpouseA := make([]core.TaxCredit, len(c.crSpouseA))
	for i, cr := range c.crSpouseA {
		crSpouseA[i] = cr.ShallowCopy()
	}

	crSpouseB := make([]core.TaxCredit, len(c.crSpouseB))
	for i, cr := range c.crSpouseB {
		crSpouseB[i] = cr.ShallowCopy()
	}

	totalCrA = append(crSpouseA, newCrA...)
	totalCrB = append(crSpouseB, newCrB...)

	c.contraFormula.FilterAndSort(&totalCrA)
	c.contraFormula.FilterAndSort(&totalCrB)

	return totalCrA, totalCrB
}

// netPayableTax returns the payable tax after using the given credits. It also
// sets new values for the given tax credits according to how they were used.
// The credits are used one by one in the same given order
func (c *Calculator) netPayableTax(taxAmount float64, credits []core.TaxCredit) float64 {

	// FIXME: this is a bit complex
	for _, cr := range credits {

		ruleType := cr.Rule().Type
		initial, used, remaining := cr.Amounts()

		if taxAmount <= 0.0 && ruleType == core.CrRuleTypeCanCarryForward {
			continue
		}

		if taxAmount <= 0.0 && ruleType == core.CrRuleTypeNotCarryForward {
			cr.SetAmounts(initial, used, 0)
			continue
		}

		if taxAmount >= remaining || ruleType == core.CrRuleTypeCashable {
			taxAmount -= remaining
			cr.SetAmounts(initial, used+remaining, 0)
			continue
		}

		// reached at most once: 0 < taxAmount < remaining
		if ruleType == core.CrRuleTypeNotCarryForward {
			cr.SetAmounts(initial, used+taxAmount, 0)
			taxAmount = 0.0
		} else if ruleType == core.CrRuleTypeCanCarryForward {
			cr.SetAmounts(initial, used+taxAmount, remaining-taxAmount)
			taxAmount = 0.0
		}

		// we only reach here if ruleType is unknown which is a nop
	}

	return taxAmount
}

// setCredits stores relevent credits from the given credits in this calculator.
// Subsequent calls to other calculator functions may or may not be influenced
// by these credits.
func (c *Calculator) setCredits(credits []core.TaxCredit) {

	c.crSpouseA = make([]core.TaxCredit, 0)
	c.crSpouseB = make([]core.TaxCredit, 0)

	for _, cr := range credits {

		if cr == nil {
			continue
		}

		if _, _, remaining := cr.Amounts(); remaining == 0 {
			continue
		}

		if cr.Region() != c.taxRegion {
			continue
		}

		if cr.Year() > c.taxYear {
			continue
		}

		ref := cr.ReferenceFinancer()
		if ref == nil {
			continue
		}

		if ref == c.finances.SpouseA() {
			c.crSpouseA = append(c.crSpouseA, cr)
		} else if ref == c.finances.SpouseB() {
			c.crSpouseB = append(c.crSpouseB, cr)
		}

	}

}

// makeTaxPayers returns dual tax payers from the given net income amounts and
// the finances stored in this calculator
func (c *Calculator) makeTaxPayers(netIncomeA, netIncomeB float64) (taxPayerA, taxPayerB *TaxPayer) {

	financesA := c.finances.SpouseA()
	financesB := c.finances.SpouseB()

	if financesA != nil {
		taxPayerA = &TaxPayer{
			Finances:        financesA,
			NetIncome:       netIncomeA,
			SpouseFinances:  financesB,
			SpouseNetIncome: netIncomeB,
			Dependents:      c.dependents,
		}
	}

	if financesB != nil {
		taxPayerB = &TaxPayer{
			Finances:        financesB,
			NetIncome:       netIncomeB,
			SpouseFinances:  financesA,
			SpouseNetIncome: netIncomeA,
			Dependents:      c.dependents,
		}
	}

	return taxPayerA, taxPayerB
}

func (c *Calculator) panicIfEqNonNilSpouses() {
	if c.finances.SpouseA() != nil && c.finances.SpouseB() != nil {
		if c.finances.SpouseA() == c.finances.SpouseB() {
			panic("household finances cannot reference the same spouse as spouse A & spouse B")
		}
	}
}
