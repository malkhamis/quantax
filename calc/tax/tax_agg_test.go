package tax

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestCalculatorAgg_Calc(t *testing.T) {

	canada2018 := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			-0.150: finance.Bracket{0, 11809},
			0.150:  finance.Bracket{0, 46605},
			0.205:  finance.Bracket{46606, 93208},
			0.260:  finance.Bracket{93209, 144489},
			0.290:  finance.Bracket{144490, 205842},
			0.330:  finance.Bracket{205843, math.Inf(1)},
		},
	}

	bc2018 := &CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			-0.0506: finance.Bracket{0, 10412},
			0.0506:  finance.Bracket{0, 39676},
			0.0770:  finance.Bracket{39676, 79353},
			0.1050:  finance.Bracket{79353, 91107},
			0.1229:  finance.Bracket{91107, 110630},
			0.1470:  finance.Bracket{110630, 150000},
			0.1680:  finance.Bracket{150000, math.Inf(1)},
		},
	}

	agg, err := NewAggregator(canada2018, bc2018)
	if err != nil {
		t.Fatal(err)
	}

	expectedTax := 91226.32
	finances := finance.IndividualFinances{
		Income: finance.IncomeBySource{finance.IncSrcEarned: 250000.0},
	}
	actualTax := agg.Calc(finances)
	if !areEqual(actualTax, expectedTax, 1e-6) {
		t.Errorf(
			"difference between actual and expected total "+
				"tax exceeds error margin\nwant: %04f\n got: %04f",
			expectedTax, actualTax,
		)
	}

}

func TestNewAggregator_Error(t *testing.T) {

	_, err := NewAggregator(nil, nil)
	if err != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

}
