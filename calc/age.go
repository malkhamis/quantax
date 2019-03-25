package calc

import (
	"github.com/pkg/errors"
)

// AgeRange represents a range of age in months, e.g. [12, 36] to indicate
// one to three years old
type AgeRange [2]uint

// NewAgeRange returns a new age range from the given min/max ages
func NewAgeRange(min, max uint) (AgeRange, error) {
	ar := AgeRange{min, max}
	return ar, ar.Validate()
}

// Validate ensures this instance is valid for the intended use. Users need to
// call this method before use only if the instance was manually created/modified
func (ar AgeRange) Validate() error {

	switch {

	case ar[0] < 0 || ar[1] < 0:
		return errors.Wrap(ErrValNeg, "invalid min/max amounts")

	case ar[0] > ar[1]:
		return errors.Wrap(ErrBoundsReversed, "invalid age range")

	default:
		return nil
	}

}

// Min returns the lower limit/bound of this age
func (ar AgeRange) Min() uint {
	return ar[0]
}

// Max returns the upper limit/bound of this age
func (ar AgeRange) Max() uint {
	return ar[1]
}

// Clone returns a copy of this bracket
func (ar AgeRange) Clone() AgeRange {
	return AgeRange{ar[0], ar[1]}
}
