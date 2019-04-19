package tax

import (
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
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

func TestAggregator_Calc(t *testing.T) {

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

	aggregator.SetFinances(finance.NewEmptyIndividualFinances(2019))
	actual, _ := aggregator.TaxPayable()
	t.Fatal("^^ missing check for credits")
	expected := (3000.0 / 2.0) * float64(len(aggregator.calculators))
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

}
