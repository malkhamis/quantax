package history

import (
	"testing"

	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

func TestGetTaxParams(t *testing.T) {

	params, err := GetTaxParams(2018, core.RegionBC)
	if err != nil {
		t.Fatal(err)
	}

	if params.Formula == nil {
		t.Fatal("formula must not be nil")
	}

	if params.IncomeRecipe == nil {
		t.Fatal("income recipe must not be nil")
	}
}

func TestGetTaxParams_Errors(t *testing.T) {

	_, err := GetTaxParams(2018, core.Region("OhCanada"))
	if errors.Cause(err) != ErrRegionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrRegionNotExist, err)
	}

	_, err = GetTaxParams(2108, core.RegionCA)
	if errors.Cause(err) != ErrParamsNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrParamsNotExist, err)
	}

}

func TestGetRRSPParams(t *testing.T) {

	_, err := GetRRSPParams(2018, core.RegionCA)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRRSPParams_Errors(t *testing.T) {

	_, err := GetRRSPParams(2018, core.Region("OhCanada"))
	if errors.Cause(err) != ErrRegionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrRegionNotExist, err)
	}

	_, err = GetRRSPParams(2108, core.RegionCA)
	if errors.Cause(err) != ErrParamsNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrParamsNotExist, err)
	}

}

func TestGetChildBenefitParams(t *testing.T) {

	_, err := GetChildBenefitParams(2018, core.RegionCA)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetChildBenefitParams_Errors(t *testing.T) {

	_, err := GetChildBenefitParams(2018, core.Region("OhCanada"))
	if errors.Cause(err) != ErrRegionNotExist {
		t.Fatalf("unexpected error\nwant: %v\n got: %v", ErrRegionNotExist, err)
	}

	_, err = GetChildBenefitParams(2108, core.RegionCA)
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
