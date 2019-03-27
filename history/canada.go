package history

import (
	"math"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
)

var (
	taxFormulasCanada = yearlyTaxFormulas{
		2018: taxFormulaCanada2018,
	}

	cbFormulasCanada = yearlyCBFormulas{
		2017: cbFormulaCanada2017,
	}
)

var taxFormulaCanada2018 = calc.WeightedBracketFormula{
	-0.150: calc.Bracket{0, 11809},
	0.150:  calc.Bracket{0, 46605},
	0.205:  calc.Bracket{46605, 93208},
	0.260:  calc.Bracket{93208, 144489},
	0.290:  calc.Bracket{144489, 205842},
	0.330:  calc.Bracket{205842, math.Inf(1)},
}

var cbFormulaCanada2017 = &benefits.MaxReducerFormula{
	BenefitClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			AgesMonths:      calc.AgeRange{0, (MonthsInYear * 6) - 1},
			AmountsPerMonth: calc.Bracket{0, 541.33},
		},
		benefits.AgeGroupBenefits{
			AgesMonths:      calc.AgeRange{MonthsInYear * 6, MonthsInYear * 17},
			AmountsPerMonth: calc.Bracket{0, 456.75},
		},
	},
	BenefitReducer: &benefits.StepReducer{
		StepFormulas: []calc.WeightedBracketFormula{
			calc.WeightedBracketFormula{}, // 0 child
			calc.WeightedBracketFormula{ // 1 child
				0.000: calc.Bracket{0, 30450},
				0.070: calc.Bracket{30450, 65976},
				0.032: calc.Bracket{65976, math.Inf(1)},
			},
			calc.WeightedBracketFormula{ // 2 children
				0.000: calc.Bracket{0, 30450},
				0.135: calc.Bracket{30450, 65976},
				0.057: calc.Bracket{65976, math.Inf(1)},
			},
			calc.WeightedBracketFormula{ // 3 children
				0.000: calc.Bracket{0, 30450},
				0.190: calc.Bracket{30450, 65976},
				0.080: calc.Bracket{65976, math.Inf(1)},
			},
		},
		AboveMaxStepFormula: calc.WeightedBracketFormula{ // 4+ children
			0.000: calc.Bracket{0, 30450},
			0.230: calc.Bracket{30450, 65976},
			0.095: calc.Bracket{65976, math.Inf(1)},
		},
	},
}
