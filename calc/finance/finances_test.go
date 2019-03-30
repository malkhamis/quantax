package finance

import "testing"

func TestFamilyFinances(t *testing.T) {

	spouse1 := IndividualFinances{Income: 10, Deductions: 20}
	spouse2 := IndividualFinances{Income: 15, Deductions: 50}
	finances := FamilyFinances{spouse1, spouse2}

	actualIncome := finances.Income()
	expectedIncome := float64(10 + 15)
	if actualIncome != expectedIncome {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedIncome, actualIncome,
		)
	}

	actualDeductions := finances.Deductions()
	expectedDeductions := float64(20 + 50)
	if actualDeductions != expectedDeductions {
		t.Errorf(
			"unexpected income total\nwant: %.2f\n got: %.2f",
			expectedDeductions, actualDeductions,
		)
	}

	actualSpouse1, actualSpouse2 := finances.Split()
	if actualSpouse1 != spouse1 {
		t.Errorf(
			"actual split finances is not equal to expected\nwant: %v\n got: %v",
			spouse1, actualSpouse1,
		)
	}
	if actualSpouse2 != spouse2 {
		t.Errorf(
			"actual split finances is not equal to expected\nwant: %v\n got: %v",
			spouse2, actualSpouse2,
		)
	}
}
