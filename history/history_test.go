package history

import (
	"testing"

	"github.com/pkg/errors"
)

func TestGetTaxParams(t *testing.T) {

	_, err := GetTaxParams(2018, BC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTaxParams_Errors(t *testing.T) {

	_, err := GetTaxParams(2018, Jurisdiction("OhCanada"))
	if errors.Cause(err) != ErrJurisdictionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrJurisdictionNotExist, err)
	}

	_, err = GetTaxParams(2108, Canada)
	if errors.Cause(err) != ErrParamsNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrParamsNotExist, err)
	}

}

func TestGetRRSPParams(t *testing.T) {

	_, err := GetRRSPParams(2018, Canada)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRRSPParams_Errors(t *testing.T) {

	_, err := GetRRSPParams(2018, Jurisdiction("OhCanada"))
	if errors.Cause(err) != ErrJurisdictionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrJurisdictionNotExist, err)
	}

	_, err = GetRRSPParams(2108, Canada)
	if errors.Cause(err) != ErrParamsNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrParamsNotExist, err)
	}

}

func TestGetChildBenefitParams(t *testing.T) {

	_, err := GetChildBenefitParams(2017, Canada)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetChildBenefitParams_Errors(t *testing.T) {

	_, err := GetChildBenefitParams(2018, Jurisdiction("OhCanada"))
	if errors.Cause(err) != ErrJurisdictionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrJurisdictionNotExist, err)
	}

	_, err = GetChildBenefitParams(2108, Canada)
	if errors.Cause(err) != ErrParamsNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrParamsNotExist, err)
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
