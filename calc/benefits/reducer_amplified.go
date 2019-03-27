package benefits

import (
	"github.com/malkhamis/quantax/calc"
)

var _ Reducer = (*AmplifiedReducer)(nil)

// AmplifiedReducer represents an amount-reducing method in which the reduction
// is amplified by a given constant. This reducer is specifically written for
// calculating the reduction for BC Early Childhood Tax Benefits. It implements
// the 'Reducer' interface defined in this package
type AmplifiedReducer calc.WeightedBracketFormula

// Reduce returns the reduced value amplified by the given multiplier for the
// given amount
func (ar AmplifiedReducer) Reduce(amount float64, multiplier float64) (reduction float64) {
	return multiplier * calc.WeightedBracketFormula(ar).Apply(amount)
}

// Clone returns a copy of this amplified reducer
func (ar AmplifiedReducer) Clone() Reducer {

	clone := calc.WeightedBracketFormula(ar).Clone()
	return AmplifiedReducer(clone)
}

// Validate ensures that this instance is valid for use. Users need to call
// this method before use only if the instance was manually created/modified
func (ar AmplifiedReducer) Validate() error {
	return calc.WeightedBracketFormula(ar).Validate()
}
