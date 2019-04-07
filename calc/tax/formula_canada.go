package tax

import "github.com/malkhamis/quantax/calc/finance"

// CanadianFormula is used to calculate Canadian federal and provincial taxes
type CanadianFormula struct {
	finance.WeightedBrackets
}

// Apply applies this formula on the given net income and returns payable tax
func (ct *CanadianFormula) Apply(netIncome float64) float64 {
	return ct.WeightedBrackets.Apply(netIncome)
}

// Clone returns a copy of this formula
func (ct *CanadianFormula) Clone() Formula {

	if ct == nil {
		return nil
	}

	clone := &CanadianFormula{
		WeightedBrackets: ct.WeightedBrackets.Clone(),
	}

	return clone
}

// Validate ensures that this formula is valid for use
func (ct *CanadianFormula) Validate() error {
	return ct.WeightedBrackets.Validate()
}
