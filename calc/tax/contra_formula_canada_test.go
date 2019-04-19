package tax

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

func TestCanadianContraFormula_Apply_Nil_finances(t *testing.T) {

	cf := new(CanadianContraFormula)
	actual := cf.Apply(nil, 0)
	if actual != nil {
		t.Errorf("expected nil []credits if finances is nil")
	}

}

func TestCanadianContraFormula_Apply(t *testing.T) {

	cf := &CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 1,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 2,
			},
		},
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 3,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 4,
			},
		},
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 5,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 6,
			},
		},
		ApplicationOrder: []CreditRule{
			{"1000", CrRuleTypeCashable},
			{"2000", CrRuleTypeCanCarryForward},
		},
	}

	err := cf.Validate()
	if err != nil {
		t.Fatal(err)
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddIncome(123, 0)
	finances.AddIncome(999, 0)
	finances.AddDeduction(456, 0)
	finances.AddDeduction(999, 0)
	finances.AddMiscAmount(123, 0)
	finances.AddMiscAmount(456, 0)

	actual := cf.Apply(finances, 0)
	expected := []*taxCredit{
		&taxCredit{amount: 1, rule: CreditRule{"1000", CrRuleTypeCashable}},
		&taxCredit{amount: 5, rule: CreditRule{"1000", CrRuleTypeCashable}},
		&taxCredit{amount: 4, rule: CreditRule{"2000", CrRuleTypeCanCarryForward}},
		&taxCredit{amount: 6, rule: CreditRule{"2000", CrRuleTypeCanCarryForward}},
	}

	// this test might fail if the order of adding income/deductions/misc
	// was changed. However, the most important thing is that the rewturned
	// credits are sorted by their source according to cf.ApplicationOrder
	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_Validate(t *testing.T) {

	cases := []struct {
		name    string
		formula *CanadianContraFormula
		err     error
	}{
		//
		{
			name: "valid",
			formula: &CanadianContraFormula{
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2000"},
				},
				CreditsFromDeduction: map[finance.DeductionSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2000"},
				},
				CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2000"},
				},
				ApplicationOrder: []CreditRule{{Source: "1000"}, {Source: "2000"}},
			},
			err: nil,
		},
		//
		{
			name: "duplicates-in-app-order",
			formula: &CanadianContraFormula{
				ApplicationOrder: []CreditRule{{Source: "1000"}, {Source: "1000"}},
			},
			err: ErrDupCreditSource,
		},
		//
		{
			name: "unknown-income-creditor",
			formula: &CanadianContraFormula{
				CreditsFromIncome: map[finance.IncomeSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2222"},
				},
				ApplicationOrder: []CreditRule{{Source: "1000"}, {Source: "2000"}},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "unknown-deduc-creditor",
			formula: &CanadianContraFormula{
				CreditsFromDeduction: map[finance.DeductionSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2222"},
				},
				ApplicationOrder: []CreditRule{{Source: "1000"}, {Source: "2000"}},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "unknown-misc-creditor",
			formula: &CanadianContraFormula{
				CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
					123: testCreditor{onSource: "1000"},
					456: testCreditor{onSource: "2222"},
				},
				ApplicationOrder: []CreditRule{{Source: "1000"}, {Source: "2000"}},
			},
			err: ErrUnknownCreditSource,
		},
		//
		{
			name: "unknown-persistent-crSrc",
			formula: &CanadianContraFormula{
				PersistentCredits: map[string]float64{"2222": 0.0},
				ApplicationOrder:  []CreditRule{{Source: "1000"}, {Source: "2000"}},
			},
			err: ErrUnknownCreditSource,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case%d-%s", i, c.name), func(t *testing.T) {

			err := c.formula.Validate()
			if errors.Cause(err) != c.err {
				t.Errorf("unexpected error\nwant: %v\n got: %v", c.err, err)
			}

		})
	}
}

func TestCanadianContraFormula_creditsFromIncSrcs(t *testing.T) {

	cf := CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 115,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 95,
			},
			111: testCreditor{
				onSource:     "3000",
				onTaxCredits: 0,
			},
		},
		ApplicationOrder: []CreditRule{
			{Source: "1000", Type: CrRuleTypeCashable},
			{Source: "2000", Type: CrRuleTypeNotCarryForward},
		},
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddIncome(123, 15000) // has creditor
	finances.AddIncome(111, 20000) // zero credits
	finances.AddIncome(999, 8000)  // no creditor

	actual := cf.creditsFromIncSrcs(finances, 0)
	expected := []*creditBySource{&creditBySource{amount: 115, source: "1000"}}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_creditsFromDeducSrcs(t *testing.T) {

	cf := CanadianContraFormula{
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 115,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 95,
			},
			111: testCreditor{
				onSource:     "3000",
				onTaxCredits: 0,
			},
		},
		ApplicationOrder: []CreditRule{
			{Source: "1000", Type: CrRuleTypeCashable},
			{Source: "2000", Type: CrRuleTypeNotCarryForward},
		},
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddDeduction(123, 15000) // has creditor
	finances.AddDeduction(111, 20000) // zero credits
	finances.AddDeduction(999, 8000)  // no creditor

	actual := cf.creditsFromDeducSrcs(finances, 0)
	expected := []*creditBySource{&creditBySource{amount: 115, source: "1000"}}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_creditsFromMiscSrcs(t *testing.T) {

	cf := CanadianContraFormula{
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			123: testCreditor{
				onSource:     "1000",
				onTaxCredits: 115,
			},
			456: testCreditor{
				onSource:     "2000",
				onTaxCredits: 95,
			},
			111: testCreditor{
				onSource:     "3000",
				onTaxCredits: 0,
			},
		},
		ApplicationOrder: []CreditRule{
			{Source: "1000", Type: CrRuleTypeCashable},
			{Source: "2000", Type: CrRuleTypeNotCarryForward},
		},
	}

	finances := finance.NewEmptyIndividualFinances(2019)
	finances.AddMiscAmount(123, 15000) // has creditor
	finances.AddMiscAmount(111, 20000) // zero credits
	finances.AddMiscAmount(999, 8000)  // no creditor

	actual := cf.creditsFromMiscSrcs(finances, 0)
	expected := []*creditBySource{&creditBySource{amount: 115, source: "1000"}}

	diff := deep.Equal(actual, expected)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

}

func TestCanadianContraFormula_checkMiscSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: "123"}
	cf := &CanadianContraFormula{
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkMiscSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[string]struct{}{"123": struct{}{}}

	cf.CreditsFromMiscAmounts[9] = nil
	err = cf.checkMiscSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromMiscAmounts, 9)
	err = cf.checkMiscSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_checkDeducSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: "123"}
	cf := &CanadianContraFormula{
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkDeducSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[string]struct{}{"123": struct{}{}}

	cf.CreditsFromDeduction[9] = nil
	err = cf.checkDeducSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromDeduction, 9)
	err = cf.checkDeducSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_checkIncSrcCreditorsInSet(t *testing.T) {

	creditor := testCreditor{onSource: "123"}
	cf := &CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			321: creditor,
		},
	}

	err := cf.checkIncSrcCreditorsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[string]struct{}{"123": struct{}{}}

	cf.CreditsFromIncome[9] = nil
	err = cf.checkIncSrcCreditorsInSet(set)
	if errors.Cause(err) != ErrNoCreditor {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrNoCreditor, err)
	}

	delete(cf.CreditsFromIncome, 9)
	err = cf.checkIncSrcCreditorsInSet(set)
	if err != nil {
		t.Errorf("unexpected error\nwant: %v\n got: %v", nil, err)
	}

}

func TestCanadianContraFormula_checkPersistentCrSrcsInSet(t *testing.T) {

	cf := &CanadianContraFormula{
		PersistentCredits: map[string]float64{"123": 0.0},
	}

	err := cf.checkPersistentCrSrcsInSet(nil)
	if errors.Cause(err) != ErrUnknownCreditSource {
		t.Errorf("unexpected error\nwant: %v\n got: %v", ErrUnknownCreditSource, err)
	}

	set := map[string]struct{}{"123": struct{}{}}

	err = cf.checkPersistentCrSrcsInSet(set)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestCanadianContraFormula_orderCredits(t *testing.T) {

	cf := &CanadianContraFormula{
		ApplicationOrder: []CreditRule{
			{Source: "1000"},
			{Source: "2000"},
			{Source: "3000"},
			{Source: "5000"},
			{Source: "4000"},
		},
	}
	err := cf.Validate()
	if err != nil {
		t.Fatal(err)
	}

	credits := []*taxCredit{
		&taxCredit{rule: CreditRule{Source: "3000"}},
		&taxCredit{rule: CreditRule{Source: "2000"}},
		&taxCredit{rule: CreditRule{Source: "5000"}},
		&taxCredit{rule: CreditRule{Source: "1000"}},
		&taxCredit{rule: CreditRule{Source: "4000"}},
	}

	cf.orderCredits(credits)

	expectedOrder := []*taxCredit{
		&taxCredit{rule: CreditRule{Source: "1000"}},
		&taxCredit{rule: CreditRule{Source: "2000"}},
		&taxCredit{rule: CreditRule{Source: "3000"}},
		&taxCredit{rule: CreditRule{Source: "5000"}},
		&taxCredit{rule: CreditRule{Source: "4000"}},
	}

	diff := deep.Equal(credits, expectedOrder)
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}

	var emptyGrp []*taxCredit
	cf.orderCredits(emptyGrp)
	diff = deep.Equal(emptyGrp, []*taxCredit(nil))
	if diff != nil {
		t.Error("actual does not match expected\n" + strings.Join(diff, "\n"))
	}
}

func TestCanadianContraFormula_Clone(t *testing.T) {

	var original *CanadianContraFormula
	if original.Clone() != nil {
		t.Error("cloning nil contra formula should return nil")
	}

	original = &CanadianContraFormula{
		CreditsFromIncome: map[finance.IncomeSource]Creditor{
			123: testCreditor{onSource: "1000"},
			456: testCreditor{onSource: "2000"},
		},
		CreditsFromDeduction: map[finance.DeductionSource]Creditor{
			123: testCreditor{onSource: "1000"},
			456: testCreditor{onSource: "2000"},
		},
		CreditsFromMiscAmounts: map[finance.MiscSource]Creditor{
			123: testCreditor{onSource: "1000"},
			456: testCreditor{onSource: "2000"},
		},
		PersistentCredits: map[string]float64{"1000": 0.0},
		ApplicationOrder: []CreditRule{
			{Source: "1000"},
			{Source: "2000"},
		},
	}

	err := original.Validate()
	if err != nil {
		t.Fatal(err)
	}

	clone := original.Clone()
	err = clone.Validate()
	if err != nil {
		t.Fatal(err)
	}

	original.ApplicationOrder[0] = CreditRule{}
	original.ApplicationOrder[1] = CreditRule{}
	original.CreditsFromIncome[123] = testCreditor{onSource: "1111"}
	original.CreditsFromDeduction[123] = testCreditor{onSource: "1111"}
	original.CreditsFromMiscAmounts[123] = testCreditor{onSource: "1111"}
	original.PersistentCredits = map[string]float64{"111": 0.0}
	err = original.Validate()
	if err == nil {
		t.Fatal("invalid formula did not return error when validated")
	}
	err = clone.Validate()
	if err != nil {
		t.Fatal("changes to original contra formula should not affect clone")
	}
}

func TestCanadianContraFormula_makeControlledTaxCreditsFrom(t *testing.T) {
	t.Skip("TODO")
}

func TestCanadianContraFormula_NumFieldsUnchanged(t *testing.T) {

	dummy := CanadianContraFormula{}
	s := reflect.ValueOf(&dummy).Elem()
	if s.NumField() != 5 {
		t.Fatal(
			"number of struct fields changed. Please update the constructor and the " +
				"clone method of this type as well as associated test. Next, update " +
				"this test with the new number of fields",
		)
	}
}
