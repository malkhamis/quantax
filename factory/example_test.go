package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc"
)

func ExampleNewIncomeTaxCalculator_Calc() {

	calcFactoryFed, err := NewTaxCalcFactory(Options{Year: 2018, Region: Canada})
	if err != nil {
		fmt.Println(err)
		return
	}
	calcFactoryBC, err := NewTaxCalcFactory(Options{Year: 2018, Region: BC})
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := calc.IndividualFinances{Income: 170000.0}

	calcFed, err := calcFactoryFed.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}
	calcBC, err := calcFactoryBC.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}

	aggTax := calcFed.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitCalculatorForBC() {

	opts := Options{
		Year:   2017,
		Region: BC,
	}
	children := []calc.Person{
		{Name: "A", AgeMonths: 3},
		{Name: "B", AgeMonths: 3},
	}

	calcFactoryBC, err := NewChildBenefitCalcFactory(opts, children...)
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := calc.FamilyFinances{
		calc.IndividualFinances{Income: 110000},
		calc.IndividualFinances{Income: 0},
	}

	calculator, err := calcFactoryBC.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}

	total := calculator.Calc()
	fmt.Printf("%.2f", total) // Output: 1056.00

}

func ExampleNewChildBenefitCalculatorForCanada() {

	opts := Options{
		Year:   2017,
		Region: Canada,
	}
	children := []calc.Person{
		{Name: "A", AgeMonths: 3},
		{Name: "B", AgeMonths: 3},
	}

	calcFactoryBC, err := NewChildBenefitCalcFactory(opts, children...)
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := calc.FamilyFinances{
		calc.IndividualFinances{Income: 110000},
		calc.IndividualFinances{Income: 0},
	}

	calculator, err := calcFactoryBC.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}

	total := calculator.Calc()
	fmt.Printf("%.2f", total) // Output: 5686.54

}
