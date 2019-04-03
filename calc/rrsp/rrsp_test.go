package rrsp

import (
	"math"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/tax"

	"github.com/pkg/errors"
)

func TestCalculator_Full(t *testing.T) {

	taxFormula := &tax.CanadianFormula{
		WeightedBrackets: finance.WeightedBrackets{
			0.150: finance.Bracket{0, 10000},
			0.300: finance.Bracket{10000, math.Inf(1)},
		},
	}

	taxCalc, err := tax.NewCalculator(taxFormula)
	if err != nil {
		t.Fatal(err)
	}

	rrspFormula := &MaxCapper{
		Rate:          0.10,
		Cap:           5000,
		IncomeSources: []finance.IncomeSource{finance.IncSrcEarned},
	}

	rrspCalc, err := NewCalculator(taxCalc, rrspFormula)
	if err != nil {
		t.Fatal(err)
	}

	finances := finance.NewEmptyIndividialFinances(2019)
	finances.Income = finance.IncomeBySource{finance.IncSrcEarned: 9000}
	rrspCalc.SetFinances(finances)

	actualContr := rrspCalc.ContributionEarned()
	expectedContr := 900.0
	if actualContr != expectedContr {
		t.Errorf(
			"unexpected earned contribution\nwant: %.2f\n got: %.2f",
			expectedContr, actualContr,
		)
	}

	actualTaxPaid := rrspCalc.TaxPaid(2000.0)
	expectedTaxPaid := (0.15 * 1000.0) + (0.30 * 1000)
	if actualTaxPaid != expectedTaxPaid {
		t.Errorf(
			"unexpected earned contribution\nwant: %.2f\n got: %.2f",
			expectedTaxPaid, actualTaxPaid,
		)
	}

	finances = finance.NewEmptyIndividialFinances(2019)
	finances.Income = finance.IncomeBySource{
		finance.IncSrcEarned: 10000,
		finance.IncSrcUCCB:   1000,
	}
	finances.RRSPContributionRoom = 2000
	rrspCalc.SetFinances(finances)

	actualTaxRefund, err := rrspCalc.TaxRefund(2000.0)
	if err != nil {
		t.Fatal(err)
	}

	expectedTaxRefund := (0.15 * 1000.0) + (0.30 * 1000)
	if actualTaxRefund != expectedTaxRefund {
		t.Errorf(
			"unexpected tax refund\nwant: %.2f\n got: %.2f",
			expectedTaxRefund, actualTaxRefund,
		)
	}
}

func TestNewCalculator_Errors(t *testing.T) {

	_, err := NewCalculator(nil, nil)
	if err != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	_, err = NewCalculator(nil, testFormula{})
	if err != ErrNoTaxCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoTaxCalc, err)
	}

	simulatedErr := errors.New("testing error")
	_, err = NewCalculator(nil, testFormula{err: simulatedErr})
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}

func TestCalculator_RefundError(t *testing.T) {

	c := &Calculator{
		finances: &finance.IndividualFinances{
			RRSPContributionRoom: 1000.0,
		},
	}

	_, err := c.TaxRefund(2000.0)
	if err != ErrNoRRSPRoom {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoRRSPRoom, err)
	}
}

type testFormula struct {
	err           error
	contribution  float64
	incomeSources []finance.IncomeSource
}

func (f testFormula) Contribution(income float64) float64 {
	return f.contribution
}
func (f testFormula) AllowedIncomeSources() []finance.IncomeSource {
	return f.incomeSources
}
func (f testFormula) Validate() error {
	return f.err
}
func (f testFormula) Clone() Formula {
	return f
}
