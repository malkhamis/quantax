package history

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/malkhamis/quantax/calc/finance/income"
)

var (
	incomeRecipeNet  income.Recipe
	incomeRecipeAFNI income.Recipe
)

// TODO: add more adjusters
func initIncomeRecipes() {

	zeroAdjuster := income.WeightedAdjuster(0.0)
	incomeRecipeNet = income.Recipe{
		IncomeAdjusters: map[finance.IncomeSource]income.Adjuster{
			finance.IncSrcTFSA:          zeroAdjuster,
			finance.IncSrcCapitalGainCA: income.WeightedAdjuster(0.5),
		},
	}

	incomeRecipeAFNI = *((&incomeRecipeNet).Clone())
	incomeRecipeAFNI.IncomeAdjusters[finance.IncSrcUCCB] = zeroAdjuster
	incomeRecipeAFNI.IncomeAdjusters[finance.IncSrcRDSP] = zeroAdjuster
}
