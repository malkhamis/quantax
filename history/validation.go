package history

import "github.com/pkg/errors"

func validateAllFormulas() error {

	err := validateAllTaxFormulas()
	if err != nil {
		return errors.Wrap(err, "invalid tax formula")
	}

	err = validateAllChildBenefitFormulas()
	if err != nil {
		return errors.Wrap(err, "invalid child benefit formula")
	}

	return nil
}

func validateAllTaxFormulas() error {

	for jursdiction, formulasAllYears := range taxFormulasAll {
		for year, formula := range formulasAllYears {

			err := formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}

func validateAllChildBenefitFormulas() error {

	for jursdiction, formulasAllYears := range cbFormulasAll {
		for year, formula := range formulasAllYears {

			err := formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}
