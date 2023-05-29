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

package main

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

func Thexint(Nthex int, Thex []THEX) {
	var i int
	var s string

	for i = 0; i < Nthex; i++ {
		if Thex[i].Cat.eh < 0.0 {
			Thex[i].Type = 't'
			Thex[i].Cat.eh = 0.0
		} else {
			Thex[i].Type = 'h'
		}

		if Thex[i].Cat.et < 0.0 {
			s = fmt.Sprintf("Name=%s catname=%s et=%f", Thex[i].Name, Thex[i].Cat.Name, Thex[i].Cat.et)
			Eprint("<Thexint>", s)
		}

		Thex[i].Xeinold = FNXtr(26.0, 50.0)
		Thex[i].Xeoutold = Thex[i].Xeinold
		Thex[i].Xoinold = Thex[i].Xeinold
		Thex[i].Xooutold = Thex[i].Xeinold
	}
}

/*  全熱交換器出口空気温湿度に関する変数割当  */
func Thexelm(NThex int, Thex []THEX) {
	var i int
	var E, E1, E2, E3 *ELOUT
	var elin, elin2 *ELIN

	for i = 0; i < NThex; i++ {
		E = Thex[i].Cmp.Elouts[0]
		E1 = Thex[i].Cmp.Elouts[1]
		E2 = Thex[i].Cmp.Elouts[2]
		E3 = Thex[i].Cmp.Elouts[3]

		// Tein variable assignment
		// E: Teout calculation, elin2: Tein
		elin2 = E.Elins[0]

		// E+2: Toout calculation, elin: Tein
		elin = E2.Elins[1]
		elin.Upo = elin2.Upo
		elin.Upv = elin2.Upv

		if Thex[i].Cat.eh > 0.0 {
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

		if Thex[i].Cat.eh > 0.0 {
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

func Thexcfv(Nthex int, Thex []THEX) {
	var Eoet, Eoot, Eoex, Eoox *ELOUT
	var etCGmin, ehGmin, Aeout, Aein, Aoout, Aoin float64
	var i int

	for i = 0; i < Nthex; i++ {
		if Thex[i].Cmp.Control != OFF_SW {
			Thex[i].ET = Thex[i].Cat.et
			Thex[i].EH = Thex[i].Cat.eh

			Eoet = Thex[i].Cmp.Elouts[0] // 排気系統（温度）
			Eoex = Thex[i].Cmp.Elouts[1] // 排気系統（エンタルピー）
			Eoot = Thex[i].Cmp.Elouts[2] // 給気系統（温度）
			Eoox = Thex[i].Cmp.Elouts[3] // 給気系統（エンタルピー）

			Thex[i].Ge = Eoet.G
			Thex[i].Go = Eoot.G

			if DEBUG {
				fmt.Printf("<Thexcfv>  %s Ge=%f Go=%f\n", Thex[i].Cmp.Name, Thex[i].Ge, Thex[i].Go)
			}

			Thex[i].CGe = Spcheat(Eoet.Fluid) * Thex[i].Ge
			Thex[i].CGo = Spcheat(Eoot.Fluid) * Thex[i].Go
			etCGmin = Thex[i].ET * math.Min(Thex[i].CGe, Thex[i].CGo)
			ehGmin = Thex[i].EH * math.Min(Thex[i].Ge, Thex[i].Go)

			Aein = Ca + Cv*Thex[i].Xeinold
			Aeout = Ca + Cv*Thex[i].Xeoutold
			Aoin = Ca + Cv*Thex[i].Xoinold
			Aoout = Ca + Cv*Thex[i].Xooutold

			// 排気系統（温度）の熱収支
			Eoet.Coeffo = Thex[i].CGe
			Eoet.Co = 0.0
			cfin := Eoet.Coeffin
			cfin[0] = etCGmin - Thex[i].CGe
			cfin[1] = -etCGmin

			// 給気系統（温度）の熱収支
			Eoot.Coeffo = Thex[i].CGo
			Eoot.Co = 0.0
			cfin = Eoot.Coeffin
			cfin[0] = etCGmin - Thex[i].CGo
			cfin[1] = -etCGmin

			if Thex[i].Type == 'h' {
				// 排気系統（エンタルピー）の熱収支
				Eoex.Coeffo = Thex[i].Ge * Ro
				Eoex.Co = 0.0
				cfin = Eoex.Coeffin
				cfin[0] = Ro * (ehGmin - Thex[i].Ge)
				cfin[1] = Aein * (ehGmin - Thex[i].Ge)
				cfin[2] = Aeout * Thex[i].Ge
				cfin[3] = -ehGmin * Aoin
				cfin[4] = -ehGmin * Ro

				// 給気系統（エンタルピー）の熱収支
				Eoox.Coeffo = Thex[i].Go * Ro
				Eoox.Co = 0.0
				cfin = Eoox.Coeffin
				cfin[0] = Ro * (ehGmin - Thex[i].Go)
				cfin[1] = Aoin * (ehGmin - Thex[i].Go)
				cfin[2] = Thex[i].Go * Aoout
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

func Thexene(Nthex int, Thex []THEX) {
	for i := 0; i < Nthex; i++ {
		Eoet := Thex[i].Cmp.Elouts[0] // 排気系統（温度）
		Eoex := Thex[i].Cmp.Elouts[1] // 排気系統（エンタルピー）
		Eoot := Thex[i].Cmp.Elouts[2] // 給気系統（温度）
		Eoox := Thex[i].Cmp.Elouts[3] // 給気系統（エンタルピー）

		Thex[i].Tein = Eoet.Elins[0].Upo.Sysv
		Thex[i].Teout = Eoet.Sysv
		Thex[i].Xein = Eoex.Elins[0].Upo.Sysv
		Thex[i].Xeout = Eoex.Sysv

		Thex[i].Toin = Eoot.Elins[0].Upo.Sysv
		Thex[i].Toout = Eoot.Sysv
		Thex[i].Xoin = Eoox.Elins[0].Upo.Sysv
		Thex[i].Xoout = Eoox.Sysv

		Thex[i].Hein = FNH(Thex[i].Tein, Thex[i].Xein)
		Thex[i].Heout = FNH(Thex[i].Teout, Thex[i].Xeout)
		Thex[i].Hoin = FNH(Thex[i].Toin, Thex[i].Xoin)
		Thex[i].Hoout = FNH(Thex[i].Toout, Thex[i].Xoout)

		if Thex[i].Cmp.Control != OFF_SW {
			// 交換熱量の計算
			Thex[i].Qes = Ca * Thex[i].Ge * (Thex[i].Teout - Thex[i].Tein)
			Thex[i].Qel = Ro * Thex[i].Ge * (Thex[i].Xeout - Thex[i].Xein)
			Thex[i].Qet = Thex[i].Qes + Thex[i].Qel

			Thex[i].Qos = Ca * Thex[i].Go * (Thex[i].Toout - Thex[i].Toin)
			Thex[i].Qol = Ro * Thex[i].Go * (Thex[i].Xoout - Thex[i].Xoin)
			Thex[i].Qot = Thex[i].Qos + Thex[i].Qol

			// 前時刻の絶対湿度の入れ替え
			Thex[i].Xeinold = Thex[i].Xein
			Thex[i].Xeoutold = Thex[i].Xeout
			Thex[i].Xoinold = Thex[i].Xoin
			Thex[i].Xooutold = Thex[i].Xoout
		} else {
			Thex[i].Qes = 0.0
			Thex[i].Qel = 0.0
			Thex[i].Qet = 0.0
			Thex[i].Qos = 0.0
			Thex[i].Qol = 0.0
			Thex[i].Qot = 0.0
			Thex[i].Ge = 0.0
			Thex[i].Tein = 0.0
			Thex[i].Teout = 0.0
			Thex[i].Xein = 0.0
			Thex[i].Xeout = 0.0
			Thex[i].Hein = 0.0
			Thex[i].Heout = 0.0
			Thex[i].Go = 0.0
			Thex[i].Toin = 0.0
			Thex[i].Toout = 0.0
			Thex[i].Xoin = 0.0
			Thex[i].Xoout = 0.0
			Thex[i].Hoin = 0.0
			Thex[i].Hoout = 0.0
		}
	}
}

func Thexprint(fo io.Writer, id, Nthex int, Thex []THEX) {
	var i int
	var el *ELOUT

	switch id {
	case 0:
		if Nthex > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, Nthex)
		}
		for i = 0; i < Nthex; i++ {
			fmt.Fprintf(fo, " %s 1 22\n", Thex[i].Name)
		}

	case 1:
		for i = 0; i < Nthex; i++ {
			fmt.Fprintf(fo, "%s_ce c c %s_Ge m f %s_Tei t f %s_Teo t f %s_xei t f %s_xeo t f\n",
				Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
			fmt.Fprintf(fo, "%s_hei h f %s_heo h f %s_Qes q f %s_Qel q f %s_Qet q f\n",
				Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)

			fmt.Fprintf(fo, "%s_co c c %s_Go m f %s_Toi t f %s_Too t f %s_xoi t f %s_xoo t f\n",
				Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
			fmt.Fprintf(fo, "%s_hoi h f %s_hoo h f %s_Qos q f %s_Qol q f %s_Qot q f\n",
				Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
		}

	default:
		for i = 0; i < Nthex; i++ {
			el = Thex[i].Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, Thex[i].Ge, Thex[i].Tein, Thex[i].Teout, Thex[i].Xein, Thex[i].Xeout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				Thex[i].Hein, Thex[i].Heout, Thex[i].Qes, Thex[i].Qel, Thex[i].Qet)

			el = Thex[i].Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, Thex[i].Go, Thex[i].Toin, Thex[i].Toout, Thex[i].Xoin, Thex[i].Xoout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				Thex[i].Hoin, Thex[i].Hoout, Thex[i].Qos, Thex[i].Qol, Thex[i].Qot)
		}
	}
}

func Thexdyint(Nthex int, Thex []THEX) {
	var i int

	for i = 0; i < Nthex; i++ {
		svdyint(&Thex[i].Teidy)
		svdyint(&Thex[i].Teody)
		svdyint(&Thex[i].Xeidy)
		svdyint(&Thex[i].Xeody)

		svdyint(&Thex[i].Toidy)
		svdyint(&Thex[i].Toody)
		svdyint(&Thex[i].Xoidy)
		svdyint(&Thex[i].Xoody)

		qdyint(&Thex[i].Qdyes)
		qdyint(&Thex[i].Qdyel)
		qdyint(&Thex[i].Qdyet)

		qdyint(&Thex[i].Qdyos)
		qdyint(&Thex[i].Qdyol)
		qdyint(&Thex[i].Qdyot)
	}
}

func Thexmonint(Nthex int, Thex []THEX) {
	var i int

	for i = 0; i < Nthex; i++ {
		svdyint(&Thex[i].MTeidy)
		svdyint(&Thex[i].MTeody)
		svdyint(&Thex[i].MXeidy)
		svdyint(&Thex[i].MXeody)

		svdyint(&Thex[i].MToidy)
		svdyint(&Thex[i].MToody)
		svdyint(&Thex[i].MXoidy)
		svdyint(&Thex[i].MXoody)

		qdyint(&Thex[i].MQdyes)
		qdyint(&Thex[i].MQdyel)
		qdyint(&Thex[i].MQdyet)

		qdyint(&Thex[i].MQdyos)
		qdyint(&Thex[i].MQdyol)
		qdyint(&Thex[i].MQdyot)
	}
}

func Thexday(Mon, Day, ttmm, Nthex int, Thex []THEX, Nday, SimDayend int) {
	var i int

	for i = 0; i < Nthex; i++ {
		// 日集計
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Tein, &Thex[i].Teidy)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Teout, &Thex[i].Teody)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Xein, &Thex[i].Xeidy)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Xeout, &Thex[i].Xeody)

		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Toin, &Thex[i].Toidy)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Toout, &Thex[i].Toody)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Xoin, &Thex[i].Xoidy)
		svdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Xoout, &Thex[i].Xoody)

		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qes, &Thex[i].Qdyes)
		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qel, &Thex[i].Qdyel)
		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qet, &Thex[i].Qdyet)

		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qos, &Thex[i].Qdyos)
		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qol, &Thex[i].Qdyol)
		qdaysum(int64(ttmm), Thex[i].Cmp.Control, Thex[i].Qot, &Thex[i].Qdyot)

		// 月集計
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Tein, &Thex[i].MTeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Teout, &Thex[i].MTeody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Xein, &Thex[i].MXeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Xeout, &Thex[i].MXeody, Nday, SimDayend)

		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Toin, &Thex[i].MToidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Toout, &Thex[i].MToody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Xoin, &Thex[i].MXoidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Xoout, &Thex[i].MXoody, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qes, &Thex[i].MQdyes, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qel, &Thex[i].MQdyel, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qet, &Thex[i].MQdyet, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qos, &Thex[i].MQdyos, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qol, &Thex[i].MQdyol, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Thex[i].Cmp.Control, Thex[i].Qot, &Thex[i].MQdyot, Nday, SimDayend)
	}
}

func Thexdyprt(fo io.Writer, id, Nthex int, Thex []THEX) {
	for i := 0; i < Nthex; i++ {
		switch id {
		case 0:
			if Nthex > 0 {
				fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, Nthex)
			}
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, " %s 1 48\n", Thex[i].Name)
			}
		case 1:
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)

				fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)

				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
			}
		default:
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					Thex[i].Teidy.Hrs, Thex[i].Teidy.M,
					Thex[i].Teidy.Mntime, Thex[i].Teidy.Mn,
					Thex[i].Teidy.Mxtime, Thex[i].Teidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
					Thex[i].Toidy.Hrs, Thex[i].Toidy.M,
					Thex[i].Toidy.Mntime, Thex[i].Toidy.Mn,
					Thex[i].Toidy.Mxtime, Thex[i].Toidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					Thex[i].Xeidy.Hrs, Thex[i].Xeidy.M,
					Thex[i].Xeidy.Mntime, Thex[i].Xeidy.Mn,
					Thex[i].Xeidy.Mxtime, Thex[i].Xeidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
					Thex[i].Xoidy.Hrs, Thex[i].Xoidy.M,
					Thex[i].Xoidy.Mntime, Thex[i].Xoidy.Mn,
					Thex[i].Xoidy.Mxtime, Thex[i].Xoidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyes.Hhr, Thex[i].Qdyes.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyes.Chr, Thex[i].Qdyes.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].Qdyes.Hmxtime, Thex[i].Qdyes.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].Qdyes.Cmxtime, Thex[i].Qdyes.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyel.Hhr, Thex[i].Qdyel.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyel.Chr, Thex[i].Qdyel.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].Qdyel.Hmxtime, Thex[i].Qdyel.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].Qdyel.Cmxtime, Thex[i].Qdyel.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyet.Hhr, Thex[i].Qdyet.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].Qdyet.Chr, Thex[i].Qdyet.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].Qdyet.Hmxtime, Thex[i].Qdyet.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].Qdyet.Cmxtime, Thex[i].Qdyet.Cmx)
			}
		}
	}
}
func Thexmonprt(fo io.Writer, id, Nthex int, Thex []THEX) {
	for i := 0; i < Nthex; i++ {
		switch id {
		case 0:
			if Nthex > 0 {
				fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, Nthex)
			}
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, " %s 1 48\n", Thex[i].Name)
			}
		case 1:
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)

				fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)

				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
				fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
					Thex[i].Name, Thex[i].Name, Thex[i].Name, Thex[i].Name)
			}
		default:
			for i := 0; i < Nthex; i++ {
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					Thex[i].MTeidy.Hrs, Thex[i].MTeidy.M,
					Thex[i].MTeidy.Mntime, Thex[i].MTeidy.Mn,
					Thex[i].MTeidy.Mxtime, Thex[i].MTeidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
					Thex[i].MToidy.Hrs, Thex[i].MToidy.M,
					Thex[i].MToidy.Mntime, Thex[i].MToidy.Mn,
					Thex[i].MToidy.Mxtime, Thex[i].MToidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					Thex[i].MXeidy.Hrs, Thex[i].MXeidy.M,
					Thex[i].MXeidy.Mntime, Thex[i].MXeidy.Mn,
					Thex[i].MXeidy.Mxtime, Thex[i].MXeidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
					Thex[i].MXoidy.Hrs, Thex[i].MXoidy.M,
					Thex[i].MXoidy.Mntime, Thex[i].MXoidy.Mn,
					Thex[i].MXoidy.Mxtime, Thex[i].MXoidy.Mx)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyes.Hhr, Thex[i].MQdyes.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyes.Chr, Thex[i].MQdyes.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].MQdyes.Hmxtime, Thex[i].MQdyes.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].MQdyes.Cmxtime, Thex[i].MQdyes.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyel.Hhr, Thex[i].MQdyel.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyel.Chr, Thex[i].MQdyel.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].MQdyel.Hmxtime, Thex[i].MQdyel.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].MQdyel.Cmxtime, Thex[i].MQdyel.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyet.Hhr, Thex[i].MQdyet.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Thex[i].MQdyet.Chr, Thex[i].MQdyet.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Thex[i].MQdyet.Hmxtime, Thex[i].MQdyet.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Thex[i].MQdyet.Cmxtime, Thex[i].MQdyet.Cmx)
			}
		}
	}
}
