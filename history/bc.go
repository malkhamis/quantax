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
		2019: TaxParams{
			Formula:       taxFormulaBC2019,
			ContraFormula: taxContraFormulaBC2019,
			IncomeRecipe:  incomeRecipeNetCA2019,
		},
		2018: TaxParams{
			Formula:       taxFormulaBC2018,
			ContraFormula: taxContraFormulaBC2018,
			IncomeRecipe:  incomeRecipeNetCA2018,
		},
	}

	cbParamsBC = yearlyCBParams{
		2018: CBParams{cbFormulaBC2018, incomeRecipeAFNICA2018},
	}
)

var taxFormulaBC2019 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.0506: core.Bracket{0, 40707},
		0.0770: core.Bracket{40707, 81416},
		0.1050: core.Bracket{81416, 93476},
		0.1229: core.Bracket{93476, 113506},
		0.1470: core.Bracket{113506, 153900},
		0.1680: core.Bracket{153900, math.Inf(1)},
	},
	TaxRegion: core.RegionBC,
	TaxYear:   2019,
}

var taxContraFormulaBC2019 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.ConstCreditor{Amount: 0.0506 * 10682, CreditDescriptor: crDescPersonalAmount},
		tax.CanadianSpouseCreditor{BaseAmount: 9147, Weight: 0.0506, CreditDescriptor: crDescCanadianSpouse},
		tax.WeightedCreditor{Weight: 0.0506, CreditDescriptor: crDescTuitionAmount},
		tax.WeightedCreditor{Weight: 1.38 * 0.12, CreditDescriptor: crDescCanadianEligibleDividends},
		tax.WeightedCreditor{Weight: 1.15 * 0.0196, CreditDescriptor: crDescCanadianNonEligibleDividends},
	},
	TaxYear:   2019,
	TaxRegion: core.RegionBC,
}

var taxFormulaBC2018 = &tax.CanadianFormula{
	WeightedBrackets: core.WeightedBrackets{
		0.0506: core.Bracket{0, 39676},
		0.0770: core.Bracket{39676, 79353},
		0.1050: core.Bracket{79353, 91107},
		0.1229: core.Bracket{91107, 110630},
		0.1470: core.Bracket{110630, 150000},
		0.1680: core.Bracket{150000, math.Inf(1)},
	},
	TaxRegion: core.RegionBC,
	TaxYear:   2018,
}

var taxContraFormulaBC2018 = &tax.CanadianContraFormula{
	OrderedCreditors: []tax.Creditor{
		tax.ConstCreditor{Amount: 0.0506 * 10412, CreditDescriptor: crDescPersonalAmount},
		tax.CanadianSpouseCreditor{BaseAmount: 8915, Weight: 0.0506, CreditDescriptor: crDescCanadianSpouse},
		tax.WeightedCreditor{Weight: 0.0506, CreditDescriptor: crDescTuitionAmount},
		tax.WeightedCreditor{Weight: 1.38 * 0.10, CreditDescriptor: crDescCanadianEligibleDividends},
		tax.WeightedCreditor{Weight: 1.16 * 0.0207, CreditDescriptor: crDescCanadianNonEligibleDividends},
	},
	TaxYear:   2018,
	TaxRegion: core.RegionBC,
}

var cbFormulaBC2018 = &benefits.BCECTBMaxReducer{
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
