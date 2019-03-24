package benefits

// Reducer is type that reduces a given amount
type Reducer interface {
	// Reduce reduces the given amount, where the reduction is controlled by the
	// given control value. It is up to the implementation to interpret the
	// control value. The returned value is the reduciton amount
	Reduce(amount float64, control float64) float64
	// Clone returns a copy of the underlying reducer
	Clone() Reducer
	// Validate ensures that the reducer is valid for use
	Validate() error
}
