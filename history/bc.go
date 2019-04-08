package history

import (
	"math"

	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
	"github.com/malkhamis/quantax/calc/tax"
)

var (
	taxParamsBC = yearlyTaxParams{
		2018: TaxParams{taxFormulaBC2018, incomeRecipeNet},
	}

	cbParamsBC = yearlyCBParams{
		2017: CBParams{cbFormulaBC2017, incomeRecipeNet},
	}
)

var taxFormulaBC2018 = &tax.CanadianFormula{
	WeightedBrackets: finance.WeightedBrackets{
		-0.0506: finance.Bracket{0, 10412},
		0.0506:  finance.Bracket{0, 39676},
		0.0770:  finance.Bracket{39676, 79353},
		0.1050:  finance.Bracket{79353, 91107},
		0.1229:  finance.Bracket{91107, 110630},
		0.1470:  finance.Bracket{110630, 150000},
		0.1680:  finance.Bracket{150000, math.Inf(1)},
	},
}

var cbFormulaBC2017 = &benefits.BCECTBMaxReducer{
	BeneficiaryClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{0, (monthsInYear * 6) - 1},
			AmountsPerMonth: finance.Bracket{0, 55},
		},
	},
	ReducerFormula: finance.WeightedBrackets{
		0.0132: finance.Bracket{100000, math.Inf(1)},
	},
}
