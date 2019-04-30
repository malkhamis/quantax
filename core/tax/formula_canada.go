package tax

import "github.com/malkhamis/quantax/core"

// CanadianFormula is used to calculate Canadian federal and provincial taxes
type CanadianFormula struct {
	// TaxYear is the tax year this contra formula is associated with
	TaxYear uint
	// TaxRegion is the tax region this contra formula is associated with
	TaxRegion core.Region
	core.WeightedBrackets
}

// Apply applies this formula on the given net income and returns payable tax
func (cf *CanadianFormula) Apply(netIncome float64) float64 {
	return cf.WeightedBrackets.Apply(netIncome)
}

// Year returns the tax year for this formula
func (cf *CanadianFormula) Year() uint {
	return cf.TaxYear
}

// Region returns the tax region for this formula
func (cf *CanadianFormula) Region() core.Region {
	return cf.TaxRegion
}

// Clone returns a copy of this formula
func (cf *CanadianFormula) Clone() Formula {

	if cf == nil {
		return nil
	}

	clone := &CanadianFormula{
		WeightedBrackets: cf.WeightedBrackets.Clone(),
		TaxYear:          cf.TaxYear,
		TaxRegion:        cf.TaxRegion,
	}

	return clone
}

// Validate ensures that this formula is valid for use
func (cf *CanadianFormula) Validate() error {
	return cf.WeightedBrackets.Validate()
}
