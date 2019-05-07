// Package benefits implements benefit calculators' interfaces defined by
// package calc
package benefits

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"

	"github.com/pkg/errors"
)

// Sentinel errors that can be wrapped and returned by this package
var (
	ErrNoFormula = errors.New("no formula given/set")
	ErrNoIncCalc = errors.New("no income calculator given")
	ErrNoCalc    = errors.New("no benefit calculator given")
)

// ChildBenefitFormula represents a method for calculating child benefits
type ChildBenefitFormula interface {
	// Apply returns the sum of benefits for all beneficiaries
	Apply(netIncome float64, children ...*human.Person) float64
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() ChildBenefitFormula
}

// CalcConfigCB is used to pass configurations to create new child benefit
// calculator
type CalcConfigCB struct {
	Formula    ChildBenefitFormula
	IncomeCalc core.IncomeCalculator
}

// validate checks if the configurations are valid for use by calc constructors
func (cfg CalcConfigCB) validate() error {

	if cfg.Formula == nil {
		return ErrNoFormula
	}

	err := cfg.Formula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula")
	}

	if cfg.IncomeCalc == nil {
		return ErrNoIncCalc
	}

	return nil
}

func getChildCount(children []*human.Person) int {
	childCount := len(children)
	for _, c := range children {
		if c == nil {
			childCount--
		}
	}
	return childCount
}
