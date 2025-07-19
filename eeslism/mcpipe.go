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

func Pipedata(cattype EqpType, s string, Pipeca *PIPECA) int {
	var st string
	var dt float64
	var id int

	if cattype == DUCT_TYPE {
		Pipeca.Type = DUCT_PDT
	} else if cattype == PIPEDUCT_TYPE {
		Pipeca.Type = PIPE_PDT
	} else {
		panic(cattype)
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

func Pipeint(Pipe []*PIPE, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT) {
	for _, pipe := range Pipe {
		if pipe.Cmp.Ivparm != nil {
			pipe.L = *pipe.Cmp.Ivparm
		} else {
			pipe.L = FNAN
		}

		if pipe.Cmp.Envname != "" {
			pipe.Tenv = envptr(pipe.Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			pipe.Room = roomptr(pipe.Cmp.Roomname, Compnt)
		}

		if pipe.Cat.Ko < 0.0 {
			Err := fmt.Sprintf("Name=%s  Ko=%.4g", pipe.Cmp.Name, pipe.Cat.Ko)
			Eprint("Pipeint", Err)
		}

		if pipe.L < 0.0 {
			Err := fmt.Sprintf("Name=%s  L=%.4g", pipe.Cmp.Name, pipe.L)
			Eprint("Pipeint", Err)
		}
	}
}

/* --------------------------- */

/*  特性式の係数  */

// [IN 1] ---> +------+ ---> [OUT 1] 空気 or 温水温度
//
//	| PIPE |
//
// [IN 2] ---> +------+ ---> [OUT 2] 湿度 (DUCT_PDTのみ)
func Pipecfv(Pipe []*PIPE) {
	for _, pipe := range Pipe {
		Te := 0.0
		if pipe.Cmp.Control != OFF_SW {
			if pipe.Cmp.Envname != "" {
				Te = *pipe.Tenv
			} else if pipe.Room != nil {
				Te = pipe.Room.Tot
			} else {
				Err := fmt.Sprintf("Undefined Pipe Environment  name=%s", pipe.Name)
				Eprint("<Pipecfv>", Err)
			}
			pipe.Ko = pipe.Cat.Ko

			Eo1 := pipe.Cmp.Elouts[0]
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			pipe.Ep = 1.0 - math.Exp(-(pipe.Ko*pipe.L)/cG)
			pipe.D1 = cG * pipe.Ep
			pipe.Do = pipe.D1 * Te
			Eo1.Coeffo = cG
			Eo1.Co = pipe.Do
			Eo1.Coeffin[0] = pipe.D1 - cG

			if pipe.Cat.Type == DUCT_PDT {
				Eo2 := pipe.Cmp.Elouts[1]
				Eo2.Coeffo = 1.0
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -1.0
			}
		}
	}
}

/* --------------------------- */

/* 取得熱量の計算 */

func Pipeene(Pipe []*PIPE) {
	for _, pipe := range Pipe {
		pipe.Tin = pipe.Cmp.Elins[0].Sysvin

		if pipe.Cmp.Control != OFF_SW {
			Eo := pipe.Cmp.Elouts[0]
			pipe.Tout = pipe.Do
			pipe.Q = pipe.Do - pipe.D1*pipe.Tin

			if pipe.Room != nil {
				pipe.Room.Qeqp += (-pipe.Q)
			}

			if pipe.Cat.Type == DUCT_PDT {
				Eo = pipe.Cmp.Elouts[1]
				pipe.Xout = Eo.Sysv
				pipe.RHout = FNRhtx(pipe.Tout, pipe.Xout)
				pipe.Hout = FNH(pipe.Tout, Eo.Sysv)
			} else {
				pipe.Hout = FNAN
			}
		} else {
			pipe.Q = 0.0
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

func pipeprint(fo io.Writer, id int, Pipe []*PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, " %s 1 5\n", pipe.Name)
		}
	case 1:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f %s_Q q f\n",
				pipe.Name, pipe.Name, pipe.Name, pipe.Name, pipe.Name)
		}
	default:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %.2f\n",
				pipe.Cmp.Elouts[0].Control, pipe.Cmp.Elouts[0].G,
				pipe.Tin, pipe.Cmp.Elouts[0].Sysv, pipe.Q)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func pipedyint(Pipe []*PIPE) {
	for _, pipe := range Pipe {
		svdyint(&pipe.Tidy)
		qdyint(&pipe.Qdy)
	}
}

func pipemonint(Pipe []*PIPE) {
	for _, pipe := range Pipe {
		svdyint(&pipe.MTidy)
		qdyint(&pipe.MQdy)
	}
}

func pipeday(Mon int, Day int, ttmm int, Pipe []*PIPE, Nday int, SimDayend int) {
	for _, pipe := range Pipe {
		// 日集計
		svdaysum(int64(ttmm), pipe.Cmp.Elouts[0].Control, pipe.Tin, &pipe.Tidy)
		qdaysum(int64(ttmm), pipe.Cmp.Elouts[0].Control, pipe.Q, &pipe.Qdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, pipe.Cmp.Elouts[0].Control, pipe.Tin, &pipe.MTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, pipe.Cmp.Elouts[0].Control, pipe.Q, &pipe.MQdy, Nday, SimDayend)
	}
}

func pipedyprt(fo io.Writer, id int, Pipe []*PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, " %s 1 14\n", pipe.Name)
		}

	case 1:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
		}

	default:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				pipe.Tidy.Hrs, pipe.Tidy.M, pipe.Tidy.Mntime,
				pipe.Tidy.Mn, pipe.Tidy.Mxtime, pipe.Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", pipe.Qdy.Hhr, pipe.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", pipe.Qdy.Chr, pipe.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", pipe.Qdy.Hmxtime, pipe.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", pipe.Qdy.Cmxtime, pipe.Qdy.Cmx)
		}
	}
}

func pipemonprt(fo io.Writer, id int, Pipe []*PIPE) {
	switch id {
	case 0:
		if len(Pipe) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PIPEDUCT_TYPE, len(Pipe))
		}
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, " %s 1 14\n", pipe.Name)
		}

	case 1:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n\n", pipe.Name, pipe.Name, pipe.Name, pipe.Name)
		}

	default:
		for _, pipe := range Pipe {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				pipe.MTidy.Hrs, pipe.MTidy.M, pipe.MTidy.Mntime,
				pipe.MTidy.Mn, pipe.MTidy.Mxtime, pipe.MTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", pipe.MQdy.Hhr, pipe.MQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", pipe.MQdy.Chr, pipe.MQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", pipe.MQdy.Hmxtime, pipe.MQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", pipe.MQdy.Cmxtime, pipe.MQdy.Cmx)
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
