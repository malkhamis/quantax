package history

import "github.com/malkhamis/quantax/facts"

// Province represents a name of a province/state
type Province string

// convenience constants for use with this package
const (
	BC  Province = "British Columbia"
	YT  Province = "Yukon Territory"
	AB  Province = "Alberta"
	NWT Province = "Northwest Territories"
	SK  Province = "Saskatchewan"
	MB  Province = "Manitoba"
	NU  Province = "Nunavut"
	ON  Province = "Ontario"
	QC  Province = "Quebec"
	NL  Province = "Newfoundland and Labrador"
	NB  Province = "New Brunswick"
	NS  Province = "Nova Scotia"
	PE  Province = "Prince Edward Island"
)

var factsAllProvinces = map[Province]map[uint]facts.FactsProv{
	BC: factsBC,
}
