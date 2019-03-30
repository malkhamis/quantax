package human

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
	if errors.Cause(err) != ErrInvalidAgeRange {
		t.Errorf("unexpected error\nwant: %v\n got: %d", ErrInvalidAgeRange, err)
	}

}
