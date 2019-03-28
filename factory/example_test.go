package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc"
)

func ExampleNewTaxCalcFactory() {

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
