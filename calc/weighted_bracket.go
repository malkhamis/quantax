package calc

import (
	"math"

	"github.com/pkg/errors"
)

var _ BasicFormula = (WeightedBracketFormula)(nil)

// Bracket represents a float number range, e.g. [47630.51, 95259.32]
type Bracket [2]float64

// Validate ensures that this bracket is within [0, +Inf]
func (b Bracket) Validate() error {

	if b.Lower() > b.Upper() {
		return errors.Wrapf(ErrBoundsReversed, "[%.2f, %.2f]", b.Lower(), b.Upper())
	}

	return nil
}

// Amount returns the difference between upper and lower bound
func (b Bracket) Amount() float64 {
	return b.Upper() - b.Lower()
}

// Lower returns the lower limit/bound of this bracket
func (b Bracket) Lower() float64 {
	return b[0]
}

// Upper returns the upper limit/bound of this bracket
func (b Bracket) Upper() float64 {
	return b[1]
}

// Clone returns a copy of this bracket
func (b Bracket) Clone() Bracket {
	return Bracket{b[0], b[1]}
}

// WeightedBracketFormula maps weights (rates) to numeric ranges, e.g. brackets
type WeightedBracketFormula map[float64]Bracket

// Apply slices the given param into this formula's brackets. Then, it applies
// the rate asscoiated with the bracket to the sliced amounts and returns the
// sum of all results
func (wb WeightedBracketFormula) Apply(param float64) (result float64) {

	for rate, bracket := range wb {

		if param < bracket.Lower() {
			continue
		}

		if param >= bracket.Upper() {
			result += rate * bracket.Amount()
			continue
		}

		result += rate * (param - bracket.Lower())
	}

	return result
}

// Validate ensures that this weighted brackets object is valid for use
func (wb WeightedBracketFormula) Validate() error {

	for rate, bracket := range wb {

		if math.IsInf(rate, 0) {
			return errors.Wrap(ErrValInf, "invalid rate")
		}

		err := bracket.Validate()
		if err != nil {
			return errors.Wrap(err, "invalid bracket")
		}

	}

	return nil
}

// Clone returns a copy of this weighted bracket instance
func (wb WeightedBracketFormula) Clone() WeightedBracketFormula {

	var clone WeightedBracketFormula

	if wb != nil {
		clone = make(WeightedBracketFormula)
	}

	for rate, bracket := range wb {
		clone[rate] = bracket.Clone()
	}

	return clone
}
