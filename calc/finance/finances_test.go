package finance

import "testing"

func TestFamilyFinances(t *testing.T) {

	spouse1 := &IndividualFinances{
		Income:     IncomeBySource{IncSrcEarned: 10},
		Deductions: DeductionBySource{DeducSrcRRSP: 20},
	}

	spouse2 := &IndividualFinances{
		Income:     IncomeBySource{IncSrcEarned: 15},
		Deductions: DeductionBySource{DeducSrcRRSP: 50},
	}

	finances := HouseholdFinances{spouse1, spouse2}

	actualIncome := finances.TotalIncome()
	expectedIncome := float64(10 + 15)
	if actualIncome != expectedIncome {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedIncome, actualIncome,
		)
	}

	actualDeductions := finances.TotalDeductions()
	expectedDeductions := float64(20 + 50)
	if actualDeductions != expectedDeductions {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedDeductions, actualDeductions,
		)
	}
}
