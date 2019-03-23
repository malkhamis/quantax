package calc

type Person struct {
	Name      string
	AgeMonths uint
}

// Person represents a dependent child for tax purposes
// type Person struct {
// 	name      string
// 	birthdate date.Date
// }
//
// // NewPerson returns a new child instance. Month and day values are normalized
// // if the exceed their limits
// func NewPerson(name string, birthdate Date) Person {
//
// 	dob := date.New(birthdate.Year, birthdate.Month, birthdate.Day)
// 	return Person{name, dob}
// }
//
// // IsOlderThan returns true if the child is older than nMonths on the given date
// func (p Person) IsOlderThan(nMonths uint, on Date) (bool, error) {
//
// 	onDate := date.New(on.Year, on.Month, on.Day)
// 	if p.birthdate.After(onDate) {
// 		return false, errors.Wrap(ErrInvalidDate, "date is before child birthdate")
// 	}
//
// 	age := int32(onDate.Sub(p.birthdate))
// 	testAge := int32(nMonths * 30)
//
// 	// TODO inaccurate?
// 	return age > testAge, nil
// }
