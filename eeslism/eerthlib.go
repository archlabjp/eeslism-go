package eeslism

import (
	"math"
)

/* 地中温度の計算 */

func Tearth(Z float64, n int, nmx int, Tgro float64, DTg float64, a float64) float64 {
	const t = 31.536e+6
	Cz := Z * math.Sqrt(math.Pi/(a*t))
	Tearth := Tgro + 0.5*DTg*math.Exp(-Cz)*math.Cos(float64(n-nmx)*0.017214-Cz)
	return Tearth
}

/* -------------------------------------------------- */

func Exsfter(day int, daymx int, Tgrav float64, DTgr float64, Exs []*EXSF, Wd *WDAT, tt int) {
	if Exs != nil {
		for i := range Exs {
			_Exs := Exs[i]
			if _Exs.Typ == EXSFType_E {
				_Exs.Tearth = Tearth(_Exs.Z, day, daymx, Tgrav, DTgr, _Exs.Erdff)
			} else if _Exs.Typ == EXSFType_e {
				_Exs.Tearth = Wd.EarthSurface[day*24+tt]
			}
		}
	}
}
