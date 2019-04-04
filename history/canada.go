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
	taxFormulasCanada = yearlyTaxFormulas{
		2018: taxFormulaCanada2018,
	}

	cbFormulasCanada = yearlyCBFormulas{
		2017: cbFormulaCanada2017,
	}

	rrspFormulasCanada = yearlyRRSPFormulas{
		2018: rrspFormulaCanada2018,
	}
)

var rrspFormulaCanada2018 = &rrsp.MaxCapper{
	Rate: 0.18,
	Cap:  26230.00,
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
	ExcludedIncome: []finance.IncomeSource{finance.IncSrcTFSA},
}

var cbFormulaCanada2017 = &benefits.CCBMaxReducer{
	BeneficiaryClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{0, (MonthsInYear * 6) - 1},
			AmountsPerMonth: finance.Bracket{0, 541.33},
		},
		benefits.AgeGroupBenefits{
			AgesMonths:      human.AgeRange{MonthsInYear * 6, MonthsInYear * 17},
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
	ExcludedIncome: []finance.IncomeSource{
		finance.IncSrcTFSA,
		finance.IncSrcUCCB,
		finance.IncSrcRDSP,
	},
}
