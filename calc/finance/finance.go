// Package finances provides the basic tools and data type needed to compute
// Canadian taxes and benefits given financial information
package finance

var (
	_ IncomeDeductor = (*IndividualFinances)(nil)
	_ IncomeDeductor = (HouseholdFinances)(nil)
)

type IncomeDeductor interface {
	// TotalIncome returns the sum of income for the given sources only. If no
	// sources given, the total income for all sources is returned
	TotalIncome(sources ...IncomeSource) float64
	// TotalDeductions returns the sum of TotalDeductionss for the given sources
	// only. If no sources given, the total deduction for all sources is returned
	TotalDeductions(sources ...DeductionSource) float64
	// MiscAmount returns the amount of the given miscellaneous sources only.
	// If no sources given, the total amount for all sources is returned
	MiscAmount(source ...MiscSource) float64
	// IncomeSources returns a set of all income sources in this instance. The
	// returned map is never nil
	IncomeSources() IncomeSourceSet
	// DeductionSources returns a set of all deduciton sources in this instance.
	// The returned map is never nil
	DeductionSources() DeductionSourceSet
	// MiscSources returns a set of all miscellaneous sources in this instance.
	// The returned map is never nil
	MiscSources() MiscSourceSet
}

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	EOY                     uint
	Cash                    float64
	Income                  IncomeBySource
	Deductions              DeductionBySource
	MiscAmounts             MiscAmountsBySource
	RRSPContributionRoom    float64
	RRSPUnclaimedDeductions float64
}

// NewEmptyIndividualFinances returns an instance whose EOY is initialized to
// endOfYear and whose maps are initialized with no income sources
func NewEmptyIndividualFinances(endOfYear uint) *IndividualFinances {
	return &IndividualFinances{
		EOY:         endOfYear,
		Income:      make(IncomeBySource),
		Deductions:  make(DeductionBySource),
		MiscAmounts: make(MiscAmountsBySource),
	}
}

// TotalIncome returns the sum of income for the given sources only. If no
// sources given, the total income for all sources is returned. If 'f' is nil
// zero is returned
func (f *IndividualFinances) TotalIncome(sources ...IncomeSource) float64 {

	if f == nil {
		return 0.0
	}

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

	if f == nil {
		return 0.0
	}

	if len(sources) == 0 {
		return f.Deductions.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.Deductions[source]
	}
	return total
}

// MiscAmount returns the amount for the given miscellaneous source
func (f *IndividualFinances) MiscAmount(sources ...MiscSource) float64 {

	if f == nil {
		return 0.0
	}

	if len(sources) == 0 {
		return f.MiscAmounts.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.MiscAmounts[source]
	}
	return total
}

// AddIncome adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddIncome(source IncomeSource, amount float64) {
	f.Income[source] += amount
}

// RemoveIncome removes the given stored income sources. If sources is empty,
// the function is noop. This operation ensures that subsequent calls to
// IncomeSources() returns a list that does not contain given income source(s)
func (f *IndividualFinances) RemoveIncome(sources ...IncomeSource) {
	for _, s := range sources {
		delete(f.Income, s)
	}
}

// AddDeduction adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddDeduction(source DeductionSource, amount float64) {
	f.Deductions[source] += amount
}

// RemoveIncome removes the given stored deduction sources. If sources is empty,
// the function is noop. This operation ensures that subsequent calls to
// DeductionSources() returns a list that does not contain given deduciton
// source(s)
func (f *IndividualFinances) RemoveDeduction(sources ...DeductionSource) {
	for _, s := range sources {
		delete(f.Deductions, s)
	}
}

// AddDeduction adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddMiscAmount(source MiscSource, amount float64) {
	f.MiscAmounts[source] += amount
}

// RemoveMiscAmount removes the given stored miscellaneous sources. If sources
// is empty, the function is noop. This operation ensures that subsequent calls
// to DeductionSources() returns a list that does not contain given misc
// source(s)
func (f *IndividualFinances) RemoveMiscAmount(sources ...MiscSource) {
	for _, s := range sources {
		delete(f.MiscAmounts, s)
	}
}

// IncomeSources returns a set of all income sources in this instance
func (f *IndividualFinances) IncomeSources() IncomeSourceSet {

	set := make(IncomeSourceSet)

	if f == nil {
		return set
	}

	for source := range f.Income {
		set[source] = struct{}{}
	}

	return set

}

// DeductionSources returns a set of all deduciton sources in this instance.
// The returned map is never nil
func (f *IndividualFinances) DeductionSources() DeductionSourceSet {

	set := make(DeductionSourceSet)

	if f == nil {
		return set
	}

	for source := range f.Deductions {
		set[source] = struct{}{}
	}

	return set
}

// MiscSourceSet returns a set of all miscellaneous sources in this instance
func (f *IndividualFinances) MiscSources() MiscSourceSet {

	set := make(MiscSourceSet)

	if f == nil {
		return set
	}

	for source := range f.MiscAmounts {
		set[source] = struct{}{}
	}

	return set
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
		MiscAmounts:             f.MiscAmounts.Clone(),
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

// MiscAMount calculate the the total miscellaneous amount of the household from
// the given sources. if no sources are given, the sum of all miscellaneous
// sources is returned
func (hf HouseholdFinances) MiscAmount(sources ...MiscSource) float64 {

	var total float64
	for _, f := range hf {
		total += f.MiscAmount(sources...)
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

// IncomeSources returns a set of all income sources in this instance. The
// returned map is never nil
func (hf HouseholdFinances) IncomeSources() IncomeSourceSet {

	set := make(IncomeSourceSet)

	for _, f := range hf {

		if f == nil {
			continue
		}

		for source := range f.Income {
			set[source] = struct{}{}
		}

	}

	return set
}

// DeductionSources returns a set of all deduciton sources in this instance.
// The returned map is never nil
func (hf HouseholdFinances) DeductionSources() DeductionSourceSet {

	set := make(DeductionSourceSet)

	for _, f := range hf {

		if f == nil {
			continue
		}

		for source := range f.Deductions {
			set[source] = struct{}{}
		}

	}

	return set
}

// MiscSources returns a set of all miscellaneous sources in this instance.
// The returned map is never nil
func (hf HouseholdFinances) MiscSources() MiscSourceSet {

	set := make(MiscSourceSet)

	for _, f := range hf {

		if f == nil {
			continue
		}

		for source := range f.MiscAmounts {
			set[source] = struct{}{}
		}

	}

	return set
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
