package history

import (
	"math"
)

// var (
// 	_ = childBenefitFormulasBC
// 	_ = cbFormulaCA2018
// )

var taxFormulasBC = yearlyTaxFormulas{
	2018: WeightedBracketFormula{
		-0.0506: Bracket{0, 10412},
		0.0506:  Bracket{0, 39676},
		0.0770:  Bracket{39677, 79353},
		0.1050:  Bracket{79354, 91107},
		0.1229:  Bracket{91108, 110630},
		0.1470:  Bracket{110631, 150000},
		0.1680:  Bracket{150001, math.Inf(1)},
	},
}

// var childBenefitFormulasBC = map[uint]calc.ChildBenefitFormula{
// 	2018: cbFormulaCA2018,
// }

// TODO
//func cbFormulaCA2018(finances calc.FamilyFinances, children calc.Children) float64 {

// We calculate the Canada child benefit (CCB) as follows:
//
//   $6,496 per year ($541.33 per month) for each eligible child under the age of six
//
//   $5,481 per year ($456.75 per month) for each eligible child aged 6 to 17
//
//   We start to reduce the amount of CCB you get when your adjusted family net income (AFNI) is over $30,450. The reduction is calculated as follows:
//
//   families with one eligible child: the reduction is 7% of the amount of AFNI between $30,450 and $65,976, plus 3.2% of the amount of AFNI over $65,976
//
//   families with two eligible children: the reduction is 13.5% of the amount of AFNI between $30,450 and $65,976, plus 5.7% of the amount of AFNI over $65,976
//
//   families with three eligible children: the reduction is 19% of the amount of AFNI between $30,450 and $65,976, plus 8% of the amount of AFNI over $65,976
//
//   families with four or more eligible children: the reduction is 23% of the amount of AFNI between $30,450 and $65,976, plus 9.5% of the amount of AFNI over $65,976

// 	return 0.
// }
