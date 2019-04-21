package benefits

import (
	"github.com/malkhamis/quantax/core/finance"
	"github.com/malkhamis/quantax/core/human"
)

type testIncomeCalculator struct {
	onNetIncome       float64
	onTotalDeductions float64
	onTotalIncome     float64
}

func (tic testIncomeCalculator) TotalIncome(_ finance.IncomeDeductor) float64 {
	return tic.onTotalIncome
}
func (tic testIncomeCalculator) TotalDeductions(_ finance.IncomeDeductor) float64 {
	return tic.onTotalDeductions
}
func (tic testIncomeCalculator) NetIncome(_ finance.IncomeDeductor) float64 {
	return tic.onNetIncome
}

type testCBFormula struct {
	onApply    float64
	onValidate error
	onClone    ChildBenefitFormula
}

func (tcb testCBFormula) Apply(_ float64, _ ...human.Person) float64 {
	return tcb.onApply
}
func (tcb testCBFormula) Validate() error {
	return tcb.onValidate
}
func (tcb testCBFormula) Clone() ChildBenefitFormula {
	return tcb.onClone
}
