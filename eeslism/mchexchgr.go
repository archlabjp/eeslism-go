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

func Hexcfv(Hex []*HEX) {
	for _, hex := range Hex {

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

func Hexene(Hex []*HEX) {
	for _, hex := range Hex {

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

func hexprint(fo io.Writer, id int, Hex []*HEX) {
	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 9\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%s_c c c %s:c_G m f %s:c_Ti t f %s:c_To t f %s:c_Q q f\n",
				hex.Name, hex.Name, hex.Name, hex.Name, hex.Name)
			fmt.Fprintf(fo, "%s:h_G m f %s:h_Ti t f %s:h_To t f %s:h_Q q f\n",
				hex.Name, hex.Name, hex.Name, hex.Name)
		}
	default:
		for _, hex := range Hex {
			eo_Tc := hex.Cmp.Elouts[0]
			eo_Th := hex.Cmp.Elouts[1]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f", hex.Cmp.Control, eo_Tc.G, hex.Tcin, eo_Tc.Sysv, hex.Qci)
			fmt.Fprintf(fo, " %6.4g %4.1f %4.1f %2.0f\n", eo_Th.G, hex.Thin, eo_Th.Sysv, hex.Qhi)
		}
	}
}

/* 日積算値に関する処理 */

func hexdyint(Hex []*HEX) {
	for _, hex := range Hex {
		svdyint(&hex.Tcidy)
		svdyint(&hex.Thidy)
		qdyint(&hex.Qcidy)
		qdyint(&hex.Qhidy)
	}
}

func hexmonint(Hex []*HEX) {
	for _, hex := range Hex {
		svdyint(&hex.MTcidy)
		svdyint(&hex.MThidy)
		qdyint(&hex.MQcidy)
		qdyint(&hex.MQhidy)
	}
}

func hexday(Mon, Day, ttmm int, Hex []*HEX, Nday, SimDayend int) {
	for _, hex := range Hex {
		// 日集計
		svdaysum(int64(ttmm), hex.Cmp.Control, hex.Tcin, &hex.Tcidy)
		svdaysum(int64(ttmm), hex.Cmp.Control, hex.Thin, &hex.Thidy)
		qdaysum(int64(ttmm), hex.Cmp.Control, hex.Qci, &hex.Qcidy)
		qdaysum(int64(ttmm), hex.Cmp.Control, hex.Qhi, &hex.Qhidy)

		// 月集計
		svmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Tcin, &hex.MTcidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Thin, &hex.MThidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Qci, &hex.MQcidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Qhi, &hex.MQhidy, Nday, SimDayend)
	}
}

func hexdyprt(fo io.Writer, id int, Hex []*HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
			}
		}
	default:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.Tcidy.Hrs, hex.Tcidy.M,
				hex.Tcidy.Mntime, hex.Tcidy.Mn,
				hex.Tcidy.Mxtime, hex.Tcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qcidy.Hhr, hex.Qcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qcidy.Chr, hex.Qcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qcidy.Hmxtime, hex.Qcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qcidy.Cmxtime, hex.Qcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.Thidy.Hrs, hex.Thidy.M,
				hex.Thidy.Mntime, hex.Thidy.Mn,
				hex.Thidy.Mxtime, hex.Thidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qhidy.Hhr, hex.Qhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qhidy.Chr, hex.Qhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qhidy.Hmxtime, hex.Qhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hex.Qhidy.Cmxtime, hex.Qhidy.Cmx)
		}
	}
}

func hexmonprt(fo io.Writer, id int, Hex []*HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
			}
		}
	default:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.MTcidy.Hrs, hex.MTcidy.M,
				hex.MTcidy.Mntime, hex.MTcidy.Mn,
				hex.MTcidy.Mxtime, hex.MTcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQcidy.Hhr, hex.MQcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQcidy.Chr, hex.MQcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQcidy.Hmxtime, hex.MQcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQcidy.Cmxtime, hex.MQcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.MThidy.Hrs, hex.MThidy.M,
				hex.MThidy.Mntime, hex.MThidy.Mn,
				hex.MThidy.Mxtime, hex.MThidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQhidy.Hhr, hex.MQhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQhidy.Chr, hex.MQhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQhidy.Hmxtime, hex.MQhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hex.MQhidy.Cmxtime, hex.MQhidy.Cmx)
		}
	}
}
