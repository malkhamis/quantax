package history

import (
	"math"

	"github.com/malkhamis/quantax/calc"
)

var taxFormulasBC = yearlyTaxFormulas{
	2018: calc.WeightedBracketFormula{
		-0.0506: calc.Bracket{0, 10412},
		0.0506:  calc.Bracket{0, 39676},
		0.0770:  calc.Bracket{39676, 79353},
		0.1050:  calc.Bracket{79353, 91107},
		0.1229:  calc.Bracket{91107, 110630},
		0.1470:  calc.Bracket{110630, 150000},
		0.1680:  calc.Bracket{150000, math.Inf(1)},
	},
}
