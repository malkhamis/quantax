package finance

import (
	"strings"
	"testing"
)

func TestIncomeType_String(t *testing.T) {

	if AFNI.String() != "AFNI" {
		t.Error("unexpected stringified iota")
	}

	actualStr := IncomeType(10000).String()
	if !strings.Contains(actualStr, "10000") {
		t.Errorf(
			"unexpected strignfied income type\nwant: %s\n got: %s",
			"IncomeType(10000)", actualStr,
		)
	}

}

func TestIncomeTypeAFNI(t *testing.T) {
	t.Skip("TODO")
}
