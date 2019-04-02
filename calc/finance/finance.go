// Package finances provides the basic tools and data type needed to compute
// Canadian taxes and benefits given financial information
package finance

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	EndOfYear uint
	Cash      float64
	Income    IncomeBySource
	Deduction DeductionBySource
	RRSP      struct {
		ContributionRoom    float64
		UnclaimedDeductions float64
	}
}

func NewEmptyIndividialFinances(endOfYear uint) *IndividualFinances {
	return &IndividualFinances{
		EndOfYear: endOfYear,
		Income:    make(IncomeBySource),
		Deduction: make(DeductionBySource),
	}
}

// HouseholdFinances represents financial data for a couple, family etc
type HouseholdFinances []IndividualFinances

// Income calculate the the total income of the household from the given income
// sources. if no sources are given, the sum of all income sources is returned
func (hf HouseholdFinances) Income(sources ...IncomeSource) float64 {

	var total float64

	if len(sources) == 0 {
		for _, individualFinances := range hf {
			total += individualFinances.Income.Sum()
		}
		return total
	}

	for _, individualFinances := range hf {
		for _, source := range sources {
			total += individualFinances.Income[source]
		}
	}

	return total
}

// Deductions calculate the the total deduciton of the household from the given
// deduction sources. if no sources are given, the sum of all deduction sources
// is returned
func (hf HouseholdFinances) Deductions(sources ...DeductionSource) float64 {

	var total float64

	if len(sources) == 0 {
		for _, individualFinances := range hf {
			total += individualFinances.Deduction.Sum()
		}
		return total
	}

	for _, individualFinances := range hf {
		for _, source := range sources {
			total += individualFinances.Deduction[source]
		}
	}

	return total
}

// Cash returns the total cash balanace of this household
func (hf HouseholdFinances) Cash() float64 {
	var total float64
	for _, f := range hf {
		total += f.Cash
	}
	return total
}
