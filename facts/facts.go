package facts

// Facts represents tax parameters set by governments for a given tax year
type Facts struct {
	Year uint
	FactsFed
	FactsProv
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

// Clone returns a copy of this FactsFed
func (f FactsFed) Clone() FactsFed {
	return FactsFed{Rates: f.Rates.Clone()}
}

// FactsProv represents tax parameters set by the provincial government
type FactsProv struct {
	Rates BracketRates
}

// Clone returns a copy of this FactsProv
func (f FactsProv) Clone() FactsProv {
	return FactsProv{Rates: f.Rates.Clone()}
}
