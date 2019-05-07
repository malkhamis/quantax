package income

import "github.com/malkhamis/quantax/core"

type testFinancer struct {
	_currentTotalAmount int // don't set
	onTotalAmount       []float64

	onIncomeSources    []core.FinancialSource
	onDeductionSources []core.FinancialSource
	onMiscSources      []core.FinancialSource
	onAllSources       []core.FinancialSource
}

func (tf *testFinancer) TotalAmount(sources ...core.FinancialSource) float64 {

	total := tf.onTotalAmount[tf._currentTotalAmount]
	tf._currentTotalAmount++
	return total
}
func (tf *testFinancer) IncomeSources() []core.FinancialSource {
	return tf.onIncomeSources
}
func (tf *testFinancer) DeductionSources() []core.FinancialSource {
	return tf.onDeductionSources
}
func (tf *testFinancer) MiscSources() []core.FinancialSource {
	return tf.onMiscSources
}
func (tf *testFinancer) AllSources() []core.FinancialSource {
	return tf.onAllSources
}
func (tf *testFinancer) Clone() core.FinanceMutator {
	return nil // not needed for this package
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
