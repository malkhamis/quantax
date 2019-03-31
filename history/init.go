package history

import "github.com/pkg/errors"

func init() {

	err := validateAllTaxFormulas()
	panicIfError(errors.Wrap(err, "invalid tax formula"))

	err = validateAllRRSPFormulas()
	panicIfError(errors.Wrap(err, "invalid RRSP formula"))

	err = validateAllChildBenefitFormulas()
	panicIfError(errors.Wrap(err, "invalid child benefit formula"))
}

func validateAllTaxFormulas() error {

	for jursdiction, formulasAllYears := range taxFormulasAll {
		for year, formula := range formulasAllYears {

			if formula == nil {
				return errors.New("history must not contain nil formulas")
			}

			err := formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}

func validateAllRRSPFormulas() error {

	for jursdiction, formulasAllYears := range rrspFormulasAll {
		for year, formula := range formulasAllYears {

			if formula == nil {
				return errors.New("history must not contain nil formulas")
			}

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

			if formula == nil {
				return errors.New("history must not contain nil formulas")
			}

			err := formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "jurisdiction %q, year %d", jursdiction, year)
			}

		}
	}

	return nil
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
