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

/*  boiler.c  */

/*  ボイラ－   */

package eeslism

/* 機器仕様入力　　　　　　*/

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func Boidata(s string, boica *BOICA) int {
	var id int

	st := strings.IndexRune(s, '=')
	if st == -1 && s[0] != '-' {
		boica.name = s
		boica.unlimcap = 'n'
		boica.ene = ' '
		boica.plf = ' '
		boica.Qo = nil
		boica.eff = 1.0
		boica.Ph = -999.0
		boica.Qmin = -999.0
		//boica.mode = 'n'
		boica.Qostr = ""
	} else if s == "-U" {
		boica.unlimcap = 'y'
	} else {
		if st >= 0 {
			s1, s2 := s[:st], s[st+1:]
			switch s1 {
			case "p":
				boica.plf = rune(s2[0])
			case "en":
				boica.ene = rune(s2[0])
			case "blwQmin":

				if s2 == "ON" {
					// 負荷が最小出力以下のときに最小出力でONとする
					boica.belowmin = ON_SW
				} else if s2 == "OFF" {
					// 負荷が最小出力以下のときにOFFとする
					boica.belowmin = OFF_SW
				} else {
					id = 1
				}
			case "Qo":
				boica.Qostr = s2
			case "Qmin", "eff", "Ph":
				dt, err := strconv.ParseFloat(s2, 64)
				if err != nil {
					id = 1
				} else {
					switch s1 {
					case "Qmin":
						boica.Qmin = dt
					case "eff":
						boica.eff = dt
					case "Ph":
						boica.Ph = dt
					}
				}
			default:
				id = 1
			}
		}
	}
	return id
}

func Boicaint(_Boica []*BOICA, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	for _, Boica := range _Boica {
		if idx, err := idsch(Boica.Qostr, Schdl.Sch, ""); err == nil {
			Boica.Qo = &Schdl.Val[idx]
		} else {
			Boica.Qo = envptr(Boica.Qostr, Simc, nil, nil, nil)
		}
	}
}

/* --------------------------- */

/*  特性式の係数  */

//
//             +-----+
// [IN 1] ---> | BOI | ---> [OUT 1] 出口温度??
//             +-----+
//
func Boicfv(Boi []*BOI) {
	var cG, Qocat, Temp float64

	if len(Boi) != len(Boi) {
		panic("len(Boi) != len(Boi)")
	}

	for _, boi := range Boi {

		Eo1 := boi.Cmp.Elouts[0]

		if boi.Cmp.Control != OFF_SW {
			Temp = math.Abs(*boi.Cat.Qo - (-999.9))
			if math.Abs(Temp) < 1e-3 {
				Qocat = 0.0
			} else {
				Qocat = *boi.Cat.Qo
			}

			if Qocat > 0.0 {
				boi.HCmode = HEATING_LOAD
			} else {
				boi.HCmode = COOLING_LOAD
			}

			boi.Do = Qocat

			if (boi.Do < 0.0 && boi.HCmode == HEATING_LOAD) || (boi.Do > 0.0 && boi.HCmode == COOLING_LOAD) || boi.HCmode == 'n' {
				fmt.Printf("<BOI> name=%s  Qo=%.4g\n", boi.Cmp.Name, boi.Do)
			}

			boi.D1 = 0.0

			cG = Spcheat(Eo1.Fluid) * Eo1.G
			boi.cG = cG
			Eo1.Coeffo = cG

			if Eo1.Control != OFF_SW {
				if Eo1.Sysld == 'y' {
					// 出口を設定温度に制御
					Eo1.Co = 0.0
					Eo1.Coeffin[0] = -cG
				} else {
					if boi.Mode == 'M' {
						// 最大能力
						Eo1.Co = boi.Do
					} else {
						// 最小能力
						Eo1.Co = boi.Cat.Qmin
					}
					Eo1.Coeffin[0] = boi.D1 - cG
				}
			}
		} else {
			// 機器が停止
			Eo1.Co = 0.0
			Eo1.Coeffo = 1.0
			Eo1.Coeffin[0] = -1.0
		}
	}
}

/* --------------------------- */

/*  供給熱量、エネルギーの計算 */

func Boiene(Boi []*BOI, BOIreset *int) {
	for i, boi := range Boi {
		boi.Tin = boi.Cmp.Elins[0].Sysvin
		Qmin := boi.Cat.Qmin
		if math.Abs(Qmin-(-999.0)) < 1.0e-5 {
			Qmin = 0.0
		}

		Eo := boi.Cmp.Elouts[0]
		reset := 0

		if Eo.Control != OFF_SW {
			boi.Q = boi.cG * (Eo.Sysv - boi.Tin)

			// 次回ループの機器制御判定用の熱量
			Qcheck := boi.Q

			// 加熱設定での冷却、冷却設定での加熱時はボイラを止める
			if (Qcheck < 0.0 && boi.HCmode == 'H') || (Qcheck > 0.0 && boi.HCmode == 'C') {
				boi.Cmp.Control = OFF_SW
				Eo.Control = ON_SW
				Eo.Emonitr.Control = ON_SW
				Eo.Sysld = 'n'

				reset = 1
			} else if Qmin > 0.0 && Qcheck < Qmin {
				// 最小出力以下はOFFにする場合
				if boi.Cat.belowmin == OFF_SW {
					boi.Cmp.Elouts[0].Control = OFF_SW
					boi.Cmp.Control = OFF_SW
					Eo.Control = ON_SW
					Eo.Emonitr.Control = ON_SW
					Eo.Sysld = 'n'
				} else {
					Eo.Control = ON_SW
					Eo.Emonitr.Control = ON_SW
					Eo.Sysld = 'n'
					boi.Mode = 'm'
				}

				reset = 1
			} else if boi.Cat.unlimcap == 'n' {
				// 過負荷状態のチェック
				Qocat := 0.0
				if math.Abs(*boi.Cat.Qo-(-999.9)) < 1.0e-3 {
					Qocat = 0.0
				} else {
					Qocat = *boi.Cat.Qo
				}

				reset0 := maxcapreset(Qcheck, Qocat, boi.HCmode, Eo)

				if reset == 0 {
					reset = reset0
				}
			}

			if reset == 1 {
				Boicfv(Boi[i : i+1])
				(*BOIreset)++
			}

			boi.E = boi.Q / boi.Cat.eff
			boi.Ph = boi.Cat.Ph
		} else {
			boi.Q = 0.0
			boi.E = 0.0
			boi.Ph = 0.0
		}
	}
}

/* --------------------------- */

/* 負荷計算指定時の設定値のポインター */

func boildptr(load *ControlSWType, key []string, Boi *BOI) (VPTR, error) {
	var err error
	var vptr VPTR

	if strings.Compare(key[1], "Tout") == 0 {
		vptr.Ptr = &Boi.Toset
		vptr.Type = VAL_CTYPE
		Boi.Load = load
	} else {
		err = errors.New("Tout expected")
	}
	return vptr, err
}

/* --------------------------- */

/* 負荷計算指定時のスケジュール設定 */

func boildschd(Boi *BOI) {
	Eo := Boi.Cmp.Elouts[0]

	if Boi.Load != nil {
		if Eo.Control != OFF_SW {
			if Boi.Toset > TEMPLIMIT {
				Eo.Control = LOAD_SW
				Eo.Sysv = Boi.Toset
			} else {
				Eo.Control = OFF_SW
			}
		}
	}
}

/* --------------------------- */

func boiprint(fo io.Writer, id int, Boi []*BOI) {
	for _, boi := range Boi {
		switch id {
		case 0:
			if len(Boi) > 0 {
				fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
			}
			fmt.Fprintf(fo, " %s 1 7\n", boi.Name)
		case 1:
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f ", boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_P e f\n", boi.Name, boi.Name, boi.Name)
		default:
			fmt.Fprintf(fo, "%c %.4g %4.2f %4.2f %4.0f %4.0f %2.0f\n",
				boi.Cmp.Elouts[0].Control, boi.Cmp.Elouts[0].G,
				boi.Tin, boi.Cmp.Elouts[0].Sysv, boi.Q, boi.E, boi.Ph)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func boidyint(Boi []*BOI) {
	for _, boi := range Boi {
		// 日集計のリセット
		svdyint(&boi.Tidy)
		qdyint(&boi.Qdy)
		edyint(&boi.Edy)
		edyint(&boi.Phdy)
	}
}

/* --------------------------- */

/* 月積算値に関する処理 */

func boimonint(Boi []*BOI) {
	for _, boi := range Boi {
		// 日集計のリセット
		svdyint(&boi.mTidy)
		qdyint(&boi.mQdy)
		edyint(&boi.mEdy)
		edyint(&boi.mPhdy)
	}
}

func boiday(Mon, Day, ttmm int, Boi []*BOI, Nday, SimDayend int) {
	var Mo, tt int

	Mo = Mon - 1
	tt = ConvertHour(ttmm)
	for _, boi := range Boi {
		// 日集計
		svdaysum(int64(ttmm), boi.Cmp.Control, boi.Tin, &boi.Tidy)
		qdaysum(int64(ttmm), boi.Cmp.Control, boi.Q, &boi.Qdy)
		edaysum(ttmm, boi.Cmp.Control, boi.E, &boi.Edy)
		edaysum(ttmm, boi.Cmp.Control, boi.Ph, &boi.Phdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Tin, &boi.mTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Q, &boi.mQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.mEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Ph, &boi.mPhdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.MtEdy[Mo][tt])
		emtsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.MtPhdy[Mo][tt])
	}
}

func boidyprt(fo io.Writer, id int, Boi []*BOI) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 22\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				boi.Tidy.Hrs, boi.Tidy.M,
				boi.Tidy.Mntime, boi.Tidy.Mn,
				boi.Tidy.Mxtime, boi.Tidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Qdy.Hhr, boi.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", boi.Qdy.Chr, boi.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Qdy.Hmxtime, boi.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Qdy.Cmxtime, boi.Qdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Edy.Hrs, boi.Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Edy.Mxtime, boi.Edy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Phdy.Hrs, boi.Phdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", boi.Phdy.Mxtime, boi.Phdy.Mx)
		}
	}
}

func boimonprt(fo io.Writer, id int, Boi []*BOI) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 22\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				boi.mTidy.Hrs, boi.mTidy.M,
				boi.mTidy.Mntime, boi.mTidy.Mn,
				boi.mTidy.Mxtime, boi.mTidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mQdy.Hhr, boi.mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", boi.mQdy.Chr, boi.mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mQdy.Hmxtime, boi.mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mQdy.Cmxtime, boi.mQdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mEdy.Hrs, boi.mEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mEdy.Mxtime, boi.mEdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mPhdy.Hrs, boi.mPhdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", boi.mPhdy.Mxtime, boi.mPhdy.Mx)
		}
	}
}

func boimtprt(fo io.Writer, id int, Boi []*BOI, Mo int, tt int) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 2\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_E E f %s_Ph E f \n", boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %.2f %.2f\n",
				boi.MtEdy[Mo-1][tt-1].D*Cff_kWh, boi.MtPhdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}
