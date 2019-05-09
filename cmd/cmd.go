package main

import (
	"fmt"
	"os"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/factory"
)

func main() {

	initAmounts := map[core.FinancialSource]float64{
		core.IncSrcEligibleDividendsCA: 119000,
		core.MiscSrcTuition:            0,
	}
	ff := factory.NewFinanceFactory()
	myFinances := ff.NewHouseholdFinancesForSingle(initAmounts)

	tf := factory.NewTaxFactory(2018, core.RegionCA, core.RegionBC)
	calc, err := tf.NewCalculator()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	calc.SetFinances(myFinances, nil)
	taxPaid, _, credits := calc.TaxPayable()
	fmt.Println(taxPaid)

	for _, cr := range credits {
		fmt.Println(cr.Source(), cr.Region())
		initial, used, remain := cr.Amounts()
		fmt.Printf("initial: %.2f\tused: %.2f\tremaining: %.2f\n\n", initial, used, remain)
	}
}
