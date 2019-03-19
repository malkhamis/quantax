// Package history provides historical tax params for various jurisdictions
package history

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type historicalRates = map[uint]calc.BracketRates

var ratesAll = map[Jurisdiction]historicalRates{
	BC:     ratesBC,
	Canada: ratesCanada,
}

func init() {

	err := validateAllRates()
	if err != nil {
		panic(err)
	}

}

// Get returns the rates for the given jurisdiction in a specific tax year
func Get(year uint, region Jurisdiction) (calc.BracketRates, error) {

	jurisdictionRates, ok := ratesAll[region]
	if !ok {
		return calc.BracketRates{}, ErrNoJurisdiction
	}

	rates, ok := jurisdictionRates[year]
	if !ok {
		err := errors.Wrap(ErrNoRates, "jurisdiction year's rates don't exist")
		return calc.BracketRates{}, err
	}

	return rates.Clone(), nil
}

func validateAllRates() error {

	for jursdiction, ratesAllYears := range ratesAll {
		for year, rates := range ratesAllYears {

			err := rates.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}
