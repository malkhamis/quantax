package income

import (
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/core"
)

func TestRecipe_Clone_Nil(t *testing.T) {

	var nilRecipe *Recipe
	clone := nilRecipe.Clone()
	if clone != nil {
		t.Fatal("expected cloned nil recipe to also be nil")
	}
}

func TestRecipe_Clone(t *testing.T) {

	original := &Recipe{
		IncomeAdjusters: map[core.FinancialSource]Adjuster{
			core.FinancialSource(1000): testAdjuster{adjusted: 250.0},
		},
		DeductionAdjusters: map[core.FinancialSource]Adjuster{
			core.FinancialSource(2000): testAdjuster{adjusted: 100.0},
		},
	}

	clone := original.Clone()
	delete(original.IncomeAdjusters, core.FinancialSource(1000))
	delete(original.DeductionAdjusters, core.FinancialSource(2000))

	_, ok := clone.IncomeAdjusters[core.FinancialSource(1000)]
	if !ok {
		t.Errorf("expected changes to original to not affect clone")
	}

	_, ok = clone.DeductionAdjusters[core.FinancialSource(2000)]
	if !ok {
		t.Errorf("expected changes to original to not affect clone")
	}
}

func TestBCECTBReducer_NumFieldsUnchanged(t *testing.T) {

	dummy := Recipe{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 2 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
