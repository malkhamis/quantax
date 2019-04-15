package history

import (
	"math"

	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
	"github.com/malkhamis/quantax/calc/rrsp"
	"github.com/malkhamis/quantax/calc/tax"
)

var (
	taxParamsCanada = yearlyTaxParams{
		2018: TaxParams{
			Formula:       taxFormulaCanada2018,
			ContraFormula: taxContraFormulaCanada2018,
			IncomeRecipe:  &incomeRecipeNet,
		},
	}

	cbParamsCanada = yearlyCBParams{
		2017: CBParams{cbFormulaCanada2017, &incomeRecipeAFNI},
	}

	rrspParamsCanada = yearlyRRSPParams{
		2018: RRSPParams{rrspFormulaCanada2018},
	}
)

var rrspFormulaCanada2018 = &rrsp.MaxCapper{
	Rate:                           0.18,
	Cap:                            26230.00,
	IncomeSources:                  []finance.IncomeSource{finance.IncSrcEarned},
	IncomeSourceForWithdrawal:      finance.IncSrcRRSP,
	DeductionSourceForContribution: finance.DeducSrcRRSP,
}

var taxFormulaCanada2018 = &tax.CanadianFormula{
	WeightedBrackets: finance.WeightedBrackets{
		-0.150: finance.Bracket{0, 11809},
		0.150:  finance.Bracket{0, 46605},
		0.205:  finance.Bracket{46605, 93208},
		0.260:  finance.Bracket{93208, 144489},
		0.290:  finance.Bracket{144489, 205842},
		0.330:  finance.Bracket{205842, math.Inf(1)},
	},
}

// TODO
var taxContraFormulaCanada2018 = tax.NopContraFormula{}

var cbFormulaCanada2017 = &benefits.CCBMaxReducer{
	BeneficiaryClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{0, (monthsInYear * 6) - 1},
			AmountsPerMonth: finance.Bracket{0, 541.33},
		},
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{monthsInYear * 6, monthsInYear * 17},
			AmountsPerMonth: finance.Bracket{0, 456.75},
		},
	},
	Reducers: []finance.WeightedBrackets{
		finance.WeightedBrackets{ // 1 child
			0.070: finance.Bracket{30450, 65976},
			0.032: finance.Bracket{65976, math.Inf(1)},
		},
		finance.WeightedBrackets{ // 2 children
			0.135: finance.Bracket{30450, 65976},
			0.057: finance.Bracket{65976, math.Inf(1)},
		},
		finance.WeightedBrackets{ // 3 children
			0.190: finance.Bracket{30450, 65976},
			0.080: finance.Bracket{65976, math.Inf(1)},
		},
		finance.WeightedBrackets{ // 4+ children
			0.230: finance.Bracket{30450, 65976},
			0.095: finance.Bracket{65976, math.Inf(1)},
		},
	},
}
