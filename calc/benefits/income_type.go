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

	// this is not exactly right as we should incorporate
	// UCCB and other shit that I don't care about for now

	switch it {
	case AFNI:
		fallthrough
	case REGU:
		fallthrough
	default:
		return finances.Income() - finances.Deductions()
	}
}
