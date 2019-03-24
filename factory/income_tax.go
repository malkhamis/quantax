package factory

import (
	"errors"

	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/benefits"
	"github.com/malkhamis/quantax/calc/tax"
	"github.com/malkhamis/quantax/history"
)

// Config is used to set up the various calculators this factory creates
type CalculatorConfig struct {
	Year   uint
	Region history.Jurisdiction
}

// NewIncomeTaxCalculator returns a tax calculator for the given parameters
func NewIncomeTaxCalculator(finances calc.IndividualFinances, cfg CalculatorConfig) (calc.TaxCalculator, error) {

	formula, err := history.GetTaxFormula(cfg.Year, cfg.Region)
	if err != nil {
		return nil, err
	}

	return tax.NewCalculator(finances, formula)
}

// ChildBenefitParams is used to pass child benefit calculation parameters to
// the relevant factory functions
type CBCalculatorConfig struct {
	Children []calc.Person
	CalculatorConfig
}

// NewChildBenefitCalculator returns a tax calculator for the given parameters
func NewChildBenefitCalculator(finances calc.FamilyFinances, cfg CBCalculatorConfig) (calc.ChildBenefitCalculator, error) {

	if len(cfg.Children) < 1 {
		return nil, errors.New("a minimum of one child is required")
	}

	formula, err := history.GetChildBenefitFormula(cfg.Year, cfg.Region)
	if err != nil {
		return nil, err
	}

	calcCfg := benefits.ConfigCCB{
		Finances: finances,
		Formula:  formula,
	}

	return benefits.NewCBCalculator(calcCfg, cfg.Children[0], cfg.Children[1:]...)
}
