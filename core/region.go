package core

// Region represent a tax jurisdiction
type Region string

// Defined tax regions which implementations should support
const (
	RegionCA Region = "Canada"
	RegionNT Region = "Northwest Territories"
	RegionBC Region = "British Columbia"
	RegionYT Region = "Yukon Territory"
	RegionAB Region = "Alberta"
	RegionSK Region = "Saskatchewan"
	RegionMB Region = "Manitoba"
	RegionNU Region = "Nunavut"
	RegionON Region = "Ontario"
	RegionQC Region = "Quebec"
	RegionNL Region = "Newfoundland and Labrador"
	RegionNB Region = "New Brunswick"
	RegionNS Region = "Nova Scotia"
	RegionPE Region = "Prince Edward Island"
)
