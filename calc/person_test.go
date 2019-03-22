package calc

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestPerson_IsOlderThan_Error(t *testing.T) {

	p := NewPerson("John", Date{1980, time.February, 1})
	_, err := p.IsOlderThan(100, Date{1970, time.January, 1})

	if errors.Cause(err) != ErrInvalidDate {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrInvalidDate, err)
	}

}

func TestPerson_IsOlderThan(t *testing.T) {

	p := NewPerson("John", Date{1980, time.January, 1})
	ok, err := p.IsOlderThan(24, Date{1980, time.December, 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("expected the child to not be older than perscribed date")
	}

	p = NewPerson("John", Date{1980, time.January, 1})
	ok, err = p.IsOlderThan(6, Date{1980, time.December, 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("expected the child to be older than perscribed date")
	}
}
