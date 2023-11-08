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

/*  hccoil.c  */

/*  冷温水コイル  */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// 冷温水コイルの機器仕様入力
// See: ../format/EQPCAT.md#HCC
func Hccdata(s string, Hccca *HCCCA) int {
	var st string
	var dt float64
	id := 0

	if stIdx := strings.IndexRune(s, '='); stIdx == -1 {
		Hccca.name = s
		Hccca.eh = 0.0
		Hccca.et = -999.0
		Hccca.KA = -999.0
	} else {
		st = s[stIdx+1:]
		dt, _ = readFloat(st)

		if s == "et" {
			// コイル温度効率
			Hccca.et = dt
		} else if s == "eh" {
			// コイルエンタルピー効率
			Hccca.eh = dt
		} else if s == "KA" {
			// コイルの熱通過率と伝熱面積の積 [W/K]
			Hccca.KA = dt
		} else {
			id = 1
		}
	}

	return id
}

/* ------------------------------------------ */
func Hccdwint(_hcc []*HCC) {
	for _, hcc := range _hcc {

		// 乾きコイルと湿りコイルの判定
		if hcc.Cat.eh > 1.0e-10 {
			hcc.Wet = 'w' // 湿りコイル
		} else {
			hcc.Wet = 'd' // 乾きコイル
		}

		// 温度効率固定タイプと変動タイプの判定
		if hcc.Cat.et > 0.0 {
			hcc.Etype = 'e' // 定格(温度効率固定タイプ)
		} else if hcc.Cat.KA > 0.0 {
			hcc.Etype = 'k' // 変動タイプ
		} else {
			fmt.Printf("Hcc %s  Undefined Character et or KA\n", hcc.Name)
			hcc.Etype = '-'
		}

		// 入口水温、入口空気絶対湿度を初期化
		//Hcc.Twin = 5.0
		//Hcc.xain = FNXtr(25.0, 50.0)
	}
}

/* ------------------------------------------ */
/*  特性式の係数  */

//
// [IN 1] ----> +-----+ ----> [OUT 1] 空気の温度
// [IN 2] ----> | HCC | ----> [OUT 2] 空気の絶対湿度
// [IN 3] ----> +-----+ ----> [OUT 3] 水の温度
//
func Hcccfv(_hcc []*HCC) {
	for _, hcc := range _hcc {
		hcc.Ga = 0.0
		hcc.Gw = 0.0
		hcc.et = 0.0
		hcc.eh = 0.0

		// 経路が停止していなければ
		if hcc.Cmp.Control == OFF_SW {
			continue
		}

		// 機器出力は3つ
		if len(hcc.Cmp.Elouts) != 3 || len(hcc.Cmp.Elins) != 0 {
			panic("HCCの機器出力数は3、機器入力は0です。")
		}

		eo_ta := hcc.Cmp.Elouts[0] // 排気温度
		eo_xa := hcc.Cmp.Elouts[1] // 排気湿度
		eo_tw := hcc.Cmp.Elouts[2] // 排水温度

		var AirSW, WaterSW ControlSWType

		// 排気量・排気熱量
		hcc.Ga = eo_ta.G                        // 排気量
		hcc.cGa = Spcheat(eo_ta.Fluid) * hcc.Ga // 排気熱量
		if hcc.Ga > 0.0 {
			AirSW = ON_SW
		} else {
			AirSW = OFF_SW
		}

		// 排水量・排水熱量
		hcc.Gw = eo_tw.G                        // 排水量
		hcc.cGw = Spcheat(eo_tw.Fluid) * hcc.Gw // 排水熱量
		if hcc.Gw > 0.0 {
			WaterSW = ON_SW
		} else {
			WaterSW = OFF_SW
		}

		// 温度効率
		if hcc.Etype == 'e' {
			// 定格温度効率
			hcc.et = hcc.Cat.et
		} else if hcc.Etype == 'k' {
			// 温度効率を計算
			hcc.et = FNhccet(hcc.cGa, hcc.cGw, hcc.Cat.KA)
		} else {
			panic(hcc.Etype)
		}

		// エンタルピ効率 [-]
		hcc.eh = hcc.Cat.eh

		// 冷温水コイルの処理熱量
		hcc.Et, hcc.Ex, hcc.Ew = wcoil(AirSW, WaterSW, hcc.Wet, hcc.Ga*hcc.et, hcc.Ga*hcc.eh, hcc.Xain, hcc.Twin)

		// 排気温度に関する係数の設定
		eo_ta.Coeffo = hcc.cGa
		eo_ta.Co = -(hcc.Et.C)
		eo_ta.Coeffin[0] = hcc.Et.T - hcc.cGa
		eo_ta.Coeffin[1] = hcc.Et.X
		eo_ta.Coeffin[2] = -(hcc.Et.W)

		// 排気湿度に関する係数の設定
		eo_xa.Coeffo = hcc.Ga
		eo_xa.Co = -(hcc.Ex.C)
		eo_xa.Coeffin[0] = hcc.Ex.T
		eo_xa.Coeffin[1] = hcc.Ex.X - hcc.Ga
		eo_xa.Coeffin[2] = -(hcc.Ex.W)

		// 排水温度に関する係数の設定
		eo_tw.Coeffo = hcc.cGw
		eo_tw.Co = hcc.Ew.C
		eo_tw.Coeffin[0] = -(hcc.Ew.T)
		eo_tw.Coeffin[1] = -(hcc.Ew.X)
		eo_tw.Coeffin[2] = hcc.Ew.W - hcc.cGw
	}
}

/* ------------------------------------------ */

// 供給熱量の計算
func Hccdwreset(Hcc []*HCC, DWreset *int) {
	for i, hcc := range Hcc {
		xain := hcc.Cmp.Elins[1].Sysvin // <給気>絶対湿度 [kg/kg]
		Twin := hcc.Cmp.Elins[2].Sysvin // <給水>温水の温度 [C]

		reset := false
		if hcc.Cat.eh > 1.0e-10 {
			Tdp := FNDp(FNPwx(xain)) // 露点温度
			if hcc.Wet == 'w' && Twin > Tdp {
				// 露点温度を上回った => 結露なし (乾きコイル)
				hcc.Wet = 'd'
				reset = true
			} else if hcc.Wet == 'd' && Twin < Tdp {
				// 露点温度を上回った => 結露あり (湿りコイル)
				hcc.Wet = 'w'
				reset = true
			}

			if reset {
				(*DWreset)++
				Hcccfv(Hcc[i : i+1])
			}
		}
	}
}

/* ------------------------------------------ */

// 冷温水コイルHccの供給熱量 Qs, Ql, Qt の計算を行う。
func Hccene(Hcc []*HCC) {
	for _, hcc := range Hcc {
		hcc.Tain = hcc.Cmp.Elins[0].Sysvin // <給気>空気温度 [C]
		hcc.Xain = hcc.Cmp.Elins[1].Sysvin // <給気>絶対湿度 [kg/kg]
		hcc.Twin = hcc.Cmp.Elins[2].Sysvin // <給水>温水の温度 [C]

		if hcc.Cmp.Control != OFF_SW {
			// <排気>空気温度 [C]
			hcc.Taout = hcc.Cmp.Elouts[0].Sysv
			hcc.Qs = hcc.cGa * (hcc.Taout - hcc.Tain)

			// <排気>空気絶対湿度 [kg/kg]
			Xaout := hcc.Cmp.Elouts[1].Sysv
			hcc.Ql = Ro * hcc.Ga * (Xaout - hcc.Xain)

			// <排水>温水の温度 [C]
			hcc.Twout = hcc.Cmp.Elouts[2].Sysv
			hcc.Qt = hcc.cGw * (hcc.Twout - hcc.Twin)
		} else {
			// 経路が停止している場合は熱供給しない
			hcc.Qs = 0.0
			hcc.Ql = 0.0
			hcc.Qt = 0.0
		}
	}
}

/* ------------------------------------------ */

// 冷温水コイルHccの状態をfoに出力する。
func hccprint(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, " %s 1 16\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_ca c c %s_Ga m f %s_Ti t f %s_To t f %s_Qs q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_cx c c %s_xi x f %s_xo x f %s_Ql q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_cw c c %s_Gw m f %s_Twi t f %s_Two t f %s_Qt q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_et m f %s_eh m f\n\n", hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			// 給排気温度に関する事項
			eo_ta := hcc.Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ", eo_ta.Control, hcc.Ga, hcc.Tain, eo_ta.Sysv, hcc.Qs)

			// 給排気湿度に関する事項
			eo_xa := hcc.Cmp.Elouts[1]
			fmt.Fprintf(fo, "%c %5.3f %5.3f %2.0f ", eo_xa.Control, hcc.Xain, eo_xa.Sysv, hcc.Ql)

			// 給排水温度に関する事項
			eo_tw := hcc.Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ", eo_tw.Control, hcc.Gw, hcc.Twin, eo_tw.Sysv, hcc.Qt)

			// 温度効率、エンタルピー
			fmt.Fprintf(fo, "%6.4g %6.4g\n", hcc.et, hcc.eh)
		}
	}
}

/* ------------------------------ */

/* 日積算値に関する処理 */

func hccdyint(Hcc []*HCC) {
	for _, hcc := range Hcc {
		svdyint(&hcc.Taidy)
		svdyint(&hcc.xaidy)
		svdyint(&hcc.Twidy)
		qdyint(&hcc.Qdys)
		qdyint(&hcc.Qdyl)
		qdyint(&hcc.Qdyt)
	}
}

func hccmonint(Hcc []*HCC) {
	for _, hcc := range Hcc {
		svdyint(&hcc.mTaidy)
		svdyint(&hcc.mxaidy)
		svdyint(&hcc.mTwidy)
		qdyint(&hcc.mQdys)
		qdyint(&hcc.mQdyl)
		qdyint(&hcc.mQdyt)
	}
}

func hccday(Mon, Day, ttmm int, Hcc []*HCC, Nday, SimDayend int) {
	for _, hcc := range Hcc {
		// 日集計
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Tain, &hcc.Taidy)
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Xain, &hcc.xaidy)
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Twin, &hcc.Twidy)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Qs, &hcc.Qdys)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Ql, &hcc.Qdyl)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Qt, &hcc.Qdyt)

		// 月集計
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Tain, &hcc.mTaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Xain, &hcc.mxaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Twin, &hcc.mTwidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Qs, &hcc.mQdys, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Ql, &hcc.mQdyl, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Qt, &hcc.mQdyt, Nday, SimDayend)
	}
}

func hccdyprt(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.Taidy.Hrs, hcc.Taidy.M,
				hcc.Taidy.Mntime, hcc.Taidy.Mn,
				hcc.Taidy.Mxtime, hcc.Taidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdys.Hhr, hcc.Qdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdys.Chr, hcc.Qdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdys.Hmxtime, hcc.Qdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdys.Cmxtime, hcc.Qdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				hcc.xaidy.Hrs, hcc.xaidy.M,
				hcc.xaidy.Mntime, hcc.xaidy.Mn,
				hcc.xaidy.Mxtime, hcc.xaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyl.Hhr, hcc.Qdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyl.Chr, hcc.Qdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyl.Hmxtime, hcc.Qdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyl.Cmxtime, hcc.Qdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.Twidy.Hrs, hcc.Twidy.M,
				hcc.Twidy.Mntime, hcc.Twidy.Mn,
				hcc.Twidy.Mxtime, hcc.Twidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyt.Hhr, hcc.Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyt.Chr, hcc.Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyt.Hmxtime, hcc.Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hcc.Qdyt.Cmxtime, hcc.Qdyt.Cmx)
		}
	}
}

func hccmonprt(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.mTaidy.Hrs, hcc.mTaidy.M,
				hcc.mTaidy.Mntime, hcc.mTaidy.Mn,
				hcc.mTaidy.Mxtime, hcc.mTaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdys.Hhr, hcc.mQdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdys.Chr, hcc.mQdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdys.Hmxtime, hcc.mQdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdys.Cmxtime, hcc.mQdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				hcc.mxaidy.Hrs, hcc.mxaidy.M,
				hcc.mxaidy.Mntime, hcc.mxaidy.Mn,
				hcc.mxaidy.Mxtime, hcc.mxaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyl.Hhr, hcc.mQdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyl.Chr, hcc.mQdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyl.Hmxtime, hcc.mQdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyl.Cmxtime, hcc.mQdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.mTwidy.Hrs, hcc.mTwidy.M,
				hcc.mTwidy.Mntime, hcc.mTwidy.Mn,
				hcc.mTwidy.Mxtime, hcc.mTwidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyt.Hhr, hcc.mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyt.Chr, hcc.mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyt.Hmxtime, hcc.mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hcc.mQdyt.Cmxtime, hcc.mQdyt.Cmx)
		}
	}
}

// 温水コイルの温度効率計算関数
// 計算モデルは向流コイル
func FNhccet(Wa, Ww, KA float64) float64 {
	Ws := Wa
	Wl := Ww

	NTU := KA / Ws
	C := Ws / Wl
	B := (1.0 - C) * NTU

	if math.Abs(Ws-Wl) < 1.0e-5 {
		return NTU / (1.0 + NTU)
	} else {
		if exB := math.Exp(-B); math.IsInf(exB, 0) {
			return 1.0 / C
		} else {
			return (1.0 - exB) / (1.0 - C*exB)
		}
	}
}
