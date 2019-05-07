package tax

import "github.com/malkhamis/quantax/core"

var (
	// compile-time check for interface implementation
	_ ContraFormula = (NopContraFormula)(NopContraFormula{})
)

type NopContraFormula struct {
	_ byte // to guarantee a new instance
}

func (NopContraFormula) Apply(*TaxPayer) []*TaxCredit {
	return nil
}
func (NopContraFormula) FilterAndSort(*[]core.TaxCredit) {}
func (NopContraFormula) Clone() ContraFormula {
	return NopContraFormula{}
}
func (NopContraFormula) Validate() error {
	return nil
}
func (NopContraFormula) Year() uint {
	return 0
}
func (NopContraFormula) Region() core.Region {
	return core.Region("")
}
