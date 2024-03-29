package benefits

import (
	"math"
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

func TestBCECTBMaxReducer_Apply(t *testing.T) {

	bracket := core.WeightedBrackets{
		0.0132: core.Bracket{100000, math.Inf(1)},
	}

	mr := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: core.Bracket{0, 55},
			},
		},
	}

	err := mr.Validate()
	if err != nil {
		t.Fatal(err)
	}

	child1, child2 := &human.Person{AgeMonths: 0}, &human.Person{AgeMonths: 6}
	income := 110000.0
	max := 2.0 * 12.0 * 55
	expected := max - (2 * 0.0132 * 10000)

	actual := mr.Apply(income, child1, child2)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 0.0
	expected = 0.0
	actual = mr.Apply(income) // no children
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}

	income = 500000.0
	expected = 0.0
	actual = mr.Apply(income, child1)
	if actual != expected {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", expected, actual)
	}
}

func TestBCECTBMaxReducer_Validate_InvalidAgeRanges(t *testing.T) {

	formula := BCECTBMaxReducer{
		BeneficiaryClasses: []AgeGroupBenefits{
			AgeGroupBenefits{
				AgesMonths:      human.AgeRange{10, 0},
				AmountsPerMonth: core.Bracket{0, 55},
			},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != human.ErrInvalidAgeRange {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", human.ErrInvalidAgeRange, err)
	}
}

func TestBCECTBMaxReducer_Validate_NilFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		BeneficiaryClasses: nil,
		ReducerFormula:     nil,
	}

	err := formula.Validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

}

func TestBCECTBMaxReducer_Validate_InvalidFormula(t *testing.T) {

	formula := BCECTBMaxReducer{
		ReducerFormula: core.WeightedBrackets{
			0.0132: core.Bracket{100000, 1},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != core.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", core.ErrBoundsReversed, err)
	}
}

func TestBCECTBMaxReducer_Clone(t *testing.T) {

	bracket := core.WeightedBrackets{
		0.0132: core.Bracket{100000, math.Inf(1)},
	}

	originalFormula := &BCECTBMaxReducer{
		ReducerFormula: bracket,
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 6*12 - 1},
				AmountsPerMonth: core.Bracket{0, 55},
			},
		},
	}

	err := originalFormula.Validate()
	if err != nil {
		t.Fatal(err)
	}

	income := 100000.0
	child1, child2 := &human.Person{AgeMonths: 0}, &human.Person{AgeMonths: 6}
	originalResults := originalFormula.Apply(income, child1, child2)

	clone := originalFormula.Clone()
	originalFormula.BeneficiaryClasses = nil
	originalFormula.ReducerFormula = nil

	actualResults := clone.Apply(income, child1, child2)
	if actualResults != originalResults {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", originalResults, actualResults)
	}

}

func TestBCECTBReducer_Clone_Nil(t *testing.T) {

	var mr *BCECTBMaxReducer
	clone := mr.Clone()
	if clone != nil {
		t.Fatal("cloning a nil formula should return nil")
	}
}

func TestBCECTBReducer_NumFieldsUnchanged(t *testing.T) {

	dummy := BCECTBMaxReducer{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 2 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
