// Package finances provides the basic tools and data type needed to compute
// Canadian taxes and benefits given financial information
package finance

// TODO: deductions should be a map[Type]float64 to allow for exclusion of
// certain deductions under different criteria
// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	Income     float64 `json:"income"`
	Deductions float64 `json:"deductions"`
	RRSPRoom   float64 `json:"rrsp-room"`
}

// FamilyFinances represents financial data for a couple
type FamilyFinances [2]IndividualFinances

// Income calculate the the total income of the couple. The calculation is
// only based on adding the income components of the couple
func (f FamilyFinances) Income() float64 {
	return f[0].Income + f[1].Income
}

// Deductions calculate the the total deductions of the couple. The calculation
// is only based on adding the deduction components of the couple
func (f FamilyFinances) Deductions() float64 {
	return f[0].Deductions + f[1].Deductions
}

// Split returns the individual finances that jointly represent this object
func (f FamilyFinances) Split() (IndividualFinances, IndividualFinances) {
	return f[0], f[1]
}
