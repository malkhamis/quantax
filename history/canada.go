package history

import (
	"math"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/benefits"
	"github.com/malkhamis/quantax/core/human"
	"github.com/malkhamis/quantax/core/rrsp"
	"github.com/malkhamis/quantax/core/tax"
)

var (
	crDescPersonalAmount = tax.CreditDescriptor{
		CreditDescription:     "the basic personal amount",
		TargetFinancialSource: core.SrcNone,
		CreditRule: core.CreditRule{
			CrSource: "personal-amount",
			Type:     core.CrRuleTypeNotCarryForward,
		},
	}

	crDescTuitionAmount = tax.CreditDescriptor{
		CreditDescription:     "credits for paid university tuition fees",
		TargetFinancialSource: core.MiscSrcTuition,
		CreditRule: core.CreditRule{
			CrSource: "tuition-amount",
			Type:     core.CrRuleTypeCanCarryForward,
		},
	}

	crDescCanadianEligibleDividends = tax.CreditDescriptor{
		CreditDescription:     "credits for recieving Canadian-sourced eligible dividends",
		TargetFinancialSource: core.IncSrcEligibleDividendsCA,
		CreditRule: core.CreditRule{
			CrSource: "canadian-eligible-dividends",
			Type:     core.CrRuleTypeNotCarryForward,
		},
	}

	crDescCanadianNonEligibleDividends = tax.CreditDescriptor{
		CreditDescription:     "credits for recieving Canadian-sourced non-eligible dividends",
		TargetFinancialSource: core.IncSrcNonEligibleDividendsCA,
		CreditRule: core.CreditRule{
			CrSource: "canadian-non-eligible-dividends",
			Type:     core.CrRuleTypeNotCarryForward,
		},
	}

	crDescCanadianSpouse = tax.CreditDescriptor{
		CreditDescription:     "credits for having a Canadian spouse",
		TargetFinancialSource: core.SrcNone,
		CreditRule: core.CreditRule{
			CrSource: "canadian-spouse-credit",
			Type:     core.CrRuleTypeNotCarryForward,
		},
	}
)

var (
	taxParamsCanada = yearlyTaxParams{
		2022: TaxParams{
			Formula:       taxFormulaCanada2022,
			ContraFormula: taxContraFormulaCanada2022,
			IncomeRecipe:  incomeRecipeNetCA2022,
		},
		2019: TaxParams{
			Formula:       taxFormulaCanada2019,
			ContraFormula: taxContraFormulaCanada2019,
			IncomeRecipe:  incomeRecipeNetCA2019,
		},
		2018: TaxParams{
			Formula:       taxFormulaCanada2018,
			ContraFormula: taxContraFormulaCanada2018,
			IncomeRecipe:  incomeRecipeNetCA2018,
		},
	}

	cbParamsCanada = yearlyCBParams{
		2018: CBParams{cbFormulaCanada2018, incomeRecipeAFNICA2018},
	}

	rrspParamsCanada = yearlyRRSPParams{
		2019: RRSPParams{rrspFormulaCanada2019},
		2018: RRSPParams{rrspFormulaCanada2018},
	}
)

/* 2022 */

var taxFormulaCanada2022 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.1500: core.Bracket{0, 50197},
		0.2050: core.Bracket{50197, 100392},
		0.2600: core.Bracket{100392, 155625},
		0.2938: core.Bracket{155625, 221708},
		0.3300: core.Bracket{221708, math.Inf(1)},
	},
	TaxRegion: core.RegionCA,
	TaxYear:   2022,
}

var taxContraFormulaCanada2022 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.ConstCreditor{Amount: 0.150 * 14398, CreditDescriptor: crDescPersonalAmount},
		tax.CanadianSpouseCreditor{BaseAmount: 14398, Weight: 0.150, CreditDescriptor: crDescCanadianSpouse},
	},
	TaxYear:   2022,
	TaxRegion: core.RegionCA,
}
/* 2019 */

var taxFormulaCanada2019 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.150: core.Bracket{0, 47630},
		0.205: core.Bracket{47630, 95259},
		0.260: core.Bracket{95259, 147667},
		0.290: core.Bracket{147667, 210371},
		0.330: core.Bracket{210371, math.Inf(1)},
	},
	TaxRegion: core.RegionCA,
	TaxYear:   2019,
}

var taxContraFormulaCanada2019 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.ConstCreditor{Amount: 0.150 * 12069, CreditDescriptor: crDescPersonalAmount},
		tax.CanadianSpouseCreditor{BaseAmount: 12069, Weight: 0.150, CreditDescriptor: crDescCanadianSpouse},
		tax.WeightedCreditor{Weight: 0.150, CreditDescriptor: crDescTuitionAmount},
		tax.WeightedCreditor{Weight: 1.38 * 0.150198, CreditDescriptor: crDescCanadianEligibleDividends},
		tax.WeightedCreditor{Weight: 1.15 * 0.090301, CreditDescriptor: crDescCanadianNonEligibleDividends},
	},
	TaxYear:   2019,
	TaxRegion: core.RegionCA,
}

/* 2018 */

var taxFormulaCanada2018 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.150: core.Bracket{0, 46605},
		0.205: core.Bracket{46605, 93208},
		0.260: core.Bracket{93208, 144489},
		0.290: core.Bracket{144489, 205842},
		0.330: core.Bracket{205842, math.Inf(1)},
	},
	TaxRegion: core.RegionCA,
	TaxYear:   2018,
}

var taxContraFormulaCanada2018 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.ConstCreditor{Amount: 0.150 * 11809, CreditDescriptor: crDescPersonalAmount},
		tax.CanadianSpouseCreditor{BaseAmount: 11809, Weight: 0.150, CreditDescriptor: crDescCanadianSpouse},
		tax.WeightedCreditor{Weight: 0.150, CreditDescriptor: crDescTuitionAmount},
		tax.WeightedCreditor{Weight: 1.38 * 0.150198, CreditDescriptor: crDescCanadianEligibleDividends},
		tax.WeightedCreditor{Weight: 1.16 * 0.100313, CreditDescriptor: crDescCanadianNonEligibleDividends},
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
