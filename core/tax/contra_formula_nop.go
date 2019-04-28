package tax

import "github.com/malkhamis/quantax/core"

var (
	// compile-time check for interface implementation
	_ ContraFormula = (NopContraFormula)(NopContraFormula{})
)

type NopContraFormula struct{}

func (NopContraFormula) Apply(*TaxPayer) []*TaxCredit {
	return nil
}
func (NopContraFormula) FilterAndSort([]core.TaxCredit) []core.TaxCredit {
	return nil
}
func (NopContraFormula) Clone() ContraFormula {
	return NopContraFormula{}
}
func (NopContraFormula) Validate() error {
	return nil
}
func (NopContraFormula) TaxInfo() core.TaxInfo {
	return core.TaxInfo{}
}
