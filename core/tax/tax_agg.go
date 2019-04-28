package tax

//
// import (
// 	"github.com/malkhamis/quantax/core"
// 	"github.com/malkhamis/quantax/core/human"
// 	"github.com/pkg/errors"
// )
//
// // compile-time check for interface implementation
// var _ core.TaxCalculator = (*Aggregator)(nil)
//
// // Aggregator is used to aggregate payable tax from multiple tax calculators
// type Aggregator struct {
// 	calculators []core.TaxCalculator
// }
//
// // NewAggregator returns a new tax aggregator for the given tax calculators
// func NewAggregator(c0, c1 core.TaxCalculator, extras ...core.TaxCalculator) (*Aggregator, error) {
//
// 	cAgg := &Aggregator{
// 		calculators: make([]core.TaxCalculator, 0, len(extras)+2),
// 	}
//
// 	for i, c := range append([]core.TaxCalculator{c0, c1}, extras...) {
// 		if c == nil {
// 			return nil, errors.Wrapf(ErrNoCalc, "index %d: invalid calculator", i)
// 		}
// 		cAgg.calculators = append(cAgg.calculators, c)
// 	}
//
// 	return cAgg, nil
//
// }
//
// // TaxPayable returns the sum of payable tax from the underlying calculators
// func (agg *Aggregator) TaxPayable() (spouseA, spouseB float64, unusedCredits []core.TaxCredit) {
//
// 	var (
// 		taxAggA float64
// 		taxAggB float64
// 		crAgg   []core.TaxCredit
// 	)
//
// 	for _, c := range agg.calculators {
// 		taxA, taxB, credits := c.TaxPayable()
// 		taxAggA += taxA
// 		taxAggB += taxB
// 		crAgg = append(crAgg, credits...)
// 	}
//
// 	return taxAggA, taxAggB, crAgg
// }
//
// // SetFinances sets the given finances in all underlying tax calculators
// func (agg *Aggregator) SetFinances(f core.HouseholdFinances) {
// 	for _, c := range agg.calculators {
// 		c.SetFinances(f)
// 	}
// }
//
// // SetCredits sets the given credits in all underlying tax calculators
// func (agg *Aggregator) SetCredits(credits []core.TaxCredit) {
// 	for _, c := range agg.calculators {
// 		c.SetCredits(credits)
// 	}
// }
//
// // SetDependents sets the dependents which the calculator might use for tax-
// // related calculations
// func (agg *Aggregator) SetDependents(dependents ...*human.Person) {
// 	for _, c := range agg.calculators {
// 		c.SetDependents(dependents)
// 	}
// }
