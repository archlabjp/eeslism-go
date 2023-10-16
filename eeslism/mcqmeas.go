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

/*  mcqmeas.c  */

/*  QMEAS */

package eeslism

import (
	"fmt"
	"io"
)

// 出入り口温湿度の割り当て

func Qmeaselm(Qmeas []QMEAS) {
	for i := range Qmeas {
		qmeas := &Qmeas[i]
		PlistTh := qmeas.PlistTh
		N := qmeas.Nelmh
		Pelm := PlistTh.Pelm[N]
		if N > 0 && (Pelm.Cmp.Eqptype != DIVERG_TYPE) && (Pelm.Cmp.Eqptype != DIVGAIR_TYPE) {
			Pelm = PlistTh.Pelm[N-1]
			qmeas.Th = &Pelm.Out.Sysv

			if PlistTh.Plistx != nil {
				qmeas.Plistxh = PlistTh.Plistx
				Plist := qmeas.Plistxh
				Pelm = Plist.Pelm[N-1]
				qmeas.Xh = &Pelm.Out.Sysv
			}
		} else {
			Pelm = PlistTh.Pelm[1]
			var sysvin *float64

			if Pelm.Out != nil {
				sysvin = &Pelm.Out.Elins[0].Sysvin
			} else {
				sysvin = &Pelm.In.Sysvin
			}

			qmeas.Th = sysvin

			if PlistTh.Plistx != nil {
				qmeas.Plistxh = PlistTh.Plistx
				Pelm = qmeas.Plistxh.Pelm[1]

				if Pelm.Out != nil {
					qmeas.Xh = &Pelm.Out.Elins[0].Sysvin
				} else {
					qmeas.Xh = &Pelm.In.Sysvin
				}
			}
		}

		PlistTc := qmeas.PlistTc
		N = qmeas.Nelmc
		Pelm = PlistTc.Pelm[N]
		if N > 0 && (Pelm.Cmp.Eqptype != DIVERG_TYPE) && (Pelm.Cmp.Eqptype != DIVGAIR_TYPE) {
			Pelm = PlistTc.Pelm[N-1]
			qmeas.Tc = &Pelm.Out.Sysv

			if PlistTc.Plistx != nil {
				qmeas.Plistxc = PlistTc.Plistx
				Plist := qmeas.Plistxc
				Pelm = Plist.Pelm[N-1]
				qmeas.Xc = &Pelm.Out.Sysv
			}
		} else {
			Pelm = PlistTc.Pelm[1]
			var sysvin *float64

			if Pelm.Out != nil {
				sysvin = &Pelm.Out.Elins[0].Sysvin
			} else {
				sysvin = &Pelm.In.Sysvin
			}

			qmeas.Tc = sysvin

			if PlistTc.Plistx != nil {
				qmeas.Plistxc = PlistTc.Plistx
				Pelm = qmeas.Plistxc.Pelm[1]

				if Pelm.Out != nil {
					qmeas.Xc = &Pelm.Out.Elins[0].Sysvin
				} else {
					qmeas.Xc = &Pelm.In.Sysvin
				}
			}
		}
	}
}

func Qmeasene(Qmeas []QMEAS) {
	for i := range Qmeas {
		qmeas := &Qmeas[i]
		PG := qmeas.PlistG
		Ph := qmeas.PlistTh
		Pc := qmeas.PlistTc

		if PG.Control != OFF_SW && Ph.Control != OFF_SW && Pc.Control != OFF_SW {
			qmeas.Qs = Spcheat(PG.Mpath.Fluid) * *qmeas.G * (*qmeas.Th - *qmeas.Tc)

			if qmeas.Plistxc != nil {
				qmeas.Ql = Ro * *qmeas.G * (*qmeas.Xh - *qmeas.Xc)
			} else {
				qmeas.Ql = 0.0
			}

			qmeas.Qt = qmeas.Qs + qmeas.Ql
		} else {
			qmeas.Qs = 0.0
			qmeas.Ql = 0.0
			qmeas.Qt = 0.0
		}
	}
}

func Qmeasprint(fo io.Writer, id int, Qmeas []QMEAS) {
	for i := range Qmeas {
		qmeas := &Qmeas[i]
		el := qmeas.Cmp.Elouts[0]

		switch id {
		case 0:
			if len(Qmeas) > 0 {
				fmt.Fprintf(fo, "%s %d\n", QMEAS_TYPE, len(Qmeas))
			}
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, " %s 1 9\n", qmeas.Name)
			} else {
				fmt.Fprintf(fo, " %s 1 5\n", qmeas.Name)
			}
		case 1:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%s_ce c c %s_G m f %s_Th t f %s_Tc t f %s_xh t f %s_xc t f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_Qs q f %s_Ql q f %s_Qt q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name)
			} else {
				fmt.Fprintf(fo, "%s_ce c c %s_G m f %s_Th t f %s_Tc t f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_Qt q f\n", qmeas.Name)
			}
		default:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %.3f %.3f ",
					el.Control, *qmeas.G, *qmeas.Th, *qmeas.Tc, *qmeas.Xh, *qmeas.Xc)
				fmt.Fprintf(fo, "%.0f %.0f %.0f\n",
					qmeas.Qs, qmeas.Ql, qmeas.Qt)
			} else {
				fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f ",
					el.Control, *qmeas.G, *qmeas.Th, *qmeas.Tc)
				fmt.Fprintf(fo, "%.0f\n", qmeas.Qt)
			}
		}
	}
}

func Qmeasdyint(Qmeas []QMEAS) {
	for i := range Qmeas {
		svdyint(&Qmeas[i].Tcdy)
		svdyint(&Qmeas[i].Thdy)
		svdyint(&Qmeas[i].xcdy)
		svdyint(&Qmeas[i].xhdy)

		qdyint(&Qmeas[i].Qdys)
		qdyint(&Qmeas[i].Qdyl)
		qdyint(&Qmeas[i].Qdyt)
	}
}

func Qmeasmonint(Qmeas []QMEAS) {
	for i := range Qmeas {
		svdyint(&Qmeas[i].mTcdy)
		svdyint(&Qmeas[i].mThdy)
		svdyint(&Qmeas[i].mxcdy)
		svdyint(&Qmeas[i].mxhdy)

		qdyint(&Qmeas[i].mQdys)
		qdyint(&Qmeas[i].mQdyl)
		qdyint(&Qmeas[i].mQdyt)
	}
}

func Qmeasday(Mon, Day, ttmm int, Qmeas []QMEAS, Nday, SimDayend int) {
	for i := range Qmeas {
		// 日次集計
		svdaysum(int64(ttmm), Qmeas[i].PlistG.Control, *Qmeas[i].Th, &Qmeas[i].Thdy) // 温度
		svdaysum(int64(ttmm), Qmeas[i].PlistG.Control, *Qmeas[i].Tc, &Qmeas[i].Tcdy) // 温度

		if Qmeas[i].Xh != nil {
			svdaysum(int64(ttmm), Qmeas[i].PlistG.Control, *Qmeas[i].Xh, &Qmeas[i].xhdy) // 出口湿度
		}

		if Qmeas[i].Xc != nil {
			svdaysum(int64(ttmm), Qmeas[i].PlistG.Control, *Qmeas[i].Xc, &Qmeas[i].xcdy) // 入口湿度
		}

		qdaysum(int64(ttmm), Qmeas[i].PlistG.Control, Qmeas[i].Qs, &Qmeas[i].Qdys) // 流量
		qdaysum(int64(ttmm), Qmeas[i].PlistG.Control, Qmeas[i].Ql, &Qmeas[i].Qdyl) // 湿度流量
		qdaysum(int64(ttmm), Qmeas[i].PlistG.Control, Qmeas[i].Qt, &Qmeas[i].Qdyt) // 総熱量

		// 月次集計
		svmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, *Qmeas[i].Th, &Qmeas[i].mThdy, Nday, SimDayend) // 温度
		svmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, *Qmeas[i].Tc, &Qmeas[i].mTcdy, Nday, SimDayend) // 温度

		if Qmeas[i].Xh != nil {
			svmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, *Qmeas[i].Xh, &Qmeas[i].mxhdy, Nday, SimDayend) // 出口湿度
		}

		if Qmeas[i].Xc != nil {
			svmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, *Qmeas[i].Xc, &Qmeas[i].mxcdy, Nday, SimDayend) // 入口湿度
		}

		qmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, Qmeas[i].Qs, &Qmeas[i].mQdys, Nday, SimDayend) // 流量
		qmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, Qmeas[i].Ql, &Qmeas[i].mQdyl, Nday, SimDayend) // 湿度流量
		qmonsum(Mon, Day, ttmm, Qmeas[i].PlistG.Control, Qmeas[i].Qt, &Qmeas[i].mQdyt, Nday, SimDayend) // 総熱量
	}
}

func Qmeasdyprt(fo io.Writer, id int, Qmeas []QMEAS) {
	for i := range Qmeas {
		switch id {
		case 0:
			if len(Qmeas) > 0 {
				fmt.Fprintf(fo, "%s %d\n", QMEAS_TYPE, len(Qmeas))
			}
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, " %s 1 24\n", Qmeas[i].Name)
			} else {
				fmt.Fprintf(fo, " %s 1 8\n", Qmeas[i].Name)
			}
		case 1:
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
			}
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
		default:
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdys.Hhr, Qmeas[i].Qdys.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdys.Chr, Qmeas[i].Qdys.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].Qdys.Hmxtime, Qmeas[i].Qdys.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].Qdys.Cmxtime, Qmeas[i].Qdys.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdyl.Hhr, Qmeas[i].Qdyl.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdyl.Chr, Qmeas[i].Qdyl.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].Qdyl.Hmxtime, Qmeas[i].Qdyl.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].Qdyl.Cmxtime, Qmeas[i].Qdyl.Cmx)
			}

			fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdyt.Hhr, Qmeas[i].Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].Qdyt.Chr, Qmeas[i].Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].Qdyt.Hmxtime, Qmeas[i].Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].Qdyt.Cmxtime, Qmeas[i].Qdyt.Cmx)
		}
	}
}

func Qmeasmonprt(fo io.Writer, id int, Qmeas []QMEAS) {
	for i := range Qmeas {
		switch id {
		case 0:
			if len(Qmeas) > 0 {
				fmt.Fprintf(fo, "%s %d\n", QMEAS_TYPE, len(Qmeas))
			}
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, " %s 1 24\n", Qmeas[i].Name)
			} else {
				fmt.Fprintf(fo, " %s 1 8\n", Qmeas[i].Name)
			}
		case 1:
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n",
					Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
			}
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name, Qmeas[i].Name)
		default:
			if Qmeas[i].Plistxc != nil && Qmeas[i].Plistxh != nil {
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdys.Hhr, Qmeas[i].mQdys.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdys.Chr, Qmeas[i].mQdys.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].mQdys.Hmxtime, Qmeas[i].mQdys.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].mQdys.Cmxtime, Qmeas[i].mQdys.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdyl.Hhr, Qmeas[i].mQdyl.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdyl.Chr, Qmeas[i].mQdyl.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].mQdyl.Hmxtime, Qmeas[i].mQdyl.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].mQdyl.Cmxtime, Qmeas[i].mQdyl.Cmx)
			}

			fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdyt.Hhr, Qmeas[i].mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Qmeas[i].mQdyt.Chr, Qmeas[i].mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Qmeas[i].mQdyt.Hmxtime, Qmeas[i].mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Qmeas[i].mQdyt.Cmxtime, Qmeas[i].mQdyt.Cmx)
		}
	}
}
