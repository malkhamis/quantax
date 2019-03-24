package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

type WBFormula = calc.WeightedBracketFormula

// StepReducer represents a lookup table for formulas such that each formula
// is mapped to an index that could represent something meaningful. For example
// this number may indicate the number of children or a range of number of tax
// days spent in a specific region
type StepReducer struct {
	StepFormulas        []calc.WeightedBracketFormula
	AboveMaxStepFormula calc.WeightedBracketFormula
}

// NewStepReducerFormula returns a step reducer instance. The above-max step
// indicates that this formula is used for steps above the maximum available
func NewStepReducer(aboveMaxStep, firstStep WBFormula, extraSteps ...WBFormula) (*StepReducer, error) {

	stepReducer := &StepReducer{
		AboveMaxStepFormula: aboveMaxStep.Clone(),
		StepFormulas:        []WBFormula{firstStep.Clone()},
	}

	for _, extra := range extraSteps {
		stepReducer.StepFormulas = append(stepReducer.StepFormulas, extra.Clone())
	}

	return stepReducer, stepReducer.Validate()
}

// Reduce returns the reduced value from the given amount using the given step.
// The first step is indexed at zero. If the given step exceeds the maximum
// available step, the formula of the above-max step is used for the reduction
func (sr *StepReducer) Reduce(amount float64, step uint) (reduction float64) {

	indexMaxAvailableStep := len(sr.StepFormulas) - 1
	if step > uint(indexMaxAvailableStep) {
		return sr.AboveMaxStepFormula.Apply(amount)
	}

	return sr.StepFormulas[step].Apply(amount)
}

// Validate ensures this step reducer is valid for use
func (sr *StepReducer) Validate() error {

	if sr.AboveMaxStepFormula == nil {
		return errors.Wrap(calc.ErrNoFormula, "above-max-step formula")
	}

	err := sr.AboveMaxStepFormula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula for the above-max step")
	}

	if len(sr.StepFormulas) < 1 {
		return errors.Wrap(calc.ErrNoFormula, "empty step formulas")
	}

	for i, formula := range sr.StepFormulas {

		if formula == nil {
			return errors.Wrapf(calc.ErrNoFormula, "step formulas %d", i)
		}

		if err = formula.Validate(); err != nil {
			return errors.Wrapf(err, "step %d: invalid formula", i)
		}
	}

	return nil
}

// Clone returns a copy of this step reducer
func (sr *StepReducer) Clone() *StepReducer {

	clone := &StepReducer{
		AboveMaxStepFormula: sr.AboveMaxStepFormula.Clone(),
	}

	if sr.StepFormulas != nil {
		clone.StepFormulas = make([]WBFormula, len(sr.StepFormulas))
	}
	for i, formula := range sr.StepFormulas {
		clone.StepFormulas[i] = formula.Clone()
	}

	return clone
}
