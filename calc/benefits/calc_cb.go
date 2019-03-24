package benefits

import (
	"github.com/malkhamis/quantax/calc"
	"github.com/pkg/errors"
)

var _ calc.ChildBenefitCalculator = (*CCBCalculator)(nil)

type FamilyFinances = calc.FamilyFinances

type CCBCalculator struct {
	children []Child
	formula  calc.ChildBenefitFormula
	calc.FamilyFinances
}

func NewCCBCalculator() (*CCBCalculator, error) {
	return nil, errors.New("not implemented")
}

// Calc returns the recievable amount of child benefits
func (c *CCBCalculator) Calc() float64 {
	// Not implemented
	return 0
}

// UpdateFinances sets the financial numbers which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set financial
// numbers to anything other than what the calculator was initialized with
func (c *CCBCalculator) UpdateFinances(newFinNums calc.FamilyFinances) {
	// Not implemented
}

// UpdateBeneficiary sets the child which the calculator will use in
// subsequent calls to Calc(). Users may call this method to set beneficiary
// to anything other than what the calculator was initialized with
func (c *CCBCalculator) UpdateBeneficiaries(newChildren []calc.Person) {
	// Not implemented
}

// UpdateForumla sets the formula for calculating the amount of benefits for
// children given family finances. Users may call this method to set the
// formula to anything other than what the calculator was initialized with
func (c *CCBCalculator) UpdateForumla(newFormula calc.ChildBenefitFormula) error {

	if newFormula == nil {
		return calc.ErrNoFormula
	}

	err := newFormula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula")
	}

	c.formula = newFormula
	return nil
}
