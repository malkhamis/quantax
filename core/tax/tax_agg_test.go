package tax

import (
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

func TestAggregator(t *testing.T) {

	c0, c1, c2 := &Calculator{}, &Calculator{}, &Calculator{}
	_, err := NewAggregator(c0, c1, c2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = NewAggregator(nil, nil, nil)
	if errors.Cause(err) != ErrNoCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCalc, err)
	}

}

func TestAggregator_TaxPayable(t *testing.T) {

	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
	formula := testTaxFormula{onApply: incCalc.TotalIncome(nil) / 2.0}

	cfg := CalcConfig{
		TaxFormula:       formula,
		ContraTaxFormula: &testTaxContraFormula{},
		IncomeCalc:       incCalc,
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

	aggregator.SetFinances(core.NewEmptyIndividualFinances())
	actual, _ := aggregator.TaxPayable()
	expected := (3000.0 / 2.0) * float64(len(aggregator.calculators))
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

func TestAggregator_SetCredits(t *testing.T) {

	c0, c1, c2 := &Calculator{}, &Calculator{}, &Calculator{}
	agg, err := NewAggregator(c0, c1, c2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tc0 := &taxCredit{owner: c0}
	tc1 := &taxCredit{owner: c1}
	tcAnonymous := &taxCredit{owner: nil}
	credits := []core.TaxCredit{tc0, tc1, tcAnonymous}

	agg.SetCredits(credits)

	if len(c0.credits) != 1 {
		t.Fatalf("expected exactly one credit in c0, got: %d", len(c0.credits))
	}
	if c0.credits[0] != tc0 {
		t.Error("expected c0 to accept owned tc0 as credit")
	}

	if len(c1.credits) != 1 {
		t.Fatalf("expected exactly one credit in c1, got: %d", len(c1.credits))
	}
	if c1.credits[0] != tc1 {
		t.Error("expected c1 to accept owned tc0 as credit")
	}

	if len(c2.credits) != 0 {
		t.Fatalf("expected zero credit in c2, got: %d", len(c2.credits))
	}

}
