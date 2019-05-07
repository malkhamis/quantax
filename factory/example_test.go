package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/finance"
	"github.com/malkhamis/quantax/core/human"
)

func ExampleNewTaxFactory() {

	f := NewTaxFactory(2018, core.RegionCA, core.RegionBC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	myFinances := finance.NewIndividualFinances()
	finances := finance.NewHouseholdFinances(myFinances, nil)
	calculator.SetFinances(finances, nil)

	myFinances.AddAmount(core.IncSrcEarned, 170000.0)
	myFinances.AddAmount(core.IncSrcCapitalGainCA, 20000)
	myFinances.AddAmount(core.IncSrcTFSA, 12000)
	myFinances.AddAmount(core.DeducSrcRRSP, 10000)

	aggTax, _, _ := calculator.TaxPayable()
	fmt.Printf("%.2f\n", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitFactory() {

	f := NewChildBenefitFactory(2018, core.RegionCA, core.RegionBC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	children := []*human.Person{
		&human.Person{Name: "A", AgeMonths: 3},
		&human.Person{Name: "B", AgeMonths: 3},
	}
	calculator.SetBeneficiaries(children)

	f1 := finance.NewIndividualFinances()
	f2 := finance.NewIndividualFinances()

	f1.AddAmount(core.IncSrcEarned, 109500.0)
	f1.AddAmount(core.IncSrcCapitalGainCA, 1000)

	f2.AddAmount(core.IncSrcEarned, 14750)
	f2.AddAmount(core.IncSrcTFSA, 32000)
	f1.AddAmount(core.IncSrcCapitalGainCA, 500)
	f2.AddAmount(core.DeducSrcRRSP, 15000)

	finances := finance.NewHouseholdFinances(f1, f2)
	calculator.SetFinances(finances)
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

	myFinances := finance.NewIndividualFinances()
	hf := finance.NewHouseholdFinances(myFinances, nil)
	myFinances.AddAmount(core.IncSrcEarned, 100000.0)

	calculator.SetFinances(hf, nil)
	taxRecievable, _ := calculator.TaxRefund(15000.0)

	fmt.Printf("%.2f", taxRecievable) // Output: 5182.74
}
