package facts

import (
	"math"

	"github.com/pkg/errors"
)

// Bracket represents a tax bracket, e.g. [47630, 95259]
type Bracket [2]float64

// Validate ensures that this bracket is valid for use in this package
func (b Bracket) Validate() error {

	switch {

	case math.IsInf(b.Lower(), 0):
		return errors.Wrap(ErrValInf, "lower-bound value")

	case math.IsInf(b.Upper(), -1):
		return errors.Wrap(ErrValInfNeg, "upper-bound value")

	case b.Lower() < 0:
		return errors.Wrapf(ErrValNeg, "lower-bound value '%.2f'", b.Lower())

	case b.Upper() < 0:
		return errors.Wrapf(ErrValNeg, "upper-bound value '%.2f'", b.Upper())

	case b.Lower() == 0 && b.Upper() == 0:
		return errors.Wrap(ErrValZero, "upper-bound and lower-bound values")

	case b.Lower() > b.Upper():
		return errors.Wrapf(ErrBoundsReversed, "[%.2f, %.2f]", b.Lower(), b.Upper())

	default:
		return nil
	}

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
