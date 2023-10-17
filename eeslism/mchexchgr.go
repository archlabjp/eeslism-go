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

/* hexchgr.c */

package eeslism

// -------------------------------------------------------------
// 熱交換器
//
// 冷風入力 [IN  1] ---> +-----+ <--- [IN  2] 温風入力
//                       | HEX |
// 冷風出力 [OUT 1] <--- +-----+ ---> [OUT 2] 温風出力
//
// -------------------------------------------------------------

/*  仕様入力  */

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func Hexdata(s string, Hexca *HEXCA) int {
	st := strings.IndexByte(s, '=')
	if st == -1 {
		Hexca.Name = s
		Hexca.eff = -999.0
		Hexca.KA = -999.0
	} else {
		s1 := s[:st]
		s2 := s[st+1:]
		if s1 == "eff" {
			e, err := strconv.ParseFloat(s2, 64)
			if err != nil {
				return 1
			}
			Hexca.eff = e
		} else if s1 == "KA" {
			ka, err := strconv.ParseFloat(s2, 64)
			if err != nil {
				return 1
			}
			Hexca.KA = ka
		} else {
			return 1
		}
	}
	return 0
}

/* --------------------------- */

/*  特性式の係数  */

func Hexcfv(Hex []HEX) {
	for i := range Hex {
		hex := &Hex[i]

		// 計算準備
		if hex.Id == 0 {
			/* 温度効率固定タイプと変動タイプの判定 */
			if hex.Cat.eff > 0.0 {
				hex.Etype = 'e'
			} else if hex.Cat.KA > 0.0 {
				hex.Etype = 'k'
			} else {
				fmt.Printf("Hex %s  Undefined Character eff or KA\n", hex.Name)
				hex.Etype = '-'
			}

			hex.Id = 1
		}

		if hex.Cmp.Control != OFF_SW {
			hex.Eff = hex.Cat.eff

			if hex.Eff < 0.0 {
				errMsg := fmt.Sprintf("Name=%s  eff=%.4g", hex.Cmp.Name, hex.Eff)
				Eprint("Hexcfv", errMsg)
			}

			eoh := hex.Cmp.Elouts[1]
			eoc := hex.Cmp.Elouts[0]
			hex.CGc = Spcheat(eoc.Fluid) * eoc.G
			hex.CGh = Spcheat(eoh.Fluid) * eoh.G

			if hex.Etype == 'k' {
				hex.Eff = FNhccet(hex.CGc, hex.CGh, hex.Cat.KA)
			}

			eCGmin := hex.Eff * math.Min(hex.CGc, hex.CGh)
			hex.ECGmin = eCGmin
			eoc.Coeffin[0] = -hex.CGc + eCGmin
			eoc.Coeffin[1] = -eCGmin
			eoc.Coeffo = hex.CGc
			eoc.Co = 0.0

			eoh.Coeffin[0] = -eCGmin
			eoh.Coeffin[1] = -hex.CGh + eCGmin
			eoh.Coeffo = hex.CGh
			eoh.Co = 0.0
		}
	}
}

/* --------------------------- */

/* 交換熱量の計算 */

func Hexene(Hex []HEX) {
	for i := range Hex {
		hex := &Hex[i]

		// 流入
		hex.Tcin = hex.Cmp.Elins[0].Sysvin
		hex.Thin = hex.Cmp.Elins[1].Sysvin

		if hex.Cmp.Control != OFF_SW {
			// 流出
			hex.Qci = hex.CGc * (hex.Cmp.Elouts[0].Sysv - hex.Tcin)
			hex.Qhi = hex.CGh * (hex.Cmp.Elouts[1].Sysv - hex.Thin)
		} else {
			hex.Qci = 0.0
			hex.Qhi = 0.0
		}
	}
}

/* --------------------------- */

func hexprint(fo io.Writer, id int, Hex []HEX) {
	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for i := range Hex {
			fmt.Fprintf(fo, " %s 1 9\n", Hex[i].Name)
		}
	case 1:
		for i := range Hex {
			fmt.Fprintf(fo, "%s_c c c %s:c_G m f %s:c_Ti t f %s:c_To t f %s:c_Q q f\n",
				Hex[i].Name, Hex[i].Name, Hex[i].Name, Hex[i].Name, Hex[i].Name)
			fmt.Fprintf(fo, "%s:h_G m f %s:h_Ti t f %s:h_To t f %s:h_Q q f\n",
				Hex[i].Name, Hex[i].Name, Hex[i].Name, Hex[i].Name)
		}
	default:
		for i := range Hex {
			Eo := Hex[i].Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f", Hex[i].Cmp.Control, Eo.G, Hex[i].Tcin, Eo.Sysv, Hex[i].Qci)
			Eo = Hex[i].Cmp.Elouts[1]
			fmt.Fprintf(fo, " %6.4g %4.1f %4.1f %2.0f\n", Eo.G, Hex[i].Thin, Eo.Sysv, Hex[i].Qhi)
		}
	}
}

/* 日積算値に関する処理 */

func hexdyint(Hex []HEX) {
	for i := range Hex {
		svdyint(&Hex[i].Tcidy)
		svdyint(&Hex[i].Thidy)
		qdyint(&Hex[i].Qcidy)
		qdyint(&Hex[i].Qhidy)
	}
}

func hexmonint(Hex []HEX) {
	for i := range Hex {
		svdyint(&Hex[i].MTcidy)
		svdyint(&Hex[i].MThidy)
		qdyint(&Hex[i].MQcidy)
		qdyint(&Hex[i].MQhidy)
	}
}

func hexday(Mon, Day, ttmm int, Hex []HEX, Nday, SimDayend int) {
	for i := range Hex {
		// 日集計
		svdaysum(int64(ttmm), Hex[i].Cmp.Control, Hex[i].Tcin, &Hex[i].Tcidy)
		svdaysum(int64(ttmm), Hex[i].Cmp.Control, Hex[i].Thin, &Hex[i].Thidy)
		qdaysum(int64(ttmm), Hex[i].Cmp.Control, Hex[i].Qci, &Hex[i].Qcidy)
		qdaysum(int64(ttmm), Hex[i].Cmp.Control, Hex[i].Qhi, &Hex[i].Qhidy)

		// 月集計
		svmonsum(Mon, Day, ttmm, Hex[i].Cmp.Control, Hex[i].Tcin, &Hex[i].MTcidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, Hex[i].Cmp.Control, Hex[i].Thin, &Hex[i].MThidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hex[i].Cmp.Control, Hex[i].Qci, &Hex[i].MQcidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Hex[i].Cmp.Control, Hex[i].Qhi, &Hex[i].MQhidy, Nday, SimDayend)
	}
}

func hexdyprt(fo io.Writer, id int, Hex []HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for i := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", Hex[i].Name)
		}
	case 1:
		for i := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
			}
		}
	default:
		for i := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hex[i].Tcidy.Hrs, Hex[i].Tcidy.M,
				Hex[i].Tcidy.Mntime, Hex[i].Tcidy.Mn,
				Hex[i].Tcidy.Mxtime, Hex[i].Tcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].Qcidy.Hhr, Hex[i].Qcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].Qcidy.Chr, Hex[i].Qcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].Qcidy.Hmxtime, Hex[i].Qcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].Qcidy.Cmxtime, Hex[i].Qcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hex[i].Thidy.Hrs, Hex[i].Thidy.M,
				Hex[i].Thidy.Mntime, Hex[i].Thidy.Mn,
				Hex[i].Thidy.Mxtime, Hex[i].Thidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].Qhidy.Hhr, Hex[i].Qhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].Qhidy.Chr, Hex[i].Qhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].Qhidy.Hmxtime, Hex[i].Qhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hex[i].Qhidy.Cmxtime, Hex[i].Qhidy.Cmx)
		}
	}
}

func hexmonprt(fo io.Writer, id int, Hex []HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for i := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", Hex[i].Name)
		}
	case 1:
		for i := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c, Hex[i].Name, c)
			}
		}
	default:
		for i := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hex[i].MTcidy.Hrs, Hex[i].MTcidy.M,
				Hex[i].MTcidy.Mntime, Hex[i].MTcidy.Mn,
				Hex[i].MTcidy.Mxtime, Hex[i].MTcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].MQcidy.Hhr, Hex[i].MQcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].MQcidy.Chr, Hex[i].MQcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].MQcidy.Hmxtime, Hex[i].MQcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].MQcidy.Cmxtime, Hex[i].MQcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Hex[i].MThidy.Hrs, Hex[i].MThidy.M,
				Hex[i].MThidy.Mntime, Hex[i].MThidy.Mn,
				Hex[i].MThidy.Mxtime, Hex[i].MThidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].MQhidy.Hhr, Hex[i].MQhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Hex[i].MQhidy.Chr, Hex[i].MQhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Hex[i].MQhidy.Hmxtime, Hex[i].MQhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Hex[i].MQhidy.Cmxtime, Hex[i].MQhidy.Cmx)
		}
	}
}
