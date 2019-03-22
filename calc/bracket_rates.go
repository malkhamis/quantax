package calc

import (
	"math"

	"github.com/pkg/errors"
)

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

// WeightedBrackets maps weights (rates) to numeric ranges, e.g. brackets
type WeightedBrackets map[float64]Bracket

// // Validate ensures that this weighted brackets object is valid for use
func (wb WeightedBrackets) Validate() error {

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
func (wb WeightedBrackets) Clone() WeightedBrackets {

	var clone WeightedBrackets

	if wb != nil {
		clone = make(WeightedBrackets)
	}

	for rate, bracket := range wb {
		clone[rate] = bracket.Clone()
	}

	return clone
}

// RateAdjBracketFormula is a formula used to slice the given parameter into
// ranges such that each range is weighted according to the mapped rate
type RateAdjBracketFormula struct {
	RateMap WeightedBrackets
	Param   float64
}

// NewRateAdjBracketFormula returns a new instance, copying the given rate map
func NewRateAdjBracketFormula(rateMap WeightedBrackets, param float64) (*RateAdjBracketFormula, error) {

	rab := &RateAdjBracketFormula{
		RateMap: rateMap.Clone(),
		Param:   param,
	}

	return rab, rab.Validate()
}

// Validate ensures that this formula is valid for use
func (r *RateAdjBracketFormula) Validate() error {
	return r.RateMap.Validate()
}

// Clone returns a copy of this formula
func (r *RateAdjBracketFormula) Clone() *RateAdjBracketFormula {

	return &RateAdjBracketFormula{
		RateMap: r.RateMap.Clone(),
		Param:   r.Param,
	}
}
