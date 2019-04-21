package history

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/income"
)

var (
	incomeRecipeNet  income.Recipe
	incomeRecipeAFNI income.Recipe
)

// TODO: add more adjusters
func initIncomeRecipes() {

	zeroAdjuster := income.WeightedAdjuster(0.0)
	incomeRecipeNet = income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcTFSA:          zeroAdjuster,
			core.IncSrcCapitalGainCA: income.WeightedAdjuster(0.5),
		},
	}

	incomeRecipeAFNI = *((&incomeRecipeNet).Clone())
	incomeRecipeAFNI.IncomeAdjusters[core.IncSrcUCCB] = zeroAdjuster
	incomeRecipeAFNI.IncomeAdjusters[core.IncSrcRDSP] = zeroAdjuster
}
