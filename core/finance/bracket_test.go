package finance

import (
	"fmt"
	"math"
	"testing"

	"github.com/pkg/errors"
)

func TestBracket_Validate(t *testing.T) {

	cases := []struct {
		name    string
		bracket Bracket
		err     error
	}{
		{
			name:    "lower-bound-greater",
			bracket: Bracket{2000, 1000},
			err:     ErrBoundsReversed,
		},
		{
			name:    "lower-bound-neg-inf",
			bracket: Bracket{math.Inf(-1), 1000},
			err:     nil,
		},
		{
			name:    "lower-bound-pos-inf",
			bracket: Bracket{math.Inf(1), 1000},
			err:     ErrBoundsReversed,
		},
		{
			name:    "upper-bound-neg-inf",
			bracket: Bracket{1000, math.Inf(-1)},
			err:     ErrBoundsReversed,
		},
		{
			name:    "upper-bound-pos-inf",
			bracket: Bracket{1000, math.Inf(1)},
			err:     nil,
		},
		{
			name:    "upper-bound-neg",
			bracket: Bracket{1000, -2000},
			err:     ErrBoundsReversed,
		},
		{
			name:    "lower-bound-neg",
			bracket: Bracket{-1000, 2000},
			err:     nil,
		},
		{
			name:    "bounds-zero",
			bracket: Bracket{0, 0},
			err:     nil,
		},
		{
			name:    "upper-bound-zero",
			bracket: Bracket{1000, 0},
			err:     ErrBoundsReversed,
		},
		{
			name:    "simple",
			bracket: Bracket{1000, 2000},
			err:     nil,
		},
	}

	for i, c := range cases {
		c, i := c, i
		t.Run(fmt.Sprintf("case-%d-%s", i, c.name), func(t *testing.T) {

			err := c.bracket.Validate()
			cause := errors.Cause(err)
			if cause != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}
		})
	}
}

func TestBracket_Amount(t *testing.T) {

	amnt := (Bracket{500, 750}).Amount()
	if amnt != 250 {
		t.Fatalf("expected bracket amount to be 250, got: %.2f", amnt)
	}
}

func TestBracket_Lower_Upper(t *testing.T) {

	b := Bracket{500, 750}
	if b.Lower() != 500 {
		t.Fatalf("expected bracket lower bound to be 500, got: %.2f", b.Lower())
	}
	if b.Upper() != 750 {
		t.Fatalf("expected bracket lower bound to be 750, got: %.2f", b.Upper())
	}
}
