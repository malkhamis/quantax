package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type WBFormula = calc.WeightedBracketFormula

// StepReducer represents a lookup table for formulas such that each formula
// is mapped to a number that could represent something meaningful. For example
// this number may indicate the number of children or the number of tax days
// spent in a specific region
type StepReducer map[int]calc.WeightedBracketFormula

// NewStepReducerFormula returns a step reducer instance. 'zeroStep' and
// 'infStep' are for the zero and the infinity step, respectively, which
// must exist in the map to indicate a step that means greater than max
func NewStepReducer(zeroStep WBFormula, infyStep WBFormula, others ...WBFormula) (StepReducer, error) {

	err := zeroStep.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula for the zero step")
	}

	err = infyStep.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid formula for the infinity step")
	}

	stepReducer := StepReducer{
		0:  zeroStep,
		-1: infyStep,
	}

	for i, next := range others {
		err = next.Validate()
		if err != nil {
			return nil, errors.Wrap(err, "invalid formula") // TODO: formula name
		}
		stepReducer[i+1] = next
	}

	return stepReducer, nil
}

// Reduce returns the reduced value from the given amount using the given step
func (sr StepReducer) Reduce(amount float64, step int) (newAmount float64) {
	// TODO not implemented
	return 0
}

// Validate ensures this step reducer is valid for use
func (sr StepReducer) Validate() error {

	var (
		existInf  bool
		existZero bool
	)

	for step, formula := range sr {

		if step < 0 {
			existInf = true
		} else if step == 0 {
			existZero = true
		}

		if err := formula.Validate(); err != nil {
			return errors.Wrapf(err, "step [%d]: invalid sub-formula", step)
		}

	}

	if !existInf {
		return errors.New(
			"a subformula mapped to a negative number " +
				"to indicate more than [max] step must exist",
		)
	}

	if !existZero {
		return errors.New("a subformula mapped to zero step must exist")
	}

	return nil
}

// Clone returns a copy of this step reducer
func (sr StepReducer) Clone() StepReducer {

	newSR := make(StepReducer)
	for step, wbf := range sr {
		newSR[step] = wbf.Clone()
	}
	return newSR
}
