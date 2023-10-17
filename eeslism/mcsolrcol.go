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

/*  solrcol.c  */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*　太陽熱集熱器

機器仕様入力　　　　*/

func Colldata(typeStr EqpType, s string, Collca *COLLCA) int {
	var st string
	id := 0

	if typeStr == COLLECTOR_TYPE {
		Collca.Type = COLLECTOR_PDT
	} else {
		Collca.Type = ACOLLECTOR_PDT
	}

	if idx := strings.Index(s, "="); idx == -1 {
		Collca.name = s
		Collca.b0 = -999.0
		Collca.b1 = -999.0
		Collca.Ac = -999.0
		Collca.Ag = -999.0
	} else {
		st = s[idx+1:]
		s = s[:idx]

		dt, err := strconv.ParseFloat(st, 64)
		if err != nil {
			panic(err)
		}

		switch s {
		case "b0":
			Collca.b0 = dt
		case "b1":
			Collca.b1 = dt
		case "Fd":
			Collca.Fd = dt
		case "Ac":
			Collca.Ac = dt
		case "Ag":
			Collca.Ag = dt
		default:
			id = 1
		}
	}

	return id
}

/* ------------------------------------- */

/*  初期設定 */

func Collint(Coll []COLL, Nexsf int, Exs []EXSF, Wd *WDAT) {
	for i := range Coll {
		Coll[i].Ta = &Wd.T
		Coll[i].sol = nil
		for j := 0; j < Nexsf; j++ {
			exs := &Exs[j]
			if Coll[i].Cmp.Exsname == exs.Name {
				Coll[i].sol = exs
			}
		}
		if Coll[i].sol == nil {
			Eprint("Collint", Coll[i].Cmp.Exsname)
		}

		if Coll[i].Cat.b0 < 0.0 {
			Err := fmt.Sprintf("Name=%s b0=%.4g", Coll[i].Cmp.Name, Coll[i].Cat.b0)
			Eprint("Collint", Err)
		}
		if Coll[i].Cat.b1 < 0.0 {
			Err := fmt.Sprintf("Name=%s b1=%.4g", Coll[i].Cmp.Name, Coll[i].Cat.b1)
			Eprint("Collint", Err)
		}
		if Coll[i].Cat.Ac < 0.0 {
			Err := fmt.Sprintf("Name=%s Ac=%.4g", Coll[i].Cmp.Name, Coll[i].Cat.Ac)
			Eprint("Collint", Err)
		}
		if Coll[i].Cat.Ag < 0.0 {
			Err := fmt.Sprintf("Name=%s Ag=%.4g", Coll[i].Cmp.Name, Coll[i].Cat.Ag)
			Eprint("Collint", Err)
		}

		// 総合熱損失係数[W/(m2･K)]の計算
		Coll[i].Cat.Ko = Coll[i].Cat.b1 / Coll[i].Cat.Fd
	}
}

// 集熱器の相当外気温度を計算する（制御用）
func CalcCollTe(Coll []COLL) {
	for i := range Coll {
		tgaKo := Coll[i].Cat.b0 / Coll[i].Cat.b1
		Coll[i].Te = scolte(tgaKo, Coll[i].sol.Cinc, Coll[i].sol.Idre, Coll[i].sol.Idf, *Coll[i].Ta)
	}
}

/* ------------------------------------- */

/*  特性式の係数   */

//
//   +------+ ---> [OUT 1]
//   | COLL |
//   +------+ ---> [OUT 2] ACOLLECTOR_PDTのみ
//
func Collcfv(Coll []COLL) {
	for i := range Coll {
		// 制御用の相当外気温度（現在時刻）は計算済みなのでここでは計算しない
		if Coll[i].Cmp.Control != OFF_SW {
			Eo1 := Coll[i].Cmp.Elouts[0]
			Kcw := Coll[i].Cat.b1
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			Coll[i].ec = 1.0 - math.Exp(-Kcw*Coll[i].Cmp.Ac/cG)
			Coll[i].D1 = cG * Coll[i].ec
			Coll[i].Do = Coll[i].D1 * Coll[i].Te

			Eo1.Coeffo = cG
			Eo1.Co = Coll[i].Do
			Eo1.Coeffin[0] = Coll[i].D1 - cG

			if Coll[i].Cat.Type == ACOLLECTOR_PDT {
				Eo2 := Coll[i].Cmp.Elouts[1]
				Eo2.Coeffo = 1.0
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -1.0
			}
		}
	}
}

/* ------------------------------------- */

/*  集熱量の計算 */

func Collene(Coll []COLL) {
	for i := range Coll {
		Coll[i].Tin = Coll[i].Cmp.Elins[0].Sysvin

		if Coll[i].Cmp.Control != OFF_SW {
			Coll[i].Q = Coll[i].Do - Coll[i].D1*Coll[i].Tin
		} else {
			Coll[i].Q = 0.0
		}

		// 集熱板温度の計算
		if Coll[i].Q > 0.0 {
			Coll[i].Tcb = Coll[i].Te - Coll[i].Q/(Coll[i].Ac*Coll[i].Cat.Ko)
		} else {
			// 集熱ポンプ停止時は相当外気温度に等しいとする
			Coll[i].Tcb = Coll[i].Te
		}

		Coll[i].Sol = Coll[i].sol.Iw * Coll[i].Cmp.Ac
	}
}

/* ------------------------------------- */

/*  集熱器到達温度　Te　　　　　　　*/

func scolte(rtgko, cinc, Idre, Idf, Ta float64) float64 {
	Cidf := 0.91
	return rtgko*(Glscid(cinc)*Idre+Cidf*Idf) + Ta
}

/* ------------------------------------- */

/*  集熱器内部変数のポインター  */

func collvptr(key []string, Coll *COLL, vptr *VPTR) int {
	err := 0

	if key[1] == "Te" {
		vptr.Ptr = &Coll.Te
		vptr.Type = VAL_CTYPE
	} else if key[1] == "Tcb" {
		vptr.Ptr = &Coll.Tcb
		vptr.Type = VAL_CTYPE
	} else {
		err = 1
	}

	return err
}

/* ------------------------------------------------------------- */

func collprint(fo io.Writer, id int, Coll []COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for i := range Coll {
			fmt.Fprintf(fo, " %s 1 7\n", Coll[i].Name)
		}
	case 1:
		for i := range Coll {
			fmt.Fprintf(fo, "%s_c c c %s_Ti t f %s_To t f %s_Te t f %s_Tcb t f %s_Q q f %s_S e f\n",
				Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
		}
	default:
		for i := range Coll {
			fmt.Fprintf(fo, "%c %4.1f %4.1f %4.1f %4.1f %3.0f %3.0f\n",
				Coll[i].Cmp.Elouts[0].Control,
				Coll[i].Tin, Coll[i].Cmp.Elouts[0].Sysv, Coll[i].Te, Coll[i].Tcb, Coll[i].Q, Coll[i].Sol)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func colldyint(Coll []COLL) {
	for i := range Coll {
		svdyint(&Coll[i].Tidy) // 温度の日次積分
		qdyint(&Coll[i].Qdy)   // 熱量の日次積分
		edyint(&Coll[i].Soldy) // 経済性指標の日次積分
	}
}

func collmonint(Coll []COLL) {
	for i := range Coll {
		svdyint(&Coll[i].mTidy) // 温度の月次積分
		qdyint(&Coll[i].mQdy)   // 熱量の月次積分
		edyint(&Coll[i].mSoldy) // 経済性指標の月次積分
	}
}

func collday(Mon, Day, ttmm int, Coll []COLL, Nday, SimDayend int) {
	var sw ControlSWType

	for i := range Coll {
		// 日次集計
		svdaysum(int64(ttmm), Coll[i].Cmp.Control, Coll[i].Tin, &Coll[i].Tidy) // 温度の日次集計
		qdaysum(int64(ttmm), Coll[i].Cmp.Control, Coll[i].Q, &Coll[i].Qdy)     // 熱量の日次集計

		if Coll[i].Sol > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		edaysum(ttmm, sw, Coll[i].Sol, &Coll[i].Soldy)

		// 月次集計
		svmonsum(Mon, Day, ttmm, Coll[i].Cmp.Control, Coll[i].Tin, &Coll[i].mTidy, Nday, SimDayend) // 温度の月次集計
		qmonsum(Mon, Day, ttmm, Coll[i].Cmp.Control, Coll[i].Q, &Coll[i].mQdy, Nday, SimDayend)     // 熱量の月次集計

		if Coll[i].Sol > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		emonsum(Mon, Day, ttmm, sw, Coll[i].Sol, &Coll[i].mSoldy, Nday, SimDayend)
	}
}

func colldyprt(fo io.Writer, id int, Coll []COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for i := range Coll {
			fmt.Fprintf(fo, " %s 1 18\n", Coll[i].Name)
		}
	case 1:
		for i := range Coll {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
		}
	default:
		for i := range Coll {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Coll[i].Tidy.Hrs, Coll[i].Tidy.M, Coll[i].Tidy.Mntime, Coll[i].Tidy.Mn,
				Coll[i].Tidy.Mxtime, Coll[i].Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].Qdy.Hhr, Coll[i].Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].Qdy.Chr, Coll[i].Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Coll[i].Qdy.Hmxtime, Coll[i].Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Coll[i].Qdy.Cmxtime, Coll[i].Qdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].Soldy.Hrs, Coll[i].Soldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", Coll[i].Soldy.Mxtime, Coll[i].Soldy.Mx)
		}
	}
}

func collmonprt(fo io.Writer, id int, Coll []COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for i := range Coll {
			fmt.Fprintf(fo, " %s 1 18\n", Coll[i].Name)
		}
	case 1:
		for i := range Coll {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", Coll[i].Name, Coll[i].Name, Coll[i].Name, Coll[i].Name)
		}
	default:
		for i := range Coll {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Coll[i].mTidy.Hrs, Coll[i].mTidy.M, Coll[i].mTidy.Mntime,
				Coll[i].mTidy.Mn, Coll[i].mTidy.Mxtime, Coll[i].mTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].mQdy.Hhr, Coll[i].mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].mQdy.Chr, Coll[i].mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Coll[i].mQdy.Hmxtime, Coll[i].mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Coll[i].mQdy.Cmxtime, Coll[i].mQdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", Coll[i].mSoldy.Hrs, Coll[i].mSoldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", Coll[i].mSoldy.Mxtime, Coll[i].mSoldy.Mx)
		}
	}
}
