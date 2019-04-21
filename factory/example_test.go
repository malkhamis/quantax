package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
)

func ExampleNewTaxFactory() {

	f := NewTaxFactory(2018, Canada, BC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := core.NewEmptyIndividualFinances()
	calculator.SetFinances(finances)

	finances.AddAmount(core.IncSrcEarned, 170000.0)
	finances.AddAmount(core.IncSrcCapitalGainCA, 20000)
	finances.AddAmount(core.IncSrcTFSA, 12000)
	finances.AddAmount(core.DeducSrcRRSP, 10000)

	aggTax, _ := calculator.TaxPayable()
	fmt.Printf("%.2f\n", aggTax) // Output: 52821.09
}

func ExampleNewChildBenefitFactory() {

	f := NewChildBenefitFactory(2017, Canada, BC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	children := []human.Person{{Name: "A", AgeMonths: 3}, {Name: "B", AgeMonths: 3}}
	calculator.SetBeneficiaries(children...)

	f1 := core.NewEmptyIndividualFinances()
	f2 := core.NewEmptyIndividualFinances()

	f1.AddAmount(core.IncSrcEarned, 109500.0)
	f1.AddAmount(core.IncSrcCapitalGainCA, 1000)

	f2.AddAmount(core.IncSrcEarned, 14750)
	f2.AddAmount(core.IncSrcTFSA, 32000)
	f1.AddAmount(core.IncSrcCapitalGainCA, 500)
	f2.AddAmount(core.DeducSrcRRSP, 15000)

	finances := core.NewHouseholdFinances(f1, f2)
	total := calculator.Calc(finances)

	fmt.Printf("%.2f", total) // Output: 6742.54
}

func ExampleNewRRSPFactory() {

	config := RRSPFactoryConfig{
		Year:       2018,
		RRSPRegion: Canada,
		TaxRegions: []Region{Canada, BC},
	}

	f := NewRRSPFactory(config)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := core.NewEmptyIndividualFinances()
	finances.AddAmount(core.IncSrcEarned, 100000.0)
	finances.SetRRSPAmounts(core.RRSPAmounts{ContributionRoom: 15000.0})

	calculator.SetFinances(finances)
	taxRecievable, err := calculator.TaxRefund(15000.0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%.2f", taxRecievable) // Output: 5182.74
}
