package calc

import (
	"testing"

	"github.com/pkg/errors"
)

func TestAgeRange_Valid(t *testing.T) {

	validAge, err := NewAgeRange(0, 10)
	if err != nil {
		t.Fatal(err)
	}

	if validAge.Min() != 0 {
		t.Errorf("unexpected minimum age\nwant: %d\n got: %d", 0, validAge.Min())
	}

	if validAge.Max() != 10 {
		t.Errorf("unexpected maximum age\nwant: %d\n got: %d", 0, validAge.Max())
	}

}

func TestAgeRange_Invalid(t *testing.T) {

	invalidAge := AgeRange{10, 0}
	err := invalidAge.Validate()
	if errors.Cause(err) != ErrBoundsReversed {
		t.Errorf("unexpected error\nwant: %v\n got: %d", ErrBoundsReversed, err)
	}

}

func TestAgeRange_Clone(t *testing.T) {

	ar, err := NewAgeRange(1, 5)
	if err != nil {
		t.Fatal(err)
	}

	clone := ar.Clone()
	if ar[0] != clone[0] || ar[1] != clone[1] {
		t.Errorf("unexpected results\nwant: %v\n got: %v", ar, clone)
	}
}
