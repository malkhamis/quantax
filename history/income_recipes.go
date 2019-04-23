package history

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/income"
)

var (
	incomeRecipeNet = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA: income.WeightedAdjuster(0.5),
			core.IncSrcTFSA:          income.WeightedAdjuster(0.0),
		},
	}
	incomeRecipeAFNI = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA: income.WeightedAdjuster(0.5),
			core.IncSrcTFSA:          income.WeightedAdjuster(0.0),
			core.IncSrcUCCB:          income.WeightedAdjuster(0.0),
			core.IncSrcRDSP:          income.WeightedAdjuster(0.0),
		},
	}
)
