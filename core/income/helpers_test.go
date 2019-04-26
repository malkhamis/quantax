package income

import "github.com/malkhamis/quantax/core"

type testFinancer struct {
	_currentTotalAmount int // don't set
	onTotalAmount       []float64

	onIncomeSources    []core.FinancialSource
	onDeductionSources []core.FinancialSource
	onMiscSources      []core.FinancialSource
	onVersion          uint64
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
func (tf *testFinancer) Clone() core.FinanceMutator {
	return nil // not needed for this package
}
func (tf *testFinancer) Version() uint64 {
	return tf.onVersion
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
