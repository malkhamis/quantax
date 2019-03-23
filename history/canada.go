package history

import (
	"math"
)

var taxFormulasCanada = yearlyTaxFormulas{
	2018: WeightedBracketFormula{
		-0.150: Bracket{0, 11809},
		0.150:  Bracket{0, 46605},
		0.205:  Bracket{46606, 93208},
		0.260:  Bracket{93209, 144489},
		0.290:  Bracket{144490, 205842},
		0.330:  Bracket{205843, math.Inf(1)},
	},
}
