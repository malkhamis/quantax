// Package benefits provides implementation for various benefit interfaces
// as defined in package calc
package benefits

import "github.com/malkhamis/quantax/calc"

//go:generate stringer -type=IncomeType -output=incometype_gen.go
type IncomeType uint

const (
	REGU IncomeType = iota // regular income (income - deductions)
	AFNI                   // adjusted family net income (income - deductions - UCCB - RDSP)
)

// Calc calculates the income according to the underlying method
// TODO fix this
func (it IncomeType) Calc(finances calc.FamilyFinances) float64 {

	switch it {
	case AFNI, REGU:
		fallthrough
	default:
		// this is not exactly right as we should incorporate
		// UCCB and other shit that I don't care about for now
		return finances.Income() - finances.Deductions()
	}

}

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(income float64, children ...calc.Person) float64
	// IncomeCalcMethod returns the method of calculating the income
	IncomeCalcMethod() IncomeType
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() ChildBenefitFormula
}
