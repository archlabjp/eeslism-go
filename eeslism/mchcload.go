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

/*  hcload.c  */

/*  空調負荷仮想機器  */

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
)

/* ------------------------------------------ */

/*  コイル出口空気温湿度に関する変数割当  */

func Hclelm(Hcload []HCLOAD) {
	for i := range Hcload {
		hc := &Hcload[i]
		// 湿りコイルの場合
		if hc.Wet {
			// 空気温度出口の計算式
			eo := hc.Cmp.Elouts[0]
			// elo:空気湿度出口の計算式
			elo := hc.Cmp.Elouts[1]
			// elini:空気湿度出口の要素方程式の2つ目の変数
			elini := elo.Elins[1]
			// 空気絶対湿度の要素方程式の2つ目の変数にupo、upvに空気出口をつなげる
			elini.Upo = eo
			elini.Upv = eo
		}

		if hc.Type == 'W' {
			eo := hc.Cmp.Elouts[0]
			elo := hc.Cmp.Elouts[2]
			elini := elo.Elins[1]
			elin := elo.Elins[3]

			elini.Upo = eo.Elins[0].Upo
			elini.Upv = eo.Elins[0].Upo
			elin.Upo = eo
			elin.Upv = eo
			elini = elo.Elins[2]
			elin = elo.Elins[4]
			eo = hc.Cmp.Elouts[1]

			elini.Upo = eo.Elins[0].Upo
			elini.Upv = eo.Elins[0].Upo
			elin.Upo = eo
			elin.Upv = eo
		}
	}
}

/* -------------------------------------------------- */

/* ルームエアコン（事業主基準モデル）機器仕様の入力処理 */

func rmacdat(Hcld *HCLOAD) {
	const (
		ERRFMT = "%s (rmacdat)"
		SCHAR  = 256
	)
	ss := Hcld.Cmp.Tparm

	Hcld.Qhmax = -999.0
	Hcld.Qh = -999.0
	Hcld.COPc = -999.0
	Hcld.COPh = -999.0
	Hcld.Qcmax = 999.0
	Hcld.Qc = 999.0

	parseFloat := func(value string) (float64, error) {
		return strconv.ParseFloat(value, 64)
	}

	for {
		var s string
		_, err := fmt.Sscanf(ss, "%s", &s)
		if err != nil || strings.Contains(string(s), "*") {
			break
		}

		ss = ss[len(s):]
		ss = strings.TrimLeftFunc(ss, unicode.IsSpace)

		keyValue := strings.SplitN(string(s), "=", 2)
		if len(keyValue) != 2 {
			Eprint("<rmacdat>", string(s))
			continue
		}

		key, value := keyValue[0], keyValue[1]
		switch key {
		case "Qc":
			if Hcld.Qc, err = parseFloat(value); err != nil {
				panic(err)
			}
		case "Qcmax":
			if Hcld.Qcmax, err = parseFloat(value); err != nil {
				panic(err)
			}
		case "Qh":
			if Hcld.Qh, err = parseFloat(value); err != nil {
				panic(err)
			}
		case "Qhmax":
			if Hcld.Qhmax, err = parseFloat(value); err != nil {
				panic(err)
			}
		case "COPc":
			if Hcld.COPc, err = parseFloat(value); err != nil {
				panic(err)
			}
		case "COPh":
			if Hcld.COPh, err = parseFloat(value); err != nil {
				panic(err)
			}
		default:
			Eprint("<rmacdat>", key)
		}
	}

	if Hcld.Qc < 0.0 {
		Hcld.rc = Hcld.Qcmax / Hcld.Qc
		Hcld.Ec = -Hcld.Qc / Hcld.COPc
	}
	if Hcld.Qh > 0.0 {
		Hcld.rh = Hcld.Qhmax / Hcld.Qh
		Hcld.Eh = Hcld.Qh / Hcld.COPh
	}
}

/* ルームエアコン（電中研モデル）機器仕様の入力処理 */

func rmacddat(Hcld *HCLOAD) {
	//Err := fmt.Sprintf(ERRFMT, "(rmacddat)")

	ss := Hcld.Cmp.Tparm
	Hcld.Qcmax, Hcld.Qc, Hcld.Qcmin = 999.0, 999.0, 999.0
	Hcld.Ecmax, Hcld.Ec, Hcld.Ecmin, Hcld.Qh, Hcld.Qhmax, Hcld.Qhmin, Hcld.Ehmax, Hcld.Eh, Hcld.Ehmin = -999.0, -999.0, -999.0, -999.0, -999.0, -999.0, -999.0, -999.0, -999.0
	Hcld.Gi, Hcld.Go = -999.0, -999.0

	for {
		var s string
		_, err := fmt.Sscanf(ss, "%s", &s)
		if err != nil {
			break
		}
		ss = ss[len(s):]
		for len(ss) > 0 && unicode.IsSpace(rune(ss[0])) {
			ss = ss[1:]
		}

		if st := strings.IndexByte(s, '='); st != -1 {
			key := s[:st]
			value := s[st+1:]

			switch key {
			case "Qc":
				Hcld.Qc, _ = strconv.ParseFloat(value, 64)
			case "Qcmax":
				Hcld.Qcmax, _ = strconv.ParseFloat(value, 64)
			case "Qcmin":
				Hcld.Qcmin, _ = strconv.ParseFloat(value, 64)
			case "Ec":
				Hcld.Ec, _ = strconv.ParseFloat(value, 64)
			case "Ecmax":
				Hcld.Ecmax, _ = strconv.ParseFloat(value, 64)
			case "Ecmin":
				Hcld.Ecmin, _ = strconv.ParseFloat(value, 64)
			case "Qh":
				Hcld.Qh, _ = strconv.ParseFloat(value, 64)
			case "Qhmax":
				Hcld.Qhmax, _ = strconv.ParseFloat(value, 64)
			case "Qhmin":
				Hcld.Qhmin, _ = strconv.ParseFloat(value, 64)
			case "Eh":
				Hcld.Eh, _ = strconv.ParseFloat(value, 64)
			case "Ehmax":
				Hcld.Ehmax, _ = strconv.ParseFloat(value, 64)
			case "Ehmin":
				Hcld.Ehmin, _ = strconv.ParseFloat(value, 64)
			case "Gi":
				Hcld.Gi, _ = strconv.ParseFloat(value, 64)
			case "Go":
				Hcld.Go, _ = strconv.ParseFloat(value, 64)
			default:
				Eprint("<rmacddat>", key)
			}
		} else {
			Eprint("<rmacddat>", s)
		}
	}

	// 機器固有値の計算
	if Hcld.Qc < 0.0 {
		Hcld.COPc = -Hcld.Qc / Hcld.Ec
		Hcld.COPcmax = -Hcld.Qcmax / Hcld.Ecmax
		Hcld.COPcmin = -Hcld.Qcmin / Hcld.Ecmin

		// JISにおける温湿度条件
		DBco := 35.0
		DBci := 27.0
		WBco := 24.0
		WBci := 19.0
		// 絶対湿度
		xco := FNXtw(DBco, WBco)
		xci := FNXtw(DBci, WBci)
		// 湿り比熱
		cao := Ca + Cv*xco
		cai := Ca + Cv*xci

		// 室内機、室外機熱交換器のバイパスファクタ
		Hcld.BFi = 0.2
		Hcld.BFo = 0.2

		// 理論効率の計算
		// 定格条件
		effthr := FNeffthc(DBco, DBci, xci, Hcld.Qc, Hcld.Ec, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)
		// 最小能力
		effthmin := FNeffthc(DBco, DBci, xci, Hcld.Qcmin, Hcld.Ecmin, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)
		// 最大能力
		effthmax := FNeffthc(DBco, DBci, xci, Hcld.Qcmax, Hcld.Ecmax, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)

		// ファン等動力の計算
		X := Hcld.COPcmin * Hcld.Qcmin / effthmin * effthr / (Hcld.Qc * Hcld.COPc)
		Hcld.Pcc = (-Hcld.Qcmin - X*(-Hcld.Qc)) / (Hcld.COPcmin - X*Hcld.COPc)

		// 定格条件、最小能力時の理論COPと実働COPの比R（両条件のRは等しいと仮定）
		Rr := (-Hcld.Qc * Hcld.COPc) / (effthr * (-Hcld.Qc - Hcld.Pcc*Hcld.COPc))
		// 最大能力時のRを計算
		Rmax := (-Hcld.Qcmax * Hcld.COPcmax) / (effthmax * (-Hcld.Qcmax - Hcld.Pcc*Hcld.COPcmax))

		// Rの回帰式係数の計算
		U := make([]float64, 9)
		Qc := make([]float64, 3)
		R := make([]float64, 3)
		Qc[0] = -Hcld.Qcmin
		Qc[1] = -Hcld.Qc
		Qc[2] = -Hcld.Qcmax
		R[0] = Rr
		R[1] = Rr
		R[2] = Rmax
		// 行列Uの作成
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				U[i*3+j] = math.Pow(Qc[i], 2.0-float64(j))
			}
		}

		// Uの逆行列の計算
		Matinv(U, 3, 3, "<rmacddat> UX")

		// 回帰係数の計算
		//Hcld.Rc = make([]float64, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				Hcld.Rc[i] += U[i*3+j] * R[j]
			}
		}
	}

	if Hcld.Qh > 0.0 {
		Hcld.COPh = Hcld.Qh / Hcld.Eh
		Hcld.COPhmax = Hcld.Qhmax / Hcld.Ehmax
		Hcld.COPhmin = Hcld.Qhmin / Hcld.Ehmin

		// JISにおける温湿度条件
		DBco := 7.0
		DBci := 20.0
		WBco := 6.0
		WBci := 15.0
		// 絶対湿度
		xco := FNXtw(DBco, WBco)
		xci := FNXtw(DBci, WBci)
		// 湿り比熱
		cao := Ca + Cv*xco
		cai := Ca + Cv*xci

		// 室内機、室外機熱交換器のバイパスファクタ
		Hcld.BFi = 0.2
		Hcld.BFo = 0.2

		// 理論効率の計算
		// 定格条件
		effthr := FNeffthh(DBco, DBci, xco, Hcld.Qh, Hcld.Eh, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)
		// 最小能力
		effthmin := FNeffthh(DBco, DBci, xco, Hcld.Qhmin, Hcld.Ehmin, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)
		// 最大能力
		effthmax := FNeffthh(DBco, DBci, xco, Hcld.Qhmax, Hcld.Ehmax, (1.0-Hcld.BFo)*Hcld.Go, (1.0-Hcld.BFi)*Hcld.Gi, cai, cao)

		// ファン等動力の計算
		X := Hcld.COPhmin * Hcld.Qhmin / effthmin * effthr / (Hcld.Qh * Hcld.COPh)
		Hcld.Pch = (Hcld.Qhmin - X*Hcld.Qh) / (Hcld.COPhmin - X*Hcld.COPh)

		// 定格条件、最小能力時の理論COPと実働COPの比R（両条件のRは等しいと仮定）
		Rr := (Hcld.Qh * Hcld.COPh) / (effthr * (Hcld.Qh - Hcld.Pch*Hcld.COPh))
		// 最大能力時のRを計算
		Rmax := (Hcld.Qhmax * Hcld.COPhmax) / (effthmax * (Hcld.Qhmax - Hcld.Pch*Hcld.COPhmax))

		// Rの回帰式係数の計算
		U := make([]float64, 9)
		Qc := make([]float64, 3)
		R := make([]float64, 3)
		Qc[0] = Hcld.Qhmin
		Qc[1] = Hcld.Qh
		Qc[2] = Hcld.Qhmax
		R[0] = Rr
		R[1] = Rr
		R[2] = Rmax
		// 行列Uの作成
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				U[i*3+j] = math.Pow(Qc[i], 2.0-float64(j))
			}
		}

		// Uの逆行列の計算
		Matinv(U, 3, 3, "<rmacddat> UX")

		// 回帰係数の計算
		//Hcld.Rh = make([]float64, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				Hcld.Rh[i] += U[i*3+j] * R[j]
			}
		}
	}
}

// 冷房時の理論COPを計算する
func FNeffthc(Tco, Tci, xci, Qc, Ec, Go, Gi, cai, cao float64) float64 {
	// 凝縮温度の計算
	Tcnd := FNTcndc(Tco, Qc, Ec, cao, Go)
	// 蒸発温度の計算
	Tevp := FNTevpc(Tci, Qc, cai, Gi, xci)

	// 理論効率の計算
	return (Tevp + 273.15) / (Tcnd - Tevp)
}

// 暖房時の理論COPを計算する
func FNeffthh(Tco, Tci, xco, Qc, Eh, Go, Gi, cai, cao float64) float64 {
	// 凝縮温度の計算
	Tcnd := FNTcndh(Tci, Qc, cai, Gi)
	// 蒸発温度の計算
	Tevp := FNTevph(Tco, Qc, Eh, cao, Go, xco)

	// 理論効率の計算
	return (Tcnd + 273.15) / (Tcnd - Tevp)
}

// 冷房時凝縮温度の計算
func FNTcndc(Tco, Qc, Ec, cao, Go float64) float64 {
	return (Tco + (-Qc+Ec)/(cao*Go))
}

// 暖房時凝縮温度の計算
func FNTcndh(Tci, Qc, cai, Gi float64) float64 {
	return (Tci + Qc/(cai*Gi))
}

// 冷房時蒸発温度の計算
func FNTevpc(Tci, Qc, cai, Gi, xci float64) float64 {

	// 蒸発温度の計算
	Tevp := Tci - (-Qc)/(cai*Gi)
	// 室内が結露するかどうかの判定（結露時は等エンタルピー変化による飽和状態とする）
	RHi := FNRhtx(Tevp, xci)
	if RHi > 100.0 {
		Tevp = FNDbrh(100.0, FNH(Tevp, xci))
	}

	return (Tevp)
}

// 暖房時蒸発温度の計算
func FNTevph(Tco, Qc, Eh, cao, Go, xco float64) float64 {
	// 蒸発温度の計算
	Tevp := Tco - (Qc-Eh)/(cao*Go)
	// 室外が結露するかどうかの判定（結露時は等エンタルピー変化による飽和状態とする）
	RHo := FNRhtx(Tevp, xco)
	if RHo > 100.0 {
		Tevp = FNDbrh(100.0, FNH(Tevp, xco))
	}
	return Tevp
}

/*  特性式の係数  */

//
// +--------+ ---> [OUT 1]
// | HCLOAD | ---> [OUT 2]
// +--------+ ---> [OUT 3] 冷温水コイル想定時のみ
//
func Hcldcfv(_Hcload []HCLOAD) {
	var f0, f1 float64

	Tout15 := 15.0
	Tout20 := 20.0

	for i := range _Hcload {
		Hcload := &_Hcload[i]
		Xout15 := FNXtr(Tout15, Hcload.RHout)
		Xout20 := FNXtr(Tout20, Hcload.RHout)
		f1 = (Xout20 - Xout15) / (Tout20 - Tout15)
		f0 = Xout15 - f1*Tout15

		Eo1 := Hcload.Cmp.Elouts[0]
		Hcload.Ga = Eo1.G

		if Eo1.Control != OFF_SW {
			Hcload.Ga = Eo1.G
			Hcload.CGa = Spcheat(Eo1.Fluid) * Hcload.Ga

			Eo1.Coeffo = Hcload.CGa
			Eo1.Co = 0.0
			Eo1.Coeffin[0] = -Hcload.CGa
		}

		Eo2 := Hcload.Cmp.Elouts[1]
		if Eo2.Control != OFF_SW {
			if Hcload.Wetmode {
				Eo2.Coeffo = 1.0
				Eo2.Co = f0
				Eo2.Coeffin[0] = 0.0
				Eo2.Coeffin[1] = -f1
			} else {
				Eo2.Coeffo = Hcload.Ga
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -Hcload.Ga
				Eo2.Coeffin[1] = 0.0
			}
		}

		// 冷温水コイル想定時
		if Hcload.Type == 'W' {
			Eo3 := Hcload.Cmp.Elouts[2]
			if Eo3.Control != OFF_SW {
				Hcload.Gw = Eo3.G
				Hcload.CGw = Spcheat(Eo3.Fluid) * Hcload.Gw
				rGa := Ro * Hcload.Ga

				Eo3.Coeffo = Hcload.CGw

				Eo3.Coeffin[0] = -Hcload.CGw
				Eo3.Coeffin[1] = -Hcload.CGa
				Eo3.Coeffin[2] = -rGa

				if Hcload.Wetmode && Hcload.Chmode == COOLING_SW {
					Eo3.Coeffin[3] = Hcload.CGa + rGa*f1
					Eo3.Coeffin[4] = 0.0
					Eo3.Co = -rGa * f0
				} else {
					Eo3.Coeffin[3] = Hcload.CGa
					Eo3.Coeffin[4] = rGa
					Eo3.Co = 0.0
				}
			}
		}
	}
}

/* ------------------------------------------ */

/* 空調負荷の計算 */

func Hcldene(_Hcload []HCLOAD, LDrest *int, Wd *WDAT) {
	var rest int
	var elo *ELOUT
	ro := 0.0
	ca := 0.0
	cv := 0.0

	rest = 0

	for i := range _Hcload {
		Hcload := &_Hcload[i]
		elo = Hcload.Cmp.Elouts[0]
		Hcload.Tain = elo.Elins[0].Sysvin
		elo = Hcload.Cmp.Elouts[1]
		Hcload.Xain = elo.Elins[0].Sysvin
		Hcload.Qfusoku = 0.0

		if Hcload.Type == 'W' {
			elo = Hcload.Cmp.Elouts[2]
			if elo.Elins[0].Upv != nil {
				Hcload.Twin = elo.Elins[0].Upv.Sysv
				Hcload.Twout = elo.Sysv
			} else {
				Hcload.Twin = -999.0
				Hcload.Twout = -999.0
			}
		}

		if Hcload.Cmp.Control != OFF_SW {
			elo = Hcload.Cmp.Elouts[0]
			if elo.Control == ON_SW && elo.Sysld == 'y' {
				Hcload.Qs = elo.Load
			} else {
				Hcload.Qs = Hcload.CGa * (elo.Sysv - Hcload.Tain)
			}

			rest = chswreset(Hcload.Qs, Hcload.Chmode, elo)

			if rest != 0 {
				(*LDrest)++
				Hcload.Cmp.Control = OFF_SW
			}

			elo = Hcload.Cmp.Elouts[1]
			if elo.Control == ON_SW && elo.Sysld == 'y' {
				Hcload.Ql = ro * elo.Load
			} else {
				Hcload.Ql = ro * Hcload.Ga * (elo.Sysv - Hcload.Xain)
			}

			if chqlreset(Hcload) != 0 {
				(*LDrest)++
			}

			Hcload.Qt = Hcload.Qs + Hcload.Ql

			if Hcload.RMACFlg == 'Y' {
				if Hcload.Qt > 0.0 {
					var qrhmax, qrhd, Temp float64
					var To_7 float64
					var Cafh, Cdf float64
					var Qhmax, Qhd float64
					var fht1, fht2, fht3 float64

					To_7 = Wd.T - 7.0
					Temp = (Hcload.rh - 1.0) / 1.8
					Cafh = 0.8

					qrhmax = -1.0e-6*(1.0+Temp)*math.Pow(To_7, 3.0) +
						2.0e-4*(1.0+Temp)*math.Pow(To_7, 2.0) +
						(0.0134+(0.0457-0.0134)*Temp)*To_7 + Hcload.rh

					Cdf = 1.0
					if Wd.T < 5.0 && Wd.RH >= 80.0 {
						Cdf = 0.9
					}

					Qhmax = qrhmax * Hcload.Qh * Cafh * Cdf

					if Qhmax > Hcload.Qt {
						Qhd = Hcload.Qt
					} else {
						Qhd = Qhmax
						Hcload.Qfusoku = Hcload.Qt - Qhmax
					}
					qrhd = Qhd / Hcload.Qh / Cafh / Cdf

					fht1 = fhtlb(Wd.T, 1.0)
					fht2 = fhtlb(Wd.T, qrhd*1.9/Hcload.rh)
					fht3 = fhtlb(Wd.T, 1.9/Hcload.rh)

					Eff := 1.0 / Cafh / Cdf * fht1 * fht2 / qrhd / fht3
					Hcload.Ele = Eff / Hcload.COPh * Qhd
					Hcload.COP = Qhd / Hcload.Ele

					Hcload.Ele = Hcload.Qt / Hcload.COP
				} else {
					var qrcmax, qrcd float64
					var To_35 float64
					var Cafc, Chm float64
					var Qcmax, Qcd float64
					var fct1, fct2, fct3 float64

					To_35 = Wd.T - 35.0
					Cafc = 0.85
					Chm = 1.0

					qrcmax = -1.0e-5*Hcload.rc*math.Pow(To_35, 3.0) +
						2.0e-4*0.5*(1.0+Hcload.rc)*math.Pow(To_35, 2.0) -
						(0.0147+0.014*(Hcload.rc-1.0))*To_35 + Hcload.rc

					Qcmax = qrcmax * Hcload.Qc * (Cafc + Chm) / 2.0

					if Qcmax < Hcload.Qt {
						Qcd = Hcload.Qt
					} else {
						Qcd = Qcmax
						Hcload.Qfusoku = Hcload.Qt - Qcmax
					}
					qrcd = Qcd / Hcload.Qc / ((Cafc + Chm) / 2.0)

					fct1 = fctlb(Wd.T, 1.0)
					fct2 = fctlb(Wd.T, qrcd*1.5/Hcload.rc)
					fct3 = fctlb(Wd.T, 1.5/Hcload.rc)

					Hcload.Ele = 1.0 / ((Cafc + Chm) / 2.0) * fct1 * fct2 / qrcd / fct3 / Hcload.COPc * Qcd
					Hcload.COP = Qcd / Hcload.Ele

					Hcload.Ele = -Hcload.Qt / Hcload.COP
				}
			} else if Hcload.RMACFlg == 'y' {
				if Hcload.Qt < 0.0 {
					var effth, Tevp, Tcnd, cai, cao, COP, COPd, R float64
					Qc := make([]float64, 3)

					cao = ca + cv*Wd.X
					cai = ca + cv*Hcload.Xain
					Tevp = FNTevpc(Hcload.Tain, Hcload.Qt, cai, (1.0-Hcload.BFi)*Hcload.Ga, Hcload.Xain)

					COPd = Hcload.COPc

					Qc[0] = Hcload.Qt * Hcload.Qt
					Qc[1] = -Hcload.Qt
					Qc[2] = 1.0

					R = 0.0
					for j := 0; j < 3; j++ {
						R += Hcload.Rc[j] * Qc[j]
					}

					for i := 0; i < 100; i++ {
						var E float64

						E = -Hcload.Qt / COPd

						Tcnd = FNTcndc(Wd.T, Hcload.Qt, E, cao, (1.0-Hcload.BFo)*Hcload.Go)

						effth = (Tevp + 273.15) / (Tcnd - Tevp)

						COP = 1.0 / (1.0/(R*effth) + Hcload.Pcc/-Hcload.Qt)

						if math.Abs(COP-COPd) < 1.0e-4 {
							Hcload.Ele = E
							Hcload.COP = COP
							break
						} else {
							COPd = COP
						}
					}
				} else {
					var effth, Tevp, Tcnd, cai, cao, COP, COPd, R float64
					Qc := make([]float64, 3)

					cao = ca + cv*Wd.X
					cai = ca + cv*Hcload.Xain
					Tcnd = FNTcndh(Hcload.Tain, Hcload.Qt, cai, (1.0-Hcload.BFi)*Hcload.Ga)

					COPd = Hcload.COPh

					Qc[0] = Hcload.Qt * Hcload.Qt
					Qc[1] = Hcload.Qt
					Qc[2] = 1.0

					R = 0.0
					for j := 0; j < 3; j++ {
						R += Hcload.Rh[j] * Qc[j]
					}

					for i := 0; i < 100; i++ {
						var E float64

						E = Hcload.Qt / COPd

						Tevp = FNTevph(Wd.T, Hcload.Qt, E, cao, (1.0-Hcload.BFo)*Hcload.Go, Wd.X)

						effth = (Tcnd + 273.15) / (Tcnd - Tevp)

						COP = 1.0 / (1.0/(R*effth) + Hcload.Pch/Hcload.Qt)

						if math.Abs(COP-COPd) < 1.0e-4 {
							Hcload.Ele = E
							Hcload.COP = COP
							break
						} else {
							COPd = COP
						}
					}
				}
			} else {
				Hcload.Qs = 0.0
				Hcload.Ql = 0.0
				Hcload.Qt = 0.0
				Hcload.Ele = 0.0
				Hcload.COP = 0.0
				Hcload.Qfusoku = 0.0
			}
		} else {
			Hcload.Qs = 0.0
			Hcload.Ql = 0.0
			Hcload.Qt = 0.0
			Hcload.Ele = 0.0
			Hcload.COP = 0.0
			Hcload.Qfusoku = 0.0
		}
	}
}

func fctlb(T, x float64) float64 {
	a := [...]float64{
		0.0148*T + 0.0089,
		-0.0153*T + 0.1429,
		0.034*T - 0.4963,
		-0.0012*T + 0.288 + 0.0322,
	}

	var Temp float64
	for i := 0; i < 4; i++ {
		Temp += a[i] * math.Pow(x, float64(i))
	}

	return Temp
}

func fhtlb(T, x float64) float64 {
	a := [...]float64{
		0.0018*T*T - 0.0424*T + 0.4554,
		-0.006*T*T + 0.1347*T - 1.56,
		0.0063*T*T - 0.1406*T + 2.2902,
		-0.002*T*T + 0.0176*T - 0.3789,
		0.0002*T*T - 0.0007*T + 0.4202,
	}

	var Temp float64
	for i := 0; i < 5; i++ {
		Temp += a[i] * math.Pow(x, float64(i))
	}

	return Temp
}

/* --------------------------- */

/* 負荷計算指定時の設定値のポインター */

func hcldptr(load *ControlSWType, key []string, Hcload *HCLOAD, idmrk *byte) (VPTR, error) {
	var err error
	var vptr VPTR
	if key[1] == "Tout" || key[1] == "Tr" || key[1] == "Tot" {
		vptr.Ptr = &Hcload.Toset
		vptr.Type = VAL_CTYPE
		Hcload.Loadt = load
		*idmrk = 't'
	} else if key[1] == "xout" {
		vptr.Ptr = &Hcload.Xoset
		vptr.Type = VAL_CTYPE
		Hcload.Loadx = load
		*idmrk = 'x'
	} else {
		err = errors.New("Tout, Tr, Tot or xout are expected")
	}
	return vptr, err
}

/* ------------------------------------------ */

/* 負荷計算指定時のスケジュール設定 */

func hcldschd(Hcload *HCLOAD) {
	Eo := Hcload.Cmp.Elouts

	if Hcload.Loadt != nil {
		if Eo[0].Control != OFF_SW {
			if Hcload.Toset > TEMPLIMIT {
				Eo[0].Control = LOAD_SW
				Eo[0].Sysv = Hcload.Toset
			} else {
				Eo[0].Control = OFF_SW

				if Hcload.Wetmode {
					Eo[1].Control = OFF_SW
				}
			}
		}
	} else if Hcload.Loadx != nil {
		if len(Eo) > 1 && Eo[1].Control != OFF_SW {
			if Hcload.Xoset > 0.0 {
				Eo[1].Control = LOAD_SW
				Eo[1].Sysv = Hcload.Xoset
			} else {
				Eo[1].Control = OFF_SW
			}
		}
	}

	if Hcload.Type == 'W' {
		if len(Eo) > 2 && Eo[0].Control == OFF_SW && Eo[1].Control == OFF_SW {
			Eo[2].Control = OFF_SW
		}
	}
}

/* ------------------------------------------ */

func hcldprint(fo io.Writer, id int, _Hcload []HCLOAD) {
	switch id {
	case 0:
		if len(_Hcload) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCLOAD_TYPE, len(_Hcload))
		}
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			if Hcload.Type == 'W' {
				fmt.Fprintf(fo, " %s 1 15\n", Hcload.Name)
			} else {
				if Hcload.RMACFlg == 'Y' || Hcload.RMACFlg == 'y' {
					fmt.Fprintf(fo, " %s 1 14\n", Hcload.Name)
				} else {
					fmt.Fprintf(fo, " %s 1 11\n", Hcload.Name)
				}
			}
		}
	case 1:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s_ca c c %s_Ga m f %s_Ti t f %s_To t f %s_Qs q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_cx c c %s_xi x f %s_xo x f %s_RHo r f %s_Ql q f",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			if Hcload.Type == 'W' {
				fmt.Fprintf(fo, "%s_cw c c %s_G m f %s_Twi t f %s_Two t f",
					Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			}
			fmt.Fprintf(fo, " %s_Qt q f\n", Hcload.Name)

			if Hcload.RMACFlg == 'Y' || Hcload.RMACFlg == 'y' {
				fmt.Fprintf(fo, " %s_Qfusoku q f %s_Ele E f %s_COP C f\n",
					Hcload.Name, Hcload.Name, Hcload.Name)
			}
		}
	default:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			el := Hcload.Cmp.Elouts[0]
			Taout := el.Sysv
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.2f ",
				el.Control, Hcload.Ga, Hcload.Tain, el.Sysv, Hcload.Qs)

			el = Hcload.Cmp.Elouts[1]
			RHout := FNRhtx(Taout, float64(el.Sysv))
			if RHout > 100.0 {
				RHout = 999
			} else if RHout < 0.0 {
				RHout = -99.0
			}

			fmt.Fprintf(fo, "%c %.4f %.4f %3.0f %.2f ",
				el.Control, Hcload.Xain, el.Sysv, RHout, Hcload.Ql)

			if Hcload.Type == 'W' {
				el = Hcload.Cmp.Elouts[2]
				fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f",
					el.Control, Hcload.Gw, Hcload.Twin, el.Sysv)
			}
			fmt.Fprintf(fo, " %.2f\n", Hcload.Qt)

			if Hcload.RMACFlg == 'Y' || Hcload.RMACFlg == 'y' {
				fmt.Fprintf(fo, " %.0f %.0f %.2f\n", Hcload.Qfusoku, Hcload.Ele, Hcload.COP)
			}
		}
	}
}

/* ------------------------------ */

/* 日積算値に関する処理 */

func hclddyint(_Hcload []HCLOAD) {
	for i := range _Hcload {
		Hcload := &_Hcload[i]
		svdyint(&Hcload.Taidy)
		svdyint(&Hcload.xaidy)

		qdyint(&Hcload.Qdys)
		qdyint(&Hcload.Qdyl)
		qdyint(&Hcload.Qdyt)
		qdyint(&Hcload.Qdyfusoku)
		qdyint(&Hcload.Edy)
	}
}

func hcldmonint(_Hcload []HCLOAD) {
	for i := range _Hcload {
		Hcload := &_Hcload[i]
		svdyint(&Hcload.mTaidy)
		svdyint(&Hcload.mxaidy)

		qdyint(&Hcload.mQdys)
		qdyint(&Hcload.mQdyl)
		qdyint(&Hcload.mQdyt)
		qdyint(&Hcload.mQdyfusoku)
		qdyint(&Hcload.mEdy)
	}
}

func hcldday(Mon, Day, ttmm, Nday, SimDayend int, _Hcload []HCLOAD) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for i := range _Hcload {
		Hcload := &_Hcload[i]

		// 日集計
		svdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Tain, &Hcload.Taidy)
		svdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Xain, &Hcload.xaidy)

		qdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Qs, &Hcload.Qdys)
		qdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Ql, &Hcload.Qdyl)
		qdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Qt, &Hcload.Qdyt)

		qdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Qfusoku, &Hcload.Qdyfusoku)
		qdaysum(int64(ttmm), Hcload.Cmp.Control, Hcload.Ele, &Hcload.Edy)

		// 月集計
		svmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Tain, &Hcload.mTaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Xain, &Hcload.mxaidy, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Qs, &Hcload.mQdys, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Ql, &Hcload.mQdyl, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Qt, &Hcload.mQdyt, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Qfusoku, &Hcload.mQdyfusoku, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Ele, &Hcload.mEdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, Hcload.Cmp.Control, Hcload.Ele, &Hcload.mtEdy[Mo][tt])
	}
}

func hclddyprt(fo io.Writer, id int, _Hcload []HCLOAD) {
	switch id {
	case 0:
		if len(_Hcload) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCLOAD_TYPE, len(_Hcload))
		}
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s 4 36 14 14 8\n", Hcload.Name)
		}

	case 1:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
		}

	default:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcload.Taidy.Hrs, Hcload.Taidy.M,
				Hcload.Taidy.Mntime, Hcload.Taidy.Mn,
				Hcload.Taidy.Mxtime, Hcload.Taidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdys.Hhr, Hcload.Qdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdys.Chr, Hcload.Qdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.Qdys.Hmxtime, Hcload.Qdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.Qdys.Cmxtime, Hcload.Qdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				Hcload.xaidy.Hrs, Hcload.xaidy.M,
				Hcload.xaidy.Mntime, Hcload.xaidy.Mn,
				Hcload.xaidy.Mxtime, Hcload.xaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdyl.Hhr, Hcload.Qdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdyl.Chr, Hcload.Qdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.Qdyl.Hmxtime, Hcload.Qdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.Qdyl.Cmxtime, Hcload.Qdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdyt.Hhr, Hcload.Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.Qdyt.Chr, Hcload.Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.Qdyt.Hmxtime, Hcload.Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hcload.Qdyt.Cmxtime, Hcload.Qdyt.Cmx)
		}
	}
}

func hcldmonprt(fo io.Writer, id int, _Hcload []HCLOAD) {
	switch id {
	case 0:
		if len(_Hcload) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCLOAD_TYPE, len(_Hcload))
		}
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s 4 36 14 14 8\n", Hcload.Name)
		}

	case 1:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				Hcload.Name, Hcload.Name, Hcload.Name, Hcload.Name)
		}

	default:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcload.mTaidy.Hrs, Hcload.mTaidy.M,
				Hcload.mTaidy.Mntime, Hcload.mTaidy.Mn,
				Hcload.mTaidy.Mxtime, Hcload.mTaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdys.Hhr, Hcload.mQdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdys.Chr, Hcload.mQdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.mQdys.Hmxtime, Hcload.mQdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.mQdys.Cmxtime, Hcload.mQdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				Hcload.mxaidy.Hrs, Hcload.mxaidy.M,
				Hcload.mxaidy.Mntime, Hcload.mxaidy.Mn,
				Hcload.mxaidy.Mxtime, Hcload.mxaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdyl.Hhr, Hcload.mQdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdyl.Chr, Hcload.mQdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.mQdyl.Hmxtime, Hcload.mQdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.mQdyl.Cmxtime, Hcload.mQdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdyt.Hhr, Hcload.mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcload.mQdyt.Chr, Hcload.mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcload.mQdyt.Hmxtime, Hcload.mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hcload.mQdyt.Cmxtime, Hcload.mQdyt.Cmx)
		}
	}
}

func hcldmtprt(fo io.Writer, id, Mo, tt int, _Hcload []HCLOAD) {
	switch id {
	case 0:
		if len(_Hcload) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCLOAD_TYPE, len(_Hcload))
		}
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, " %s 1 1\n", Hcload.Name)
		}

	case 1:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, "%s_E E f \n", Hcload.Name)
		}

	default:
		for i := range _Hcload {
			Hcload := &_Hcload[i]
			fmt.Fprintf(fo, " %.2f \n", Hcload.mtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

func hcldswptr(key []string, Hcload *HCLOAD) (VPTR, error) {
	if key[1] == "chmode" {
		return VPTR{
			Ptr:  &Hcload.Chmode,
			Type: SW_CTYPE,
		}, nil
	}

	return VPTR{}, errors.New("hcldswptr error")
}

func chhcldswreset(Qload, Ql float64, chmode ControlSWType, Eo *ELOUT) int {
	if (chmode == HEATING_SW && Qload < 0.0) ||
		(chmode == COOLING_SW && Ql > 0.0) ||
		(chmode == COOLING_SW && Qload > 0.0) {
		Eo.Control = ON_SW
		Eo.Sysld = 'n'
		Eo.Emonitr.Control = ON_SW

		return 1
	} else {
		return 0
	}
}

func hcldwetmdreset(Eqsys *EQSYS) {
	Hcload := Eqsys.Hcload

	for i := range Hcload {
		Hcload[i].Wetmode = Hcload[i].Wet
	}
}
