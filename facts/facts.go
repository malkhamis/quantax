// Package facts provides convenience types for building applications meant
// for bracket-based tax systems
package facts

import "github.com/pkg/errors"

// Facts represents tax parameters set by governments for a given tax year
type Facts struct {
	Year uint
	FactsFed
	FactsProv
}

// Validate ensures the given federal facts are valid for use
func (f Facts) Validate() error {

	err := f.FactsFed.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid federal facts")
	}

	err = f.FactsProv.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid provincial facts")
	}

	return nil
}

// Clone returns a copy of this Facts
func (f Facts) Clone() Facts {
	return Facts{
		Year:      f.Year,
		FactsFed:  f.FactsFed.Clone(),
		FactsProv: f.FactsProv.Clone(),
	}
}

// FactsFed represents tax parameters set by the federal government
type FactsFed struct {
	Rates BracketRates
}

// Validate ensures the given federal facts are valid for use
func (f FactsFed) Validate() error {
	return f.Rates.Validate()
}

// Clone returns a copy of this FactsFed
func (f FactsFed) Clone() FactsFed {
	return FactsFed{Rates: f.Rates.Clone()}
}

// FactsProv represents tax parameters set by the provincial government
type FactsProv struct {
	Rates BracketRates
}

// Validate ensures the given provincial facts are valid for use
func (f FactsProv) Validate() error {
	return f.Rates.Validate()
}

// Clone returns a copy of this FactsProv
func (f FactsProv) Clone() FactsProv {
	return FactsProv{Rates: f.Rates.Clone()}
}
