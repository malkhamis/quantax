package history

import (
	"math"

	"github.com/malkhamis/quantax/calc"
)

var ratesCanada = historicalRates{
	2018: calc.BracketRates{
		-0.150: calc.Bracket{0, 11809},
		0.150:  calc.Bracket{0, 46605},
		0.205:  calc.Bracket{46606, 93208},
		0.260:  calc.Bracket{93209, 144489},
		0.290:  calc.Bracket{144490, 205842},
		0.330:  calc.Bracket{205843, math.Inf(1)},
	},
}
