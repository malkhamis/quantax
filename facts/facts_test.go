package facts

import (
	"reflect"
	"testing"
)

func TestFacts_Clone_NewCopy(t *testing.T) {

	original := testFacts(t)
	cloned := original.Clone()
	cloned.FactsFed.Rates[-1000.123] = Bracket{123456789, 987654321}
	cloned.FactsProv.Rates[-1000.123] = Bracket{123456789, 987654321}

	_, ok := original.FactsFed.Rates[-1000.123]
	if ok {
		t.Errorf("expected the cloned Facts to not point to original")
	}

	_, ok = original.FactsProv.Rates[-1000.123]
	if ok {
		t.Errorf("expected the cloned Facts to not point to original")
	}
}

func TestFacts_NumFieldsUnchanged(t *testing.T) {

	dummy := Facts{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 3 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor of " +
				"type. Next, update this test with the new number of fields",
		)
	}
}

func TestFactsProv_NumFieldsUnchanged(t *testing.T) {

	dummy := FactsProv{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor of " +
				"type. Next, update this test with the new number of fields",
		)
	}
}

func TestFactsFed_NumFieldsUnchanged(t *testing.T) {

	dummy := FactsFed{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 1 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor of " +
				"type. Next, update this test with the new number of fields",
		)
	}
}
