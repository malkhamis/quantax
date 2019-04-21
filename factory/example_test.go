package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/core/finance"
	"github.com/malkhamis/quantax/core/human"
)

func ExampleNewTaxFactory() {

	f := NewTaxFactory(2018, Canada, BC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := finance.NewEmptyIndividualFinances(2018)
	calculator.SetFinances(finances)

	finances.AddIncome(finance.IncSrcEarned, 170000.0)
	finances.AddIncome(finance.IncSrcCapitalGainCA, 20000)
	finances.AddIncome(finance.IncSrcTFSA, 12000)
	finances.AddDeduction(finance.DeducSrcRRSP, 10000)

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

	f1 := finance.NewEmptyIndividualFinances(2017)
	f2 := finance.NewEmptyIndividualFinances(2017)

	f1.AddIncome(finance.IncSrcEarned, 109500.0)
	f1.AddIncome(finance.IncSrcCapitalGainCA, 1000)

	f2.AddIncome(finance.IncSrcEarned, 14750)
	f2.AddIncome(finance.IncSrcTFSA, 32000)
	f1.AddIncome(finance.IncSrcCapitalGainCA, 500)
	f2.AddDeduction(finance.DeducSrcRRSP, 15000)

	finances := finance.NewHouseholdFinances(f1, f2)
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

	finances := finance.NewEmptyIndividualFinances(2018)
	finances.AddIncome(finance.IncSrcEarned, 100000.0)
	finances.RRSPContributionRoom = 15000.0

	calculator.SetFinances(finances)
	taxRecievable, err := calculator.TaxRefund(15000.0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%.2f", taxRecievable) // Output: 5182.74
}
