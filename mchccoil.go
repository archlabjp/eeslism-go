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

package main

import (
	"fmt"
	"io"
	"math"
	"os"
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
func Hccdwint(Nhcc int, _Hcc []HCC) {
	for i := 0; i < Nhcc; i++ {
		Hcc := &_Hcc[i] // Get the address of the current element
		// 乾きコイルと湿りコイルの判定
		if Hcc.Cat.eh > 1.0e-10 {
			Hcc.Wet = 'w'
		} else {
			Hcc.Wet = 'd'
		}

		// 温度効率固定タイプと変動タイプの判定
		if Hcc.Cat.et > 0.0 {
			Hcc.Etype = 'e'
		} else if Hcc.Cat.KA > 0.0 {
			Hcc.Etype = 'k'
		} else {
			fmt.Printf("Hcc %s  Undefined Character et or KA\n", Hcc.Name)
			Hcc.Etype = '-'
		}

		// 入口水温、入口空気絶対湿度を初期化
		//Hcc.Twin = 5.0
		//Hcc.xain = FNXtr(25.0, 50.0)
	}
}

/* ------------------------------------------ */
/*  特性式の係数  */

func Hcccfv(Nhcc int, Hcc []HCC) {
	for i := 0; i < Nhcc; i++ {
		hcc := &Hcc[i] // Get the address of the current element

		hcc.Ga = 0.0
		hcc.Gw = 0.0
		hcc.et = 0.0
		hcc.eh = 0.0

		if hcc.Cmp.Control != OFF_SW {
			var Eo *ELOUT
			var cfin []float64 // Use a slice instead of a pointer
			var AirSW, WaterSW ControlSWType

			Eo = hcc.Cmp.Elouts[0]
			hcc.Ga = Eo.G
			hcc.cGa = Spcheat(Eo.Fluid) * hcc.Ga

			AirSW = OFF_SW
			if hcc.Ga > 0.0 {
				AirSW = ON_SW
			}

			Eo = hcc.Cmp.Elouts[2]
			hcc.Gw = Eo.G
			hcc.cGw = Spcheat(Eo.Fluid) * hcc.Gw

			WaterSW = OFF_SW
			if hcc.Gw > 0.0 {
				WaterSW = ON_SW
			}

			Eo = hcc.Cmp.Elouts[0]

			if hcc.Etype == 'e' {
				hcc.et = hcc.Cat.et
			} else if hcc.Etype == 'k' {
				hcc.et = FNhccet(hcc.cGa, hcc.cGw, hcc.Cat.KA)
			}

			hcc.eh = hcc.Cat.eh

			wcoil(AirSW, WaterSW, hcc.Wet, hcc.Ga*hcc.et, hcc.Ga*hcc.eh, hcc.Xain, hcc.Twin, &hcc.Et, &hcc.Ex, &hcc.Ew)

			Eo.Coeffo = hcc.cGa
			Eo.Co = -(hcc.Et.C)
			cfin = Eo.Coeffin[:]
			cfin[0] = hcc.Et.T - hcc.cGa
			cfin = cfin[1:]
			cfin[0] = hcc.Et.X
			cfin = cfin[1:]
			cfin[0] = -(hcc.Et.W)

			Eo = hcc.Cmp.Elouts[1]
			Eo.Coeffo = hcc.Ga
			Eo.Co = -(hcc.Ex.C)
			cfin = Eo.Coeffin[:]
			cfin[0] = hcc.Ex.T
			cfin = cfin[1:]
			cfin[0] = hcc.Ex.X - hcc.Ga
			cfin = cfin[1:]
			cfin[0] = -(hcc.Ex.W)

			Eo = hcc.Cmp.Elouts[2]
			Eo.Coeffo = hcc.cGw
			Eo.Co = hcc.Ew.C
			cfin = Eo.Coeffin[:]
			cfin[0] = -(hcc.Ew.T)
			cfin = cfin[1:]
			cfin[0] = -(hcc.Ew.X)
			cfin = cfin[1:]
			cfin[0] = hcc.Ew.W - hcc.cGw
		}
	}
}

/* ------------------------------------------ */

/* 供給熱量の計算 */

func Hccdwreset(Nhcc int, Hcc []HCC, DWreset *int) {
	for i := 0; i < Nhcc; i++ {
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
				Hcccfv(1, Hcc)
			}
		}
	}
}

/* ------------------------------------------ */

/* 供給熱量の計算 */

func Hccene(Nhcc int, Hcc []HCC) {
	for i := 0; i < Nhcc; i++ {
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

func hccprint(fo io.Writer, id int, Nhcc int, Hcc []HCC) {
	switch id {
	case 0:
		if Nhcc > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, Nhcc)
		}
		for i := 0; i < Nhcc; i++ {
			fmt.Fprintf(fo, " %s 1 16\n", Hcc[i].Name)
		}
	case 1:
		for i := 0; i < Nhcc; i++ {
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
		for i := 0; i < Nhcc; i++ {
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

func hccdyint(Nhcc int, Hcc []HCC) {
	for i := 0; i < Nhcc; i++ {
		svdyint(&Hcc[i].Taidy)
		svdyint(&Hcc[i].xaidy)
		svdyint(&Hcc[i].Twidy)
		qdyint(&Hcc[i].Qdys)
		qdyint(&Hcc[i].Qdyl)
		qdyint(&Hcc[i].Qdyt)
	}
}

func hccmonint(Nhcc int, Hcc []HCC) {
	for i := 0; i < Nhcc; i++ {
		svdyint(&Hcc[i].mTaidy)
		svdyint(&Hcc[i].mxaidy)
		svdyint(&Hcc[i].mTwidy)
		qdyint(&Hcc[i].mQdys)
		qdyint(&Hcc[i].mQdyl)
		qdyint(&Hcc[i].mQdyt)
	}
}

func hccday(Mon, Day, ttmm, Nhcc int, Hcc []HCC, Nday, SimDayend int) {
	for i := 0; i < Nhcc; i++ {
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

func hccdyprt(fo io.Writer, id, Nhcc int, Hcc []HCC) {
	switch id {
	case 0:
		if Nhcc > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, Nhcc)
		}
		for i := 0; i < Nhcc; i++ {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", Hcc[i].Name)
		}
	case 1:
		for i := 0; i < Nhcc; i++ {
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
		for i := 0; i < Nhcc; i++ {
			fmt.Fprintf(fo, "%1ld %3.1f %1ld %3.1f %1ld %3.1f ",
				Hcc[i].Taidy.Hrs, Hcc[i].Taidy.M,
				Hcc[i].Taidy.Mntime, Hcc[i].Taidy.Mn,
				Hcc[i].Taidy.Mxtime, Hcc[i].Taidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdys.Hhr, Hcc[i].Qdys.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdys.Chr, Hcc[i].Qdys.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].Qdys.Hmxtime, Hcc[i].Qdys.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].Qdys.Cmxtime, Hcc[i].Qdys.Cmx)

			fmt.Fprintf(fo, "%1ld %5.4f %1ld %5.4f %1ld %5.4f ",
				Hcc[i].xaidy.Hrs, Hcc[i].xaidy.M,
				Hcc[i].xaidy.Mntime, Hcc[i].xaidy.Mn,
				Hcc[i].xaidy.Mxtime, Hcc[i].xaidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdyl.Hhr, Hcc[i].Qdyl.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdyl.Chr, Hcc[i].Qdyl.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].Qdyl.Hmxtime, Hcc[i].Qdyl.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].Qdyl.Cmxtime, Hcc[i].Qdyl.Cmx)

			fmt.Fprintf(fo, "%1ld %3.1f %1ld %3.1f %1ld %3.1f ",
				Hcc[i].Twidy.Hrs, Hcc[i].Twidy.M,
				Hcc[i].Twidy.Mntime, Hcc[i].Twidy.Mn,
				Hcc[i].Twidy.Mxtime, Hcc[i].Twidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdyt.Hhr, Hcc[i].Qdyt.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].Qdyt.Chr, Hcc[i].Qdyt.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].Qdyt.Hmxtime, Hcc[i].Qdyt.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f\n", Hcc[i].Qdyt.Cmxtime, Hcc[i].Qdyt.Cmx)
		}
	}
}

func hccmonprt(fo *os.File, id int, Nhcc int, Hcc []HCC) {
	switch id {
	case 0:
		if Nhcc > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, Nhcc)
		}
		for i := 0; i < Nhcc; i++ {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", Hcc[i].Name)
		}
	case 1:
		for i := 0; i < Nhcc; i++ {
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
		for i := 0; i < Nhcc; i++ {
			fmt.Fprintf(fo, "%1ld %3.1f %1ld %3.1f %1ld %3.1f ",
				Hcc[i].mTaidy.Hrs, Hcc[i].mTaidy.M,
				Hcc[i].mTaidy.Mntime, Hcc[i].mTaidy.Mn,
				Hcc[i].mTaidy.Mxtime, Hcc[i].mTaidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdys.Hhr, Hcc[i].mQdys.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdys.Chr, Hcc[i].mQdys.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].mQdys.Hmxtime, Hcc[i].mQdys.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].mQdys.Cmxtime, Hcc[i].mQdys.Cmx)

			fmt.Fprintf(fo, "%1ld %5.4f %1ld %5.4f %1ld %5.4f ",
				Hcc[i].mxaidy.Hrs, Hcc[i].mxaidy.M,
				Hcc[i].mxaidy.Mntime, Hcc[i].mxaidy.Mn,
				Hcc[i].mxaidy.Mxtime, Hcc[i].mxaidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdyl.Hhr, Hcc[i].mQdyl.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdyl.Chr, Hcc[i].mQdyl.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].mQdyl.Hmxtime, Hcc[i].mQdyl.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].mQdyl.Cmxtime, Hcc[i].mQdyl.Cmx)

			fmt.Fprintf(fo, "%1ld %3.1f %1ld %3.1f %1ld %3.1f ",
				Hcc[i].mTwidy.Hrs, Hcc[i].mTwidy.M,
				Hcc[i].mTwidy.Mntime, Hcc[i].mTwidy.Mn,
				Hcc[i].mTwidy.Mxtime, Hcc[i].mTwidy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdyt.Hhr, Hcc[i].mQdyt.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", Hcc[i].mQdyt.Chr, Hcc[i].mQdyt.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", Hcc[i].mQdyt.Hmxtime, Hcc[i].mQdyt.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f\n", Hcc[i].mQdyt.Cmxtime, Hcc[i].mQdyt.Cmx)
		}
	}
}

/* 温水コイルの温度効率計算関数　計算モデルは向流コイル */
func FNhccet(Wa, Ww, KA float64) float64 {
	var NTU, B, Ws, Wl, exB, C float64

	Ws = Wa
	Wl = Ww

	NTU = KA / Ws
	C = Ws / Wl
	B = (1.0 - C) * NTU

	if math.Abs(Ws-Wl) < 1.0e-5 {
		return NTU / (1.0 + NTU)
	} else {
		if exB = math.Exp(-B); math.IsInf(exB, 0) {
			return 1.0 / C
		} else {
			return (1.0 - exB) / (1.0 - C*exB)
		}
	}
}
