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

	taxParamsProv := IncomeTaxParams{
		Year:   2018,
		Region: history.BC,
	}

	finNums := calc.Finances{TaxableAmount: 170000.0}

	c, err := NewIncomeTaxCalculatorAgg(finNums, taxParamsProv, taxParamsFed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aggTax := c.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 52819.71
}

func ExampleNewIncomeTaxCalculatorAgg_Update() {

	taxParamsFed := IncomeTaxParams{
		Year:   2018,
		Region: history.Canada,
	}

	taxParamsProv := IncomeTaxParams{
		Year:   2018,
		Region: history.BC,
	}

	finNums := calc.Finances{TaxableAmount: 170000.0}

	c, err := NewIncomeTaxCalculatorAgg(finNums, taxParamsProv, taxParamsFed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newFinNums := calc.Finances{TaxableAmount: 20000.0}
	c.Update(newFinNums)

	aggTax := c.Calc()
	fmt.Printf("%.2f", aggTax) // Output: 1713.80

}
