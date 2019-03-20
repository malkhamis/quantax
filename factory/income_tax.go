package factory

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/malkhamis/quantax/calc/incometax"
	"github.com/malkhamis/quantax/history"
)

// IncomeTaxParams is used to pass tax parameters to relevant factory functions
type IncomeTaxParams struct {
	Year   uint
	Region history.Jurisdiction
}

// NewIncomeTaxCalculator returns a tax calculator for the given parameters
func NewIncomeTaxCalculator(finNums calc.Finances, params IncomeTaxParams) (calc.TaxCalculator, error) {

	rates, err := history.Get(params.Year, params.Region)
	if err != nil {
		return nil, err
	}

	return incometax.NewCalculator(finNums, rates)
}

// incomeTaxCalculatorAgg wraps multiple calc.Calculator's and presents an
// implementation of calc.TaxCalculator as a single calculator. Results of
// the underlying calculators are aggregated
type incomeTaxCalculatorAgg struct {
	calculators []calc.TaxCalculator
}

// Calc returns the sum of results of calling Calc() on underlying calculators
func (agg *incomeTaxCalculatorAgg) Calc() float64 {
	var sum float64
	for _, c := range agg.calculators {
		sum += c.Calc()
	}
	return sum
}

// Update sets the given financial numbers in all underlying calculators
func (agg *incomeTaxCalculatorAgg) Update(finNums calc.Finances) {
	for _, c := range agg.calculators {
		c.Update(finNums)
	}
}

// NewIncomeTaxCalculatorAgg returns a wrapper arounf multiple tax calculators.
// Calls to Calc return the sum of the results from all underlying calculators.
// This also applies to Update(), where all underlying calculators are updated
// to the same given financial numbers. The returned calculator is useful for
// the cases in which income tax is calculated for two or more jurisdications
func NewIncomeTaxCalculatorAgg(finNums calc.Finances, params1, params2 IncomeTaxParams, extras ...IncomeTaxParams) (calc.TaxCalculator, error) {

	first, err := NewIncomeTaxCalculator(finNums, params1)
	if err != nil {
		return nil, err
	}

	second, err := NewIncomeTaxCalculator(finNums, params2)
	if err != nil {
		return nil, err
	}

	agg := &incomeTaxCalculatorAgg{
		calculators: []calc.TaxCalculator{first, second},
	}

	for _, params := range extras {
		c, err := NewIncomeTaxCalculator(finNums, params)
		if err != nil {
			return nil, err
		}
		agg.calculators = append(agg.calculators, c)
	}

	return agg, nil
}
