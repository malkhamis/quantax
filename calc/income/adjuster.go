package income

// compile-time check for interface implementatino
var (
	_ Adjuster = WeightedAdjuster(0.0)
)

// Adjuster is a type that adjusts any given amount according to some logic
type Adjuster interface {
	// Adjusted returns an adjusted amount from the given finances
	Adjusted(float64) float64
	// Clone returns a copy of this adjuster
	Clone() Adjuster
}

// WeightedAdjuster multiplies itself by a given amount
type WeightedAdjuster float64

// Adjusted multiplies 'wa' by 'amount' and returns the result
func (wa WeightedAdjuster) Adjusted(amount float64) float64 {
	return amount * float64(wa)
}

// Clone returns a copy of this instance
func (wa WeightedAdjuster) Clone() Adjuster {
	return wa
}
