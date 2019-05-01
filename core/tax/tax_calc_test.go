package tax

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewCalculator(t *testing.T) {

	cfg := CalcConfig{
		ContraTaxFormula: &testContraTaxFormula{},
		TaxFormula:       &testTaxFormula{},
		IncomeCalc:       &testIncomeCalculator{},
	}

	c, err := NewCalculator(cfg)
	if err != nil {
		t.Fatalf("expected no errors creating new calculator with valid configs, got: %v", err)
	}

	if c == nil {
		t.Fatal("expected a non-nil instance if no error returned")
	}
}

func TestCalculator_SetFinances(t *testing.T) {

	f := core.NewHouseholdFinancesNop()
	crSpouseA := &testTaxCredit{
		onAmounts:           [3]float64{0, 0, 1500},
		onReferenceFinancer: f.SpouseA(),
	}
	crSpouseB := &testTaxCredit{
		onAmounts:           [3]float64{0, 0, 3000},
		onReferenceFinancer: f.SpouseB(),
	}
	credits := []core.TaxCredit{crSpouseA, crSpouseB}

	calc := &Calculator{}
	calc.SetFinances(f, credits...)

	diff := deep.Equal(calc.crSpouseA[0], crSpouseA)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	diff = deep.Equal(calc.crSpouseB[0], crSpouseB)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_SetFinances_Nil(t *testing.T) {

	c := &Calculator{}

	c.SetFinances(nil)
	if c.finances == nil {
		t.Fatal("expected empty finances to be set")
	}

	c.SetFinances(&testHouseholdFinances{})
	if c.finances == nil {
		t.Fatal("expected empty finances to be set")
	}
}

func TestCalculator_setCredits_Nil_Credit(t *testing.T) {

	calc := &Calculator{}
	credits := []core.TaxCredit{nil}
	calc.setCredits(credits)

	if len(calc.crSpouseA) != 0 {
		t.Error("expected no credits to be set")
	}
	if len(calc.crSpouseB) != 0 {
		t.Error("expected no credits to be set")
	}
}

func TestCalculator_setCredits_Zero_Credit(t *testing.T) {

	calc := &Calculator{}
	credits := []core.TaxCredit{
		&testTaxCredit{onAmounts: [3]float64{0, 0, 0}},
	}
	calc.setCredits(credits)

	if len(calc.crSpouseA) != 0 {
		t.Error("expected no credits to be set")
	}
	if len(calc.crSpouseB) != 0 {
		t.Error("expected no credits to be set")
	}
}

func TestCalculator_setCredits_different_tax_region(t *testing.T) {

	calc := &Calculator{taxRegion: core.Region(t.Name())}
	credits := []core.TaxCredit{
		&testTaxCredit{
			onRegion:  "something-else",
			onAmounts: [3]float64{0, 0, 1000},
		},
	}
	calc.setCredits(credits)

	if len(calc.crSpouseA) != 0 {
		t.Error("expected no credits to be set")
	}
	if len(calc.crSpouseB) != 0 {
		t.Error("expected no credits to be set")
	}
}

func TestCalculator_setCredits_future_credits(t *testing.T) {

	calc := &Calculator{taxRegion: core.Region(t.Name()), taxYear: 2019}
	credits := []core.TaxCredit{
		&testTaxCredit{
			onRegion:  core.Region(t.Name()),
			onAmounts: [3]float64{0, 0, 1000},
			onYear:    2020,
		},
	}
	calc.setCredits(credits)

	if len(calc.crSpouseA) != 0 {
		t.Error("expected no credits to be set")
	}
	if len(calc.crSpouseB) != 0 {
		t.Error("expected no credits to be set")
	}
}

func TestCalculator_setCredits_different_finances(t *testing.T) {

	calc := &Calculator{
		finances:  core.NewHouseholdFinancesNop(),
		taxRegion: core.Region(t.Name()),
		taxYear:   2019,
	}
	credits := []core.TaxCredit{
		&testTaxCredit{
			onRegion:            core.Region(t.Name()),
			onAmounts:           [3]float64{0, 0, 1000},
			onYear:              2019,
			onReferenceFinancer: core.NewFinancerNop(),
		},
	}
	calc.setCredits(credits)

	if len(calc.crSpouseA) != 0 {
		t.Error("expected no credits to be set")
	}
	if len(calc.crSpouseB) != 0 {
		t.Error("expected no credits to be set")
	}
}

func TestCalculator_SetDependents(t *testing.T) {

	dep := &human.Person{AgeMonths: 10, Name: t.Name()}
	calc := new(Calculator)
	calc.SetDependents(nil, nil, dep)
	diff := deep.Equal(calc.dependents, []*human.Person{dep})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_Year(t *testing.T) {

	calc := &Calculator{taxYear: 2019}
	if calc.Year() != 2019 {
		t.Errorf("expected call to return 2019, got: %d", calc.Year())
	}
}

func TestCalculator_Regions(t *testing.T) {

	calc := &Calculator{taxRegion: core.Region(t.Name())}
	diff := deep.Equal(calc.Regions(), []core.Region{core.Region(t.Name())})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

//
// func TestCalculator_TaxPayable(t *testing.T) {
//
// 	incCalc := testIncomeCalculator{onTotalIncome: 3000.0}
// 	formula := testTaxFormula{onApply: incCalc.TotalIncome() / 2.0}
// 	cformula := &testContraTaxFormula{
// 		onApply: []*taxCredit{
// 			&taxCredit{
// 				amount: 50,
// 				rule:   CreditRule{Source: "tuition", Type: CrRuleTypeCashable},
// 			},
// 		},
// 	}
//
// 	cfg := CalcConfig{
// 		TaxFormula:       formula,
// 		ContraTaxFormula: cformula,
// 		IncomeCalc:       incCalc,
// 	}
//
// 	c, err := NewCalculator(cfg)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	c.SetFinances(core.NewEmptyIndividualFinances())
// 	actualTax, actualCr := c.TaxPayable()
//
// 	expectedTax := formula.onApply - 50.0
// 	expectedCr := []core.TaxCredit{
// 		&taxCredit{
// 			owner:  c,
// 			amount: 0.0,
// 			rule:   CreditRule{Source: "tuition", Type: CrRuleTypeCashable},
// 		},
// 	}
//
// 	if actualTax != expectedTax {
// 		t.Fatalf("unexpected tax\nwant: %.2f\n got: %.2f", expectedTax, actualTax)
// 	}
//
// 	diff := deep.Equal(actualCr, expectedCr)
// 	if diff != nil {
// 		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
// 	}
// }
//
// func TestCalculator_netPayableTax(t *testing.T) {
//
// 	crGroup := []*taxCredit{
// 		&taxCredit{amount: 5000, rule: CreditRule{Source: "1", Type: CrRuleTypeCashable}},
// 		&taxCredit{amount: 4000, rule: CreditRule{Source: "2", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 2000, rule: CreditRule{Source: "3", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 1000, rule: CreditRule{Source: "4", Type: CrRuleTypeCashable}},
// 		&taxCredit{amount: 500, rule: CreditRule{Source: "5", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 500, rule: CreditRule{Source: "6", Type: CrRuleTypeCanCarryForward}},
// 	}
//
// 	actualNetTax, actualRemainingCrs := (&Calculator{}).netPayableTax(10000, crGroup)
// 	expectedNetTax := -1000.0
// 	expectedRemainingCrs := []*taxCredit{
// 		&taxCredit{amount: 0, rule: CreditRule{Source: "1", Type: CrRuleTypeCashable}},
// 		&taxCredit{amount: 0, rule: CreditRule{Source: "2", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 0, rule: CreditRule{Source: "3", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 0, rule: CreditRule{Source: "4", Type: CrRuleTypeCashable}},
// 		&taxCredit{amount: 0, rule: CreditRule{Source: "5", Type: CrRuleTypeNotCarryForward}},
// 		&taxCredit{amount: 500, rule: CreditRule{Source: "6", Type: CrRuleTypeCanCarryForward}},
// 	}
//
// 	if actualNetTax != expectedNetTax {
// 		t.Errorf(
// 			"actual net tax does not match expected\nwant: %.2f\ngot: %.2f",
// 			expectedNetTax, actualNetTax,
// 		)
// 	}
//
// 	diff := deep.Equal(actualRemainingCrs, expectedRemainingCrs)
// 	if diff != nil {
// 		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
// 	}
//
// }
//

func TestCalculator_makeTaxPayers_couple(t *testing.T) {

	calc := &Calculator{
		finances:   core.NewHouseholdFinancesNop(),
		dependents: []*human.Person{{AgeMonths: 10, Name: t.Name()}},
	}
	taxPayerA, taxPayerB := calc.makeTaxPayers(1000, 2000)

	expectedA := &TaxPayer{
		Finances:        calc.finances.SpouseA(),
		NetIncome:       1000,
		SpouseFinances:  calc.finances.SpouseB(),
		SpouseNetIncome: 2000,
		Dependents:      calc.dependents,
	}

	expectedB := &TaxPayer{
		Finances:        calc.finances.SpouseB(),
		NetIncome:       2000,
		SpouseFinances:  calc.finances.SpouseA(),
		SpouseNetIncome: 1000,
		Dependents:      calc.dependents,
	}

	diff := deep.Equal(taxPayerA, expectedA)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	diff = deep.Equal(taxPayerB, expectedB)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}

func TestCalculator_makeTaxPayers_singleA(t *testing.T) {

	calc := &Calculator{
		finances:   &testHouseholdFinances{onSpouseA: core.NewFinancerNop()},
		dependents: []*human.Person{{AgeMonths: 10, Name: t.Name()}},
	}
	taxPayerA, taxPayerB := calc.makeTaxPayers(1000, 2000)

	expectedA := &TaxPayer{
		Finances:        calc.finances.SpouseA(),
		NetIncome:       1000,
		SpouseFinances:  nil,
		SpouseNetIncome: 2000,
		Dependents:      calc.dependents,
	}

	diff := deep.Equal(taxPayerA, expectedA)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	if taxPayerB != nil {
		t.Fatal("expected taxPayerB to be nil for unmarried taxpayer")
	}

}

func TestCalculator_makeTaxPayers_singleB(t *testing.T) {

	calc := &Calculator{
		finances:   &testHouseholdFinances{onSpouseB: core.NewFinancerNop()},
		dependents: []*human.Person{{AgeMonths: 10, Name: t.Name()}},
	}
	taxPayerA, taxPayerB := calc.makeTaxPayers(1000, 2000)

	expectedB := &TaxPayer{
		Finances:        calc.finances.SpouseB(),
		NetIncome:       2000,
		SpouseFinances:  nil,
		SpouseNetIncome: 1000,
		Dependents:      calc.dependents,
	}

	diff := deep.Equal(taxPayerB, expectedB)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	if taxPayerA != nil {
		t.Fatal("expected taxPayerA to be nil for unmarried taxpayer")
	}

}

func TestNewCalculator_Error(t *testing.T) {

	cfg := CalcConfig{
		TaxFormula:       &testTaxFormula{},
		ContraTaxFormula: &testContraTaxFormula{},
		IncomeCalc:       nil,
	}
	_, err := NewCalculator(cfg)
	if errors.Cause(err) != ErrNoIncCalc {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoIncCalc, err)
	}

}
