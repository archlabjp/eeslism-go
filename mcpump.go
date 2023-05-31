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

/*  pump.c  */

/*  ポンプ   */

package main

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/* 機器仕様入力　　　　　　*/

func Pumpdata(cattype string, s string, Pumpca *PUMPCA, pfcmp []PFCMP) int {
	st := strings.IndexByte(s, '=')
	var dt float64
	var id int
	var i int

	if st == -1 {
		Pumpca.name = s
		Pumpca.Type = ""
		Pumpca.Wo = -999.0
		Pumpca.Go = -999.0
		Pumpca.qef = -999.0
		Pumpca.val = nil

		if cattype == PUMP_TYPE {
			Pumpca.pftype = PUMP_PF
		} else if cattype == FAN_TYPE {
			Pumpca.pftype = FAN_PF
		} else {
			Pumpca.pftype = rune(OFF_SW)
		}
	} else {
		s1, s2 := s[:st], s[st+1:]

		if s1 == "type" {
			Pumpca.Type = s2
			if Pumpca.Type == "P" {
				Pumpca.val = make([]float64, 4)
			}

			for i = 0; i < len(pfcmp); i++ {
				pfc := &pfcmp[i]
				if Pumpca.pftype == pfc.pftype && Pumpca.Type == pfc.Type {
					Pumpca.pfcmp = pfc
					break
				}
			}
		} else {
			dt, _ = strconv.ParseFloat(s[st+1:], 64)
			if s == "qef" {
				Pumpca.qef = dt
			} else {
				if Pumpca.Type != "P" {
					switch s {
					case "Go":
						Pumpca.Go = dt
					case "Wo":
						Pumpca.Wo = dt
					default:
						id = 1
					}
				} else {
					switch s {
					case "a0":
						Pumpca.val[0] = dt
					case "a1":
						Pumpca.val[1] = dt
					case "a2":
						Pumpca.val[2] = dt
					case "Ic":
						Pumpca.val[3] = dt
					default:
						id = 1
					}
				}
			}
		}
	}
	return id
}

/* --------------------------- */

/* 太陽電池ポンプの太陽電池パネルの方位設定　*/

func Pumpint(Npump int, Pump []PUMP, Nexsf int, Exs []EXSF) {
	for i := 0; i < Npump; i++ {
		p := &Pump[i]
		if p.Cat.Type == "P" {
			p.Sol = nil
			for j := 0; j < Nexsf; j++ {
				if p.Cmp.Exsname == Exs[j].Name {
					p.Sol = &Exs[j]
					break
				}
			}
			if p.Sol == nil {
				Eprint("Pumpint", p.Cmp.Exsname)
			}
		}
	}
}

/* --------------------------- */

/* ポンプ流量設定（太陽電池ポンプのみ） */

func Pumpflow(Npump int, Pump []PUMP) {
	for i := 0; i < Npump; i++ {
		p := &Pump[i]
		if p.Cat.Type == "P" {
			S := p.Sol.Iw

			if DEBUG {
				fmt.Printf("<Pumpflow> i=%d S=%f Ic=%f a0=%f a1=%e a2=%e\n",
					i, S, p.Cat.val[3], p.Cat.val[0],
					p.Cat.val[1], p.Cat.val[2])
			}

			if S > p.Cat.val[3] {
				p.G = p.Cat.val[0] + (p.Cat.val[1]+p.Cat.val[2]*S)*S
			} else {
				p.G = -999.0
			}

			p.E = 0
		} else {
			if p.Cmp.Control != OFF_SW {
				p.G = p.Cat.Go
				p.E = p.Cat.Wo
			} else {
				p.G = -999.0
				p.E = 0.0
			}

			if DEBUG {
				fmt.Printf("<Pumpflow>  control=%c G=%f E=%f\n",
					p.Cmp.Control, p.G, p.E)
			}
		}
	}
}

/* --------------------------- */

/*  特性式の係数  */

func Pumpcfv(Npump int, Pump []PUMP) {
	for i := 0; i < Npump; i++ {
		p := &Pump[i]
		if p.Cmp.Control != OFF_SW {
			Eo := p.Cmp.Elouts[0]
			cG := Spcheat(Eo.Fluid) * Eo.G
			p.CG = cG
			Eo.Coeffo = cG
			p.PLC = PumpFanPLC(Eo.G/p.G, p)
			Eo.Co = p.Cat.qef * p.E * p.PLC
			Eo.Coeffin[0] = -cG

			if p.Cat.pftype == FAN_PF {
				Eo = p.Cmp.Elouts[1]
				Eo.Coeffo = Eo.G
				Eo.Co = 0.0
				Eo.Coeffin[0] = -Eo.G
			}
		} else {
			p.G = 0.0
			p.E = 0.0
		}
	}
}

// -----------------------------------------
//
//  ポンプ、ファンの部分負荷特性曲線
//

func PumpFanPLC(XQ float64, Pump *PUMP) float64 {
	var Buff, dQ float64
	var i int
	cat := Pump.Cat

	dQ = math.Min(1.0, math.Max(XQ, 0.25))

	if cat.pfcmp == nil {
		Err := fmt.Sprintf("<PumpFanPLC>  PFtype=%c  type=%s", cat.pftype, cat.Type)
		Eprint("PUMP oir FAN", string(Err[:]))
		Buff = 0.0
	} else {
		Buff = 0.0

		for i = 0; i < 5; i++ {
			Buff += cat.pfcmp.dblcoeff[i] * math.Pow(dQ, float64(i))
		}
	}
	return Buff
}

/* --------------------------- */

/*  供給熱量、エネルギーの計算 */

func Pumpene(Npump int, Pump []PUMP) {
	for i := 0; i < Npump; i++ {
		Pump[i].Tin = Pump[i].Cmp.Elins[0].Sysvin
		Eo := Pump[i].Cmp.Elouts[0]

		if Eo.Control != OFF_SW {
			Pump[i].Q = Pump[i].CG * (Eo.Sysv - Pump[i].Tin)
		} else {
			Pump[i].Q = 0.0
		}
	}
}

/* --------------------------- */

func pumpprint(fo io.Writer, id int, Npump int, Pump []PUMP) {
	var G float64

	switch id {
	case 0:
		if Npump > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, Npump)
		}
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, " %s 1 6\n", Pump[i].Name)
		}
	case 1:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%s_c c c %s_Ti t f %s_To t f ", Pump[i].Name, Pump[i].Name, Pump[i].Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_G m f\n", Pump[i].Name, Pump[i].Name, Pump[i].Name)
		}
	default:
		for i := 0; i < Npump; i++ {
			if Pump[i].Cmp.Elouts[0].G > 0.0 && Pump[i].Cmp.Elouts[0].Control != OFF_SW {
				G = Pump[i].Cmp.Elouts[0].G
			} else {
				G = 0.0
			}
			fmt.Fprintf(fo, "%c %4.1f %4.1f %4.0f %4.0f %.5g\n", Pump[i].Cmp.Elouts[0].Control,
				Pump[i].Tin, Pump[i].Cmp.Elouts[0].Sysv, Pump[i].Q, Pump[i].E*Pump[i].PLC, G)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func pumpdyint(Npump int, Pump []PUMP) {
	for i := 0; i < Npump; i++ {
		edyint(&Pump[i].Qdy)
		edyint(&Pump[i].Edy)
		edyint(&Pump[i].Gdy)
	}
}

func pumpmonint(Npump int, Pump []PUMP) {
	for i := 0; i < Npump; i++ {
		edyint(&Pump[i].MQdy)
		edyint(&Pump[i].MEdy)
		edyint(&Pump[i].MGdy)
	}
}

func pumpday(Mon, Day, ttmm, Npump int, Pump []PUMP, Nday, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for i := 0; i < Npump; i++ {
		// 日集計
		edaysum(ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].Q, &Pump[i].Qdy)
		edaysum(ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].E, &Pump[i].Edy)
		edaysum(ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].G, &Pump[i].Gdy)

		// 月集計
		emonsum(Mon, Day, ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].Q, &Pump[i].MQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].E, &Pump[i].MEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].G, &Pump[i].MGdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, Pump[i].Cmp.Elouts[0].Control, Pump[i].E, &Pump[i].MtEdy[Mo][tt])
	}
}

func pumpdyprt(fo io.Writer, id, Npump int, Pump []PUMP) {
	switch id {
	case 0:
		if Npump > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, Npump)
		}
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, " %s 1 12\n", Pump[i].Name)
		}
	case 1:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
		}
	default:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].Qdy.Hrs, Pump[i].Qdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f ", Pump[i].Qdy.Mxtime, Pump[i].Qdy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].Edy.Hrs, Pump[i].Edy.D)
			fmt.Fprintf(fo, "%1ld %2.0f ", Pump[i].Edy.Mxtime, Pump[i].Edy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].Gdy.Hrs, Pump[i].Gdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f\n", Pump[i].Gdy.Mxtime, Pump[i].Gdy.Mx)
		}
	}
}

func pumpmonprt(fo io.Writer, id, Npump int, Pump []PUMP) {
	switch id {
	case 0:
		if Npump > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, Npump)
		}
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, " %s 1 12\n", Pump[i].Name)
		}
	case 1:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				Pump[i].Name, Pump[i].Name, Pump[i].Name, Pump[i].Name)
		}
	default:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].MQdy.Hrs, Pump[i].MQdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f ", Pump[i].MQdy.Mxtime, Pump[i].MQdy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].MEdy.Hrs, Pump[i].MEdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f ", Pump[i].MEdy.Mxtime, Pump[i].MEdy.Mx)
			fmt.Fprintf(fo, "%1ld %3.1f ", Pump[i].MGdy.Hrs, Pump[i].MGdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f\n", Pump[i].MGdy.Mxtime, Pump[i].MGdy.Mx)
		}
	}
}
func pumpmtprt(fo io.Writer, id, Npump int, Pump []PUMP, Mo, tt int) {
	switch id {
	case 0:
		if Npump > 0 {
			fmt.Fprintf(fo, "%s %d\n", PUMP_TYPE, Npump)
		}
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, " %s 1 1\n", Pump[i].Name)
		}
	case 1:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, "%s_E E f \n", Pump[i].Name)
		}
	default:
		for i := 0; i < Npump; i++ {
			fmt.Fprintf(fo, " %.2f \n", Pump[i].MtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

func PFcmpInit(Npfcmp int, Pfcmp []PFCMP) {
	for i := 0; i < Npfcmp; i++ {
		Pfcmp[i].pftype = ' '
		Pfcmp[i].Type = ""
		matinit(Pfcmp[i].dblcoeff[:], 5)
	}
}

func PFcmpdata(fl io.Reader, Pfcmp *[]PFCMP) {
	var s string
	var c byte
	var i int

	for {
		_, err := fmt.Fscanf(fl, "%s", &s)
		if err != nil || s[0] == '*' {
			break
		}

		if s == "!" {
			for {
				_, err = fmt.Fscanf(fl, "%c", &c)
				if err != nil || c == '\n' {
					break
				}
			}
		} else {
			pfcmp := PFCMP{}

			if s == PUMP_TYPE {
				pfcmp.pftype = PUMP_PF
			} else if s == FAN_TYPE {
				pfcmp.pftype = FAN_PF
			} else {
				Eprint("<pumpfanlst.efl>", s)
			}

			_, err = fmt.Fscanf(fl, "%s", &s)
			if err != nil {
				break
			}

			pfcmp.Type = s

			i = 0
			for {
				_, err = fmt.Fscanf(fl, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}

				var err error
				pfcmp.dblcoeff[i], err = strconv.ParseFloat(s, 64)
				if err != nil {
					panic(err)
				}
				i++
			}

			*Pfcmp = append(*Pfcmp, pfcmp)
		}
	}
}
