package eeslism

import (
	"fmt"
)

// 水、空気の比熱
func Spcheat(fluid FliudType) float64 {
	if fluid == WATER_FLD {
		// 水の比熱
		return Cw
	} else if fluid == AIRa_FLD {
		// 空気の比熱
		return Ca
	} else {
		s := fmt.Sprintf("xxx fluid='%c'", fluid)
		Eprint("<spcheat>", s)
		return -9999.0
	}
}
