package finance

import (
	"github.com/malkhamis/quantax/core"
)

// compile-time check for interface implementation
var _ core.Financer = (*IndividualFinances)(nil)

// IndividualFinances represents the financial data of an individual
type IndividualFinances struct {
	income      map[core.FinancialSource]float64
	deductions  map[core.FinancialSource]float64
	miscAmounts map[core.FinancialSource]float64
	version     uint64
}

// NewEmptyIndividualFinances returns an empty instance with version equals one
func NewIndividualFinances() *IndividualFinances {
	return &IndividualFinances{
		income:      make(map[core.FinancialSource]float64),
		deductions:  make(map[core.FinancialSource]float64),
		miscAmounts: make(map[core.FinancialSource]float64),
		version:     1,
	}
}

// TotalAmount returns the sum of of the given sources only. If no sources
// given or f is nil, it returns zero. Source that don't exist in f are skipped
func (f *IndividualFinances) TotalAmount(sources ...core.FinancialSource) float64 {

	if f == nil {
		return 0.0
	}

	var total float64
	for _, src := range sources {
		switch {

		case src.IsIncomeSource():
			total += f.income[src]

		case src.IsDeductionSource():
			total += f.deductions[src]

		default:
			total += f.miscAmounts[src]
		}
	}

	return total
}

// AddAmount adds the given amount to the finances. If the source is within the
// range of recognized income sources, it is added as an income. If it within
// the range of recognized deduction sources, it is added as a deduction. Other
// sources, including unknown ones are added as misc amounts. If 'f' is nil,
// this method panics
func (f *IndividualFinances) AddAmount(source core.FinancialSource, amount float64) {

	switch {

	case source.IsIncomeSource():
		f.income[source] += amount

	case source.IsDeductionSource():
		f.deductions[source] += amount

	default:
		f.miscAmounts[source] += amount
	}

	// TODO: this is too naive
	f.version++
}

// SetAmount force-set the given amount for the given source. If the source
// is within the range of recognized income sources, it is set as an income
// source.If it is within the range of recognized deduction sources, it is
// added as a deduction. Other sources, including unknown ones are added as
// misc amounts. If 'f' is nil, this method panics
func (f *IndividualFinances) SetAmount(source core.FinancialSource, amount float64) {

	switch {

	case source.IsIncomeSource():
		f.income[source] = amount

	case source.IsDeductionSource():
		f.deductions[source] = amount

	default:
		f.miscAmounts[source] = amount
	}

	// TODO: this is too naive
	f.version++
}

// RemoveAmounts removes the given stored sources. If sources is empty, the call
// is noop. This operation ensures that subsequent calls to IncomeSources,
// DeductionSources, and MiscSources returns a list that does not contain given
// source(s). If 'f' is nil, this method panics
func (f *IndividualFinances) RemoveAmounts(sources ...core.FinancialSource) {

	for _, src := range sources {
		switch {

		case src.IsIncomeSource():
			delete(f.income, src)

		case src.IsDeductionSource():
			delete(f.deductions, src)

		default:
			delete(f.miscAmounts, src)
		}
	}

	// TODO: this is too naive
	f.version++
}

// IncomeSources returns a set of all income sources in this instance
func (f *IndividualFinances) IncomeSources() []core.FinancialSource {

	if f == nil {
		return nil
	}

	sources := make([]core.FinancialSource, 0, len(f.income))
	for src := range f.income {
		sources = append(sources, src)
	}

	return sources

}

// DeductionSources returns a set of all deduciton sources in this instance.
// The returned map is never nil
func (f *IndividualFinances) DeductionSources() []core.FinancialSource {

	if f == nil {
		return nil
	}

	sources := make([]core.FinancialSource, 0, len(f.deductions))
	for src := range f.deductions {
		sources = append(sources, src)
	}

	return sources
}

// MiscSourceSet returns a set of all miscellaneous sources in this instance
func (f *IndividualFinances) MiscSources() []core.FinancialSource {

	if f == nil {
		return nil
	}

	sources := make([]core.FinancialSource, 0, len(f.miscAmounts))
	for src := range f.miscAmounts {
		sources = append(sources, src)
	}

	return sources
}

// AllSources returns the income, deduction, and misc sources in this instance
func (f *IndividualFinances) AllSources() []core.FinancialSource {

	all := f.IncomeSources()
	all = append(all, f.DeductionSources()...)
	all = append(all, f.MiscSources()...)
	return all
}

// Version returns the version of this instance. If 'f' is nil, it returns zero
func (f *IndividualFinances) Version() uint64 {
	if f == nil {
		return 0
	}
	return f.version
}

// Clone returns a copy of this instance. If 'f' is nil, it returns nil
func (f *IndividualFinances) Clone() core.FinanceMutator {
	if f == nil {
		return nil
	}
	return core.FinanceMutator(f.clone())
}

// clone returns a copy of this instance
func (f *IndividualFinances) clone() *IndividualFinances {

	if f == nil {
		return nil
	}

	clone := &IndividualFinances{
		income:      amountBySource(f.income).Clone(),
		deductions:  amountBySource(f.deductions).Clone(),
		miscAmounts: amountBySource(f.miscAmounts).Clone(),
		version:     f.version,
	}

	return clone
}

// amountBySource is a helper type used to encapsulate the logic of map cloning
type amountBySource map[core.FinancialSource]float64

// Clone returns a copy of this instance
func (s amountBySource) Clone() amountBySource {

	var clone amountBySource

	if s != nil {
		clone = make(amountBySource)
	}

	for source, amount := range s {
		clone[source] = amount
	}

	return clone
}
