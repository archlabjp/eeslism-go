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
	"errors"
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

func Collint(Coll []*COLL, Exs []*EXSF, Wd *WDAT) {
	for _, coll := range Coll {
		coll.Ta = &Wd.T
		coll.sol = nil
		for _, exs := range Exs {
			if coll.Cmp.Exsname == exs.Name {
				coll.sol = exs
			}
		}
		if coll.sol == nil {
			Eprint("Collint", coll.Cmp.Exsname)
		}

		if coll.Cat.b0 < 0.0 {
			Err := fmt.Sprintf("Name=%s b0=%.4g", coll.Cmp.Name, coll.Cat.b0)
			Eprint("Collint", Err)
		}
		if coll.Cat.b1 < 0.0 {
			Err := fmt.Sprintf("Name=%s b1=%.4g", coll.Cmp.Name, coll.Cat.b1)
			Eprint("Collint", Err)
		}
		if coll.Cat.Ac < 0.0 {
			Err := fmt.Sprintf("Name=%s Ac=%.4g", coll.Cmp.Name, coll.Cat.Ac)
			Eprint("Collint", Err)
		}
		if coll.Cat.Ag < 0.0 {
			Err := fmt.Sprintf("Name=%s Ag=%.4g", coll.Cmp.Name, coll.Cat.Ag)
			Eprint("Collint", Err)
		}

		// 総合熱損失係数[W/(m2･K)]の計算
		coll.Cat.Ko = coll.Cat.b1 / coll.Cat.Fd
	}
}

// 集熱器の相当外気温度を計算する（制御用）
func CalcCollTe(Coll []*COLL) {
	for _, coll := range Coll {
		tgaKo := coll.Cat.b0 / coll.Cat.b1
		coll.Te = scolte(tgaKo, coll.sol.Cinc, coll.sol.Idre, coll.sol.Idf, *coll.Ta)
	}
}

/* ------------------------------------- */

/*  特性式の係数   */

//
//   +------+ ---> [OUT 1]
//   | COLL |
//   +------+ ---> [OUT 2] ACOLLECTOR_PDTのみ
//
func Collcfv(Coll []*COLL) {
	for _, coll := range Coll {
		// 制御用の相当外気温度（現在時刻）は計算済みなのでここでは計算しない
		if coll.Cmp.Control != OFF_SW {
			Eo1 := coll.Cmp.Elouts[0]
			Kcw := coll.Cat.b1
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			coll.ec = 1.0 - math.Exp(-Kcw*coll.Cmp.Ac/cG)
			coll.D1 = cG * coll.ec
			coll.Do = coll.D1 * coll.Te

			Eo1.Coeffo = cG
			Eo1.Co = coll.Do
			Eo1.Coeffin[0] = coll.D1 - cG

			if coll.Cat.Type == ACOLLECTOR_PDT {
				Eo2 := coll.Cmp.Elouts[1]
				Eo2.Coeffo = 1.0
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -1.0
			}
		}
	}
}

/* ------------------------------------- */

/*  集熱量の計算 */

func Collene(Coll []*COLL) {
	for _, coll := range Coll {
		coll.Tin = coll.Cmp.Elins[0].Sysvin

		if coll.Cmp.Control != OFF_SW {
			coll.Q = coll.Do - coll.D1*coll.Tin
		} else {
			coll.Q = 0.0
		}

		// 集熱板温度の計算
		if coll.Q > 0.0 {
			coll.Tcb = coll.Te - coll.Q/(coll.Ac*coll.Cat.Ko)
		} else {
			// 集熱ポンプ停止時は相当外気温度に等しいとする
			coll.Tcb = coll.Te
		}

		coll.Sol = coll.sol.Iw * coll.Cmp.Ac
	}
}

/* ------------------------------------- */

/*  集熱器到達温度　Te　　　　　　　*/

func scolte(rtgko, cinc, Idre, Idf, Ta float64) float64 {
	Cidf := 0.91
	return rtgko*(Glscid(cinc)*Idre+Cidf*Idf) + Ta
}

/* ------------------------------------- */

// 集熱器内部変数のポインターの作成
func collvptr(key []string, Coll *COLL) (VPTR, error) {
	if key[1] == "Te" {
		return VPTR{Ptr: &Coll.Te, Type: VAL_CTYPE}, nil
	} else if key[1] == "Tcb" {
		return VPTR{Ptr: &Coll.Tcb, Type: VAL_CTYPE}, nil
	}

	return VPTR{}, errors.New("collvptr error")
}

/* ------------------------------------------------------------- */

func collprint(fo io.Writer, id int, Coll []*COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for _, coll := range Coll {
			fmt.Fprintf(fo, " %s 1 7\n", coll.Name)
		}
	case 1:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%s_c c c %s_Ti t f %s_To t f %s_Te t f %s_Tcb t f %s_Q q f %s_S e f\n",
				coll.Name, coll.Name, coll.Name, coll.Name, coll.Name, coll.Name, coll.Name)
		}
	default:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%c %4.1f %4.1f %4.1f %4.1f %3.0f %3.0f\n",
				coll.Cmp.Elouts[0].Control,
				coll.Tin, coll.Cmp.Elouts[0].Sysv, coll.Te, coll.Tcb, coll.Q, coll.Sol)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func colldyint(Coll []*COLL) {
	for _, coll := range Coll {
		svdyint(&coll.Tidy) // 温度の日次積分
		qdyint(&coll.Qdy)   // 熱量の日次積分
		edyint(&coll.Soldy) // 経済性指標の日次積分
	}
}

func collmonint(Coll []*COLL) {
	for _, coll := range Coll {
		svdyint(&coll.mTidy) // 温度の月次積分
		qdyint(&coll.mQdy)   // 熱量の月次積分
		edyint(&coll.mSoldy) // 経済性指標の月次積分
	}
}

func collday(Mon, Day, ttmm int, Coll []*COLL, Nday, SimDayend int) {
	var sw ControlSWType

	for _, coll := range Coll {
		// 日次集計
		svdaysum(int64(ttmm), coll.Cmp.Control, coll.Tin, &coll.Tidy) // 温度の日次集計
		qdaysum(int64(ttmm), coll.Cmp.Control, coll.Q, &coll.Qdy)     // 熱量の日次集計

		if coll.Sol > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		edaysum(ttmm, sw, coll.Sol, &coll.Soldy)

		// 月次集計
		svmonsum(Mon, Day, ttmm, coll.Cmp.Control, coll.Tin, &coll.mTidy, Nday, SimDayend) // 温度の月次集計
		qmonsum(Mon, Day, ttmm, coll.Cmp.Control, coll.Q, &coll.mQdy, Nday, SimDayend)     // 熱量の月次集計

		if coll.Sol > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		emonsum(Mon, Day, ttmm, sw, coll.Sol, &coll.mSoldy, Nday, SimDayend)
	}
}

func colldyprt(fo io.Writer, id int, Coll []*COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for _, coll := range Coll {
			fmt.Fprintf(fo, " %s 1 18\n", coll.Name)
		}
	case 1:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", coll.Name, coll.Name, coll.Name, coll.Name)
		}
	default:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				coll.Tidy.Hrs, coll.Tidy.M, coll.Tidy.Mntime, coll.Tidy.Mn,
				coll.Tidy.Mxtime, coll.Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.Qdy.Hhr, coll.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.Qdy.Chr, coll.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", coll.Qdy.Hmxtime, coll.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", coll.Qdy.Cmxtime, coll.Qdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.Soldy.Hrs, coll.Soldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", coll.Soldy.Mxtime, coll.Soldy.Mx)
		}
	}
}

func collmonprt(fo io.Writer, id int, Coll []*COLL) {
	switch id {
	case 0:
		if len(Coll) > 0 {
			fmt.Fprintf(fo, "%s %d\n", COLLECTOR_TYPE, len(Coll))
		}
		for _, coll := range Coll {
			fmt.Fprintf(fo, " %s 1 18\n", coll.Name)
		}
	case 1:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n", coll.Name, coll.Name, coll.Name, coll.Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", coll.Name, coll.Name, coll.Name, coll.Name)
		}
	default:
		for _, coll := range Coll {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				coll.mTidy.Hrs, coll.mTidy.M, coll.mTidy.Mntime,
				coll.mTidy.Mn, coll.mTidy.Mxtime, coll.mTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.mQdy.Hhr, coll.mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.mQdy.Chr, coll.mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", coll.mQdy.Hmxtime, coll.mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", coll.mQdy.Cmxtime, coll.mQdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", coll.mSoldy.Hrs, coll.mSoldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", coll.mSoldy.Mxtime, coll.mSoldy.Mx)
		}
	}
}
