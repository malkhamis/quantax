package income

import "github.com/malkhamis/quantax/calc/finance"

type testFormula struct {
	incomeAdjusters map[finance.IncomeSource]Adjuster
	deducAdjusters  map[finance.DeductionSource]Adjuster
	err             error
}

func (tf testFormula) IncomeAdjusters() map[finance.IncomeSource]Adjuster {
	return tf.incomeAdjusters
}

func (tf testFormula) DeductionAdjusters() map[finance.DeductionSource]Adjuster {
	return tf.deducAdjusters
}

func (tf testFormula) Clone() Formula {
	return tf
}

func (tf testFormula) Validate() error {
	return tf.err
}

type testIncomeDeductor struct {
	totalIncome     float64
	totalDeductions float64
	incomeSrcs      finance.IncomeSourceSet
	deducSrcs       finance.DeductionSourceSet
}

func (tid testIncomeDeductor) TotalIncome(sources ...finance.IncomeSource) float64 {
	return tid.totalIncome
}
func (tid testIncomeDeductor) TotalDeductions(sources ...finance.DeductionSource) float64 {
	return tid.totalDeductions
}
func (tid testIncomeDeductor) IncomeSources() finance.IncomeSourceSet {
	return tid.incomeSrcs
}
func (tid testIncomeDeductor) DeductionSources() finance.DeductionSourceSet {
	return tid.deducSrcs
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
