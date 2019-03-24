package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/history"
)

func ExampleNewIncomeTaxCalculator_Calc() {

	cfgFed := CalculatorConfig{
		Year:   2018,
		Region: history.Canada,
	}

	cfgBC := CalculatorConfig{
		Year:   2018,
		Region: history.BC,
	}

	finances := calc.IndividualFinances{Income: 170000.0}

	calcFed, err := NewIncomeTaxCalculator(finances, cfgFed)
	if err != nil {
		fmt.Println(err)
		return
	}

	calcBC, err := NewIncomeTaxCalculator(finances, cfgBC)
	if err != nil {
		fmt.Println(err)
		return
	}

	aggTax := calcFed.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 52819.71
}

func ExampleNewIncomeTaxCalculator_UpdateFinances() {

	cfgFed := CalculatorConfig{
		Year:   2018,
		Region: history.Canada,
	}

	cfgBC := CalculatorConfig{
		Year:   2018,
		Region: history.BC,
	}

	finances := calc.IndividualFinances{Income: 170000.0}

	calcFed, err := NewIncomeTaxCalculator(finances, cfgFed)
	if err != nil {
		fmt.Println(err)
		return
	}

	calcBC, err := NewIncomeTaxCalculator(finances, cfgBC)
	if err != nil {
		fmt.Println(err)
		return
	}

	newFinNums := calc.IndividualFinances{Income: 20000.0}
	calcFed.UpdateFinances(newFinNums)
	calcBC.UpdateFinances(newFinNums)

	aggTax := calcFed.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 1713.80

}

func ExampleNewChildBenefitCalculator() {

	cfg := CBCalculatorConfig{
		CalculatorConfig: CalculatorConfig{
			Year:   2018,
			Region: history.Canada,
		},
		Children: []calc.Person{
			{Name: "Hafedh", AgeMonths: 12*6 - 2},
			{Name: "Zaid", AgeMonths: 12 * 3},
		},
	}

	finances := calc.FamilyFinances{
		calc.IndividualFinances{Income: 83000},
		calc.IndividualFinances{Income: 0},
	}

	ccb, err := NewChildBenefitCalculator(finances, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	total := ccb.Calc()
	fmt.Printf("%.2f", total) // Output: 8809.20

}
