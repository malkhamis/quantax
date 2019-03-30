package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc"
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

	finances := calc.IndividualFinances{Income: 170000.0}
	aggTax := calculator.Calc(finances)
	fmt.Printf("%.2f", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitCalcFactory() {

	children := []calc.Person{
		{Name: "A", AgeMonths: 3},
		{Name: "B", AgeMonths: 3},
	}

	opts := Options{Year: 2017, Region: Canada}
	calcFactoryCanada, err := NewChildBenefitCalcFactory(opts, children...)
	if err != nil {
		fmt.Println(err)
		return
	}

	opts = Options{Year: 2017, Region: BC}
	calcFactoryBC, err := NewChildBenefitCalcFactory(opts, children...)
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := calc.FamilyFinances{
		calc.IndividualFinances{Income: 110000},
		calc.IndividualFinances{Income: 0},
	}

	calcCanada, err := calcFactoryCanada.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}

	calcBC, err := calcFactoryBC.NewCalculator(finances)
	if err != nil {
		fmt.Println(err)
		return
	}

	total := calcCanada.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", total) // Output: 6742.54

}
