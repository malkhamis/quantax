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
	finances    core.HouseholdFinances
	children    []*human.Person
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

// BenefitRecievable returns the aggregate recievable amount of child benefits
func (agg *ChildBenfitAggregator) BenefitRecievable() float64 {

	var total float64
	for _, c := range agg.calculators {
		agg.setupChildBenefitCalculator(c)
		total += c.BenefitRecievable()
	}
	return total
}

// SetBeneficiaries sets the children which the calculator will compute the
// benefits for in subsequent calls to BenefitRecievable()
func (agg *ChildBenfitAggregator) SetBeneficiaries(children []*human.Person) {
	agg.children = children
}

// SetFinances stores the given financial data in this calculator. Subsequent
// calls to other calculator functions will be based on the the given finances.
// Changes to the given finances after calling this function will affect future
// calculations. If finances is nil, a non-nil, empty finances is set
func (agg *ChildBenfitAggregator) SetFinances(finances core.HouseholdFinances) {
	agg.finances = finances
}

// setupChildBenefitCalculator sets up the given calculator with the finances as
// well as dependents stored in this aggregator
func (agg *ChildBenfitAggregator) setupChildBenefitCalculator(c core.ChildBenefitCalculator) {
	c.SetBeneficiaries(agg.children)
	c.SetFinances(agg.finances)
}
