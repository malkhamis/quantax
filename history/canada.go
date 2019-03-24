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
		2018: cbFormulaCanada2018,
	}
)

var taxFormulaCanada2018 = calc.WeightedBracketFormula{
	-0.150: calc.Bracket{0, 11809},
	0.150:  calc.Bracket{0, 46605},
	0.205:  calc.Bracket{46606, 93208},
	0.260:  calc.Bracket{93209, 144489},
	0.290:  calc.Bracket{144490, 205842},
	0.330:  calc.Bracket{205843, math.Inf(1)},
}

var cbFormulaCanada2018 = &benefits.CCBFormula{
	BenefitClasses: []benefits.AgeGroupBenefits{
		benefits.AgeGroupBenefits{
			Ages:    calc.AgeRange{0, (MonthsInYear * 6) - 1},
			Amounts: calc.Bracket{0, 541.33},
		},
		benefits.AgeGroupBenefits{
			Ages:    calc.AgeRange{MonthsInYear * 6, MonthsInYear * 17},
			Amounts: calc.Bracket{0, 456.75},
		},
	},
	Reducers: &benefits.StepReducer{
		StepFormulas: []calc.WeightedBracketFormula{
			calc.WeightedBracketFormula{ // no children
				0.000: calc.Bracket{0.00, math.Inf(1)},
			},
			calc.WeightedBracketFormula{ // 1 child
				0.000: calc.Bracket{0.00, 30450.00},
				0.070: calc.Bracket{30450.01, 65976.00},
				0.032: calc.Bracket{65976.01, math.Inf(1)},
			},
			calc.WeightedBracketFormula{ // 2 children
				0.000: calc.Bracket{0.00, 30450.00},
				0.135: calc.Bracket{30450.01, 65976.00},
				0.057: calc.Bracket{65976.01, math.Inf(1)},
			},
			calc.WeightedBracketFormula{ // 3 children
				0.000: calc.Bracket{0.00, 30450.00},
				0.190: calc.Bracket{30450.01, 65976.00},
				0.080: calc.Bracket{65976.01, math.Inf(1)},
			},
		},
		AboveMaxStepFormula: calc.WeightedBracketFormula{ // 4+ children
			0.000: calc.Bracket{0.00, 30450.00},
			0.230: calc.Bracket{30450.01, 65976.00},
			0.095: calc.Bracket{65976.01, math.Inf(1)},
		},
	},
}
