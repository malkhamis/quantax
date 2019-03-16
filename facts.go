package tax

import "math"

// FactsFed represents tax parameters set by the federal government
type FactsFed struct {
	Rates BracketRates
}

// FactsProv represents tax parameters set by the provincial government
type FactsProv struct {
	Rates BracketRates
}

// Facts represents tax parameters set by governments for a given tax year
type Facts struct {
	Year uint
	FactsFed
	FactsProv
}

func init() {

	for year, fact := range facts {

		if year != fact.Year {
			panic("facts' year mapped to a different year")
		}

		err := fact.FactsFed.Rates.Validate()
		if err != nil {
			panic(err)
		}

		err = fact.FactsProv.Rates.Validate()
		if err != nil {
			panic(err)
		}

	}
}

// facts is a convenience var that returns tax facts for a given year
var facts = map[uint]Facts{

	2018: Facts{
		Year: 2018,
		FactsFed: FactsFed{
			Rates: BracketRates{
				-0.150: Bracket{0, 11809},
				0.150:  Bracket{0, 46605},
				0.205:  Bracket{46606, 93208},
				0.260:  Bracket{93209, 144489},
				0.290:  Bracket{144490, 205842},
				0.330:  Bracket{205843, math.Inf(1)},
			},
		},
		FactsProv: FactsProv{
			Rates: BracketRates{
				-0.0506: Bracket{0, 10412},
				0.0506:  Bracket{0, 39676},
				0.0770:  Bracket{39677, 79353},
				0.1050:  Bracket{79354, 91107},
				0.1229:  Bracket{91108, 110630},
				0.1470:  Bracket{110631, 150000},
				0.1680:  Bracket{150001, math.Inf(1)},
			},
		},
	},
}
