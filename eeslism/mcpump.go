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

package eeslism

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

/* 機器仕様入力　　　　　　*/

func Pumpdata(cattype EqpType, s string, Pumpca *PUMPCA, pfcmp []*PFCMP) int {
	st := strings.IndexByte(s, '=')
	var dt float64
	var id int

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

			for _, pfc := range pfcmp {
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

func Pumpint(Pump []*PUMP, Exs []*EXSF) {
	for _, p := range Pump {
		if p.Cat.Type == "P" {
			p.Sol = nil
			for j := 0; j < len(Exs); j++ {
				if p.Cmp.Exsname == Exs[j].Name {
					p.Sol = Exs[j]
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

// ポンプ流量設定（太陽電池ポンプのみ
func (eqsys *EQSYS) Pumpflow() {
	for i, p := range eqsys.Pump {
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

//
//  +------+ ---> [OUT 1] 空気 or 温水温度 ?
//  | PUMP |
//  +------+ ---> [OUT 2] 湿度? (FAN_PFのみ)
//
func Pumpcfv(Pump []*PUMP) {
	for _, p := range Pump {
		if p.Cmp.Control != OFF_SW {
			Eo1 := p.Cmp.Elouts[0]
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			p.CG = cG
			Eo1.Coeffo = cG
			p.PLC = PumpFanPLC(Eo1.G/p.G, p)
			Eo1.Co = p.Cat.qef * p.E * p.PLC
			Eo1.Coeffin[0] = -cG

			if p.Cat.pftype == FAN_PF {
				Eo2 := p.Cmp.Elouts[1]
				Eo2.Coeffo = Eo2.G
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -Eo2.G
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

func Pumpene(Pump []*PUMP) {
	for _, p := range Pump {
		p.Tin = p.Cmp.Elins[0].Sysvin
		Eo := p.Cmp.Elouts[0]

		if Eo.Control != OFF_SW {
			p.Q = p.CG * (Eo.Sysv - p.Tin)
		} else {
			p.Q = 0.0
		}
	}
}

/* --------------------------- */

func pumpprint(fo io.Writer, id int, Pump []*PUMP) {
	var G float64

	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 6\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_c c c %s_Ti t f %s_To t f ", p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_G m f\n", p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			if p.Cmp.Elouts[0].G > 0.0 && p.Cmp.Elouts[0].Control != OFF_SW {
				G = p.Cmp.Elouts[0].G
			} else {
				G = 0.0
			}
			fmt.Fprintf(fo, "%c %4.1f %4.1f %4.0f %4.0f %.5g\n", p.Cmp.Elouts[0].Control,
				p.Tin, p.Cmp.Elouts[0].Sysv, p.Q, p.E*p.PLC, G)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func pumpdyint(Pump []*PUMP) {
	for _, p := range Pump {
		edyint(&p.Qdy)
		edyint(&p.Edy)
		edyint(&p.Gdy)
	}
}

func pumpmonint(Pump []*PUMP) {
	for _, p := range Pump {
		edyint(&p.MQdy)
		edyint(&p.MEdy)
		edyint(&p.MGdy)
	}
}

func pumpday(Mon, Day, ttmm int, Pump []*PUMP, Nday, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for _, p := range Pump {
		// 日集計
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.Q, &p.Qdy)
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.E, &p.Edy)
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.G, &p.Gdy)

		// 月集計
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.Q, &p.MQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.E, &p.MEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.G, &p.MGdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.E, &p.MtEdy[Mo][tt])
	}
}

func pumpdyprt(fo io.Writer, id int, Pump []*PUMP) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 12\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				p.Name, p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%1d %3.1f ", p.Qdy.Hrs, p.Qdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.Qdy.Mxtime, p.Qdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.Edy.Hrs, p.Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.Edy.Mxtime, p.Edy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.Gdy.Hrs, p.Gdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", p.Gdy.Mxtime, p.Gdy.Mx)
		}
	}
}

func pumpmonprt(fo io.Writer, id int, Pump []*PUMP) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 12\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				p.Name, p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%1d %3.1f ", p.MQdy.Hrs, p.MQdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.MQdy.Mxtime, p.MQdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.MEdy.Hrs, p.MEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.MEdy.Mxtime, p.MEdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.MGdy.Hrs, p.MGdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", p.MGdy.Mxtime, p.MGdy.Mx)
		}
	}
}
func pumpmtprt(fo io.Writer, id int, Pump []*PUMP, Mo, tt int) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 1\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_E E f \n", p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, " %.2f \n", p.MtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

func NewPFCMP() *PFCMP {
	Pfcmp := new(PFCMP)
	Pfcmp.pftype = ' '
	Pfcmp.Type = ""
	matinit(Pfcmp.dblcoeff[:], 5)
	return Pfcmp
}

func PFcmpdata() []*PFCMP {
	var s string
	var c byte
	var i int

	fl, err := os.Open("pumpfanlst.efl")
	if err != nil {
		Eprint(" file ", "pumpfanlst.efl")
	}
	Pfcmp := make([]*PFCMP, 0)

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
			pfcmp := NewPFCMP()

			if s == string(PUMP_TYPE) {
				pfcmp.pftype = PUMP_PF
			} else if s == string(FAN_TYPE) {
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

			Pfcmp = append(Pfcmp, pfcmp)
		}
	}

	fl.Close()

	return Pfcmp
}
