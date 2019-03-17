package facts

import (
	"math"

	"github.com/pkg/errors"
)

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

// Clone returns a copy of this bracket rates object
func (b BracketRates) Clone() BracketRates {

	clone := make(BracketRates)

	for rate, bracket := range b {
		clone[rate] = bracket.Clone()
	}

	return clone
}
