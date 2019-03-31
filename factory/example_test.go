package factory

import (
	"fmt"

	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/human"
)

func ExampleNewTaxFactory() {

	f := NewTaxFactory(2018, Canada, BC)
	calculator, err := f.NewCalculator()
	if err != nil {
		fmt.Println(err)
		return
	}

	finances := finance.IndividualFinances{Income: 170000.0}
	aggTax := calculator.Calc(finances)
	fmt.Printf("%.2f", aggTax) // Output: 52821.09
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

	finances := finance.FamilyFinances{
		finance.IndividualFinances{Income: 110000},
		finance.IndividualFinances{Income: 0},
	}
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

	finances := finance.IndividualFinances{
		Income:   100000.0,
		RRSPRoom: 15000.0,
	}
	calculator.SetFinances(finances)
	taxRecievable, err := calculator.TaxRefund(15000.0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%.2f", taxRecievable) // Output: 5182.74
}
