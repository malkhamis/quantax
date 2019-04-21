package income

import "github.com/malkhamis/quantax/calc/finance"

type testIncomeDeductor struct {
	totalIncome     float64
	totalDeductions float64
	totalMiscAmount float64
	incomeSrcs      finance.IncomeSourceSet
	deducSrcs       finance.DeductionSourceSet
	miscSrcs        finance.MiscSourceSet
}

func (tid testIncomeDeductor) TotalIncome(sources ...finance.IncomeSource) float64 {
	return tid.totalIncome
}
func (tid testIncomeDeductor) TotalDeductions(sources ...finance.DeductionSource) float64 {
	return tid.totalDeductions
}
func (tid testIncomeDeductor) MiscAmount(sources ...finance.MiscSource) float64 {
	return tid.totalMiscAmount
}
func (tid testIncomeDeductor) IncomeSources() finance.IncomeSourceSet {
	return tid.incomeSrcs
}
func (tid testIncomeDeductor) DeductionSources() finance.DeductionSourceSet {
	return tid.deducSrcs
}
func (tid testIncomeDeductor) MiscSources() finance.MiscSourceSet {
	return tid.miscSrcs
}

type testAdjuster struct {
	adjusted float64
}

func (ta testAdjuster) Adjusted(_ float64) float64 {
	return ta.adjusted
}

func (ta testAdjuster) Clone() Adjuster {
	return ta
}
