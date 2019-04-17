package tax

import "github.com/malkhamis/quantax/calc/finance"

var (
	// compile-time check for interface implementation
	_ ContraFormula = (NopContraFormula)(NopContraFormula{})
)

type NopContraFormula struct{}

func (NopContraFormula) Apply(finances *finance.IndividualFinances, netIncome float64) []TaxCredit {
	return nil
}

func (NopContraFormula) Clone() ContraFormula {
	return NopContraFormula{}
}

func (NopContraFormula) Validate() error {
	return nil
}
