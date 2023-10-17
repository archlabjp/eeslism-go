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
	"strconv"
	"strings"
)

/*  機器仕様入力 */

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
		dt, _ = strconv.ParseFloat(st, 64)

		if s == "et" {
			Hccca.et = dt
		} else if s == "eh" {
			Hccca.eh = dt
		} else if s == "KA" {
			Hccca.KA = dt
		} else {
			id = 1
		}
	}

	return id
}

/* ------------------------------------------ */
func Hccdwint(_hcc []HCC) {
	for i := range _hcc {
		hcc := &_hcc[i]

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
func Hcccfv(_hcc []HCC) {
	for i := range _hcc {
		hcc := &_hcc[i]

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
		Eo1 := hcc.Cmp.Elouts[0]
		Eo2 := hcc.Cmp.Elouts[1]
		Eo3 := hcc.Cmp.Elouts[2]

		var AirSW, WaterSW ControlSWType

		// 水の流量?
		hcc.Ga = Eo1.G                        // 水の流量?
		hcc.cGa = Spcheat(Eo1.Fluid) * hcc.Ga // 水の熱量?
		if hcc.Ga > 0.0 {
			AirSW = ON_SW
		} else {
			AirSW = OFF_SW
		}

		// 空気の流量?
		hcc.Gw = Eo3.G                        // 空気の流量?
		hcc.cGw = Spcheat(Eo3.Fluid) * hcc.Gw // 空気の熱量?
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
		}

		// エンタルピ効率
		hcc.eh = hcc.Cat.eh

		// 冷温水コイルの処理熱量
		hcc.Et, hcc.Ex, hcc.Ew = wcoil(AirSW, WaterSW, hcc.Wet, hcc.Ga*hcc.et, hcc.Ga*hcc.eh, hcc.Xain, hcc.Twin)

		// 空気の温度
		Eo1.Coeffo = hcc.cGa
		Eo1.Co = -(hcc.Et.C)
		Eo1.Coeffin[0] = hcc.Et.T - hcc.cGa
		Eo1.Coeffin[1] = hcc.Et.X
		Eo1.Coeffin[2] = -(hcc.Et.W)

		// 空気の絶対湿度
		Eo2.Coeffo = hcc.Ga
		Eo2.Co = -(hcc.Ex.C)
		Eo2.Coeffin[0] = hcc.Ex.T
		Eo2.Coeffin[1] = hcc.Ex.X - hcc.Ga
		Eo2.Coeffin[2] = -(hcc.Ex.W)

		// 水の温度
		Eo3.Coeffo = hcc.cGw
		Eo3.Co = hcc.Ew.C
		Eo3.Coeffin[0] = -(hcc.Ew.T)
		Eo3.Coeffin[1] = -(hcc.Ew.X)
		Eo3.Coeffin[2] = hcc.Ew.W - hcc.cGw
	}
}

/* ------------------------------------------ */

/* 供給熱量の計算 */

func Hccdwreset(Hcc []HCC, DWreset *int) {
	for i := range Hcc {
		hcc := &Hcc[i] // Get the address of the current element

		xain := hcc.Cmp.Elins[1].Sysvin
		Twin := hcc.Cmp.Elins[2].Sysvin

		reset := 0
		if hcc.Cat.eh > 1.0e-10 {
			Tdp := FNDp(FNPwx(xain))
			if hcc.Wet == 'w' && Twin > Tdp {
				hcc.Wet = 'd'
				reset = 1
			} else if hcc.Wet == 'd' && Twin < Tdp {
				hcc.Wet = 'w'
				reset = 1
			}

			if reset != 0 {
				(*DWreset)++
				Hcccfv(Hcc[i : i+1])
			}
		}
	}
}

/* ------------------------------------------ */

/* 供給熱量の計算 */

func Hccene(Hcc []HCC) {
	for i := range Hcc {
		hcc := &Hcc[i] // Get the address of the current element

		hcc.Tain = hcc.Cmp.Elins[0].Sysvin
		hcc.Xain = hcc.Cmp.Elins[1].Sysvin
		hcc.Twin = hcc.Cmp.Elins[2].Sysvin

		if hcc.Cmp.Control != OFF_SW {
			elo := hcc.Cmp.Elouts[0]
			hcc.Taout = elo.Sysv
			hcc.Qs = hcc.cGa * (elo.Sysv - hcc.Tain)

			elo = hcc.Cmp.Elouts[1]
			hcc.Ql = Ro * hcc.Ga * (elo.Sysv - hcc.Xain)

			elo = hcc.Cmp.Elouts[2]
			hcc.Twout = elo.Sysv
			hcc.Qt = hcc.cGw * (elo.Sysv - hcc.Twin)
		} else {
			hcc.Qs = 0.0
			hcc.Ql = 0.0
			hcc.Qt = 0.0
		}
	}
}

/* ------------------------------------------ */

func hccprint(fo io.Writer, id int, Hcc []HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for i := range Hcc {
			fmt.Fprintf(fo, " %s 1 16\n", Hcc[i].Name)
		}
	case 1:
		for i := range Hcc {
			fmt.Fprintf(fo, "%s_ca c c %s_Ga m f %s_Ti t f %s_To t f %s_Qs q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_cx c c %s_xi x f %s_xo x f %s_Ql q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_cw c c %s_Gw m f %s_Twi t f %s_Two t f %s_Qt q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_et m f %s_eh m f\n\n",
				Hcc[i].Name, Hcc[i].Name)
		}
	default:
		for i := range Hcc {
			el := Hcc[i].Cmp.Elouts[0] // Get the address of the first element
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ",
				el.Control, Hcc[i].Ga, Hcc[i].Tain, el.Sysv, Hcc[i].Qs)
			el = Hcc[i].Cmp.Elouts[1] // Get the address of the second element
			fmt.Fprintf(fo, "%c %5.3f %5.3f %2.0f ",
				el.Control, Hcc[i].Xain, el.Sysv, Hcc[i].Ql)
			el = Hcc[i].Cmp.Elouts[2] // Get the address of the third element
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ",
				el.Control, Hcc[i].Gw, Hcc[i].Twin, el.Sysv, Hcc[i].Qt)

			fmt.Fprintf(fo, "%6.4g %6.4g\n",
				Hcc[i].et, Hcc[i].eh)
		}
	}
}

/* ------------------------------ */

/* 日積算値に関する処理 */

func hccdyint(Hcc []HCC) {
	for i := range Hcc {
		svdyint(&Hcc[i].Taidy)
		svdyint(&Hcc[i].xaidy)
		svdyint(&Hcc[i].Twidy)
		qdyint(&Hcc[i].Qdys)
		qdyint(&Hcc[i].Qdyl)
		qdyint(&Hcc[i].Qdyt)
	}
}

func hccmonint(Hcc []HCC) {
	for i := range Hcc {
		svdyint(&Hcc[i].mTaidy)
		svdyint(&Hcc[i].mxaidy)
		svdyint(&Hcc[i].mTwidy)
		qdyint(&Hcc[i].mQdys)
		qdyint(&Hcc[i].mQdyl)
		qdyint(&Hcc[i].mQdyt)
	}
}

func hccday(Mon, Day, ttmm int, Hcc []HCC, Nday, SimDayend int) {
	for i := range Hcc {
		// 日集計
		svdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Tain, &Hcc[i].Taidy)
		svdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Xain, &Hcc[i].xaidy)
		svdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Twin, &Hcc[i].Twidy)
		qdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Qs, &Hcc[i].Qdys)
		qdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Ql, &Hcc[i].Qdyl)
		qdaysum(int64(ttmm), Hcc[i].Cmp.Control, Hcc[i].Qt, &Hcc[i].Qdyt)

		// 月集計
		svmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Tain, &Hcc[i].mTaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Xain, &Hcc[i].mxaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Twin, &Hcc[i].mTwidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Qs, &Hcc[i].mQdys, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Ql, &Hcc[i].mQdyl, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hcc[i].Cmp.Control, Hcc[i].Qt, &Hcc[i].mQdyt, Nday, SimDayend)
	}
}

func hccdyprt(fo io.Writer, id int, Hcc []HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for i := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", Hcc[i].Name)
		}
	case 1:
		for i := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
		}
	default:
		for i := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcc[i].Taidy.Hrs, Hcc[i].Taidy.M,
				Hcc[i].Taidy.Mntime, Hcc[i].Taidy.Mn,
				Hcc[i].Taidy.Mxtime, Hcc[i].Taidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdys.Hhr, Hcc[i].Qdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdys.Chr, Hcc[i].Qdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].Qdys.Hmxtime, Hcc[i].Qdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].Qdys.Cmxtime, Hcc[i].Qdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				Hcc[i].xaidy.Hrs, Hcc[i].xaidy.M,
				Hcc[i].xaidy.Mntime, Hcc[i].xaidy.Mn,
				Hcc[i].xaidy.Mxtime, Hcc[i].xaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdyl.Hhr, Hcc[i].Qdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdyl.Chr, Hcc[i].Qdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].Qdyl.Hmxtime, Hcc[i].Qdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].Qdyl.Cmxtime, Hcc[i].Qdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcc[i].Twidy.Hrs, Hcc[i].Twidy.M,
				Hcc[i].Twidy.Mntime, Hcc[i].Twidy.Mn,
				Hcc[i].Twidy.Mxtime, Hcc[i].Twidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdyt.Hhr, Hcc[i].Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].Qdyt.Chr, Hcc[i].Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].Qdyt.Hmxtime, Hcc[i].Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hcc[i].Qdyt.Cmxtime, Hcc[i].Qdyt.Cmx)
		}
	}
}

func hccmonprt(fo io.Writer, id int, Hcc []HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for i := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", Hcc[i].Name)
		}
	case 1:
		for i := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				Hcc[i].Name, Hcc[i].Name, Hcc[i].Name, Hcc[i].Name)
		}
	default:
		for i := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcc[i].mTaidy.Hrs, Hcc[i].mTaidy.M,
				Hcc[i].mTaidy.Mntime, Hcc[i].mTaidy.Mn,
				Hcc[i].mTaidy.Mxtime, Hcc[i].mTaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdys.Hhr, Hcc[i].mQdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdys.Chr, Hcc[i].mQdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].mQdys.Hmxtime, Hcc[i].mQdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].mQdys.Cmxtime, Hcc[i].mQdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				Hcc[i].mxaidy.Hrs, Hcc[i].mxaidy.M,
				Hcc[i].mxaidy.Mntime, Hcc[i].mxaidy.Mn,
				Hcc[i].mxaidy.Mxtime, Hcc[i].mxaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdyl.Hhr, Hcc[i].mQdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdyl.Chr, Hcc[i].mQdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].mQdyl.Hmxtime, Hcc[i].mQdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].mQdyl.Cmxtime, Hcc[i].mQdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hcc[i].mTwidy.Hrs, Hcc[i].mTwidy.M,
				Hcc[i].mTwidy.Mntime, Hcc[i].mTwidy.Mn,
				Hcc[i].mTwidy.Mxtime, Hcc[i].mTwidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdyt.Hhr, Hcc[i].mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hcc[i].mQdyt.Chr, Hcc[i].mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hcc[i].mQdyt.Hmxtime, Hcc[i].mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hcc[i].mQdyt.Cmxtime, Hcc[i].mQdyt.Cmx)
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
