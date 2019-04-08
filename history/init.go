package history

import "github.com/pkg/errors"

func init() {

	err := validateAllTaxParams()
	panicIfError(errors.Wrap(err, "invalid tax params"))

	err = validateAllRRSPParams()
	panicIfError(errors.Wrap(err, "invalid RRSP params"))

	err = validateAllCBParams()
	panicIfError(errors.Wrap(err, "invalid child benefit params"))

	initIncomeRecipes()
}

func validateAllTaxParams() error {

	for jursdiction, paramsAllYears := range taxParamsAll {
		for year, params := range paramsAllYears {

			if params.Formula == nil {
				return errors.Wrapf(errNilFormula, "%s[%d]", jursdiction, year)
			}

			err := params.Formula.Validate()
			if err != nil {
				return errors.Wrapf(err, "%s[%d]", jursdiction, year)
			}

		}
	}

	return nil
}

func validateAllRRSPParams() error {

	for jursdiction, paramsAllYears := range rrspParamsAll {
		for year, params := range paramsAllYears {

			if params.Formula == nil {
				return errors.Wrapf(errNilFormula, "%s[%d]", jursdiction, year)
			}

			err := params.Formula.Validate()
			if err != nil {
				return errors.Wrapf(errNilFormula, "%s[%d]", jursdiction, year)
			}

		}
	}

	return nil
}

func validateAllCBParams() error {

	for jursdiction, paramsAllYears := range cbParamsAll {
		for year, params := range paramsAllYears {

			if params.Formula == nil {
				return errors.Wrapf(errNilFormula, "%s[%d]", jursdiction, year)
			}

			err := params.Formula.Validate()
			if err != nil {
				return errors.Wrapf(errNilFormula, "%s[%d]", jursdiction, year)
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
