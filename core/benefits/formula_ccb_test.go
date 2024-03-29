package benefits

import (
	"math"
	"reflect"
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

func TestCCBMaxReducer_Apply(t *testing.T) {

	childCount1 := core.WeightedBrackets{
		0.000: core.Bracket{0, 10000},
		0.030: core.Bracket{10000, 50000},
		0.070: core.Bracket{50000, math.Inf(1)},
	}
	childCount2 := core.WeightedBrackets{
		0.000: core.Bracket{0, 10000},
		0.050: core.Bracket{10000, 50000},
		0.100: core.Bracket{50000, math.Inf(1)},
	}

	mr := &CCBMaxReducer{
		Reducers: []core.WeightedBrackets{childCount1, childCount2},
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 11},
				AmountsPerMonth: core.Bracket{0, 500},
			},
			{
				AgesMonths:      human.AgeRange{12, 23},
				AmountsPerMonth: core.Bracket{0, 250},
			},
		},
	}

	err := mr.Validate()
	if err != nil {
		t.Fatal(err)
	}

	child1, child2 := &human.Person{AgeMonths: 0}, &human.Person{AgeMonths: 6}
	max := (12.0 * 500) + (6.0*500 + 6.0*250)

	income := 100000.0
	expected := max - (0.050 * 40000) - (0.100 * 50000)
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

func TestCCBMaxReducer_Validate_InvalidAgeRanges(t *testing.T) {

	formula := CCBMaxReducer{
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

func TestCCBMaxReducer_Validate_NilFormula(t *testing.T) {

	formula := CCBMaxReducer{
		BeneficiaryClasses: nil,
		Reducers:           nil,
	}

	err := formula.Validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	formula = CCBMaxReducer{
		BeneficiaryClasses: nil,
		Reducers:           []core.WeightedBrackets{nil},
	}

	err = formula.Validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

}

func TestCCBMaxReducer_Validate_InvalidFormula(t *testing.T) {

	formula := CCBMaxReducer{
		Reducers: []core.WeightedBrackets{
			{0.0132: core.Bracket{100000, 1}},
		},
	}

	err := formula.Validate()
	if errors.Cause(err) != core.ErrBoundsReversed {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", core.ErrBoundsReversed, err)
	}
}

func TestCCBMaxReducer_Clone(t *testing.T) {

	childCount1 := core.WeightedBrackets{
		0.000: core.Bracket{0, 10000},
		0.030: core.Bracket{10000, 50000},
		0.070: core.Bracket{50000, math.Inf(1)},
	}
	childCount2 := core.WeightedBrackets{
		0.000: core.Bracket{0, 10000},
		0.050: core.Bracket{10000, 50000},
		0.100: core.Bracket{50000, math.Inf(1)},
	}

	originalFormula := &CCBMaxReducer{
		Reducers: []core.WeightedBrackets{childCount1, childCount2},
		BeneficiaryClasses: []AgeGroupBenefits{
			{
				AgesMonths:      human.AgeRange{0, 11},
				AmountsPerMonth: core.Bracket{0, 500},
			},
			{
				AgesMonths:      human.AgeRange{12, 23},
				AmountsPerMonth: core.Bracket{0, 250},
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
	originalFormula.Reducers = nil

	actualResults := clone.Apply(income, child1, child2)
	if actualResults != originalResults {
		t.Errorf("unexpected results\nwant: %.2f\n got: %.2f", originalResults, actualResults)
	}

}

func TestCCBMaxReducer_Clone_Nil(t *testing.T) {

	var mr *CCBMaxReducer
	clone := mr.Clone()
	if clone != nil {
		t.Fatal("cloning a nil formula should return nil")
	}
}

func TestCCBMaxReducer_NumFieldsUnchanged(t *testing.T) {

	dummy := CCBMaxReducer{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 2 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
