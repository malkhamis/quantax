package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
)

func ExampleNewTaxFactory() {

	initAmounts := map[core.FinancialSource]float64{
		core.IncSrcEarned:        170000.0,
		core.IncSrcCapitalGainCA: 10000,
		core.IncSrcTFSA:          12000,
		core.DeducSrcRRSP:        10000,
	}
	myFinances := NewFinanceFactory().NewHouseholdFinancesForSingle(initAmounts)

	taxCalcFactory := NewTaxFactory(2018, core.RegionCA, core.RegionBC)
	calculator, err := taxCalcFactory.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}
	calculator.SetFinances(myFinances, nil)

	// we can mutate finances later without the need to set finances again
	myFinances.MutableSpouseA().AddAmount(core.IncSrcCapitalGainCA, 10000)
	aggTax, _, _ := calculator.TaxPayable()
	fmt.Printf("%.2f\n", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitFactory() {

	children := []*human.Person{
		&human.Person{Name: "A", AgeMonths: 0},
		&human.Person{Name: "B", AgeMonths: 3},
	}

	initAmountsA := map[core.FinancialSource]float64{
		core.IncSrcEarned:        109500.0,
		core.IncSrcCapitalGainCA: 1000,
	}

	initAmountsB := map[core.FinancialSource]float64{
		core.IncSrcEarned: 14750,
		core.IncSrcTFSA:   32000,
		core.DeducSrcRRSP: 15000,
	}

	finFactory := NewFinanceFactory()
	finances := finFactory.NewHouseholdFinancesForCouple(initAmountsA, initAmountsB)

	f := NewChildBenefitFactory(2018, core.RegionCA, core.RegionBC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}
	calculator.SetBeneficiaries(children)
	calculator.SetFinances(finances)

	// we can mutate finances later without the need to set finances again
	finances.MutableSpouseA().AddAmount(core.IncSrcCapitalGainCA, 500)
	total := calculator.BenefitRecievable()

	fmt.Printf("%.2f", total) // Output: 6742.54
}

func ExampleNewRRSPFactory() {

	config := RRSPFactoryConfig{
		Year:       2018,
		RRSPRegion: core.RegionCA,
		TaxRegions: []core.Region{core.RegionCA, core.RegionBC},
	}

	f := NewRRSPFactory(config)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	initAmountsA := map[core.FinancialSource]float64{
		core.IncSrcEarned: 100000.0,
	}
	initAmountsB := map[core.FinancialSource]float64{
		core.IncSrcEarned: 50000.0,
	}

	ff := NewFinanceFactory()
	myFinances := ff.NewHouseholdFinancesForCouple(initAmountsA, initAmountsB)

	calculator.SetFinances(myFinances, nil)
	taxRecievable, _ := calculator.TaxRefund(15000.0)
	fmt.Printf("%.2f\n", taxRecievable)

	calculator.SetTargetSpouseB()
	taxRecievable, _ = calculator.TaxRefund(15000.0)
	fmt.Printf("%.2f\n", taxRecievable)

	// Output:
	// 5182.74
	// 3468.28
}
