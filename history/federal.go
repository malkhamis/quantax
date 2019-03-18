package history

import (
	"math"

	"github.com/malkhamis/tax/facts"
)

var factsFederal = map[uint]facts.FactsFed{
	2018: facts.FactsFed{
		Rates: facts.BracketRates{
			-0.150: facts.Bracket{0, 11809},
			0.150:  facts.Bracket{0, 46605},
			0.205:  facts.Bracket{46606, 93208},
			0.260:  facts.Bracket{93209, 144489},
			0.290:  facts.Bracket{144490, 205842},
			0.330:  facts.Bracket{205843, math.Inf(1)},
		},
	},
}
