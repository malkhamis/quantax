package history

import (
	"math"

	"github.com/malkhamis/tax/facts"
)

var factsBC = map[uint]facts.FactsProv{
	2018: facts.FactsProv{
		Rates: facts.BracketRates{
			-0.0506: facts.Bracket{0, 10412},
			0.0506:  facts.Bracket{0, 39676},
			0.0770:  facts.Bracket{39677, 79353},
			0.1050:  facts.Bracket{79354, 91107},
			0.1229:  facts.Bracket{91108, 110630},
			0.1470:  facts.Bracket{110631, 150000},
			0.1680:  facts.Bracket{150001, math.Inf(1)},
		},
	},
}
