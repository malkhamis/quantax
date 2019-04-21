package tax

import "github.com/malkhamis/quantax/core"

var (
	// compile-time check for interface implementation
	_ ContraFormula = (NopContraFormula)(NopContraFormula{})
)

type NopContraFormula struct{}

func (NopContraFormula) Apply(finances *core.IndividualFinances, netIncome float64) []*taxCredit {
	return nil
}

func (NopContraFormula) Clone() ContraFormula {
	return NopContraFormula{}
}

func (NopContraFormula) Validate() error {
	return nil
}
