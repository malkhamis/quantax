package core

var (
	_ Financer = (*IndividualFinances)(nil)
	_ Financer = (HouseholdFinances)(nil)
)

type Financer interface {
	// TotalIncome returns the sum of income for the given sources only. If no
	// sources given, the total income for all sources is returned
	TotalIncome(sources ...FinancialSource) float64
	// TotalDeductions returns the sum of TotalDeductionss for the given sources
	// only. If no sources given, the total deduction for all sources is returned
	TotalDeductions(sources ...FinancialSource) float64
	// MiscAmount returns the amount of the given miscellaneous sources only.
	// If no sources given, the total amount for all sources is returned
	MiscAmount(source ...FinancialSource) float64
	// IncomeSources returns a set of all income sources in this instance. The
	// returned map is never nil
	IncomeSources() map[FinancialSource]struct{}
	// DeductionSources returns a set of all deduciton sources in this instance.
	// The returned map is never nil
	DeductionSources() map[FinancialSource]struct{}
	// MiscSources returns a set of all miscellaneous sources in this instance.
	// The returned map is never nil
	MiscSources() map[FinancialSource]struct{}
	// TODO: Version()
}

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	EOY                     uint
	Cash                    float64
	Income                  AmountBySource
	Deductions              AmountBySource
	MiscAmounts             AmountBySource
	RRSPContributionRoom    float64
	RRSPUnclaimedDeductions float64
	// TODO: version
}

// NewEmptyIndividualFinances returns an instance whose EOY is initialized to
// endOfYear and whose maps are initialized with no amounts/sources
func NewEmptyIndividualFinances(endOfYear uint) *IndividualFinances {
	return &IndividualFinances{
		EOY:         endOfYear,
		Income:      make(AmountBySource),
		Deductions:  make(AmountBySource),
		MiscAmounts: make(AmountBySource),
	}
}

// TotalIncome returns the sum of income for the given sources only. If no
// sources given, the total income for all sources is returned. If 'f' is nil
// zero is returned
func (f *IndividualFinances) TotalIncome(sources ...FinancialSource) float64 {

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

// TotalDeductions returns the sum of TotalDeductions for the given sources
// only. If no sources given, the total deduction for all sources is returned.
// If 'f' is nil, zero is returned
func (f *IndividualFinances) TotalDeductions(sources ...FinancialSource) float64 {

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

// MiscAmount returns the amount for the given miscellaneous source. If 'f' is
// nil zero is returned
func (f *IndividualFinances) MiscAmount(sources ...FinancialSource) float64 {

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

// TODO: replace all these AddXXXX with AddAmount(Source, Float64). Depending
// on where the source is in the iota, we mutate the right map
// AddIncome adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddIncome(source FinancialSource, amount float64) {
	f.Income[source] += amount
}

// RemoveIncome removes the given stored income sources. If sources is empty,
// the function is noop. This operation ensures that subsequent calls to
// IncomeSources() returns a list that does not contain given income source(s)
func (f *IndividualFinances) RemoveIncome(sources ...FinancialSource) {
	for _, s := range sources {
		delete(f.Income, s)
	}
}

// AddDeduction adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddDeduction(source FinancialSource, amount float64) {
	f.Deductions[source] += amount
}

// RemoveDeduction removes the given stored deduction sources. If sources is
// empty, the function is nop. This operation ensures that subsequent calls to
// DeductionSources() returns a list that does not contain given deduction
// source(s)
func (f *IndividualFinances) RemoveDeduction(sources ...FinancialSource) {
	for _, s := range sources {
		delete(f.Deductions, s)
	}
}

// AddMiscAmount adds the given amount to the stored amount of the given source
func (f *IndividualFinances) AddMiscAmount(source FinancialSource, amount float64) {
	f.MiscAmounts[source] += amount
}

// RemoveMiscAmount removes the given stored miscellaneous sources. If sources
// is empty, the function is noop. This operation ensures that subsequent calls
// to MiscSources() returns a list that does not contain given misc source(s)
func (f *IndividualFinances) RemoveMiscAmount(sources ...FinancialSource) {
	for _, s := range sources {
		delete(f.MiscAmounts, s)
	}
}

// IncomeSources returns a set of all income sources in this instance
func (f *IndividualFinances) IncomeSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

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
func (f *IndividualFinances) DeductionSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

	if f == nil {
		return set
	}

	for source := range f.Deductions {
		set[source] = struct{}{}
	}

	return set
}

// MiscSourceSet returns a set of all miscellaneous sources in this instance
func (f *IndividualFinances) MiscSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

	if f == nil {
		return set
	}

	for source := range f.MiscAmounts {
		set[source] = struct{}{}
	}

	return set
}

// TODO update docstrings once version is added and set version of the clone to zero
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
func (hf HouseholdFinances) TotalIncome(sources ...FinancialSource) float64 {

	var total float64
	for _, f := range hf {
		total += f.TotalIncome(sources...)
	}
	return total
}

// Deductions calculate the the total deduciton of the household from the given
// deduction sources. if no sources are given, the sum of all deduction sources
// is returned
func (hf HouseholdFinances) TotalDeductions(sources ...FinancialSource) float64 {

	var total float64
	for _, f := range hf {
		total += f.TotalDeductions(sources...)
	}
	return total
}

// MiscAMount calculate the the total miscellaneous amount of the household from
// the given sources. if no sources are given, the sum of all miscellaneous
// sources is returned
func (hf HouseholdFinances) MiscAmount(sources ...FinancialSource) float64 {

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
func (hf HouseholdFinances) IncomeSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

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

// DeductionSources returns a set of all deduction sources in this instance.
// The returned map is never nil
func (hf HouseholdFinances) DeductionSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

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
func (hf HouseholdFinances) MiscSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

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
