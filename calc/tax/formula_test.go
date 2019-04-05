package tax

import (
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestFormula_Clone(t *testing.T) {

	original := &CanadianFormula{
		ExcludedIncome:     []finance.IncomeSource{finance.IncSrcTFSA},
		ExcludedDeductions: []finance.DeductionSource{},
		IncomeAdjusters: map[finance.IncomeSource]Adjuster{
			finance.IncSrcCapitalGainCA: CanadianCapitalGainAdjuster{},
		},
		DeductionAdjusters: map[finance.DeductionSource]Adjuster{
			finance.DeducSrcMedical: CanadianCapitalGainAdjuster{},
		},
		WeightedBrackets: finance.WeightedBrackets{0.1: finance.Bracket{0, 10}},
	}

	clone := original.Clone()
	originalResults := original.Apply(5)
	original.WeightedBrackets[0.1] = finance.Bracket{100, 1000}
	cloneResults := clone.Apply(5)
	if originalResults != cloneResults {
		t.Errorf("expected clone results to be equal to results of original formula prior to mutation")
	}

	original.ExcludedIncome = nil
	original.ExcludedDeductions = nil
	original.IncomeAdjusters = nil
	original.DeductionAdjusters = nil

	typed, ok := clone.(*CanadianFormula)
	if !ok {
		t.Fatal("expected typecasting to succeed")
	}

	if typed.ExcludedIncome == nil {
		t.Fatal("expected changes to original to not be reflected in clone")
	}
	if typed.ExcludedDeductions == nil {
		t.Fatal("expected changes to original to not be reflected in clone")
	}
	if typed.IncomeAdjusters == nil {
		t.Fatal("expected changes to original to not be reflected in clone")
	}
	if typed.DeductionAdjusters == nil {
		t.Fatal("expected changes to original to not be reflected in clone")
	}

}

func TestFormula_Clone_Nil(t *testing.T) {

	var mr *CanadianFormula
	clone := mr.Clone()
	if clone != nil {
		t.Fatal("cloning a nil formula should return nil")
	}
}

func TestFormula_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianFormula{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 5 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
