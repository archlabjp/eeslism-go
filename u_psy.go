/* ==================================================================

PSYLIB

湿り空気の状態値計算用ライブラリ－
（宇田川、パソコンによる空気調和計算法、プログラム3.1の C 言語版, ANSI C 版）

--------------------------------------------------------------------- */

package main

import (
	"fmt"
	"math"
)

var _R0, _Ca, _Cv, _Rc, _Cc, _Cw, _Pcnv, _P float64 = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0

func Psyint() {
	if UNIT == "SI" {
		_P = 101.325
		_R0 = 2501000.0
		_Ca = 1005.0
		_Cv = 1846.0
		_Rc = 333600.0
		_Cc = 2093.0
		_Cw = 4186.0
		_Pcnv = 1.0
	} else {
		_P = 760.0
		_R0 = 597.5
		_Ca = 0.24
		_Cv = 0.441
		_Rc = 79.7
		_Cc = 0.5
		_Cw = 1.0
		_Pcnv = 7.50062
	}
}

func Poset(Po float64) {
	_P = Po
}

func FNPo() float64 {
	return _P
}

func FNPws(T float64) float64 {
	var Tabs, Pws, Temp float64
	Tabs = T + 273.15

	if math.Abs(Tabs) < 1e-5 {
		fmt.Printf("xxxx ゼロ割が発生しています Tabs=%f\n", Tabs)
	}

	if T > 0 {
		Temp = 6.5459673*math.Log(Tabs) - 5800.2206/Tabs + 1.3914993 + Tabs*(-0.048640239+
			Tabs*(4.1764768e-5-1.4452093e-8*Tabs))
		Pws = math.Exp(Temp)
	} else {
		Pws = math.Exp(-5674.5359/Tabs+6.3925247+Tabs*(-9.677843e-3+
			Tabs*(6.2215701e-7+Tabs*(2.0747825e-9-9.484024e-13*Tabs)))) + 4.1635019*math.Log(Tabs)
	}

	//fmt.Printf("Tabs=%f Temp=%f Pws=%f\n", Tabs, Temp, Pws)
	return _Pcnv * Pws / 1000.0
}

func FNDp(Pw float64) float64 {
	var Pwx, Y float64
	Pwx = Pw * 1000.0 / _Pcnv
	Y = math.Log(Pwx)

	if Pwx >= 611.2 {
		return -77.199 + Y*(13.198+Y*(-0.63772+0.071098*Y))
	} else {
		return -60.662 + Y*(7.4624+Y*(0.20594+0.016321*Y))
	}
}

func FNDbrp(Rh, Pw float64) float64 {
	return FNDp(100.0 / Rh * Pw)
}

func FNDbxr(X, Rh float64) float64 {
	return FNDbrp(Rh, FNPwx(X))
}

func FNDbxh(X, H float64) float64 {
	return (H - _R0*X) / (_Ca + _Cv*X)
}
func FNDbxw(X, Twb float64) float64 {
	Hc := FNHc(Twb)
	return ((_Ca*Twb + (_Cv*Twb+_R0-Hc)*FNXp(FNPws(Twb)) - (_R0-Hc)*X) / (_Ca + _Cv*X))
}

func FNDbrh(Rh, H float64) float64 {
	var T0, F, Fd, Dbrh float64
	T0 = math.Min(FNDbxh(0., H), 30.)
	for I := 1; I <= 10; I++ {
		F = H - FNH(T0, FNXtr(T0, Rh))
		Fd = (H - FNH(T0+.1, FNXtr(T0+.1, Rh)) - F) / .1
		Dbrh = T0 - F/Fd
		if math.Abs(Dbrh-T0) <= .02 {
			return Dbrh
		}
		T0 = Dbrh
	}
	fmt.Printf("XXX FNDbrh  (T-T0)=%f\n", Dbrh-T0)
	return Dbrh
}

func FNDbrw(Rh, Twb float64) float64 {
	var T0, F, Fd, Dbrw float64
	T0 = Twb
	for I := 1; I <= 10; I++ {
		F = T0 - FNDbxw(FNXtr(T0, Rh), Twb)
		Fd = (T0 + .1 - FNDbxw(FNXtr(T0+.1, Rh), Twb) - F) / .1
		Dbrw = T0 - F/Fd
		if math.Abs(Dbrw-T0) <= .02 {
			return Dbrw
		}
		T0 = Dbrw
	}
	fmt.Printf("XXX FNDbrw  (T-T0)=%f\n", Dbrw-T0)
	return Dbrw
}

func FNXp(Pw float64) float64 {
	P := 101.325
	if math.Abs(P-Pw) < 1.0e-4 {
		fmt.Printf("xxxxx ゼロ割が発生しています P=%f Pw=%f\n", P, Pw)
	}
	return 0.62198 * Pw / (P - Pw)
}

func FNXtr(T, Rh float64) float64 {
	return FNXp(FNPwtr(T, Rh))
}

func FNXth(T, H float64) float64 {
	R0, Ca, Cv := 2501000.0, 1005.0, 1846.0
	return (H - Ca*T) / (Cv*T + R0)
}

func FNXtw(T, Twb float64) float64 {
	Hc := FNHc(Twb)
	return ((_R0+_Cv*Twb-Hc)*FNXp(FNPws(Twb)) - _Ca*(T-Twb)) / (_Cv*T + _R0 - Hc)
}

func FNPwx(X float64) float64 {
	P := 101.325
	return (X * P / (X + 0.62198))
}

func FNPwtr(T, Rh float64) float64 {
	return (Rh * FNPws(T) / 100.0)
}

func FNRhtp(T, Pw float64) float64 {
	return (100.0 * Pw / FNPws(T))
}

func FNRhtx(T, X float64) float64 {
	return (FNRhtp(T, FNPwx(X)))
}

func FNH(T, X float64) float64 {
	R0, Ca, Cv := 2501000.0, 1005.0, 1846.0
	return (Ca*T + (Cv*T+R0)*X)
}

func FNWbtx(T float64, X float64) float64 {
	var Tw0, H, Xs, Xss, F, Fd, Wbtx float64
	Tw0 = T
	H = FNH(T, X)
	for I := 1; I <= 10; I++ {
		Xs = FNXp(FNPws(Tw0))
		F = FNH(Tw0, Xs) - H - (Xs-X)*FNHc(Tw0)
		Xss = FNXp(FNPws(Tw0 + .1))
		Fd = (FNH(Tw0+.1, Xss) - H - (Xss-X)*FNHc(Tw0+.1) - F) / .1
		Wbtx = Tw0 - F/Fd
		if math.Abs(Wbtx-Tw0) <= .02 {
			return Wbtx
		}
		Tw0 = Wbtx
	}
	fmt.Printf("XXX FNWbtx  (Twb-Tw0)=%f\n", Wbtx-Tw0)
	return Wbtx
}

func FNHc(Twb float64) float64 {
	var Hc float64 // declare and initialize Hc variable
	if Twb >= 0.0 {
		Hc = _Cw * Twb
	} else {
		Hc = -_Rc + _Cc*Twb
	}
	return Hc
}
