package history

import (
	"math"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/benefits"
	"github.com/malkhamis/quantax/core/human"
	"github.com/malkhamis/quantax/core/rrsp"
	"github.com/malkhamis/quantax/core/tax"
)

const (
	crSrcPersonalAmount = "basic-personal-amount"
)

var (
	taxParamsCanada = yearlyTaxParams{
		2019: TaxParams{
			Formula:       taxFormulaCanada2019,
			ContraFormula: taxContraFormulaCanada2019,
			IncomeRecipe:  incomeRecipeNet,
		},
		2018: TaxParams{
			Formula:       taxFormulaCanada2018,
			ContraFormula: taxContraFormulaCanada2018,
			IncomeRecipe:  incomeRecipeNet,
		},
	}

	cbParamsCanada = yearlyCBParams{
		2018: CBParams{cbFormulaCanada2018, incomeRecipeAFNI},
	}

	rrspParamsCanada = yearlyRRSPParams{
		2019: RRSPParams{rrspFormulaCanada2019},
		2018: RRSPParams{rrspFormulaCanada2018},
	}
)

var rrspFormulaCanada2019 = &rrsp.MaxCapper{
	Rate:                           0.18,
	Cap:                            26500,
	IncomeSources:                  []core.FinancialSource{core.IncSrcEarned},
	IncomeSourceForWithdrawal:      core.IncSrcRRSP,
	DeductionSourceForContribution: core.DeducSrcRRSP,
}

var rrspFormulaCanada2018 = &rrsp.MaxCapper{
	Rate:                           0.18,
	Cap:                            26230.00,
	IncomeSources:                  []core.FinancialSource{core.IncSrcEarned},
	IncomeSourceForWithdrawal:      core.IncSrcRRSP,
	DeductionSourceForContribution: core.DeducSrcRRSP,
}

var taxFormulaCanada2019 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.150: core.Bracket{0, 47630},
		0.205: core.Bracket{47630, 95259},
		0.260: core.Bracket{95259, 147667},
		0.290: core.Bracket{147667, 210371},
		0.330: core.Bracket{210371, math.Inf(1)},
	},
}

var taxFormulaCanada2018 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.150: core.Bracket{0, 46605},
		0.205: core.Bracket{46605, 93208},
		0.260: core.Bracket{93208, 144489},
		0.290: core.Bracket{144489, 205842},
		0.330: core.Bracket{205842, math.Inf(1)},
	},
}

var taxContraFormulaCanada2019 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.CreditorConst{
			Amount:                0.150 * 12069,
			CreditDescription:     crSrcPersonalAmount,
			TargetFinancialSource: core.SrcNone,
			CreditRule: core.CreditRule{
				CrSource: crSrcPersonalAmount,
				Type:     core.CrRuleTypeNotCarryForward,
			},
		},
	},
	TaxYear:   2019,
	TaxRegion: core.RegionCA,
}

var taxContraFormulaCanada2018 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.CreditorConst{
			Amount:                0.150 * 11809,
			CreditDescription:     crSrcPersonalAmount,
			TargetFinancialSource: core.SrcNone,
			CreditRule: core.CreditRule{
				CrSource: crSrcPersonalAmount,
				Type:     core.CrRuleTypeNotCarryForward,
			},
		},
	},
	TaxYear:   2018,
	TaxRegion: core.RegionCA,
}

var cbFormulaCanada2018 = &benefits.CCBMaxReducer{
	BeneficiaryClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{0, (monthsInYear * 6) - 1},
			AmountsPerMonth: core.Bracket{0, 541.33},
		},
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{monthsInYear * 6, monthsInYear * 17},
			AmountsPerMonth: core.Bracket{0, 456.75},
		},
	},
	Reducers: []core.WeightedBrackets{
		core.WeightedBrackets{ // 1 child
			0.070: core.Bracket{30450, 65976},
			0.032: core.Bracket{65976, math.Inf(1)},
		},
		core.WeightedBrackets{ // 2 children
			0.135: core.Bracket{30450, 65976},
			0.057: core.Bracket{65976, math.Inf(1)},
		},
		core.WeightedBrackets{ // 3 children
			0.190: core.Bracket{30450, 65976},
			0.080: core.Bracket{65976, math.Inf(1)},
		},
		core.WeightedBrackets{ // 4+ children
			0.230: core.Bracket{30450, 65976},
			0.095: core.Bracket{65976, math.Inf(1)},
		},
	},
}
