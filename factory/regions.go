package factory

import "github.com/malkhamis/quantax/history"

type Region int

const (
	Canada Region = iota
	NWT
	BC
	YT
	AB
	SK
	MB
	NU
	ON
	QC
	NL
	NB
	NS
	PE
)

var knownRegions = map[Region]history.Jurisdiction{
	Canada: history.Canada,
	NWT:    history.NWT,
	BC:     history.BC,
	YT:     history.YT,
	AB:     history.AB,
	SK:     history.SK,
	MB:     history.MB,
	NU:     history.NU,
	ON:     history.ON,
	QC:     history.QC,
	NL:     history.NL,
	NB:     history.NB,
	NS:     history.NS,
	PE:     history.PE,
}
