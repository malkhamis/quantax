package income

import "github.com/malkhamis/quantax/core"

// Recipe describes required adjustment on finances before incorporating them
// in the calculation of net income
type Recipe struct {
	// IncomeAdjusters represents the income sources that require adjustment
	// before incorporating them in the net income calculation
	IncomeAdjusters map[core.FinancialSource]Adjuster
	// DeductionAdjusters represents the deduction sources that require
	// adjustment before incorporating them in the net income
	DeductionAdjusters map[core.FinancialSource]Adjuster
}

// Clone returns a copy of this recipe
func (r *Recipe) Clone() *Recipe {

	if r == nil {
		return nil
	}

	clone := new(Recipe)

	if r.IncomeAdjusters != nil {
		clone.IncomeAdjusters = make(map[core.FinancialSource]Adjuster)
		for src, adj := range r.IncomeAdjusters {
			clone.IncomeAdjusters[src] = adj.Clone()
		}
	}

	if r.DeductionAdjusters != nil {
		clone.DeductionAdjusters = make(map[core.FinancialSource]Adjuster)
		for src, adj := range r.DeductionAdjusters {
			clone.DeductionAdjusters[src] = adj.Clone()
		}
	}

	return clone
}
