package eeslism

import (
	"fmt"
)

/*  水、空気の比熱 */
/* ---------------------- */

func Spcheat(fluid rune) float64 {
	if fluid == 'W' {
		return Cw
	} else if fluid == 'a' {
		return Ca
	} else {
		s := fmt.Sprintf("xxx fluid='%c'", fluid)
		Eprint("<spcheat>", s)
		return -9999.0
	}
}
