/* ================================================================

 SUNLIB

  太陽位置および日射量計算用ライブラリ－
  （宇田川、パソコンによる空気調和計算法、プログラム4.1の C 言語版, ANSI C 版）

---------------------------------------------------------------- */

package main

import (
	"math"
)

func Sunint() {
	var Rd float64 = math.Pi / 180.0
	Slat = math.Sin(Lat * Rd)
	Clat = math.Cos(Lat * Rd)
	Tlat = math.Tan(Lat * Rd)
	if UNIT == "SI" {
		Isc = 1370.0
	} else {
		Isc = 1178.0
	}
}

func FNDecl(N int) float64 {
	return math.Asin(0.397949 * math.Sin(2.0*math.Pi*(float64(N)-81.0)/365.0))
}

func FNE(N int) float64 {
	var B float64 = 2.0 * math.Pi * (float64(N) - 81.0) / 365.0
	return 0.1645*math.Sin(2.0*B) - 0.1255*math.Cos(B) - 0.025*math.Sin(B)
}

func FNSro(N int) float64 {
	return Isc * (1.0 + 0.033*math.Cos(2.0*math.Pi*float64(N)/365.0))
}

func FNTtas(Tt float64, E float64) float64 {
	return Tt + E + (Lon-Ls)/15.0
}

func FNTt(Ttas float64, E float64) float64 {
	return Ttas - E - (Lon-Ls)/15.0
}

func FNTtd(Decl float64) float64 {
	var Tlat float64
	var Cws, Ttd float64
	Cws = -Tlat * math.Tan(Decl)
	if 1.0 > Cws && Cws > -1.0 {
		Ttd = 7.6394 * math.Acos(Cws)
	} else {
		if Cws >= 1.0 {
			Ttd = 0.0
		} else {
			Ttd = 24.0
		}
	}
	return Ttd
}

var __Solpos_Sdecl, __Solpos_Sld, __Solpos_Cld float64
var __Solpos_Ttprev float64 = 25.0

func Solpos(Ttas float64, Decl float64, Sh *float64, Sw *float64, Ss *float64, solh *float64, solA *float64) {
	const PI float64 = math.Pi
	// const Slat float64 = 0.0
	// const Clat float64 = 1.0
	var Ch, Ca, Sa, W float64

	if Ttas < __Solpos_Ttprev {
		__Solpos_Sdecl = math.Sin(Decl)
		__Solpos_Sld = Slat * __Solpos_Sdecl
		__Solpos_Cld = Clat * math.Cos(Decl)
	}

	W = (Ttas - 12.0) * 0.2618
	*Sh = __Solpos_Sld + __Solpos_Cld*math.Cos(W)
	*solh = math.Asin(*Sh) / PI * 180.0

	if *Sh > 0.0 {
		Ch = math.Sqrt(1.0 - *Sh**Sh)
		Ca = (*Sh*Slat - __Solpos_Sdecl) / (Ch * Clat)
		var fW0 float64
		if W > 0.0 {
			fW0 = 1.0
		} else {
			fW0 = 0.0
		}
		*solA = fW0*1.0 + (1.0-fW0)*(-1.0)*math.Acos(Ca)/PI*180.0
		Sa = (W / math.Abs(W)) * math.Sqrt(1.0-Ca*Ca)
		*Sw = Ch * Sa
		*Ss = Ch * Ca
	} else {
		*Sh = 0.0
		*Sw = 0.0
		*Ss = 0.0
		*solh = 0.0
		*solA = 0.0
	}

	__Solpos_Ttprev = Ttas
}

func Srdclr(Io float64, P float64, Sh float64, Idn *float64, Isky *float64) {
	if Sh > 0.001 {
		*Idn = Io * math.Pow(P, 1.0/Sh)
		*Isky = Sh * (Io - *Idn) * (0.66 - 0.32*Sh) * (0.5 + (0.4-0.3*P)*Sh)
	} else {
		*Idn = 0.0
		*Isky = 0.0
	}
}

func Dnsky(Io float64, Ihol float64, Sh float64, Idn *float64, Isky *float64) {
	if Sh > 0.001 {
		Kt := Ihol / (Io * Sh)
		if Kt >= 0.5163+(0.333+0.00803*Sh)*Sh {
			*Idn = (-0.43 + 1.43*Kt) * Io
		} else {
			*Idn = (2.277 + (-1.258+0.2396*Sh)*Sh) * (Kt * Kt * Kt) * Io
		}
		*Isky = Ihol - *Idn*Sh
	} else {
		*Idn = 0.0
		*Isky = Ihol
	}
}
