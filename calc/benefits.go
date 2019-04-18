package calc

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
)

// ChildBenefitCalculator is used to calculate recievable child benefits for
// families with dependent children.
type ChildBenefitCalculator interface {
	// Calc returns the recievable amount of child benefits for the given
	// finances and the children set in the calculator
	Calc(finance.IncomeDeductor) float64
	// SetBeneficiaries sets the children which the calculator will compute the
	// benefits for in subsequent calls to Calc()
	SetBeneficiaries(...human.Person)
}
