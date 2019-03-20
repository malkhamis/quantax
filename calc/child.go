package calc

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rickb777/date"
)

// Introduce DeductionsType iota?
//
// In root pkg:
//
// define types:
//
//  type IncomeAdjustmentMethod int
//  const(
//    NetAdjFamily = iota
// )
//
//  FinancesFamily [2]Finances
//
//   (calc) AdjustedIncome(adj AdjustmentType)
//
//  define interfaces:
//    type ChildBenefitCalculator {
//      Calc()
//      UpdateFinances(FinancesFamily)
//      UpdateBeneficiary(Child)
//      UpdateRule(Rule)
//    }
//
//    type ChildBenefitRule iface {
//
//   }
//
//

// Child represents a dependent child for tax purposes
type Child struct {
	name      string
	birthdate date.Date
}

// DOB represents a calendar date without the context of locales
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewChild returns a new child instance. Month and day values are normalized
// if the exceed their limits
func NewChild(name string, birthdate Date) Child {

	dob := date.New(birthdate.Year, birthdate.Month, birthdate.Day)
	return Child{name, dob}
}

// IsOlderThan returns true if the child is older than nMonths on the given date
func (c Child) IsOlderThan(nMonths uint, on Date) (bool, error) {

	onDate := date.New(on.Year, on.Month, on.Day)
	if c.birthdate.After(onDate) {
		return false, errors.Wrap(ErrInvalidDate, "date is before child birthdate")
	}

	age := int32(onDate.Sub(c.birthdate))
	testAge := int32(nMonths * 30)

	// TODO inaccurate?
	return age > testAge, nil
}
