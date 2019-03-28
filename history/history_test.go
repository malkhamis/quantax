package history

import (
	"testing"

	"github.com/pkg/errors"
)

func TestGetTaxFormula(t *testing.T) {

	_, err := GetTaxFormula(2018, BC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTaxFormula_Errors(t *testing.T) {

	_, err := GetTaxFormula(2018, Jurisdiction("OhCanada"))
	if errors.Cause(err) != ErrJurisdictionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrJurisdictionNotExist, err)
	}

	_, err = GetTaxFormula(2108, Canada)
	if errors.Cause(err) != ErrFormulaNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFormulaNotExist, err)
	}

}

func TestGetChildBenefitFormula(t *testing.T) {

	_, err := GetChildBenefitFormula(2017, Canada)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetChildBenefitFormula_Errors(t *testing.T) {

	_, err := GetChildBenefitFormula(2018, Jurisdiction("OhCanada"))
	if errors.Cause(err) != ErrJurisdictionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrJurisdictionNotExist, err)
	}

	_, err = GetChildBenefitFormula(2108, Canada)
	if errors.Cause(err) != ErrFormulaNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrFormulaNotExist, err)
	}

}

func TestPanicIfError(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	panicIfError(errors.New(""))
}
