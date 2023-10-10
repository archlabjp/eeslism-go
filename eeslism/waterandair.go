//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

package eeslism

import "math"

// 水、空気の物性値の計算
// パソコンによる空気調和計算法より作成

// 空気の密度　 dblT[℃]、出力[kg/m3]
func FNarow(dblT float64) float64 {
	return 1.293 / (1.0 + dblT/273.15)
}

// 空気の比熱　 dblT[℃]、出力[J/kgK]
func FNac() float64 {
	return 1005.0
}

// 空気の熱伝導率　 dblT[℃]、出力[W/mK]
func FNalam(dblT float64) float64 {
	var dblTemp float64

	if dblT > -50.0 && dblT < 100.0 {
		dblTemp = 0.0241 + 0.000077*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 空気の粘性係数　 dblT[℃]、出力[Pa s]
func FNamew(dblT float64) float64 {
	return (0.0074237 / (dblT + 390.15) * math.Pow((dblT+273.15)/293.15, 1.5))
}

// 空気の動粘性係数　 dblT[℃]、出力[m2/s]
func FNanew(dblT float64) float64 {
	return FNamew(dblT) / FNarow(dblT)
}

// 空気の膨張率　 dblT[℃]、出力[1/K]
func FNabeta(dblT float64) float64 {
	return 1.0 / (dblT + 273.15)
}

// 水の密度　 dblT[℃]、出力[kg/m3]
func FNwrow(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 1000.5 - 0.068737*dblT - 0.0035781*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 1008.7 - 0.28735*dblT - 0.0021643*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 水の比熱　 dblT[℃]、出力[J/kgK]
func FNwc(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 4210.4 - 1.356*dblT + 0.014588*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 4306.8 - 2.7913*dblT + 0.018773*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 水の熱伝導率　 dblT[℃]、出力[W/mK]
func FNwlam(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 0.56871 + 0.0018421*dblT - 7.0427e-6*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 0.60791 + 0.0012032*dblT - 4.7025e-6*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 水の動粘性係数　 dblT[℃]、出力[m2/s]
func FNwnew(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 50.0 {
		dblTemp = math.Exp(-13.233 - 0.032516*dblT + 0.000068997*dblT*dblT +
			0.0000069513*dblT*dblT*dblT - 0.00000009386*dblT*dblT*dblT*dblT)
	} else if dblT < 100.0 {
		dblTemp = math.Exp(-13.618 - 0.015499*dblT - 0.000022461*dblT*dblT +
			0.00000036334*dblT*dblT*dblT)
	} else if dblT < 200.0 {
		dblTemp = math.Exp(-13.698 - 0.016782*dblT + 0.000034425*dblT*dblT)
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 水の膨張率　 dblT[℃]、出力[1/K]
func FNwbeta(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 50.0 {
		dblTemp = -0.060159 + 0.018725*dblT - 0.00045278*dblT*dblT +
			0.0000098148*dblT*dblT*dblT - 0.000000083333*dblT*dblT*dblT*dblT
	} else if dblT < 100.0 {
		dblTemp = -0.46048 + 0.03104*dblT - 0.000325*dblT*dblT +
			0.0000013889*dblT*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 0.33381 + 0.002847*dblT + 0.000016154*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// 水の粘性係数　 dblT[℃]、出力[Pa s]
func FNwmew(dblT float64) float64 {
	return FNwnew(dblT) / FNwrow(dblT)
}

// 空気の熱拡散率　 dblT[℃]、出力[m2/s]
func FNaa(dblT float64) float64 {
	return FNalam(dblT) / FNac() / FNarow(dblT)
}

// 水の熱拡散率　 dblT[℃]、出力[m2/s]
func FNwa(dblT float64) float64 {
	return FNwlam(dblT) / FNwc(dblT) / FNwrow(dblT)
}

// プラントル数の計算
func FNPr(strF byte, dblT float64) float64 {
	var dblTemp float64

	if strF == 'A' || strF == 'a' {
		dblTemp = FNanew(dblT) / FNaa(dblT)
	} else if strF == 'W' || strF == 'w' {
		dblTemp = FNwnew(dblT) / FNwa(dblT)
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

// グラスホフ数の計算
//   dblTs:表面温度[℃]
//   dblTa:主流温度[℃]
//   dblx :代表長さ[m]
func FNGr(strF byte, dblTs, dblTa, dblx float64) float64 {
	const dblg = 9.80665

	dblT := (dblTs + dblTa) / 2.

	// 温度差の計算
	dbldT := math.Max(math.Abs(dblTs-dblTa), 0.1)

	var dblBeta, n float64
	if strF == 'A' || strF == 'a' {
		dblBeta = FNabeta(dblT)
		n = FNanew(dblT)
	} else if strF == 'W' || strF == 'w' {
		dblBeta = FNwbeta(dblT)
		n = FNwnew(dblT)
	} else {
		dblBeta = -999.0
		n = -999.0
	}

	return dblg * dblBeta * dbldT * math.Pow(dblx, 3.0) / (n * n)
}

// 各種定数を入力するとヌセルト数が計算される。
//   Nu = C * (Pr Gr)^m
func FNCNu(dblC, dblm, dblPrGr float64) float64 {
	return dblC * math.Pow(dblPrGr, dblm)
}

// 管内の強制対流熱伝達率の計算（流体は水のみ）
// dbld:配管内径[m]
// dblL:管長[m]
// dblv:管内流速[m/s]
// dblT:流体と壁面の平均温度[℃]
func FNhinpipe(dbld, dblL, dblv, dblT float64) float64 {
	dblnew := FNwnew(dblT)
	dblRe := dblv * dbld / dblnew
	dblPr := FNPr('W', dblT)
	dbldL := dbld / dblL
	dblld := FNwlam(dblT) / dbld

	var dblTemp float64
	if dblRe < 2200. {
		dblTemp = (3.66 + 0.0668*dbldL*dblRe*dblPr/(1.+0.04*math.Pow(dbldL*dblRe*dblPr, 2./3.))) * dblld
	} else {
		dblTemp = 0.023 * math.Pow(dblRe, 0.8) * math.Pow(dblPr, 0.4) * dblld
	}

	return dblTemp
}

// 円管外部の自然対流熱伝達率の計算（流体は水のみ）
// dbld:配管内径[m]
// dblT:流体と壁面の平均温度[℃]
func FNhoutpipe(dbld, dblTs, dblTa float64) float64 {
	dblC := 0.5
	dbln := 0.25

	dblPr := FNPr('W', (dblTs+dblTa)/2.)
	dblGr := FNGr('W', dblTs, dblTa, dbld)

	dblNu := dblC * math.Pow(dblPr*dblGr, dbln)

	return dblNu / dbld
}
