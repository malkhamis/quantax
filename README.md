# quantax
quantax is a quantitative-based solution for maximizing tax returns and benefits, or that was the aim of it :)
Quantax was a pet project and I was supposed to add another layer on the top of fundamental building blocks, which would apply a Monte Carlo(-like?) analysis. Said analysis would tell me how much I should contribute to the RRSP to EFFICIENTLY maximize the amount I get from tax refund AND child benefits; and so on.. However, that extra layer never saw the daylight and I had to shift gears to something else.

Working on this project taught me a great deal about the Canadian tax system, and how amusing it is to write software that deals with taxes given all the complex rules and twisted, non-sensical dependencies between such and such! I kept the project for my future self to introspect and say "did I write that?!?!?"

# Sample Driver Program
```
package main

import (
	"fmt"
	"log"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/factory"
)

func main() {

	ff := factory.NewFinanceFactory()

	rcf := factory.NewRRSPFactory(factory.RRSPFactoryConfig{
		Year: 2022,
		RRSPRegion: core.RegionCA,
		TaxRegions: []core.Region{core.RegionBC, core.RegionCA},
	})

	rrspCalc, err := rcf.NewCalculator()
	if err != nil {
		log.Fatal(err)
	}

		initAmounts := map[core.FinancialSource]float64{core.IncSrcEarned: 300000}
		myFinances := ff.NewHouseholdFinancesForCouple(initAmounts, map[core.FinancialSource]float64{})
		rrspCalc.SetFinances(myFinances, nil)
		refund, _ := rrspCalc.TaxRefund(43300.)
		fmt.Println(refund)
}
```
