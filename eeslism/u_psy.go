/* ==================================================================

PSYLIB

湿り空気の状態値計算用ライブラリ－
（宇田川、パソコンによる空気調和計算法、プログラム3.1の C 言語版, ANSI C 版）

--------------------------------------------------------------------- */

package eeslism

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

// ---- 露点温度 ----

// 湿り空気を冷却していくと、やがて飽和空気となる。このときの温度を露点温度という。
// 飽和空気にはこれ以上水分を含むことができず、これ以上冷却すると結露が生じる。
// 飽和空気の全圧のうち、水蒸気が占める圧力(水蒸気分圧)である飽和水蒸気圧 Pws [kPa] を求める。
// 計算にはウェククスラー・ハイランド(Wexler-Hyland)による式を用いる。
// See: パソコンによる空気調和計算法 P.27
func FNPws(T float64) float64 {

	// 絶対温度 Tabs
	Tabs := T + 273.15
	if math.Abs(Tabs) < 1e-5 {
		fmt.Printf("xxxx ゼロ割が発生しています Tabs=%f\n", Tabs)
	}

	// 飽和水蒸気分圧 Pws
	var Pws float64
	if T > 0.0 {
		// 0から200℃の水と接する場合
		Pws = math.Exp(6.5459673*math.Log(Tabs)-5800.2206/Tabs+1.3914993+Tabs*(-0.048640239+
			Tabs*(4.1764768e-5-1.4452093e-8*Tabs))) / 1000.0
	} else {
		// -100から0℃の氷と接する場合
		Pws = math.Exp(-5674.5359/Tabs+6.3925247+Tabs*(-9.677843e-3+
			Tabs*(6.2215701e-7+Tabs*(2.0747825e-9-9.484024e-13*Tabs)))+4.1635019*math.Log(Tabs)) / 1000.0
	}

	// 単位変換
	return _Pcnv * Pws
}

// 水蒸気分圧 Pw [kPa] から露点温度を求める
// NOTE:
// - 611.2 Paは、水の飽和蒸気圧に関連する値であり、これはおおよそ0℃の飽和蒸気圧に相当します。
//
// See: パソコンによる空気調和計算法 P.28
func FNDp(Pw float64) float64 {
	// 水蒸気分圧の単位をPaに変換
	Pwx := Pw * 1000.0 / _Pcnv

	Y := math.Log(Pwx)

	// 近似式を用いて水蒸気分圧から露点温度を求める
	if Pwx >= 611.2 {
		// 0から50℃のとき
		// NOTE: 611.2[Pa]はおおよそ0℃の飽和蒸気圧である.
		return -77.199 + Y*(13.198+Y*(-0.63772+0.071098*Y))
	} else {
		// -50から0℃のとき
		return -60.662 + Y*(7.4624+Y*(0.20594+0.016321*Y))
	}
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

// --- 相対湿度、比較湿度 ---

// 湿り空気の水蒸気分圧 Pw [kPa] から絶対湿度 x [kg/kg]を求める。
// 絶対湿度とは、湿り空気の水蒸気と乾き空気の質量の比である。
// See: パソコンによる空気調和計算法 P.29
func FNXp(Pw float64) float64 {
	// 標準大気圧 P [kPa]
	P := 101.325

	if math.Abs(P-Pw) < 1.0e-4 {
		fmt.Printf("xxxxx ゼロ割が発生しています P=%f Pw=%f\n", P, Pw)
	}

	// 絶対湿度 x [kg/kg]
	x := 0.62198 * Pw / (P - Pw)

	return x
}

// 温度 T [C] および 相対湿度 Rh [%] から絶対湿度 x [kg/kg] を求める
func FNXtr(T, Rh float64) float64 {
	return FNXp(FNPwtr(T, Rh))
}

func FNXtw(T, Twb float64) float64 {
	Hc := FNHc(Twb)
	return ((_R0+_Cv*Twb-Hc)*FNXp(FNPws(Twb)) - _Ca*(T-Twb)) / (_Cv*T + _R0 - Hc)
}

// 絶対湿度 x [kg/kg] から 水蒸気分圧 Pw [kPa] を求める。
// See: パソコンによる空気調和計算法 P.29
func FNPwx(X float64) float64 {
	// 標準大気圧 P [kPa]
	P := 101.325

	// 水蒸気分圧 Pw [kPa]
	Pw := (X * P / (X + 0.62198))

	return Pw
}

// 温度 T [C] および 湿り空気の水蒸気分圧 Pw [kPa] から 相対湿度 φ [%] を求める。
// 相対湿度は水蒸気分圧を飽和水蒸気圧の百分率である。
// See: パソコンによる空気調和計算法 P.29
func FNRhtp(T, Pw float64) float64 {
	return 100.0 * Pw / FNPws(T)
}

// 温度 T [C] および 相対湿度 Rh [%] から 湿り空気の水蒸気分圧 Pw [kPa] を求める。
// See: パソコンによる空気調和計算法 P.29
func FNPwtr(T, Rh float64) float64 {
	return (Rh * FNPws(T) / 100.0)
}

// 相対湿度φ [%]と水蒸気分圧Pwから 乾球温度 Tを求める。
func FNDbrp(Rh, Pw float64) float64 {
	return FNDp(100.0 / Rh * Pw)
}

// 温度 T [C] および 湿り空気の水蒸気分圧 Pw [kPa] から 相対湿度 φ [%] を求める。
func FNRhtx(T, X float64) float64 {
	return FNRhtp(T, FNPwx(X))
}

// 乾燥空気の温度 t [C] と絶対湿度 x [kg/kg] から、湿り空気のエンタルピ h [J/kg] を求める。
// ここで、Tは乾燥空気の温度、Xは水蒸気の質量分率（乾燥空気に対する水蒸気の質量比）である。
// エンタルピーは、乾燥空気の比熱を用いた温度の項と、水蒸気の比熱及び蒸発熱を用いた温度と絶対湿度の項の和として計算される。
// See: パソコンによる空気調和計算法 P.29
func FNH(T, X float64) float64 {
	// 定数
	Ca := 1005.0    // 乾き空気の定格比熱, 1005 J/kgK (0.240 kcal/kgC)
	Cv := 1846.0    // 水蒸気の低圧比熱, 1846 J/kgK (0.441 kcal/kgC)
	r0 := 2501000.0 // 0℃の水の蒸発潜熱, 2501x10^3 J/kg (597.5 kcak/kg)

	// エンタルピー h
	h := Ca*T + (Cv*T+r0)*X

	return h
}

// 乾燥空気の温度 t [C]とエンタルピー h [J/kg]から
func FNXth(T, h float64) float64 {
	// 定数
	Ca := 1005.0    // 乾き空気の定格比熱, 1005 J/kgK (0.240 kcal/kgC)
	Cv := 1846.0    // 水蒸気の低圧比熱, 1846 J/kgK (0.441 kcal/kgC)
	r0 := 2501000.0 // 0℃の水の蒸発潜熱, 2501x10^3 J/kg (597.5 kcak/kg)

	return (h - Ca*T) / (Cv*T + r0)
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
