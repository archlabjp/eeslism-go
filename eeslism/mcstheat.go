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

/*  mcstheat.c  */
/*  電気蓄熱式暖房器 */

package eeslism

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

/* 機器仕様入力　　　　　　*/

/*---- Satoh Debug 電気蓄熱式暖房器 2001/1/21 ----*/
func Stheatdata(s string, stheatca *STHEATCA) int {
	var id int

	if st := strings.IndexRune(s, '='); st == -1 {
		stheatca.Name = s
		stheatca.Eff = -999.0
		stheatca.Q = -999.0
		stheatca.Hcap = -999.0
		stheatca.KA = -999.0
	} else {
		sval := s[st+1:]

		if s == "PCM" {
			stheatca.PCMName = sval
		} else {
			dt, err := strconv.ParseFloat(sval, 64)
			if err != nil {
				panic(err)
			}

			if s == "Q" {
				stheatca.Q = dt
			} else if s == "KA" {
				stheatca.KA = dt
			} else if s == "eff" {
				stheatca.Eff = dt
			} else if s == "Hcap" {
				stheatca.Hcap = dt
			} else {
				id = 1
			}
		}
	}
	return id
}

/* --------------------------- */

/*  管長・ダクト長、周囲温度設定 */

func Stheatint(_stheat []STHEAT, Simc *SIMCONTL, Compnt []COMPNT, Wd *WDAT, Npcm int, _PCM []PCM) {
	for i := range _stheat {
		stheat := &_stheat[i]
		if stheat.Cmp.Envname != "" {
			stheat.Tenv = envptr(stheat.Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			stheat.Room = roomptr(stheat.Cmp.Roomname, Compnt)
		}

		if stheat.Cat.PCMName != "" {
			var j int
			for j = 0; j < Npcm; j++ {
				if stheat.Cat.PCMName == _PCM[j].Name {
					stheat.Pcm = &_PCM[j]
				}
			}
			if stheat.Pcm == nil {
				Err := fmt.Sprintf("STHEAT %s のPCM=%sが見つかりません", stheat.Name, stheat.Cat.PCMName)
				Eprint(Err, "<Stheatint>")
				os.Exit(1)
			}
		}

		st := stheat.Cat

		if st.Q < 0.0 {
			Err := fmt.Sprintf("Name=%s  Q=%.4g", stheat.Name, st.Q)
			Eprint("Stheatinit", Err)
		}
		if stheat.Pcm == nil && st.Hcap < 0.0 {
			Err := fmt.Sprintf("Name=%s  Hcap=%.4g", stheat.Name, st.Hcap)
			Eprint("Stheatinit", Err)
		}
		if st.KA < 0.0 {
			Err := fmt.Sprintf("Name=%s  KA=%.4g", stheat.Name, st.KA)
			Eprint("Stheatinit", Err)
		}
		if st.Eff < 0.0 {
			Err := fmt.Sprintf("Name=%s  eff=%.4g", stheat.Name, st.Eff)
			Eprint("Stheatinit", Err)
		}

		var err error
		stheat.Tsold, err = strconv.ParseFloat(stheat.Cmp.Tparm, 64)
		if err != nil {
			panic(err)
		}

		// 内臓PCMの質量
		stheat.MPCM = stheat.Cmp.MPCM
	}
}

/* --------------------------- */

/*  特性式の係数  */

//
//    +--------+ --> [OUT 1]
//    | STHEAT |
//    +--------+ --> [OUT 2]
//

func Stheatcfv(_stheat []STHEAT) {
	for i := range _stheat {
		stheat := &_stheat[i]

		// 作用温度 ?
		var Te float64
		if stheat.Cmp.Envname != "" {
			Te = *(stheat.Tenv)
		} else {
			Te = stheat.Room.Tot
		}

		Eo1 := stheat.Cmp.Elouts[0]
		eff := stheat.Cat.Eff
		stheat.CG = Spcheat(Eo1.Fluid) * Eo1.G
		KA := stheat.Cat.KA
		Tsold := stheat.Tsold
		pcm := stheat.Pcm
		if pcm != nil {
			//NOTE: FNPCMState のシグネチャがヘッダと一致しない。。。
			// stheat.Hcap = stheat.MPCM *
			// 	FNPCMState(pcm.Cros, pcm.Crol, pcm.Ql, pcm.Ts, pcm.Tl, Tsold, nil)
			panic("Cannot call FNPCMState")
		} else {
			stheat.Hcap = stheat.Cat.Hcap
		}
		cG := stheat.CG

		d := stheat.Hcap/DTM + eff*cG + KA
		if stheat.Cmp.Control != OFF_SW {
			stheat.E = stheat.Cat.Q
		} else {
			stheat.E = 0.0
		}

		//  空気が流れていれば出入口温度の関係式係数を作成する
		if Eo1.Control != OFF_SW {
			Eo1.Coeffo = 1.0
			Eo1.Co = eff * (stheat.Hcap/DTM*Tsold + KA*Te + stheat.E) / d
			Eo1.Coeffin[0] = eff - 1.0 - eff*eff*cG/d

			Eo2 := stheat.Cmp.Elouts[1]
			Eo2.Coeffo = 1.0
			Eo2.Co = 0.0
			Eo2.Coeffin[0] = -1.0
		} else {
			Eo1.Coeffo = 1.0
			Eo1.Co = 0.0
			Eo1.Coeffin[0] = -1.0

			Eo2 := stheat.Cmp.Elouts[1]
			Eo2.Coeffo = 1.0
			Eo2.Co = 0.0
			Eo2.Coeffin[0] = -1.0
		}
	}
}

/* --------------------------- */

/* 取得熱量の計算 */

func Stheatene(_stheat []STHEAT) {
	var elo *ELOUT
	var cat *STHEATCA
	var Te float64

	for i := range _stheat {
		stheat := &_stheat[i]
		elo = stheat.Cmp.Elouts[0]
		stheat.Tin = elo.Elins[0].Sysvin

		cat = stheat.Cat

		if stheat.Cmp.Envname != "" {
			Te = *(stheat.Tenv)
		} else {
			Te = stheat.Room.Tot
		}

		stheat.Tout = elo.Sysv
		stheat.Ts = (stheat.Hcap/DTM*stheat.Tsold +
			cat.Eff*stheat.CG*stheat.Tin +
			cat.KA*Te + stheat.E) /
			(stheat.Hcap/DTM + cat.Eff*stheat.CG + cat.KA)

		stheat.Q = stheat.CG * (stheat.Tout - stheat.Tin)

		stheat.Qls = stheat.Cat.KA * (Te - stheat.Ts)

		stheat.Qsto = stheat.Hcap / DTM * (stheat.Ts - stheat.Tsold)

		stheat.Tsold = stheat.Ts

		if stheat.Room != nil {
			stheat.Room.Qeqp += (-stheat.Qls)
		}
	}
}

func stheatvptr(key []string, Stheat *STHEAT, vptr *VPTR, vpath *VPTR) int {
	var err int

	if key[1] == "Ts" {
		vptr.Ptr = &Stheat.Tsold
		vptr.Type = VAL_CTYPE
	} else if key[1] == "control" {
		vpath.Type = 's'
		vpath.Ptr = Stheat
		vptr.Ptr = &Stheat.Cmp.Control
		vptr.Type = SW_CTYPE
	} else {
		err = 1
	}

	return err
}

/* ---------------------------*/

func stheatprint(fo io.Writer, id int, stheat []STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 11\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ts t f %s_Ti t f %s_To t f %s_Q q f %s_Qsto q f ",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls q f %s_Ec c c %s_E e f ",
				stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hcap q f\n", stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %4.1f %2.0f %.4g ",
				stheat[i].Cmp.Elouts[0].Control, stheat[i].Cmp.Elouts[0].G,
				stheat[i].Ts,
				stheat[i].Tin, stheat[i].Tout, stheat[i].Q, stheat[i].Qsto)
			fmt.Fprintf(fo, "%.4g %c %2.0f ",
				stheat[i].Qls, stheat[i].Cmp.Control, stheat[i].E)
			fmt.Fprintf(fo, "%.0f\n", stheat[i].Hcap)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

/*******************/
func stheatdyint(stheat []STHEAT) {
	for i := range stheat {
		stheat[i].Qlossdy = 0.0
		stheat[i].Qstody = 0.0

		svdyint(&stheat[i].Tidy)
		svdyint(&stheat[i].Tsdy)
		svdyint(&stheat[i].Tody)
		qdyint(&stheat[i].Qdy)
		edyint(&stheat[i].Edy)
	}
}

func stheatmonint(stheat []STHEAT) {
	for i := range stheat {
		stheat[i].MQlossdy = 0.0
		stheat[i].MQstody = 0.0

		svdyint(&stheat[i].MTidy)
		svdyint(&stheat[i].MTsdy)
		svdyint(&stheat[i].MTody)
		qdyint(&stheat[i].MQdy)
		edyint(&stheat[i].MEdy)
	}
}

func stheatday(Mon, Day, ttmm int, stheat []STHEAT, Nday, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for i := range stheat {
		// 日集計
		stheat[i].Qlossdy += stheat[i].Qls
		stheat[i].Qstody += stheat[i].Qsto
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Tin, &stheat[i].Tidy)
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Tout, &stheat[i].Tody)
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Ts, &stheat[i].Tsdy)
		qdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Q, &stheat[i].Qdy)
		edaysum(ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].Edy)

		// 月集計
		stheat[i].MQlossdy += stheat[i].Qls
		stheat[i].MQstody += stheat[i].Qsto
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Tin, &stheat[i].MTidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Tout, &stheat[i].MTody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Ts, &stheat[i].MTsdy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Q, &stheat[i].MQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].MEdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].MtEdy[Mo][tt])
	}
}

func stheatdyprt(fo io.Writer, id int, stheat []STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 32\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n",
				stheat[i].Name, stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tidy.Hrs, stheat[i].Tidy.M,
				stheat[i].Tidy.Mntime, stheat[i].Tidy.Mn,
				stheat[i].Tidy.Mxtime, stheat[i].Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tody.Hrs, stheat[i].Tody.M,
				stheat[i].Tody.Mntime, stheat[i].Tody.Mn,
				stheat[i].Tody.Mxtime, stheat[i].Tody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tsdy.Hrs, stheat[i].Tsdy.M,
				stheat[i].Tsdy.Mntime, stheat[i].Tsdy.Mn,
				stheat[i].Tsdy.Mxtime, stheat[i].Tsdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Qdy.Hhr, stheat[i].Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Qdy.Chr, stheat[i].Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Qdy.Hmxtime, stheat[i].Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Qdy.Cmxtime, stheat[i].Qdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Edy.Hrs, stheat[i].Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Edy.Mxtime, stheat[i].Edy.Mx)
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stheat[i].Qlossdy*Cff_kWh, stheat[i].Qstody*Cff_kWh)
		}
	}
}

func stheatmonprt(fo io.Writer, id int, stheat []STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 32\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n",
				stheat[i].Name, stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTidy.Hrs, stheat[i].MTidy.M,
				stheat[i].MTidy.Mntime, stheat[i].MTidy.Mn,
				stheat[i].MTidy.Mxtime, stheat[i].MTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTody.Hrs, stheat[i].MTody.M,
				stheat[i].MTody.Mntime, stheat[i].MTody.Mn,
				stheat[i].MTody.Mxtime, stheat[i].MTody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTsdy.Hrs, stheat[i].MTsdy.M,
				stheat[i].MTsdy.Mntime, stheat[i].MTsdy.Mn,
				stheat[i].MTsdy.Mxtime, stheat[i].MTsdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MQdy.Hhr, stheat[i].MQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MQdy.Chr, stheat[i].MQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MQdy.Hmxtime, stheat[i].MQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MQdy.Cmxtime, stheat[i].MQdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MEdy.Hrs, stheat[i].MEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MEdy.Mxtime, stheat[i].MEdy.Mx)
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stheat[i].MQlossdy*Cff_kWh, stheat[i].MQstody*Cff_kWh)
		}
	}
}

func stheatmtprt(fo io.Writer, id int, stheat []STHEAT, Mo, tt int) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 1\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_E E f \n", stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, " %.2f\n", stheat[i].MtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}
