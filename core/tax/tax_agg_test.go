package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

func TestNewAggregator(t *testing.T) {

	c0, c1, c2 := &Calculator{taxYear: 2000}, &Calculator{}, &Calculator{}
	_, err := NewAggregator(c0, c1, c2)
	if errors.Cause(err) != ErrTooManyYears {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrTooManyYears, err)
	}

	_, err = NewAggregator(nil, nil, nil)
	if errors.Cause(err) != ErrNoCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCalc, err)
	}

	c0, c1, c2 = &Calculator{}, &Calculator{}, &Calculator{}
	_, err = NewAggregator(c0, c1, c2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestAggregator(t *testing.T) {

	finances := core.NewHouseholdFinancesNop()

	incCalc := &testIncomeCalculator{onTotalIncome: 3000.0}
	formula := &testTaxFormula{onApply: incCalc.TotalIncome() / 2.0}

	cfg := CalcConfig{
		TaxFormula: formula,
		ContraTaxFormula: &testContraTaxFormula{
			onApply: []*TaxCredit{
				&TaxCredit{
					AmountRemaining: 100,
					CrRule:          core.CreditRule{Type: core.CrRuleTypeCashable},
				},
			},
		},
		IncomeCalc: incCalc,
	}

	c0, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	c1, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	aggregator, err := NewAggregator(c0, c1, c2)
	if err != nil {
		t.Fatal(err)
	}

	aggregator.SetFinances(finances)
	actualA, actualB, cr := aggregator.TaxPayable()
	expected := (3000.0 / 2.0) * float64(len(aggregator.calculators))
	expected -= 300.0
	if actualA != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actualA)
	}
	if actualB != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actualB)
	}
	if len(cr) != 6 {
		t.Errorf("expected 6 tax credits, got: %d", len(cr))
	}
}

func TestAggregator_Year(t *testing.T) {

	c0, c1 := &Calculator{taxYear: 2019}, &Calculator{taxYear: 2019}
	agg, err := NewAggregator(c0, c1)
	if err != nil {
		t.Fatal(err)
	}

	year := agg.Year()
	if year != 2019 {
		t.Errorf("expected tax year to be 2019, got: %d", year)
	}
}

func TestAggregator_Regions(t *testing.T) {

	c0, c1 := &Calculator{taxRegion: "BC"}, &Calculator{taxRegion: "Canada"}
	agg, err := NewAggregator(c0, c1)
	if err != nil {
		t.Fatal(err)
	}

	regions := agg.Regions()
	diff := deep.Equal(regions, []core.Region{"BC", "Canada"})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestAggregator_SetDependents(t *testing.T) {

	c0, c1 := new(Calculator), new(Calculator)
	agg, err := NewAggregator(c0, c1)
	if err != nil {
		t.Fatal(err)
	}

	deps := []*human.Person{
		&human.Person{Name: "test1"},
		&human.Person{Name: "test2"},
	}

	agg.SetDependents(deps...)

	diff := deep.Equal(c0.dependents, deps)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
	diff = deep.Equal(c1.dependents, deps)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestAggregator_SetFinances(t *testing.T) {

	c0 := &Calculator{taxYear: 2019, taxRegion: "BC"}
	c1 := &Calculator{taxYear: 2019, taxRegion: "Canada"}
	agg, err := NewAggregator(c0, c1)
	if err != nil {
		t.Fatal(err)
	}

	finances := core.NewHouseholdFinancesNop()
	crA := &testTaxCredit{
		onReferenceFinancer: finances.SpouseA(),
		onAmounts:           [3]float64{2000, 1000, 1000},
		onRegion:            "BC",
		onYear:              2019,
	}
	crB := &testTaxCredit{
		onReferenceFinancer: finances.SpouseB(),
		onAmounts:           [3]float64{4000, 2000, 2000},
		onRegion:            "Canada",
		onYear:              2019,
	}

	agg.SetFinances(finances, crA, crB)

	diff := deep.Equal(c0.crSpouseA, []core.TaxCredit{crA})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	diff = deep.Equal(c1.crSpouseB, []core.TaxCredit{crB})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	if c0.finances != finances {
		t.Errorf("expected first calculator to be set with given household finances")
	}
	if c1.finances != finances {
		t.Errorf("expected first calculator to be set with given household finances")
	}
}
