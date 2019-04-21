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

	finances := core.NewEmptyIndividualFinances(2018)
	calculator.SetFinances(finances)

	finances.AddIncome(core.IncSrcEarned, 170000.0)
	finances.AddIncome(core.IncSrcCapitalGainCA, 20000)
	finances.AddIncome(core.IncSrcTFSA, 12000)
	finances.AddDeduction(core.DeducSrcRRSP, 10000)

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

	f1 := core.NewEmptyIndividualFinances(2017)
	f2 := core.NewEmptyIndividualFinances(2017)

	f1.AddIncome(core.IncSrcEarned, 109500.0)
	f1.AddIncome(core.IncSrcCapitalGainCA, 1000)

	f2.AddIncome(core.IncSrcEarned, 14750)
	f2.AddIncome(core.IncSrcTFSA, 32000)
	f1.AddIncome(core.IncSrcCapitalGainCA, 500)
	f2.AddDeduction(core.DeducSrcRRSP, 15000)

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

	finances := core.NewEmptyIndividualFinances(2018)
	finances.AddIncome(core.IncSrcEarned, 100000.0)
	finances.RRSPContributionRoom = 15000.0

	calculator.SetFinances(finances)
	taxRecievable, err := calculator.TaxRefund(15000.0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%.2f", taxRecievable) // Output: 5182.74
}
