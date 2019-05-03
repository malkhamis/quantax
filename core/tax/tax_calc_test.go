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
	calc.SetFinances(f, credits)

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

	c.SetFinances(nil, nil)
	if c.finances == nil {
		t.Fatal("expected empty finances to be set")
	}

	c.SetFinances(&testHouseholdFinances{}, nil)
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

func TestCalculator_setCredits_unreferenced_credit(t *testing.T) {

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
			onReferenceFinancer: nil,
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
	calc.SetDependents([]*human.Person{nil, nil, dep})
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

func TestCalculator_TaxPayable(t *testing.T) {

	dummyCr := &TaxCredit{AmountInitial: 123, AmountUsed: 123}
	incCalc := &testIncomeCalculator{onTotalIncome: 3000.0}
	formula := &testTaxFormula{onApply: incCalc.TotalIncome() / 2.0}
	cformula := &testContraTaxFormula{onApply: []*TaxCredit{dummyCr}}

	cfg := CalcConfig{
		TaxFormula:       formula,
		ContraTaxFormula: cformula,
		IncomeCalc:       incCalc,
	}

	c, err := NewCalculator(cfg)
	if err != nil {
		t.Fatal(err)
	}

	actualA, actualB, actualCr := c.TaxPayable()
	if actualA != 1500.0 {
		t.Errorf("actual does not match expected\nwant: %.2f\n got: %.2f", 1500.0, actualA)
	}
	if actualB != 1500.0 {
		t.Errorf("actual does not match expected\nwant: %.2f\n got: %.2f", 1500.0, actualB)
	}

	diff := deep.Equal(actualCr, []core.TaxCredit{dummyCr, dummyCr})
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestCalculator_netPayableTax(t *testing.T) {

	crBefore := []core.TaxCredit{
		&TaxCredit{AmountInitial: 5000, AmountRemaining: 5000, CrRule: core.CreditRule{Type: core.CrRuleTypeCashable}},
		&TaxCredit{AmountInitial: 4000, AmountRemaining: 4000, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 2000, AmountRemaining: 2000, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 1000, AmountRemaining: 1000, CrRule: core.CreditRule{Type: core.CrRuleTypeCashable}},
		&TaxCredit{AmountInitial: 500, AmountRemaining: 500, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 500, AmountRemaining: 500, CrRule: core.CreditRule{Type: core.CrRuleTypeCanCarryForward}},
	}

	crAfter := []core.TaxCredit{
		&TaxCredit{AmountInitial: 5000, AmountRemaining: 0, AmountUsed: 5000, CrRule: core.CreditRule{Type: core.CrRuleTypeCashable}},
		&TaxCredit{AmountInitial: 4000, AmountRemaining: 0, AmountUsed: 4000, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 2000, AmountRemaining: 0, AmountUsed: 1000, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 1000, AmountRemaining: 0, AmountUsed: 1000, CrRule: core.CreditRule{Type: core.CrRuleTypeCashable}},
		&TaxCredit{AmountInitial: 500, AmountRemaining: 0, AmountUsed: 0, CrRule: core.CreditRule{Type: core.CrRuleTypeNotCarryForward}},
		&TaxCredit{AmountInitial: 500, AmountRemaining: 500, AmountUsed: 0, CrRule: core.CreditRule{Type: core.CrRuleTypeCanCarryForward}},
	}

	actualNetTax := (&Calculator{}).netPayableTax(10000, crBefore)

	expectedNetTax := -1000.0
	if actualNetTax != expectedNetTax {
		t.Errorf(
			"actual net tax does not match expected\nwant: %.2f\ngot: %.2f",
			expectedNetTax, actualNetTax,
		)
	}

	diff := deep.Equal(crBefore, crAfter)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}

func TestCalculator_netPayableTax_CanCarryForward_amounts(t *testing.T) {

	crBefore := []core.TaxCredit{
		&TaxCredit{
			AmountInitial:   500,
			AmountRemaining: 500,
			CrRule: core.CreditRule{
				Type: core.CrRuleTypeCanCarryForward,
			},
		},
	}

	crAfter := []core.TaxCredit{
		&TaxCredit{
			AmountInitial:   500,
			AmountRemaining: 250,
			AmountUsed:      250,
			CrRule: core.CreditRule{
				Type: core.CrRuleTypeCanCarryForward,
			},
		},
	}

	actualNetTax := (&Calculator{}).netPayableTax(250, crBefore)

	expectedNetTax := 0.0
	if actualNetTax != expectedNetTax {
		t.Errorf(
			"actual net tax does not match expected\nwant: %.2f\ngot: %.2f",
			expectedNetTax, actualNetTax,
		)
	}

	diff := deep.Equal(crBefore, crAfter)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}

func TestCalculator_netIncome(t *testing.T) {

	calc := &Calculator{
		finances:         core.NewHouseholdFinancesNop(),
		incomeCalculator: &testIncomeCalculator{onNetIncome: 1000},
	}

	actualA, actualB := calc.netIncome()
	expected := 1000.0
	if actualA != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f", expected, actualA)
	}
	if actualB != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f", expected, actualB)
	}
}

func TestCalculator_totalTax(t *testing.T) {

	calc := &Calculator{
		finances: core.NewHouseholdFinancesNop(),
		formula:  &testTaxFormula{onApply: 1000.0},
	}

	actualA, actualB := calc.totalTax(0, 0)
	expected := 1000.0
	if actualA != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f", expected, actualA)
	}
	if actualB != expected {
		t.Errorf(
			"actual does not match expected\nwant: %.2f\n got: %.2f", expected, actualB)
	}
}

func TestCalculator_totalCredits(t *testing.T) {

	crSpouseA := &testTaxCredit{onRegion: core.Region(t.Name())}
	crSpouseB := &testTaxCredit{onRegion: core.Region("another")}
	simulatedCr := []*TaxCredit{{}, {}}
	calc := &Calculator{
		finances:  core.NewHouseholdFinancesNop(),
		crSpouseA: []core.TaxCredit{crSpouseA},
		crSpouseB: []core.TaxCredit{crSpouseB},
		contraFormula: &testContraTaxFormula{
			onApply: simulatedCr,
		},
	}

	actualA, actualB := calc.totalCredits(0, 0)

	expectedA := []core.TaxCredit{crSpouseA, simulatedCr[0], simulatedCr[1]}
	diff := deep.Equal(actualA, expectedA)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

	expectedB := []core.TaxCredit{crSpouseB, simulatedCr[0], simulatedCr[1]}
	diff = deep.Equal(actualB, expectedB)
	if diff != nil {
		t.Error("actual does not match expected\n", strings.Join(diff, "\n"))
	}

}

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

func TestCalculator_panicIfEqNonNilSpouses(t *testing.T) {

	spouseFinances := core.NewFinancerNop()
	c := &Calculator{
		finances: &testHouseholdFinances{
			onSpouseA: spouseFinances,
			onSpouseB: spouseFinances,
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic if household finances contain the same spouse twice")
		}
	}()

	c.panicIfEqNonNilSpouses()
}
