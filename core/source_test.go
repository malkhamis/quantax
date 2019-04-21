package core

import "testing"

func TestFinancialSource(t *testing.T) {

	cases := []struct {
		source   FinancialSource
		isInc    bool
		isDeduc  bool
		isMisc   bool
		isUnkown bool
	}{
		{source: IncSrcRRSP, isInc: true},
		{source: DeducSrcRRSP, isDeduc: true},
		{source: MiscSrcMedical, isMisc: true},
		{source: SrcUnknown, isUnkown: true},
	}

	for _, c := range cases {

		actual := c.source.IsIncomeSource()
		if actual != c.isInc {
			t.Errorf(
				"%v: actual '%v' does not match expected '%v'",
				c, actual, c.isInc,
			)
		}

		actual = c.source.IsDeductionSource()
		if actual != c.isDeduc {
			t.Errorf(
				"%v: actual '%v' does not match expected '%v'",
				c, actual, c.isDeduc,
			)
		}

		actual = c.source.IsMiscSource()
		if actual != c.isMisc {
			t.Errorf(
				"%v: actual '%v' does not match expected '%v'",
				c, actual, c.isMisc,
			)
		}

		actual = c.source.IsUnknownSource()
		if actual != c.isUnkown {
			t.Errorf(
				"%v: actual '%v' does not match expected '%v'",
				c, actual, c.isUnkown,
			)
		}

	}
}
