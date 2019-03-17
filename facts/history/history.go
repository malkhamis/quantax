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

// Get returns a clone of the federal and the given provincial facts for a
// specific tax year
func Get(year uint, province Province) (facts.Facts, error) {

	factsFed, ok := factsFederal[year]
	if !ok {
		return facts.Facts{}, errors.Wrap(ErrNoFacts, "federal year's facts don't exist")
	}

	factsProvAllYears, ok := factsAllProvinces[province]
	if !ok {
		return facts.Facts{}, ErrNoProvince
	}

	factsProv, ok := factsProvAllYears[year]
	if !ok {
		return facts.Facts{}, errors.Wrap(ErrNoFacts, "provincial year's facts don't exist")
	}

	f := facts.Facts{
		Year:      year,
		FactsFed:  factsFed.Clone(),
		FactsProv: factsProv.Clone(),
	}
	return f, nil
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
