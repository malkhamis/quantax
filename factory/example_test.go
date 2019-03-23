package factory

import (
	"fmt"
	"os"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/history"
)

func ExampleNewIncomeTaxCalculatorAgg_Calc() {

	taxParamsFed := IncomeTaxParams{
		Year:   2018,
		Region: history.Canada,
	}

	taxParamsBC := IncomeTaxParams{
		Year:   2018,
		Region: history.BC,
	}

	finNums := calc.IndividualFinances{Income: 170000.0}

	calcFed, err := NewIncomeTaxCalculator(finNums, taxParamsFed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	calcBC, err := NewIncomeTaxCalculator(finNums, taxParamsBC)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aggTax := calcFed.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 52819.71
}

func ExampleNewIncomeTaxCalculator_Update() {

	taxParamsFed := IncomeTaxParams{
		Year:   2018,
		Region: history.Canada,
	}

	taxParamsBC := IncomeTaxParams{
		Year:   2018,
		Region: history.BC,
	}

	finNums := calc.IndividualFinances{Income: 170000.0}

	calcFed, err := NewIncomeTaxCalculator(finNums, taxParamsFed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	calcBC, err := NewIncomeTaxCalculator(finNums, taxParamsBC)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newFinNums := calc.IndividualFinances{Income: 20000.0}
	calcFed.UpdateFinances(newFinNums)
	calcBC.UpdateFinances(newFinNums)

	aggTax := calcFed.Calc() + calcBC.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 1713.80

}
