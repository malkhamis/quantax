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
	// Version returns the version, which is the number of changes made
	Version() uint64
}

// RRSPAmounts represent amounts related to RRSP
type RRSPAmounts struct {
	ContributionRoom    float64
	UnclaimedDeductions float64
}

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	cash        float64
	income      amountBySource
	deductions  amountBySource
	miscAmounts amountBySource
	rrspAmounts RRSPAmounts
	version     uint64
}

// NewEmptyIndividualFinances returns an empty instance with version zero
func NewEmptyIndividualFinances() *IndividualFinances {
	return &IndividualFinances{
		income:      make(amountBySource),
		deductions:  make(amountBySource),
		miscAmounts: make(amountBySource),
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
		return f.income.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.income[source]
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
		return f.deductions.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.deductions[source]
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
		return f.miscAmounts.Sum()
	}

	var total float64
	for _, source := range sources {
		total += f.miscAmounts[source]
	}
	return total
}

// AddAmount adds the given amount to the finances. If the source is within the
// range of recognized income sources, it is added as an income. If it within
// the range of recognized deduction sources, it is added as a deduction. Other
// sources, including unknown ones are added as misc amounts
func (f *IndividualFinances) AddAmount(source FinancialSource, amount float64) {
	switch {
	case source.IsIncomeSource():
		f.income[source] += amount
	case source.IsDeductionSource():
		f.deductions[source] += amount
	default:
		f.miscAmounts[source] += amount
	}
}

// RemoveAmounts removes the given stored sources. If sources is empty, the call
// is noop. This operation ensures that subsequent calls to IncomeSources,
// DeductionSources, and MiscSources returns a list that does not contain given
// source(s)
func (f *IndividualFinances) RemoveAmounts(sources ...FinancialSource) {
	for _, s := range sources {
		delete(f.income, s)
		delete(f.deductions, s)
		delete(f.miscAmounts, s)
	}
}

// IncomeSources returns a set of all income sources in this instance
func (f *IndividualFinances) IncomeSources() map[FinancialSource]struct{} {

	set := make(map[FinancialSource]struct{})

	if f == nil {
		return set
	}

	for source := range f.income {
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

	for source := range f.deductions {
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

	for source := range f.miscAmounts {
		set[source] = struct{}{}
	}

	return set
}

// RRSPAmounts returns the RRSP information in this instance
func (f *IndividualFinances) RRSPAmounts() RRSPAmounts {
	return f.rrspAmounts
}

// SetRRSPAmounts sets the RRSP information in this instance
func (f *IndividualFinances) SetRRSPAmounts(amounts RRSPAmounts) {
	f.rrspAmounts = amounts
}

// Cash returns the free cash balance set in this instance
func (f *IndividualFinances) Cash() float64 {
	return f.cash
}

// SetCash sets the free cash balance to the given value
func (f *IndividualFinances) SetCash(amount float64) {
	f.cash = amount
}

// Version returns the version of this instance
func (f *IndividualFinances) Version() uint64 {
	return f.version
}

// Clone returns a copy of this instance
func (f *IndividualFinances) Clone() *IndividualFinances {

	if f == nil {
		return nil
	}

	clone := &IndividualFinances{
		cash:        f.cash,
		income:      f.income.Clone(),
		deductions:  f.deductions.Clone(),
		miscAmounts: f.miscAmounts.Clone(),
		rrspAmounts: f.rrspAmounts,
		version:     f.version,
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
		total += f.Cash()
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

		for source := range f.income {
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

		for source := range f.deductions {
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

		for source := range f.miscAmounts {
			set[source] = struct{}{}
		}

	}

	return set
}

// Version returns the version of this instance
func (hf HouseholdFinances) Version() uint64 {
	var v uint64
	for _, f := range hf {
		v += f.Version()
	}
	return v
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
