// Package finances provides the basic tools and data type needed to compute
// Canadian taxes and benefits given financial information
package finance

// IncomeDeductor stores financial information such that amounts for income and
// deduction sources for a tax payer can be retrieved by their source
type IncomeDeductor interface {
	// TotalIncome returns the sum of income for the given sources only.
	// If no sources given, the total income for all sources is returned
	TotalIncome(sources ...IncomeSource) float64
	// TotalDeductions returns the sum of TotalDeductionss for the given
	// sources only. If no sources given, the total deduction for all sources
	// is returned
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

// NewEmptyIndividialFinances returns an instance whose EOY is initialized to
// endOfYear and whose maps are initialized with no income sources
func NewEmptyIndividialFinances(endOfYear uint) *IndividualFinances {
	return &IndividualFinances{
		EOY:        endOfYear,
		Income:     make(IncomeBySource),
		Deductions: make(DeductionBySource),
	}
}

// TotalIncome returns the sum of income for the given sources only. If no
// sources given, the total income for all sources is returned
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

// TotalDeductions returns the sum of TotalDeductionss for the given sources
// only. If no sources given, the total deduction for all sources is returned
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

// AddIncome adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddIncome(source IncomeSource, amount float64) {
	f.Income[source] += amount
}

// AddDeduction adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddDeduction(source DeductionSource, amount float64) {
	f.Deductions[source] += amount
}

// Clone returns a copy of this instance
func (f *IndividualFinances) Clone() *IndividualFinances {

	if f == nil {
		return nil
	}

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

// NewHouseholdFinances returns a new instance, appending the given non-nil
// individual finances. The return instance is never nil
func NewHouseholdFinances(finances ...*IndividualFinances) HouseholdFinances {

	hf := make(HouseholdFinances, 0, len(finances))

	for _, f := range finances {
		if f != nil {
			hf = append(hf, f)
		}
	}

	return hf
}

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

// Clone returns a copy of this instance
func (hf HouseholdFinances) Clone() HouseholdFinances {

	var clone HouseholdFinances

	if hf != nil {
		clone = make(HouseholdFinances, len(hf))
	}

	for i, f := range hf {
		clone[i] = f.Clone()
	}

	return clone
}
