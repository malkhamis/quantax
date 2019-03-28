package tax

// Formula computes payable taxes on the given income
type Formula interface {
	// Apply applies the formula on the income
	Apply(income float64) float64
	// Validate checks if the formula is valid for use
	Validate() error
}
