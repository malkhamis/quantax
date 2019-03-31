package finance

// TODO: this needs to be accurate and cleaned up

//go:generate stringer -type=IncomeType -output=incometype_gen.go
type IncomeType uint

const (
	UNKOWN IncomeType = iota // regular income (income - deductions)
	REGULAR
	AFNI // adjusted family net income (income - deductions - UCCB - RDSP)
	EARNED
)

// Calc calculates the income according to the underlying method
// TODO fix this
func (it IncomeType) Calc(finances FamilyFinances) float64 {

	// this is not exactly right as we should incorporate
	// UCCB and other shit that I don't care about for now

	switch it {
	case AFNI:
		fallthrough
	case REGULAR:
		fallthrough
	default:
		return finances.Income() - finances.Deductions()
	}
}

func (it IncomeType) CalcForIndividials(finances IndividualFinances) float64 {

	// this is not exactly right as we should incorporate
	// UCCB and other shit that I don't care about for now

	switch it {
	case EARNED:
		return finances.Income
	case REGULAR:
		fallthrough
	default:
		return finances.Income - finances.Deductions
	}
}
