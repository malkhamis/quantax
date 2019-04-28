package tax

import "github.com/malkhamis/quantax/core"

// CanadianFormula is used to calculate Canadian federal and provincial taxes
type CanadianFormula struct {
	RelatedTaxInfo core.TaxInfo
	core.WeightedBrackets
}

// Apply applies this formula on the given net income and returns payable tax
func (cf *CanadianFormula) Apply(netIncome float64) float64 {
	return cf.WeightedBrackets.Apply(netIncome)
}

// TaxInfo TODO
func (cf *CanadianFormula) TaxInfo() core.TaxInfo {
	return cf.RelatedTaxInfo
}

// Clone returns a copy of this formula
func (cf *CanadianFormula) Clone() Formula {

	if cf == nil {
		return nil
	}

	clone := &CanadianFormula{
		WeightedBrackets: cf.WeightedBrackets.Clone(),
		RelatedTaxInfo:   cf.RelatedTaxInfo,
	}

	return clone
}

// Validate ensures that this formula is valid for use
func (cf *CanadianFormula) Validate() error {
	return cf.WeightedBrackets.Validate()
}
