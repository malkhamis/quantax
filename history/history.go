// Package history provides historical tax params for the Canadian tax system
package history

import (
	"github.com/malkhamis/tax/facts"
	"github.com/pkg/errors"
)

func init() {

	err := validateFactsFederal()
	if err != nil {
		panic(err)
	}

	err = validateFactsProvincial()
	if err != nil {
		panic(err)
	}

}

// Get returns the facts for the given province in a specific tax year
func GetProvincial(year uint, province Province) (facts.FactsProv, error) {

	factsProvAllYears, ok := factsAllProvinces[province]
	if !ok {
		return facts.FactsProv{}, ErrNoProvince
	}

	factsProv, ok := factsProvAllYears[year]
	if !ok {
		err := errors.Wrap(ErrNoFacts, "provincial year's facts don't exist")
		return facts.FactsProv{}, err
	}

	return factsProv.Clone(), nil
}

// Get returns the facts for the given province in a specific tax year
func GetFederal(year uint) (facts.FactsFed, error) {

	factsFed, ok := factsFederal[year]
	if !ok {
		err := errors.Wrap(ErrNoFacts, "federal year's facts don't exist")
		return facts.FactsFed{}, err
	}

	return factsFed.Clone(), nil
}

func validateFactsFederal() error {

	for year, f := range factsFederal {

		err := f.Rates.Validate()
		if err != nil {
			return errors.Wrapf(err, "year %d", year)
		}

		err = f.Rates.Validate()
		if err != nil {
			return errors.Wrapf(err, "year %d", year)
		}
	}

	return nil
}

func validateFactsProvincial() error {

	for prov, factsAllYears := range factsAllProvinces {
		for year, f := range factsAllYears {

			err := f.Rates.Validate()
			if err != nil {
				return errors.Wrapf(err, "Province %q, year %d", prov, year)
			}

			err = f.Rates.Validate()
			if err != nil {
				return errors.Wrapf(err, "Province %q, year %d", prov, year)
			}

		}
	}

	return nil
}
