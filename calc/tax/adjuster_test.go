package tax

import (
	"reflect"
	"testing"
)

func TestAdjuster_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianCapitalGainAdjuster{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type. Next, update this test with the new " +
				"number of fields",
		)
	}
}
