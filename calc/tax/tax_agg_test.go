package tax

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc"
)

func TestCalculatorAgg_Calc(t *testing.T) {

	canada2018 := CanadianFormula{
		-0.150: calc.Bracket{0, 11809},
		0.150:  calc.Bracket{0, 46605},
		0.205:  calc.Bracket{46606, 93208},
		0.260:  calc.Bracket{93209, 144489},
		0.290:  calc.Bracket{144490, 205842},
		0.330:  calc.Bracket{205843, math.Inf(1)},
	}

	bc2018 := CanadianFormula{
		-0.0506: calc.Bracket{0, 10412},
		0.0506:  calc.Bracket{0, 39676},
		0.0770:  calc.Bracket{39676, 79353},
		0.1050:  calc.Bracket{79353, 91107},
		0.1229:  calc.Bracket{91107, 110630},
		0.1470:  calc.Bracket{110630, 150000},
		0.1680:  calc.Bracket{150000, math.Inf(1)},
	}

	agg, err := NewCalculatorAgg(
		calc.IndividualFinances{Income: 250000.0},
		[]Formula{canada2018, bc2018},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedTax := 91226.32
	actualTax := agg.Calc()
	if !areEqual(actualTax, expectedTax, 1e-6) {
		t.Errorf(
			"difference between actual and expected total "+
				"tax exceeds error margin\nwant: %04f\n got: %04f",
			expectedTax, actualTax,
		)
	}

	agg.UpdateFinances(calc.IndividualFinances{Income: 0})
	actualTax = agg.Calc()
	if actualTax != 0.0 {
		t.Errorf(
			"difference between actual and expected total "+
				"tax exceeds error margin\nwant: %04f\n got: %04f",
			0.0, actualTax,
		)
	}

}

func TestNewCalculatorAgg_Error(t *testing.T) {

	_, err := NewCalculatorAgg(calc.IndividualFinances{}, nil)
	if err != calc.ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}

	_, err = NewCalculatorAgg(calc.IndividualFinances{}, []Formula{nil})
	if err != calc.ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", calc.ErrNoFormula, err)
	}

}
