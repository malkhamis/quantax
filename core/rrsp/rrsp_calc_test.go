package rrsp

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

func TestNewCalculator(t *testing.T) {

	cfg := CalcConfig{Formula: new(testFormula), TaxCalc: new(testTaxCalculator)}
	c, err := NewCalculator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil calculator if no error")
	}
}

func TestNewCalculator_Error(t *testing.T) {

	_, err := NewCalculator(CalcConfig{nil, nil})
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}
}

func TestCalculator_SetTargetSpouse(t *testing.T) {

	c := &Calculator{}

	c.SetTargetSpouseB()
	if c.isTargetSpouseB != true {
		t.Error("expected target spouse to indicate spouseB")
	}

	c.SetTargetSpouseA()
	if c.isTargetSpouseB != false {
		t.Error("expected target spouse to indicate spouseA")
	}
}

func TestCalculator_SetDependents(t *testing.T) {

	c := &Calculator{}

	var f core.HouseholdFinances
	c.SetFinances(f, nil)
	if c.householdFinances == nil {
		t.Fatal("expected a noop finances to be set when method is called with nil")
	}

	c.householdFinances = nil
	c.taxCredits = nil

	f = core.NewHouseholdFinancesNop()
	credits := []core.TaxCredit{}
	c.SetFinances(f, credits)

	if c.householdFinances != f {
		t.Error("expected passed finances to be set in calculator")
	}

	if c.taxCredits == nil {
		t.Error("expected passed credits to be set in calculator")
	}

}

func TestCalculator_SetFinances(t *testing.T) {

	c := &Calculator{}
	deps := []*human.Person{{AgeMonths: 9, Name: t.Name()}}
	c.SetDependents(deps)

	diff := deep.Equal(deps, c.dependents)
	if diff != nil {
		t.Fatal("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_TaxPaid(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onTaxPayableSpouseA: []float64{1000, 1300},
		onTaxPayableSpouseB: []float64{0, 0},
		onTaxPayableCredits: [][]core.TaxCredit{nil, []core.TaxCredit{}},
	}

	c := &Calculator{
		formula:           &testFormula{},
		taxCalculator:     taxCalc,
		householdFinances: core.NewHouseholdFinancesNop(),
	}

	actualTax, actualCr := c.TaxPaid(0)
	expected := 1300.0 - 1000
	if actualTax != expected {
		t.Errorf(
			"actual tax difference does not match expected\nwant: %.2f\n got: %.2f",
			expected, actualTax,
		)
		if actualCr == nil {
			t.Errorf("expected non-nil credits")
		}
	}
}

func TestCalculator_TaxRefund(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onTaxPayableSpouseA: []float64{1000, 750},
		onTaxPayableSpouseB: []float64{0, 0},
		onTaxPayableCredits: [][]core.TaxCredit{nil, []core.TaxCredit{}},
	}

	c := &Calculator{
		formula:           &testFormula{},
		taxCalculator:     taxCalc,
		householdFinances: core.NewHouseholdFinancesNop(),
	}

	actualTax, actualCr := c.TaxRefund(0)
	expected := 1000.0 - 750.0
	if actualTax != expected {
		t.Errorf(
			"actual tax difference does not match expected\nwant: %.2f\n got: %.2f",
			expected, actualTax,
		)
		if actualCr == nil {
			t.Errorf("expected non-nil credits")
		}
	}
}

func TestCalculator_taxDiffForTargetSpouse_A(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onTaxPayableSpouseA: []float64{1000, 750},
		onTaxPayableSpouseB: []float64{0, 0},
		onTaxPayableCredits: [][]core.TaxCredit{nil, []core.TaxCredit{}},
	}

	c := &Calculator{taxCalculator: taxCalc, householdFinances: core.NewHouseholdFinancesNop()}
	actualDiff, actualCr := c.taxDiffForTargetSpouse(0, 0)
	expected := 1000.0 - 750.0
	if actualDiff != expected {
		t.Errorf(
			"actual tax difference does not match expected\nwant: %.2f\n got: %.2f",
			expected, actualDiff,
		)
		if actualCr == nil {
			t.Errorf("expected non-nil credits")
		}
	}
}

func TestCalculator_ContributedEarned(t *testing.T) {

	c := &Calculator{
		householdFinances: core.NewHouseholdFinancesNop(),
		formula:           &testFormula{onContributionEarned: 1000},
	}

	actual, expected := c.ContributionEarned(), 1000.0
	if actual != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expected, actual,
		)
	}
}

func TestCalculator_ContributedEarned_NilSpouse(t *testing.T) {

	c := &Calculator{
		householdFinances: &testHouseholdFinances{},
		formula:           &testFormula{onContributionEarned: 1000},
	}

	actual, expected := c.ContributionEarned(), 0.0
	if actual != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f",
			expected, actual,
		)
	}
}

func TestCalculator_taxDiffForTargetSpouse_B(t *testing.T) {

	taxCalc := &testTaxCalculator{
		onTaxPayableSpouseA: []float64{0, 0},
		onTaxPayableSpouseB: []float64{1000, 750},
		onTaxPayableCredits: [][]core.TaxCredit{nil, []core.TaxCredit{}},
	}

	c := &Calculator{
		taxCalculator:     taxCalc,
		householdFinances: core.NewHouseholdFinancesNop(),
		isTargetSpouseB:   true,
	}
	actualDiff, actualCr := c.taxDiffForTargetSpouse(0, 0)
	expected := 1000.0 - 750.0
	if actualDiff != expected {
		t.Errorf(
			"actual tax difference does not match expected\nwant: %.2f\n got: %.2f",
			expected, actualDiff,
		)
		if actualCr == nil {
			t.Errorf("expected non-nil credits")
		}
	}
}

func TestCalculator_taxDiffForTargetSpouse_Nil(t *testing.T) {

	c := &Calculator{
		taxCalculator:     &testTaxCalculator{},
		householdFinances: &testHouseholdFinances{},
	}

	actualDiff, actualCr := c.taxDiffForTargetSpouse(0, 0)
	if actualDiff != 0 {
		t.Errorf("expected zero tax difference for nil spouse, got: %.2f", actualDiff)
		if actualCr != nil {
			t.Errorf("expected nil credits")
		}
	}
}

func TestCalculator_targetSpouse(t *testing.T) {

	c := &Calculator{
		householdFinances: core.NewHouseholdFinancesNop(),
	}

	if target := c.targetSpouse(); target != c.householdFinances.SpouseA() {
		t.Errorf("expected default target to be spouse A")
	}

	c.isTargetSpouseB = true
	if target := c.targetSpouse(); target != c.householdFinances.SpouseB() {
		t.Errorf("expected target to be spouse B")
	}
}

func TestCalculator_cloneFinancesAndGetTargetRef(t *testing.T) {

	c := &Calculator{
		householdFinances: core.NewHouseholdFinancesNop(),
	}

	cloneHF, cloneA := c.cloneFinancesAndGetTargetRef()
	if cloneHF == c.householdFinances {
		t.Error("expected clone to not equal cached finances")
	}
	if cloneA == c.householdFinances.SpouseA() {
		t.Error("expected clone to not equal cached finances")
	}

	c.isTargetSpouseB = true
	cloneHF, cloneB := c.cloneFinancesAndGetTargetRef()
	if cloneHF == c.householdFinances {
		t.Error("expected clone to not equal cached finances")
	}
	if cloneB == c.householdFinances.SpouseB() {
		t.Error("expected clone to not equal cached finances")
	}

	if cloneA == cloneB {
		t.Error("expected clone of spouse A to not equal clone of spouse B")
	}
}

func TestCalculator_setupTaxCalculator(t *testing.T) {

	taxCalclMocker := &testTaxCalculator{}
	c := &Calculator{
		taxCalculator: taxCalclMocker,
		taxCredits:    []core.TaxCredit{},
		dependents:    []*human.Person{},
	}
	f := core.NewHouseholdFinancesNop()

	c.setupTaxCalculator(f)
	if taxCalclMocker.financesPassedOnSetFinances[0] != f {
		t.Error("expected household finances to be passed to tax calculator")
	}

	if taxCalclMocker.depsPassedOnSetDependents[0] == nil {
		t.Error("expected non-nil dependents set in calculator to be passed to tax calculator")
	}

	if taxCalclMocker.creditsPassedOnSetFinances[0] == nil {
		t.Error("expected non-nil credits set in calculator to be passed to tax calculator")
	}

}

func TestCalcConfig_validate(t *testing.T) {

	formula := &testFormula{}
	err := CalcConfig{formula, nil}.validate()
	if errors.Cause(err) != ErrNoTaxCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoTaxCalc, err)
	}

	taxCalc := &testTaxCalculator{}
	err = CalcConfig{formula, taxCalc}.validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = CalcConfig{nil, nil}.validate()
	if errors.Cause(err) != ErrNoFormula {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoFormula, err)
	}

	simulatedErr := errors.New("test error")
	formula = &testFormula{onValidate: simulatedErr}
	err = CalcConfig{formula, nil}.validate()
	if errors.Cause(err) != simulatedErr {
		t.Errorf("unexpected error\nwant: %v\n got: %v", simulatedErr, err)
	}

}
