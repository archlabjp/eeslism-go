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

/*  pipe.c  */

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*  配管・ダクト 仕様入力 */

func Pipedata(cattype string, s string, Pipeca *PIPECA) int {
	var st string
	var dt float64
	var id int

	if cattype == DUCT_TYPE {
		Pipeca.Type = DUCT_PDT
	} else {
		Pipeca.Type = PIPE_PDT
	}

	st = strings.Split(s, "=")[1]

	var err error
	dt, err = strconv.ParseFloat(st, 64)
	if err != nil {
		panic("Failed to parse float: " + err.Error())
	}

	if strings.HasPrefix(s, "Ko") {
		Pipeca.Ko = dt
	} else {
		id = 1
	}

	return id
}

/* --------------------------- */

/*  管長・ダクト長、周囲温度設定 */

func Pipeint(Pipe []PIPE, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT) {
	for i := range Pipe {
		if Pipe[i].Cmp.Ivparm != nil {
			Pipe[i].L = *Pipe[i].Cmp.Ivparm
		} else {
			Pipe[i].L = -999.0
		}

		if Pipe[i].Cmp.Envname != "" {
			Pipe[i].Tenv = envptr(Pipe[i].Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			Pipe[i].Room = roomptr(Pipe[i].Cmp.Roomname, Compnt)
		}

		if Pipe[i].Cat.Ko < 0.0 {
			Err := fmt.Sprintf("Name=%s  Ko=%.4g", Pipe[i].Cmp.Name, Pipe[i].Cat.Ko)
			Eprint("Pipeint", Err)
		}

		if Pipe[i].L < 0.0 {
			Err := fmt.Sprintf("Name=%s  L=%.4g", Pipe[i].Cmp.Name, Pipe[i].L)
			Eprint("Pipeint", Err)
		}
	}
}

/* --------------------------- */

/*  特性式の係数  */

//
// [IN 1] ---> +------+ ---> [OUT 1] 空気 or 温水温度
//             | PIPE |
// [IN 2] ---> +------+ ---> [OUT 2] 湿度 (DUCT_PDTのみ)
//
func Pipecfv(Pipe []PIPE) {
	for i := range Pipe {
		Te := 0.0
		if Pipe[i].Cmp.Control != OFF_SW {
			if Pipe[i].Cmp.Envname != "" {
				Te = *Pipe[i].Tenv
			} else if Pipe[i].Room != nil {
				Te = Pipe[i].Room.Tot
			} else {
				Err := fmt.Sprintf("Undefined Pipe Environment  name=%s", Pipe[i].Name)
				Eprint("<Pipecfv>", Err)
			}
			Pipe[i].Ko = Pipe[i].Cat.Ko

			Eo1 := Pipe[i].Cmp.Elouts[0]
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			Pipe[i].Ep = 1.0 - math.Exp(-(Pipe[i].Ko*Pipe[i].L)/cG)
			Pipe[i].D1 = cG * Pipe[i].Ep
			Pipe[i].Do = Pipe[i].D1 * Te
			Eo1.Coeffo = cG
			Eo1.Co = Pipe[i].Do
			Eo1.Coeffin[0] = Pipe[i].D1 - cG

			if Pipe[i].Cat.Type == DUCT_PDT {
				Eo2 := Pipe[i].Cmp.Elouts[1]
				Eo2.Coeffo = 1.0
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -1.0
			}
		}
	}
}

/* --------------------------- */

/* 取得熱量の計算 */

func Pipeene(Pipe []PIPE) {
	for i := range Pipe {
		Pipe[i].Tin = Pipe[i].Cmp.Elins[0].Sysvin

		if Pipe[i].Cmp.Control != OFF_SW {
			Eo := Pipe[i].Cmp.Elouts[0]
			Pipe[i].Tout = Pipe[i].Do
			Pipe[i].Q = Pipe[i].Do - Pipe[i].D1*Pipe[i].Tin

			if Pipe[i].Room != nil {
				Pipe[i].Room.Qeqp += (-Pipe[i].Q)
			}

			if Pipe[i].Cat.Type == DUCT_PDT {
				Eo = Pipe[i].Cmp.Elouts[1]
				Pipe[i].Xout = Eo.Sysv
				Pipe[i].RHout = FNRhtx(Pipe[i].Tout, Pipe[i].Xout)
				Pipe[i].Hout = FNH(Pipe[i].Tout, Eo.Sysv)
			} else {
				Pipe[i].Hout = -999.0
			}
		} else {
			Pipe[i].Q = 0.0
		}
	}
}

/* --------------------------- */

/* 負荷計算用設定値のポインター */

func pipeldsptr(load *ControlSWType, key []string, Pipe *PIPE, idmrk *byte) (VPTR, error) {
	var err error
	var vptr VPTR

	if key[1] == "Tout" {
		vptr.Ptr = &Pipe.Toset
		vptr.Type = VAL_CTYPE
		Pipe.Loadt = load
		*idmrk = 't'
	} else if Pipe.Cat.Type == DUCT_PDT && key[1] == "xout" {
		vptr.Ptr = &Pipe.Xoset
		vptr.Type = VAL_CTYPE
		Pipe.Loadx = load
		*idmrk = 'x'
	} else {
		err = errors.New("Tout or xout expected")
	}

	return vptr, err
}

/* ------------------------------------------ */

/* 負荷計算用設定値のスケジュール設定 */

func pipeldsschd(Pipe *PIPE) {
	Eo := Pipe.Cmp.Elouts[0]

	if Pipe.Loadt != nil {
		if Eo.Control != OFF_SW {
			if Pipe.Toset > TEMPLIMIT {
				Eo.Control = LOAD_SW
				Eo.Sysv = Pipe.Toset
			} else {
				Eo.Control = OFF_SW
			}
		}
	}

	if Pipe.Cat.Type == DUCT_PDT && Pipe.Loadx != nil {
		if len(Pipe.Cmp.Elouts) > 1 {
			Eo = Pipe.Cmp.Elouts[1]
			if Eo.Control != OFF_SW {
				if Pipe.Xoset > 0.0 {
					Eo.Control = LOAD_SW
					Eo.Sysv = Pipe.Xoset
				} else {
					Eo.Control = OFF_SW
				}
			}
		}
	}
}

/* --------------------------- */

func pipeprint(fo io.Writer, id int, Pipe []PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for i := range Pipe {
			fmt.Fprintf(fo, " %s 1 5\n", Pipe[i].Name)
		}
	case 1:
		for i := range Pipe {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f %s_Q q f\n",
				Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
		}
	default:
		for i := range Pipe {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %.2f\n",
				Pipe[i].Cmp.Elouts[0].Control, Pipe[i].Cmp.Elouts[0].G,
				Pipe[i].Tin, Pipe[i].Cmp.Elouts[0].Sysv, Pipe[i].Q)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func pipedyint(Pipe []PIPE) {
	for i := range Pipe {
		svdyint(&Pipe[i].Tidy)
		qdyint(&Pipe[i].Qdy)
	}
}

func pipemonint(Pipe []PIPE) {
	for i := range Pipe {
		svdyint(&Pipe[i].MTidy)
		qdyint(&Pipe[i].MQdy)
	}
}

func pipeday(Mon int, Day int, ttmm int, Pipe []PIPE, Nday int, SimDayend int) {
	for i := range Pipe {
		// 日集計
		svdaysum(int64(ttmm), Pipe[i].Cmp.Elouts[0].Control, Pipe[i].Tin, &Pipe[i].Tidy)
		qdaysum(int64(ttmm), Pipe[i].Cmp.Elouts[0].Control, Pipe[i].Q, &Pipe[i].Qdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, Pipe[i].Cmp.Elouts[0].Control, Pipe[i].Tin, &Pipe[i].MTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Pipe[i].Cmp.Elouts[0].Control, Pipe[i].Q, &Pipe[i].MQdy, Nday, SimDayend)
	}
}

func pipedyprt(fo io.Writer, id int, Pipe []PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for i := range Pipe {
			fmt.Fprintf(fo, " %s 1 14\n", Pipe[i].Name)
		}

	case 1:
		for i := range Pipe {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
		}

	default:
		for i := range Pipe {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Pipe[i].Tidy.Hrs, Pipe[i].Tidy.M, Pipe[i].Tidy.Mntime,
				Pipe[i].Tidy.Mn, Pipe[i].Tidy.Mxtime, Pipe[i].Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Pipe[i].Qdy.Hhr, Pipe[i].Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Pipe[i].Qdy.Chr, Pipe[i].Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Pipe[i].Qdy.Hmxtime, Pipe[i].Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Pipe[i].Qdy.Cmxtime, Pipe[i].Qdy.Cmx)
		}
	}
}

func pipemonprt(fo io.Writer, id int, Pipe []PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for i := range Pipe {
			fmt.Fprintf(fo, " %s 1 14\n", Pipe[i].Name)
		}

	case 1:
		for i := range Pipe {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n\n", Pipe[i].Name, Pipe[i].Name, Pipe[i].Name, Pipe[i].Name)
		}

	default:
		for i := range Pipe {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Pipe[i].MTidy.Hrs, Pipe[i].MTidy.M, Pipe[i].MTidy.Mntime,
				Pipe[i].MTidy.Mn, Pipe[i].MTidy.Mxtime, Pipe[i].MTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Pipe[i].MQdy.Hhr, Pipe[i].MQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Pipe[i].MQdy.Chr, Pipe[i].MQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Pipe[i].MQdy.Hmxtime, Pipe[i].MQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Pipe[i].MQdy.Cmxtime, Pipe[i].MQdy.Cmx)
		}
	}
}

// 配管、ダクト内部変数のポインターを作成します
func pipevptr(key []string, Pipe *PIPE) (VPTR, error) {
	var err error
	var vptr VPTR
	switch key[1] {
	case "Tout":
		vptr.Ptr = &Pipe.Tout
		vptr.Type = VAL_CTYPE
	case "hout":
		vptr.Ptr = &Pipe.Hout
		vptr.Type = VAL_CTYPE
	case "xout":
		vptr.Ptr = &Pipe.Xout
		vptr.Type = VAL_CTYPE
	case "RHout":
		vptr.Ptr = &Pipe.RHout
		vptr.Type = VAL_CTYPE
	default:
		err = errors.New("Tout, hout, xout or RHout is expected")
	}

	return vptr, err
}
