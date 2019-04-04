// Package finances provides the basic tools and data type needed to compute
// Canadian taxes and benefits given financial information
package finance

type IncomeDeductor interface {
	TotalIncome(sources ...IncomeSource) float64
	TotalDeductions(sources ...DeductionSource) float64
}

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	EOY                     uint
	Cash                    float64
	Income                  IncomeBySource
	Deductions              DeductionBySource
	RRSPContributionRoom    float64
	RRSPUnclaimedDeductions float64
}

func NewEmptyIndividialFinances(endOfYear uint) *IndividualFinances {
	return &IndividualFinances{
		EOY:        endOfYear,
		Income:     make(IncomeBySource),
		Deductions: make(DeductionBySource),
	}
}

func (f *IndividualFinances) TotalIncome(sources ...IncomeSource) float64 {

	if len(sources) == 0 {
		return f.Income.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.Income[source]
	}
	return total
}

func (f *IndividualFinances) TotalDeductions(sources ...DeductionSource) float64 {

	if len(sources) == 0 {
		return f.Deductions.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.Deductions[source]
	}
	return total
}

func (f *IndividualFinances) AddIncome(source IncomeSource, amount float64) {
	f.Income[source] += amount
}

func (f *IndividualFinances) Clone() *IndividualFinances {

	clone := &IndividualFinances{
		EOY:                     f.EOY,
		Cash:                    f.Cash,
		Income:                  f.Income.Clone(),
		Deductions:              f.Deductions.Clone(),
		RRSPContributionRoom:    f.RRSPContributionRoom,
		RRSPUnclaimedDeductions: f.RRSPUnclaimedDeductions,
	}

	return clone
}

// HouseholdFinances represents financial data for a couple, family etc
type HouseholdFinances []*IndividualFinances

// Income calculate the the total income of the household from the given income
// sources. if no sources are given, the sum of all income sources is returned
func (hf HouseholdFinances) TotalIncome(sources ...IncomeSource) float64 {

	var total float64
	for _, f := range hf {
		total += f.TotalIncome(sources...)
	}
	return total
}

// Deductions calculate the the total deduciton of the household from the given
// deduction sources. if no sources are given, the sum of all deduction sources
// is returned
func (hf HouseholdFinances) TotalDeductions(sources ...DeductionSource) float64 {

	var total float64
	for _, f := range hf {
		total += f.TotalDeductions(sources...)
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
