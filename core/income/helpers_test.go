package income

import "github.com/malkhamis/quantax/core"

type testIncomeDeductor struct {
	totalIncome     float64
	totalDeductions float64
	totalMiscAmount float64
	incomeSrcs      map[core.FinancialSource]struct{}
	deducSrcs       map[core.FinancialSource]struct{}
	miscSrcs        map[core.FinancialSource]struct{}
	version         uint64
}

func (tid testIncomeDeductor) TotalIncome(sources ...core.FinancialSource) float64 {
	return tid.totalIncome
}
func (tid testIncomeDeductor) TotalDeductions(sources ...core.FinancialSource) float64 {
	return tid.totalDeductions
}
func (tid testIncomeDeductor) MiscAmount(sources ...core.FinancialSource) float64 {
	return tid.totalMiscAmount
}
func (tid testIncomeDeductor) IncomeSources() map[core.FinancialSource]struct{} {
	return tid.incomeSrcs
}
func (tid testIncomeDeductor) DeductionSources() map[core.FinancialSource]struct{} {
	return tid.deducSrcs
}
func (tid testIncomeDeductor) MiscSources() map[core.FinancialSource]struct{} {
	return tid.miscSrcs
}
func (tid testIncomeDeductor) Version() uint64 {
	return tid.version
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
