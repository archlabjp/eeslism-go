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

/*　太陽光発電システム

機器仕様入力　　　　*/

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func PVcadata(s string, PVca *PVCA) int {
	dt := 0.0
	id := 0

	var st int
	if st = strings.IndexRune(s, '='); st == -1 {
		PVca.Name = s
		PVca.PVcap = -999.0
		PVca.Area = -999.0
		PVca.KHD = 1.0
		PVca.KPD = 0.95
		PVca.KPM = 0.94
		PVca.KPA = 0.97
		PVca.effINO = 0.9
		PVca.A = -999.0
		PVca.B = -999.0
		PVca.apmax = -0.41
	} else {
		s1 := s[:st]
		s2 := s[st+1:]

		var err error
		dt, err = strconv.ParseFloat(s2, 64)
		if err != nil {
			panic(err)
		}

		switch {
		case strings.HasPrefix(s1, "KHD"):
			// 日射量年変動補正係数
			PVca.KHD = dt
		case strings.HasPrefix(s1, "KPD"):
			// 経時変化補正係数
			PVca.KPD = dt
		case strings.HasPrefix(s1, "KPM"):
			// アレイ負荷整合補正係数
			PVca.KPM = dt
		case strings.HasPrefix(s1, "KPA"):
			// アレイ回路補正係数
			PVca.KPA = dt
		case strings.HasPrefix(s1, "EffInv"):
			// インバータ実行効率
			PVca.effINO = dt
		case strings.HasPrefix(s1, "apmax"):
			// 最大出力温度係数
			PVca.apmax = dt
		case s1 == "InstallType":
			PVca.InstallType = rune(s2[0])
			switch PVca.InstallType {
			case 'A':
				PVca.A = 46.0
				PVca.B = 0.41
			case 'B':
				PVca.A = 50.0
				PVca.B = 0.38
			case 'C':
				PVca.A = 57.0
				PVca.B = 0.33
			}
		case strings.HasPrefix(s1, "PVcap"):
			// 太陽電池容量
			PVca.PVcap = dt
		case strings.HasPrefix(s1, "Area"):
			// 太陽電池容量
			PVca.Area = dt
		default:
			id = 1
		}
	}
	return id
}

/* ------------------------------------- */

/*  初期設定 */

func PVint(PV []PV, Nexsf int, Exs []EXSF, Wd *WDAT) {
	Err := ""

	for i := range PV {
		PV[i].Ta = &Wd.T
		PV[i].V = &Wd.Wv

		for j := 0; j < Nexsf; j++ {
			exs := &Exs[j]
			if PV[i].Cmp.Exsname == exs.Name {
				PV[i].Sol = exs
				PV[i].I = &PV[i].Sol.Iw
				break
			}
		}

		if PV[i].Sol == nil {
			Eprint("PVint", PV[i].Cmp.Exsname)
		}

		if PV[i].Cat.KHD < 0.0 {
			Err = fmt.Sprintf("Name=%s KHD=%.4g", PV[i].Cmp.Name, PV[i].Cat.KHD)
			Eprint("PVint", Err)
		}

		if PV[i].Cat.KPD < 0.0 {
			Err = fmt.Sprintf("Name=%s KHD=%.4g", PV[i].Cmp.Name, PV[i].Cat.KPD)
			Eprint("PVint", Err)
		}

		if PV[i].Cat.KPM < 0.0 {
			Err = fmt.Sprintf("Name=%s KPM=%.4g", PV[i].Cmp.Name, PV[i].Cat.KPM)
			Eprint("PVint", Err)
		}

		if PV[i].Cat.KPA < 0.0 {
			Err = fmt.Sprintf("Name=%s KPA=%.4g", PV[i].Cmp.Name, PV[i].Cat.KPA)
			Eprint("PVint", Err)
		}

		if PV[i].Cat.effINO < 0.0 {
			Err = fmt.Sprintf("Name=%s EffInv=%.4g", PV[i].Cmp.Name, PV[i].Cat.effINO)
			Eprint("PVint", Err)
		}

		if PV[i].Cat.apmax > 0.0 {
			Err = fmt.Sprintf("Name=%s apmax=%.4g", PV[i].Cmp.Name, PV[i].Cat.apmax)
			Eprint("PVint", Err)
		}

		if PV[i].PVcap < 0.0 {
			Err = fmt.Sprintf("Name=%s PVcap=%.4g", PV[i].Cmp.Name, PV[i].PVcap)
			Eprint("PVint", Err)
		}

		if PV[i].Area < 0.0 {
			Err = fmt.Sprintf("Name=%s Area=%.4g", PV[i].Cmp.Name, PV[i].Area)
			Eprint("PVint", Err)
		}

		// 計算途中で変化しない各種補正係数の積
		PV[i].KConst = PV[i].Cat.KHD * PV[i].Cat.KPD * PV[i].Cat.KPM * PV[i].Cat.KPA * PV[i].Cat.effINO
	}
}

// 太陽光発電の発電量計算
/* ------------------------------------- */

/*  集熱量の計算 */

func PVene(PV []PV) {
	for i := range PV {
		// 太陽電池アレイの計算（JIS C 8907:2005　P21による）
		PV[i].TPV = *PV[i].Ta + (PV[i].Cat.A/(PV[i].Cat.B*math.Pow(*PV[i].V, 0.8)+1.0)+2.0)**PV[i].I/1000.0 - 2.0
		PV[i].KPT = FNKPT(PV[i].TPV, PV[i].Cat.apmax)
		PV[i].KTotal = PV[i].KConst * PV[i].KPT

		// 太陽電池入社日射量の計算
		PV[i].Iarea = *PV[i].I * PV[i].Area

		// 発電量の計算
		PV[i].Power = PV[i].KTotal * *PV[i].I / 1000.0 * PV[i].PVcap

		// 発電効率の計算
		PV[i].Eff = 0.0
		if PV[i].Iarea > 0.0 {
			PV[i].Eff = PV[i].Power / PV[i].Iarea
		}
	}
}

/* ------------------------------------------------------------- */

func PVprint(fo io.Writer, id int, PV []PV) {
	switch id {
	case 0:
		if len(PV) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PV_TYPE, len(PV))
		}
		for i := range PV {
			fmt.Fprintf(fo, " %s 1 4\n", PV[i].Name)
		}
		break

	case 1:
		for i := range PV {
			fmt.Fprintf(fo, " %s_TPV t f %s_I e f %s_P e f %s_Eff r f \n",
				PV[i].Name, PV[i].Name, PV[i].Name, PV[i].Name)
		}
		break

	default:
		for i := range PV {
			fmt.Fprintf(fo, " %4.1f %4.0f %3.0f %.3f\n",
				PV[i].TPV, PV[i].Iarea, PV[i].Power, PV[i].Eff)
		}
		break
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func PVdyint(PV []PV) {
	for i := range PV {
		qdyint(&PV[i].Edy)
		edyint(&PV[i].Soldy)
	}
}

func PVmonint(PV []PV) {
	for i := range PV {
		qdyint(&PV[i].mEdy)
		edyint(&PV[i].mSoldy)
	}
}

func PVday(Mon int, Day int, ttmm int, PV []PV, Nday int, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for i := range PV {
		var sw ControlSWType
		if PV[i].Power > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}

		// 日間集計
		qdaysum(int64(ttmm), sw, PV[i].Power, &PV[i].Edy)

		if *PV[i].I > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		// 時間集計
		edaysum(ttmm, sw, *PV[i].I, &PV[i].Soldy)

		// 月間集計
		sw = OFF_SW
		if PV[i].Power > 0.0 {
			sw = ON_SW
		}
		qmonsum(Mon, Day, ttmm, sw, PV[i].Power, &PV[i].mEdy, Nday, SimDayend)

		if *PV[i].I > 0.0 {
			sw = ON_SW
		} else {
			sw = OFF_SW
		}
		// 時間集計
		emonsum(Mon, Day, ttmm, sw, *PV[i].I, &PV[i].mSoldy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, sw, PV[i].Power, &PV[i].mtEdy[Mo][tt])
	}
}

func PVdyprt(fo io.Writer, id int, PV []PV) {
	switch id {
	case 0:
		if len(PV) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PV_TYPE, len(PV))
		}
		for i := range PV {
			fmt.Fprintf(fo, " %s 1 8\n", PV[i].Name)
		}
	case 1:
		for i := range PV {
			fmt.Fprintf(fo, "%s_Hh H d %s_E E f\n", PV[i].Name, PV[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_Em e f\n", PV[i].Name, PV[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", PV[i].Name, PV[i].Name, PV[i].Name, PV[i].Name)
		}
	default:
		for i := range PV {
			fmt.Fprintf(fo, "%1d %3.1f ", PV[i].Edy.Hhr, PV[i].Edy.H)
			fmt.Fprintf(fo, "%1d %2.0f\n", PV[i].Edy.Hmxtime, PV[i].Edy.Hmx)

			fmt.Fprintf(fo, "%1d %3.1f ", PV[i].Soldy.Hrs, PV[i].Soldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", PV[i].Soldy.Mxtime, PV[i].Soldy.Mx)
		}
	}
}

func PVmonprt(fo io.Writer, id int, PV []PV) {
	switch id {
	case 0:
		if len(PV) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PV_TYPE, len(PV))
		}
		for i := range PV {
			fmt.Fprintf(fo, " %s 1 8\n", PV[i].Name)
		}
	case 1:
		for i := range PV {
			fmt.Fprintf(fo, "%s_Hh H d %s_E E f\n", PV[i].Name, PV[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_Em e f\n", PV[i].Name, PV[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_S E f %s_te h d %s_Sm e f\n\n", PV[i].Name, PV[i].Name, PV[i].Name, PV[i].Name)
		}
	default:
		for i := range PV {
			fmt.Fprintf(fo, "%1d %3.1f ", PV[i].mEdy.Hhr, PV[i].mEdy.H)
			fmt.Fprintf(fo, "%1d %2.0f\n", PV[i].mEdy.Hmxtime, PV[i].mEdy.Hmx)

			fmt.Fprintf(fo, "%1d %3.1f ", PV[i].mSoldy.Hrs, PV[i].mSoldy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", PV[i].mSoldy.Mxtime, PV[i].mSoldy.Mx)
		}
	}
}

func PVmtprt(fo io.Writer, id int, PV []PV, Mo int, tt int) {
	switch id {
	case 0:
		if len(PV) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PV_TYPE, len(PV))
		}
		for i := range PV {
			fmt.Fprintf(fo, " %s 1 1\n", PV[i].Name)
		}
	case 1:
		for i := range PV {
			fmt.Fprintf(fo, "%s_E E f\n", PV[i].Name)
		}
	default:
		for i := range PV {
			fmt.Fprintf(fo, " %.2f\n", PV[i].mtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}
