package calc

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rickb777/date"
)

// DOB represents a calendar date without the context of locales
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// Child represents a dependent child for tax purposes
type Child struct {
	name      string
	birthdate date.Date
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

// Children is a convenience type for handling one or more child
type Children []Child

// NewChildren returns Children type from a multiple Child instances
func NewChildren(first Child, others ...Child) Children {
	all := []Child{first}
	all = append(all, others...)
	return all
}

// Count returns the number of children
func (c Children) Count() int {
	return len(c)
}
