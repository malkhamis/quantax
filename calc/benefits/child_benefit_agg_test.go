package benefits

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"

	"github.com/pkg/errors"
)

func TestNewChildBenefitAggregator_Error(t *testing.T) {

	_, err := NewChildBenefitAggregator(&BCECTBMaxReducer{}, &BCECTBMaxReducer{}, nil)
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected errpr\nwant: %v\n got: %v", ErrNoFormula, err)
	}
}

func TestNewChildBenefitAggregator_Full(t *testing.T) {

	formulaBC := &BCECTBMaxReducer{
		ReducerFormula: finance.WeightedBrackets{
			0.0132: finance.Bracket{100000, math.Inf(1)},
		},
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: finance.Bracket{0, 55},
			},
		},
		ExcludedIncome: []finance.IncomeSource{
			finance.IncSrcTFSA, finance.IncSrcRRSP,
		},
	}

	formulaCanada := &CCBMaxReducer{
		BeneficiaryClasses: []AgeGroupBenefits{
			AgeGroupBenefits{
				AgesMonths:      human.AgeRange{0, (12 * 6) - 1},
				AmountsPerMonth: finance.Bracket{0, 541.33},
			},
			AgeGroupBenefits{
				AgesMonths:      human.AgeRange{12 * 6, 12 * 17},
				AmountsPerMonth: finance.Bracket{0, 456.75},
			},
		},
		Reducers: []finance.WeightedBrackets{
			finance.WeightedBrackets{ // 1 child
				0.070: finance.Bracket{30450, 65976},
				0.032: finance.Bracket{65976, math.Inf(1)},
			},
			finance.WeightedBrackets{ // 2 children
				0.135: finance.Bracket{30450, 65976},
				0.057: finance.Bracket{65976, math.Inf(1)},
			},
			finance.WeightedBrackets{ // 3 children
				0.190: finance.Bracket{30450, 65976},
				0.080: finance.Bracket{65976, math.Inf(1)},
			},
			finance.WeightedBrackets{ // 4+ children
				0.230: finance.Bracket{30450, 65976},
				0.095: finance.Bracket{65976, math.Inf(1)},
			},
		},
		ExcludedIncome: []finance.IncomeSource{
			finance.IncSrcTFSA, finance.IncSrcRRSP,
		},
		ExcludedDeductions: []finance.DeductionSource{finance.DeducSrcUnknown},
	}

	calculator, err := NewChildBenefitAggregator(formulaCanada, formulaBC)
	if err != nil {
		t.Fatal(err)
	}

	children := []human.Person{{AgeMonths: 0}, {AgeMonths: (17 * 12) - 2}}
	calculator.SetBeneficiaries(children...)

	f1 := finance.NewEmptyIndividialFinances(2020)
	f1.AddIncome(finance.IncSrcEarned, 180000)
	f1.AddIncome(finance.IncSrcTFSA, 12000)
	f1.AddDeduction(finance.DeducSrcRRSP, 10000)

	f2 := finance.NewEmptyIndividialFinances(2020)
	f2.AddIncome(finance.IncSrcEarned, 20000)
	f2.AddIncome(finance.IncSrcRRSP, 5000)
	f2.AddDeduction(finance.DeducSrcRRSP, 20000)
	f2.AddDeduction(finance.DeducSrcUnknown, 1000)

	finances := finance.NewHouseholdFinances(f1, f2)
	actual := calculator.Calc(finances)

	expectedBC := 0.0
	expectedCanada := 0.0
	expected := expectedBC + expectedCanada
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.SetBeneficiaries()
	actual = calculator.Calc(finances)
	expected = 0.0
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	calculator.SetBeneficiaries(human.Person{AgeMonths: 0})
	actual = calculator.Calc(finance.HouseholdFinances{{}, {}})
	expected = (55 * 12) + (541.33 * 12)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	children = []human.Person{{AgeMonths: 0}, {AgeMonths: (6 * 12) - 2}}
	calculator.SetBeneficiaries(children...)
	finances = finance.HouseholdFinances{
		{
			Income: finance.IncomeBySource{
				finance.IncSrcEarned: 100000,
			},
		},
	}
	actual = calculator.Calc(finances)
	expectedCanada = (541.33*12 + 541.33*2 + 456.75*10) - (0.135 * (65976.0 - 30450.0)) - (0.057 * (100000.0 - 65976.0))
	expectedBC = (55 * 12) + (55 * 2)
	expected = expectedCanada + expectedBC
	if !areEqual(actual, expected, 1e-9) {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

// areEqual returns true if the difference between floor(actual) and
// floor(expected) is within the given +/- error margin of expcted. Negative
// error margins are converted to a positive number
func areEqual(actual, expected, errMargin float64) bool {

	actual, expected = math.Floor(actual), math.Floor(expected)
	allowedDiff := math.Abs(errMargin * expected)
	actualDiff := math.Abs(actual - expected)

	return actualDiff <= allowedDiff
}
