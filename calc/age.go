package calc

import "errors"

// Sentinel errors that can be wrapped and returned
var (
	ErrInvalidAge      = errors.New("invalid age")
	ErrInvalidAgeRange = errors.New("invalid age range")
)

// AgeRange represents a range of age
type AgeRange [2]uint

// NewAgeRange returns a new age range from the given min/max ages
func NewAgeRange(min, max uint) (AgeRange, error) {
	ar := AgeRange{min, max}
	return ar, ar.Validate()
}

// Validate ensures this instance is valid for the intended use. Users need to
// call this method before use only if the instance was manually created/modified
func (ar AgeRange) Validate() error {

	if ar[0] > ar[1] {
		return ErrInvalidAgeRange
	}
	return nil
}

// Min returns the lower limit/bound of this age
func (ar AgeRange) Min() uint {
	return ar[0]
}

// Max returns the upper limit/bound of this age
func (ar AgeRange) Max() uint {
	return ar[1]
}
