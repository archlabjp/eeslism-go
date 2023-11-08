package eeslism

import (
	"fmt"
)

var mapCx = map[FliudType]float64{
	WATER_FLD: Cw, // 水の比熱
	AIRa_FLD:  Ca, // 空気の比熱
}

// 水、空気の比熱
func Spcheat(fluid FliudType) float64 {
	C, ok := mapCx[fluid]
	if !ok {
		s := fmt.Sprintf("xxx fluid='%c'", fluid)
		Eprint("<spcheat>", s)
		return -9999.0
	}
	return C
}
