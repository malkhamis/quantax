package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitFormula = (*CCBFormula)(nil)

type Child = calc.Person
type ChildCount = int

// CCBFormula computes Canada Child Benefits amounts for child(ren)
type CCBFormula struct {
	// the [min, max] dollar amounts for given age groups
	Benefits []AgeGroupBenefits
	// the sub-formulas to reduce the maximum benefits. Negative child
	// count indicate any number of children above the maximum positve
	Reducers StepReducer
}

func (cbf *CCBFormula) Apply(income float64, first Child, others ...Child) float64 {
	// TODO not implemented
	return 0
}

// TODO reduction from net total vs reduction from each child benefits
// func (CCBFormula) Calc(finances calc.FamilyFinances, children []Child) float64 {
//
// 	var total float64
//
// 	for _, c := range children {
// 		for _, ba := range mr.MaxAmnts {
// 			if c.IsOlderThan(ba.Ages.Min(), mr.AgeAt) && !c.IsOlderThan(ba.Ages.Max(), mr.AgeAt) {
//         total += // reduce per reduction rule
// 			}
// 		}
// 	}
// }

func (cbf *CCBFormula) Validate() error {

	for _, ageGroupBenefit := range cbf.Benefits {
		if err := ageGroupBenefit.Validate(); err != nil {
			return errors.Wrap(err, "invalid age group benefits")
		}
	}

	if err := cbf.Reducers.Validate(); err != nil {
		return errors.Wrap(err, "invalid step reducers")
	}

	return nil
}
