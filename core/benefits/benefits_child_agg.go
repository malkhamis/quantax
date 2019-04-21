package benefits

import (
	"github.com/malkhamis/quantax/core"
	"github.com/malkhamis/quantax/core/human"
	"github.com/pkg/errors"
)

// ChildBenfitAggregator is used to calculate recievable child benefits. This
// type can aggregate the benefits from multiple child benefit calculators. it
// implements the following interface:
//  core.ChildBenefitCalculator
type ChildBenfitAggregator struct {
	calculators []core.ChildBenefitCalculator
}

// compile-time check for interface implementation
var _ core.ChildBenefitCalculator = (*ChildBenfitAggregator)(nil)

// NewChildBenfitAggregator returns a new child benefit calculator which can
// aggregate the benefits from all the given child benefit calculators. If one
// calculator is nil, it returns wrapped(ErrNoCalc)
func NewChildBenefitAggregator(c0, c1 core.ChildBenefitCalculator, extras ...core.ChildBenefitCalculator) (*ChildBenfitAggregator, error) {

	cAgg := &ChildBenfitAggregator{
		calculators: make([]core.ChildBenefitCalculator, 0, len(extras)+2),
	}

	for i, c := range append([]core.ChildBenefitCalculator{c0, c1}, extras...) {
		if c == nil {
			return nil, errors.Wrapf(ErrNoCalc, "index %d: invalid calculator", i)
		}
		cAgg.calculators = append(cAgg.calculators, c)
	}

	return cAgg, nil
}

// Calc returns the aggregate recievable amount of child benefits
func (agg *ChildBenfitAggregator) Calc(finances core.Financer) float64 {

	var total float64
	for _, c := range agg.calculators {
		total += c.Calc(finances)
	}
	return total
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to Calc()
func (agg *ChildBenfitAggregator) SetBeneficiaries(children ...human.Person) {
	for _, c := range agg.calculators {
		c.SetBeneficiaries(children...)
	}
}
