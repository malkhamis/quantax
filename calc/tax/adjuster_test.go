package tax

import (
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/calc/finance"
)

func TestAdjuster_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianCapitalGainAdjuster{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}

type testAdjuster struct {
	adjustment Adjustment
}

func (ta testAdjuster) Adjusted(_ *finance.IndividualFinances) Adjustment {
	return ta.adjustment
}

func (ta testAdjuster) Clone() Adjuster {
	return ta
}
