package history

import (
	"math"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/benefits"
	"github.com/malkhamis/quantax/core/human"
	"github.com/malkhamis/quantax/core/tax"
)

var (
	taxParamsBC = yearlyTaxParams{
		2018: TaxParams{
			Formula:       taxFormulaBC2018,
			ContraFormula: taxContraFormulaBC2018,
			IncomeRecipe:  &incomeRecipeNet,
		},
	}

	cbParamsBC = yearlyCBParams{
		2017: CBParams{cbFormulaBC2017, &incomeRecipeNet},
	}
)

var taxFormulaBC2018 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.0506: core.Bracket{0, 39676},
		0.0770: core.Bracket{39676, 79353},
		0.1050: core.Bracket{79353, 91107},
		0.1229: core.Bracket{91107, 110630},
		0.1470: core.Bracket{110630, 150000},
		0.1680: core.Bracket{150000, math.Inf(1)},
	},
}

const (
	crSrcPersonalAmountBC = "BC-basic-personal-amount"
)

var taxContraFormulaBC2018 = &tax.CanadianContraFormula{
	PersistentCredits: map[string]float64{
		crSrcPersonalAmountBC: 0.0506 * 10412,
	},
	ApplicationOrder: []tax.CreditRule{
		{
			Source: crSrcPersonalAmountBC,
			Type:   tax.CrRuleTypeNotCarryForward,
		},
	},
}

var cbFormulaBC2017 = &benefits.BCECTBMaxReducer{
	BeneficiaryClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{0, (monthsInYear * 6) - 1},
			AmountsPerMonth: core.Bracket{0, 55},
		},
	},
	ReducerFormula: core.WeightedBrackets{
		0.0132: core.Bracket{100000, math.Inf(1)},
	},
}
