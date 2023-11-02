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

/* mcthex.c */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*  全熱交換器  */

/*  仕様入力  */

func Thexdata(s string, Thexca *THEXCA) int {
	var st int
	var dt float64
	var id int

	if st = strings.IndexRune(s, '='); st == -1 {
		Thexca.Name = s
	} else {
		stval := strings.Replace(s[st:], "=", "", 1)
		dt, _ = strconv.ParseFloat(stval, 64)

		if s == "et" {
			Thexca.et = dt
		} else if s == "eh" {
			Thexca.eh = dt
		} else {
			id = 1
		}
	}

	return id
}

/* ------------------------------------------------------ */

func Thexint(Thex []*THEX) {
	for _, thex := range Thex {
		if thex.Cat.eh < 0.0 {
			thex.Type = 't'
			thex.Cat.eh = 0.0
		} else {
			thex.Type = 'h'
		}

		if thex.Cat.et < 0.0 {
			s := fmt.Sprintf("Name=%s catname=%s et=%f", thex.Name, thex.Cat.Name, thex.Cat.et)
			Eprint("<Thexint>", s)
		}

		thex.Xeinold = FNXtr(26.0, 50.0)
		thex.Xeoutold = thex.Xeinold
		thex.Xoinold = thex.Xeinold
		thex.Xooutold = thex.Xeinold
	}
}

/*  全熱交換器出口空気温湿度に関する変数割当  */
func Thexelm(Thex []*THEX) {
	var E, E1, E2, E3 *ELOUT
	var elin, elin2 *ELIN

	for _, thex := range Thex {
		E = thex.Cmp.Elouts[0]
		E1 = thex.Cmp.Elouts[1]
		E2 = thex.Cmp.Elouts[2]
		E3 = thex.Cmp.Elouts[3]

		// Tein variable assignment
		// E: Teout calculation, elin2: Tein
		elin2 = E.Elins[0]

		// E+2: Toout calculation, elin: Tein
		elin = E2.Elins[1]
		elin.Upo = elin2.Upo
		elin.Upv = elin2.Upv

		if thex.Cat.eh > 0.0 {
			// E+1: xeout calculation, elin:
			elin = E1.Elins[1]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upv

			elin = E3.Elins[3]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upv
		}

		// Toin variable assignment
		elin2 = E.Elins[1]

		elin = E2.Elins[0]
		elin.Upo = elin2.Upo
		elin.Upv = elin2.Upv

		if thex.Cat.eh > 0.0 {
			elin = E1.Elins[3]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upv

			elin = E3.Elins[1]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upv

			// Teout variable assignment
			elin = E1.Elins[2]
			elin.Upo = E
			elin.Upv = E2

			// Toout assignment
			elin = E3.Elins[2]
			elin.Upo = E2
			elin.Upv = E2

			// xein assignment
			elin = E1.Elins[0]
			elin2 = E3.Elins[4]
			elin2.Upo = elin.Upo
			elin2.Upv = elin.Upv

			// xoin assignment
			elin = E1.Elins[4]
			elin2 = E3.Elins[0]
			elin2.Upo = elin.Upo
			elin2.Upv = elin.Upv
		}
	}
}

/* ------------------------------------------------------ */

//
//  [IN 1] --(E)-->  +------+ --(E)--> [OUT 1] 排気系統（温度）
//  [IN 2] --(e)-->  |      | --(e)--> [OUT 2] 排気系統（エンタルピー）
//                   | THEX |
//  [IN 3] --(O)-->  |      | --(O)--> [OUT 3] 給気系統（温度）
//  [IN 4] --(o)-->  +------+ --(o)--> [OUT 4] 給気系統（エンタルピー）
//
func Thexcfv(Thex []*THEX) {
	var Eoet, Eoot, Eoex, Eoox *ELOUT
	var etCGmin, ehGmin, Aeout, Aein, Aoout, Aoin float64

	for _, thex := range Thex {
		if thex.Cmp.Control != OFF_SW {
			thex.ET = thex.Cat.et
			thex.EH = thex.Cat.eh

			Eoet = thex.Cmp.Elouts[0] // 排気系統（温度）
			Eoex = thex.Cmp.Elouts[1] // 排気系統（エンタルピー）
			Eoot = thex.Cmp.Elouts[2] // 給気系統（温度）
			Eoox = thex.Cmp.Elouts[3] // 給気系統（エンタルピー）

			thex.Ge = Eoet.G
			thex.Go = Eoot.G

			if DEBUG {
				fmt.Printf("<Thexcfv>  %s Ge=%f Go=%f\n", thex.Cmp.Name, thex.Ge, thex.Go)
			}

			thex.CGe = Spcheat(Eoet.Fluid) * thex.Ge
			thex.CGo = Spcheat(Eoot.Fluid) * thex.Go
			etCGmin = thex.ET * math.Min(thex.CGe, thex.CGo)
			ehGmin = thex.EH * math.Min(thex.Ge, thex.Go)

			Aein = Ca + Cv*thex.Xeinold
			Aeout = Ca + Cv*thex.Xeoutold
			Aoin = Ca + Cv*thex.Xoinold
			Aoout = Ca + Cv*thex.Xooutold

			// 排気系統（温度）の熱収支
			Eoet.Coeffo = thex.CGe
			Eoet.Co = 0.0
			cfin := Eoet.Coeffin
			cfin[0] = etCGmin - thex.CGe
			cfin[1] = -etCGmin

			// 給気系統（温度）の熱収支
			Eoot.Coeffo = thex.CGo
			Eoot.Co = 0.0
			cfin = Eoot.Coeffin
			cfin[0] = etCGmin - thex.CGo
			cfin[1] = -etCGmin

			if thex.Type == 'h' {
				// 排気系統（エンタルピー）の熱収支
				Eoex.Coeffo = thex.Ge * Ro
				Eoex.Co = 0.0
				cfin = Eoex.Coeffin
				cfin[0] = Ro * (ehGmin - thex.Ge)
				cfin[1] = Aein * (ehGmin - thex.Ge)
				cfin[2] = Aeout * thex.Ge
				cfin[3] = -ehGmin * Aoin
				cfin[4] = -ehGmin * Ro

				// 給気系統（エンタルピー）の熱収支
				Eoox.Coeffo = thex.Go * Ro
				Eoox.Co = 0.0
				cfin = Eoox.Coeffin
				cfin[0] = Ro * (ehGmin - thex.Go)
				cfin[1] = Aoin * (ehGmin - thex.Go)
				cfin[2] = thex.Go * Aoout
				cfin[3] = -ehGmin * Aein
				cfin[4] = -ehGmin * Ro
			} else {
				Eoex.Coeffo = 1.0
				Eoex.Coeffin[0] = -1.0

				Eoox.Coeffo = 1.0
				Eoox.Coeffin[0] = -1.0
			}
		}
	}
}

func Thexene(Thex []*THEX) {
	for _, thex := range Thex {
		Eoet := thex.Cmp.Elouts[0] // 排気系統（温度）
		Eoex := thex.Cmp.Elouts[1] // 排気系統（エンタルピー）
		Eoot := thex.Cmp.Elouts[2] // 給気系統（温度）
		Eoox := thex.Cmp.Elouts[3] // 給気系統（エンタルピー）

		thex.Tein = Eoet.Elins[0].Upo.Sysv
		thex.Teout = Eoet.Sysv
		thex.Xein = Eoex.Elins[0].Upo.Sysv
		thex.Xeout = Eoex.Sysv

		thex.Toin = Eoot.Elins[0].Upo.Sysv
		thex.Toout = Eoot.Sysv
		thex.Xoin = Eoox.Elins[0].Upo.Sysv
		thex.Xoout = Eoox.Sysv

		thex.Hein = FNH(thex.Tein, thex.Xein)
		thex.Heout = FNH(thex.Teout, thex.Xeout)
		thex.Hoin = FNH(thex.Toin, thex.Xoin)
		thex.Hoout = FNH(thex.Toout, thex.Xoout)

		if thex.Cmp.Control != OFF_SW {
			// 交換熱量の計算
			thex.Qes = Ca * thex.Ge * (thex.Teout - thex.Tein)
			thex.Qel = Ro * thex.Ge * (thex.Xeout - thex.Xein)
			thex.Qet = thex.Qes + thex.Qel

			thex.Qos = Ca * thex.Go * (thex.Toout - thex.Toin)
			thex.Qol = Ro * thex.Go * (thex.Xoout - thex.Xoin)
			thex.Qot = thex.Qos + thex.Qol

			// 前時刻の絶対湿度の入れ替え
			thex.Xeinold = thex.Xein
			thex.Xeoutold = thex.Xeout
			thex.Xoinold = thex.Xoin
			thex.Xooutold = thex.Xoout
		} else {
			thex.Qes = 0.0
			thex.Qel = 0.0
			thex.Qet = 0.0
			thex.Qos = 0.0
			thex.Qol = 0.0
			thex.Qot = 0.0
			thex.Ge = 0.0
			thex.Tein = 0.0
			thex.Teout = 0.0
			thex.Xein = 0.0
			thex.Xeout = 0.0
			thex.Hein = 0.0
			thex.Heout = 0.0
			thex.Go = 0.0
			thex.Toin = 0.0
			thex.Toout = 0.0
			thex.Xoin = 0.0
			thex.Xoout = 0.0
			thex.Hoin = 0.0
			thex.Hoout = 0.0
		}
	}
}

func Thexprint(fo io.Writer, id int, Thex []*THEX) {
	var el *ELOUT

	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 22\n", thex.Name)
		}

	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_ce c c %s_Ge m f %s_Tei t f %s_Teo t f %s_xei t f %s_xeo t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_hei h f %s_heo h f %s_Qes q f %s_Qel q f %s_Qet q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_co c c %s_Go m f %s_Toi t f %s_Too t f %s_xoi t f %s_xoo t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_hoi h f %s_hoo h f %s_Qos q f %s_Qol q f %s_Qot q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
		}

	default:
		for _, thex := range Thex {
			el = thex.Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, thex.Ge, thex.Tein, thex.Teout, thex.Xein, thex.Xeout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				thex.Hein, thex.Heout, thex.Qes, thex.Qel, thex.Qet)

			el = thex.Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, thex.Go, thex.Toin, thex.Toout, thex.Xoin, thex.Xoout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				thex.Hoin, thex.Hoout, thex.Qos, thex.Qol, thex.Qot)
		}
	}
}

func Thexdyint(Thex []*THEX) {
	for _, thex := range Thex {
		svdyint(&thex.Teidy)
		svdyint(&thex.Teody)
		svdyint(&thex.Xeidy)
		svdyint(&thex.Xeody)

		svdyint(&thex.Toidy)
		svdyint(&thex.Toody)
		svdyint(&thex.Xoidy)
		svdyint(&thex.Xoody)

		qdyint(&thex.Qdyes)
		qdyint(&thex.Qdyel)
		qdyint(&thex.Qdyet)

		qdyint(&thex.Qdyos)
		qdyint(&thex.Qdyol)
		qdyint(&thex.Qdyot)
	}
}

func Thexmonint(Thex []*THEX) {
	for _, thex := range Thex {
		svdyint(&thex.MTeidy)
		svdyint(&thex.MTeody)
		svdyint(&thex.MXeidy)
		svdyint(&thex.MXeody)

		svdyint(&thex.MToidy)
		svdyint(&thex.MToody)
		svdyint(&thex.MXoidy)
		svdyint(&thex.MXoody)

		qdyint(&thex.MQdyes)
		qdyint(&thex.MQdyel)
		qdyint(&thex.MQdyet)

		qdyint(&thex.MQdyos)
		qdyint(&thex.MQdyol)
		qdyint(&thex.MQdyot)
	}
}

func Thexday(Mon, Day, ttmm int, Thex []*THEX, Nday, SimDayend int) {
	for _, thex := range Thex {
		// 日集計
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Tein, &thex.Teidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Teout, &thex.Teody)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xein, &thex.Xeidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xeout, &thex.Xeody)

		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Toin, &thex.Toidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Toout, &thex.Toody)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xoin, &thex.Xoidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xoout, &thex.Xoody)

		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qes, &thex.Qdyes)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qel, &thex.Qdyel)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qet, &thex.Qdyet)

		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qos, &thex.Qdyos)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qol, &thex.Qdyol)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qot, &thex.Qdyot)

		// 月集計
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Tein, &thex.MTeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Teout, &thex.MTeody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xein, &thex.MXeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xeout, &thex.MXeody, Nday, SimDayend)

		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Toin, &thex.MToidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Toout, &thex.MToody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xoin, &thex.MXoidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xoout, &thex.MXoody, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qes, &thex.MQdyes, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qel, &thex.MQdyel, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qet, &thex.MQdyet, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qos, &thex.MQdyos, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qol, &thex.MQdyol, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qot, &thex.MQdyot, Nday, SimDayend)
	}
}

func Thexdyprt(fo io.Writer, id int, Thex []*THEX) {
	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 48\n", thex.Name)
		}
	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
		}
	default:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.Teidy.Hrs, thex.Teidy.M,
				thex.Teidy.Mntime, thex.Teidy.Mn,
				thex.Teidy.Mxtime, thex.Teidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.Toidy.Hrs, thex.Toidy.M,
				thex.Toidy.Mntime, thex.Toidy.Mn,
				thex.Toidy.Mxtime, thex.Toidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.Xeidy.Hrs, thex.Xeidy.M,
				thex.Xeidy.Mntime, thex.Xeidy.Mn,
				thex.Xeidy.Mxtime, thex.Xeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.Xoidy.Hrs, thex.Xoidy.M,
				thex.Xoidy.Mntime, thex.Xoidy.Mn,
				thex.Xoidy.Mxtime, thex.Xoidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyes.Hhr, thex.Qdyes.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyes.Chr, thex.Qdyes.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyes.Hmxtime, thex.Qdyes.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyes.Cmxtime, thex.Qdyes.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyel.Hhr, thex.Qdyel.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyel.Chr, thex.Qdyel.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyel.Hmxtime, thex.Qdyel.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyel.Cmxtime, thex.Qdyel.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyet.Hhr, thex.Qdyet.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyet.Chr, thex.Qdyet.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyet.Hmxtime, thex.Qdyet.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyet.Cmxtime, thex.Qdyet.Cmx)
		}
	}
}
func Thexmonprt(fo io.Writer, id int, Thex []*THEX) {
	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 48\n", thex.Name)
		}
	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
		}
	default:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.MTeidy.Hrs, thex.MTeidy.M,
				thex.MTeidy.Mntime, thex.MTeidy.Mn,
				thex.MTeidy.Mxtime, thex.MTeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.MToidy.Hrs, thex.MToidy.M,
				thex.MToidy.Mntime, thex.MToidy.Mn,
				thex.MToidy.Mxtime, thex.MToidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.MXeidy.Hrs, thex.MXeidy.M,
				thex.MXeidy.Mntime, thex.MXeidy.Mn,
				thex.MXeidy.Mxtime, thex.MXeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.MXoidy.Hrs, thex.MXoidy.M,
				thex.MXoidy.Mntime, thex.MXoidy.Mn,
				thex.MXoidy.Mxtime, thex.MXoidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyes.Hhr, thex.MQdyes.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyes.Chr, thex.MQdyes.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyes.Hmxtime, thex.MQdyes.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyes.Cmxtime, thex.MQdyes.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyel.Hhr, thex.MQdyel.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyel.Chr, thex.MQdyel.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyel.Hmxtime, thex.MQdyel.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyel.Cmxtime, thex.MQdyel.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyet.Hhr, thex.MQdyet.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyet.Chr, thex.MQdyet.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyet.Hmxtime, thex.MQdyet.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyet.Cmxtime, thex.MQdyet.Cmx)
		}
	}
}
