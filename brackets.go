package tax

import (
	"math"

	"github.com/pkg/errors"
)

// Bracket represents a tax bracket, e.g. [47630, 95259]
type Bracket [2]float64

// Validate ensures that this bracket is valid for use in this package
func (b *Bracket) Validate() error {

	switch {

	case math.IsInf(b.Lower(), 0):
		return errors.Wrap(ErrValInf, "lower-bound value")

	case math.IsInf(b.Upper(), -1):
		return errors.Wrap(ErrValInfNeg, "upper-bound value")

	case b.Lower() < 0:
		return errors.Wrapf(ErrValNeg, "lower-bound value '%f02'", b.Lower())

	case b.Upper() < 0:
		return errors.Wrapf(ErrValNeg, "upper-bound value '%f02'", b.Upper())

	case b.Lower() == 0 && b.Upper() == 0:
		return errors.Wrap(ErrValZero, "upper-bound and lower-bound values")

	case b.Lower() > b.Upper():
		return errors.Errorf(
			"lower-bound is greater than upper-bound [%f02 > %f02]",
			b.Lower(), b.Upper(),
		)

	default:
		return nil
	}

}

// Amount returns the difference between upper and lower bound
func (b *Bracket) Amount() float64 {
	return b.Upper() - b.Lower()
}

// Lower returns the lower limit/bound of this bracket
func (b *Bracket) Lower() float64 {
	return b[0]
}

// Upper returns the upper limit/bound of this bracket
func (b *Bracket) Upper() float64 {
	return b[1]
}

// BracketRates maps tax rates to their income brackets. Negative rates can
// be used to indicate tax credits (e.g. personal basic amount, line 300)
type BracketRates map[float64]Bracket

// Validate ensures that these brackets are valid for use in this package
func (b BracketRates) Validate() error {

	for rate, bracket := range b {

		if math.IsInf(rate, 0) {
			return errors.Wrap(ErrValInf, "invalid tax rate")
		}

		if err := bracket.Validate(); err != nil {
			return errors.Wrap(err, "invalid tax bracket")
		}

	}

	return nil
}
