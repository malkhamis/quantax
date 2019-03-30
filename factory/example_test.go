package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/finance"
)

func ExampleNewTaxCalcFactory() {

	f, err := NewTaxCalcFactory(2018, Canada, BC)
	if err != nil {
		fmt.Println(err)
		return
	}

	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := finance.IndividualFinances{Income: 170000.0}
	aggTax := calculator.Calc(finances)
	fmt.Printf("%.2f", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitCalcFactory() {

	f, err := NewChildBenefitCalcFactory(2017, Canada, BC)
	if err != nil {
		fmt.Println(err)
		return
	}

	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	children := []calc.Person{{Name: "A", AgeMonths: 3}, {Name: "B", AgeMonths: 3}}
	calculator.SetBeneficiaries(children...)

	finances := finance.FamilyFinances{
		finance.IndividualFinances{Income: 110000},
		finance.IndividualFinances{Income: 0},
	}
	total := calculator.Calc(finances)

	fmt.Printf("%.2f", total) // Output: 6742.54
}
