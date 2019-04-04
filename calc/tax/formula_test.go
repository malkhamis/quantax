package tax

import (
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestFormula_Clone(t *testing.T) {

	original := CanadianFormula{
		ExcludedIncome:     []finance.IncomeSource{finance.IncSrcTFSA},
		ExcludedDeductions: []finance.DeductionSource{},
		WeightedBrackets:   finance.WeightedBrackets{0.1: finance.Bracket{0, 10}},
	}
	clone := original.Clone()
	originalResults := original.Apply(5)
	original.WeightedBrackets[0.1] = finance.Bracket{100, 1000}
	cloneResults := clone.Apply(5)
	if originalResults != cloneResults {
		t.Errorf("expected clone results to be equal to results of original formula prior to mutation")
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
	if s.NumField() != 3 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type. Next, update this test with the new " +
				"number of fields",
		)
	}
}
