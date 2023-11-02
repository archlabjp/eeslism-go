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

func Qmeaselm(Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
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

func Qmeasene(Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
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

func Qmeasprint(fo io.Writer, id int, Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
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

func Qmeasdyint(Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
		svdyint(&qmeas.Tcdy)
		svdyint(&qmeas.Thdy)
		svdyint(&qmeas.xcdy)
		svdyint(&qmeas.xhdy)

		qdyint(&qmeas.Qdys)
		qdyint(&qmeas.Qdyl)
		qdyint(&qmeas.Qdyt)
	}
}

func Qmeasmonint(Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
		svdyint(&qmeas.mTcdy)
		svdyint(&qmeas.mThdy)
		svdyint(&qmeas.mxcdy)
		svdyint(&qmeas.mxhdy)

		qdyint(&qmeas.mQdys)
		qdyint(&qmeas.mQdyl)
		qdyint(&qmeas.mQdyt)
	}
}

func Qmeasday(Mon, Day, ttmm int, Qmeas []*QMEAS, Nday, SimDayend int) {
	for _, qmeas := range Qmeas {
		// 日次集計
		svdaysum(int64(ttmm), qmeas.PlistG.Control, *qmeas.Th, &qmeas.Thdy) // 温度
		svdaysum(int64(ttmm), qmeas.PlistG.Control, *qmeas.Tc, &qmeas.Tcdy) // 温度

		if qmeas.Xh != nil {
			svdaysum(int64(ttmm), qmeas.PlistG.Control, *qmeas.Xh, &qmeas.xhdy) // 出口湿度
		}

		if qmeas.Xc != nil {
			svdaysum(int64(ttmm), qmeas.PlistG.Control, *qmeas.Xc, &qmeas.xcdy) // 入口湿度
		}

		qdaysum(int64(ttmm), qmeas.PlistG.Control, qmeas.Qs, &qmeas.Qdys) // 流量
		qdaysum(int64(ttmm), qmeas.PlistG.Control, qmeas.Ql, &qmeas.Qdyl) // 湿度流量
		qdaysum(int64(ttmm), qmeas.PlistG.Control, qmeas.Qt, &qmeas.Qdyt) // 総熱量

		// 月次集計
		svmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, *qmeas.Th, &qmeas.mThdy, Nday, SimDayend) // 温度
		svmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, *qmeas.Tc, &qmeas.mTcdy, Nday, SimDayend) // 温度

		if qmeas.Xh != nil {
			svmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, *qmeas.Xh, &qmeas.mxhdy, Nday, SimDayend) // 出口湿度
		}

		if qmeas.Xc != nil {
			svmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, *qmeas.Xc, &qmeas.mxcdy, Nday, SimDayend) // 入口湿度
		}

		qmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, qmeas.Qs, &qmeas.mQdys, Nday, SimDayend) // 流量
		qmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, qmeas.Ql, &qmeas.mQdyl, Nday, SimDayend) // 湿度流量
		qmonsum(Mon, Day, ttmm, qmeas.PlistG.Control, qmeas.Qt, &qmeas.mQdyt, Nday, SimDayend) // 総熱量
	}
}

func Qmeasdyprt(fo io.Writer, id int, Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
		switch id {
		case 0:
			if len(Qmeas) > 0 {
				fmt.Fprintf(fo, "%s %d\n", QMEAS_TYPE, len(Qmeas))
			}
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, " %s 1 24\n", qmeas.Name)
			} else {
				fmt.Fprintf(fo, " %s 1 8\n", qmeas.Name)
			}
		case 1:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
			}
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
		default:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdys.Hhr, qmeas.Qdys.H)
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdys.Chr, qmeas.Qdys.C)
				fmt.Fprintf(fo, "%1d %2.0f ", qmeas.Qdys.Hmxtime, qmeas.Qdys.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.Qdys.Cmxtime, qmeas.Qdys.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdyl.Hhr, qmeas.Qdyl.H)
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdyl.Chr, qmeas.Qdyl.C)
				fmt.Fprintf(fo, "%1d %2.0f ", qmeas.Qdyl.Hmxtime, qmeas.Qdyl.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.Qdyl.Cmxtime, qmeas.Qdyl.Cmx)
			}

			fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdyt.Hhr, qmeas.Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", qmeas.Qdyt.Chr, qmeas.Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", qmeas.Qdyt.Hmxtime, qmeas.Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.Qdyt.Cmxtime, qmeas.Qdyt.Cmx)
		}
	}
}

func Qmeasmonprt(fo io.Writer, id int, Qmeas []*QMEAS) {
	for _, qmeas := range Qmeas {
		switch id {
		case 0:
			if len(Qmeas) > 0 {
				fmt.Fprintf(fo, "%s %d\n", QMEAS_TYPE, len(Qmeas))
			}
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, " %s 1 24\n", qmeas.Name)
			} else {
				fmt.Fprintf(fo, " %s 1 8\n", qmeas.Name)
			}
		case 1:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
				fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n",
					qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
			}
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				qmeas.Name, qmeas.Name, qmeas.Name, qmeas.Name)
		default:
			if qmeas.Plistxc != nil && qmeas.Plistxh != nil {
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdys.Hhr, qmeas.mQdys.H)
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdys.Chr, qmeas.mQdys.C)
				fmt.Fprintf(fo, "%1d %2.0f ", qmeas.mQdys.Hmxtime, qmeas.mQdys.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.mQdys.Cmxtime, qmeas.mQdys.Cmx)

				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdyl.Hhr, qmeas.mQdyl.H)
				fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdyl.Chr, qmeas.mQdyl.C)
				fmt.Fprintf(fo, "%1d %2.0f ", qmeas.mQdyl.Hmxtime, qmeas.mQdyl.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.mQdyl.Cmxtime, qmeas.mQdyl.Cmx)
			}

			fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdyt.Hhr, qmeas.mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", qmeas.mQdyt.Chr, qmeas.mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", qmeas.mQdyt.Hmxtime, qmeas.mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", qmeas.mQdyt.Cmxtime, qmeas.mQdyt.Cmx)
		}
	}
}
