package facts

import "testing"

func testFacts(t *testing.T) Facts {

	t.Helper()

	newFacts := Facts{
		Year: 2018,
		FactsFed: FactsFed{
			BracketRates{
				0.10: Bracket{100, 200},
			},
		},
		FactsProv: FactsProv{
			BracketRates{
				0.20: Bracket{1000, 2000},
			},
		},
	}

	return newFacts
}
