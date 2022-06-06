package history

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/income"
)

var (
	incomeRecipeNetCA2022 = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA:          income.WeightedAdjuster(0.5),
			core.IncSrcTFSA:                   income.WeightedAdjuster(0.0),
		},
	}

	incomeRecipeNetCA2019 = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA:          income.WeightedAdjuster(0.5),
			core.IncSrcEligibleDividendsCA:    income.WeightedAdjuster(1.38),
			core.IncSrcNonEligibleDividendsCA: income.WeightedAdjuster(1.15),
			core.IncSrcTFSA:                   income.WeightedAdjuster(0.0),
		},
	}

	incomeRecipeNetCA2018 = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA:          income.WeightedAdjuster(0.5),
			core.IncSrcEligibleDividendsCA:    income.WeightedAdjuster(1.38),
			core.IncSrcNonEligibleDividendsCA: income.WeightedAdjuster(1.16),
			core.IncSrcTFSA:                   income.WeightedAdjuster(0.0),
		},
	}

	incomeRecipeAFNICA2019 = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA:          income.WeightedAdjuster(0.5),
			core.IncSrcEligibleDividendsCA:    income.WeightedAdjuster(1.38),
			core.IncSrcNonEligibleDividendsCA: income.WeightedAdjuster(1.15),
			core.IncSrcTFSA:                   income.WeightedAdjuster(0.0),
			core.IncSrcUCCB:                   income.WeightedAdjuster(0.0),
			core.IncSrcRDSP:                   income.WeightedAdjuster(0.0),
		},
	}

	incomeRecipeAFNICA2018 = &income.Recipe{
		IncomeAdjusters: map[core.FinancialSource]income.Adjuster{
			core.IncSrcCapitalGainCA:          income.WeightedAdjuster(0.5),
			core.IncSrcEligibleDividendsCA:    income.WeightedAdjuster(1.38),
			core.IncSrcNonEligibleDividendsCA: income.WeightedAdjuster(1.16),
			core.IncSrcTFSA:                   income.WeightedAdjuster(0.0),
			core.IncSrcUCCB:                   income.WeightedAdjuster(0.0),
			core.IncSrcRDSP:                   income.WeightedAdjuster(0.0),
		},
	}
)
