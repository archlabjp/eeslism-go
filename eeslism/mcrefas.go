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

/*  refas.c */

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*  圧縮式冷凍機

機器仕様入力          */
var __Refadata_hpch *HPCH

func Refadata(s string, Refaca *REFACA, Rfcmp []*RFCMP) int {
	var c ControlSWType
	var dt float64
	var id int

	if !strings.ContainsAny(s, "=-") {
		// イコールとハイフンが含まれていない場合の処理
		Refaca.name = s
		Refaca.Nmode = 0
		Refaca.unlimcap = 'n'
		Refaca.cool = nil
		Refaca.heat = nil
		Refaca.awtyp = ' '
		Refaca.plf = ' '
		Refaca.rfc = nil
		Refaca.Ph = -999.0
	} else {
		if s[0] == 'a' {
			Refaca.awtyp = rune(s[1])
		} else if strings.Compare(s, "-U") == 0 {
			Refaca.unlimcap = 'y'
		} else if s[0] == 'c' {
			var i int
			Nrfcmp := len(Rfcmp)
			for i = 0; i < Nrfcmp; i++ {
				rfc := Rfcmp[i]
				if strings.Compare(s[1:], rfc.cname) == 0 {
					Refaca.rfc = rfc
					break
				}
			}
			if i == Nrfcmp {
				id = 1
			}
		} else if s[0] == 'p' {
			Refaca.plf = rune(s[1])
		} else if s[0] == 'm' {
			c = ControlSWType(s[1])
			Refaca.mode[Refaca.Nmode] = c
			if c == COOLING_SW {
				Refaca.cool = new(HPCH)
				__Refadata_hpch = Refaca.cool
			} else if c == HEATING_SW {
				Refaca.heat = new(HPCH)
				__Refadata_hpch = Refaca.heat
			}
			Refaca.Nmode++
		} else {
			dt, _ = strconv.ParseFloat(s[1:], 64)
			switch {
			case strings.HasPrefix(s, "Qo"):
				__Refadata_hpch.Qo = dt
			case strings.HasPrefix(s, "Go"):
				__Refadata_hpch.Go = dt
			case strings.HasPrefix(s, "Two"):
				__Refadata_hpch.Two = dt
			case strings.HasPrefix(s, "eo"):
				__Refadata_hpch.eo = dt
			case strings.HasPrefix(s, "Qex"):
				__Refadata_hpch.Qex = dt
			case strings.HasPrefix(s, "Gex"):
				__Refadata_hpch.Gex = dt
			case strings.HasPrefix(s, "Tex"):
				__Refadata_hpch.Tex = dt
			case strings.HasPrefix(s, "eex"):
				__Refadata_hpch.eex = dt
			case s[0] == 'W':
				__Refadata_hpch.Wo = dt
			case strings.HasPrefix(s, "Ph"):
				Refaca.Ph = dt
			default:
				id = 1
			}
		}
	}
	return id
}

/* -------------------------------------------- */

/*  冷凍機／ヒ－トポンプの圧縮機特性設定   */

func Refaint(Refa []*REFA, Wd *WDAT, Compnt []*COMPNT) {
	var Cmp *RFCMP
	var Teo, Tco, cGex, Qeo, Qco float64
	var Qes, Qcs, Ws, ke, kc, kw, E float64
	var i int

	for _, refa := range Refa {
		refa.Ta = &Wd.T

		if refa.Cat.awtyp != 'a' {
			fmt.Printf("Refcfi   awtyp=%c\n", refa.Cat.awtyp)
		}
		if refa.Cmp.Roomname != "" {
			refa.Room = roomptr(refa.Cmp.Roomname, Compnt)
			fmt.Printf("RefaRoom=%s\n", refa.Cmp.Roomname)
		}

		for m := 0; m < refa.Cat.Nmode; m++ {
			if refa.Cat.mode[m] == COOLING_SW {
				cGex = Ca * refa.Cat.cool.Gex
				E = (1.0 - refa.Cat.cool.eo) / refa.Cat.cool.eo
				Qeo = refa.Cat.cool.Qo
				Qco = refa.Cat.cool.Qex
				Teo = Qeo*E/(Cw*refa.Cat.cool.Go) + refa.Cat.cool.Two
				Tco = Qco/(refa.Cat.cool.eex*cGex) + refa.Cat.cool.Tex
			} else if refa.Cat.mode[m] == HEATING_SW {
				cGex = Ca * refa.Cat.heat.Gex
				E = (1.0 - refa.Cat.heat.eo) / refa.Cat.heat.eo
				Qeo = refa.Cat.heat.Qex
				Qco = refa.Cat.heat.Qo
				Tco = Qco*E/(Cw*refa.Cat.heat.Go) + refa.Cat.heat.Two
				Teo = Qeo/(refa.Cat.heat.eex*cGex) + refa.Cat.heat.Tex
			}

			Cmp = refa.Cat.rfc
			Qes = Cmp.e[0] + Cmp.e[1]*Teo + Cmp.e[2]*Tco + Cmp.e[3]*Teo*Tco
			Qcs = Cmp.d[0] + Cmp.d[1]*Teo + Cmp.d[2]*Tco + Cmp.d[3]*Teo*Tco
			Ws = Cmp.w[0] + Cmp.w[1]*Teo + Cmp.w[2]*Tco + Cmp.w[3]*Teo*Tco
			ke = Qeo / Qes
			kc = Qco / Qcs
			if refa.Cat.mode[m] == COOLING_SW {
				kw = refa.Cat.cool.Wo / Ws
			} else if refa.Cat.mode[m] == HEATING_SW {
				kw = refa.Cat.heat.Wo / Ws
			}

			if refa.Cat.mode[m] == COOLING_SW {
				for i = 0; i < 4; i++ {
					refa.c_e[i] = ke * Cmp.e[i]
					refa.c_d[i] = kc * Cmp.d[i]
					refa.c_w[i] = kw * Cmp.w[i]
				}
			} else if refa.Cat.mode[m] == HEATING_SW {
				for i = 0; i < 4; i++ {
					refa.h_e[i] = ke * Cmp.e[i]
					refa.h_d[i] = kc * Cmp.d[i]
					refa.h_w[i] = kw * Cmp.w[i]
				}
			}
		}
		if refa.Cat.Nmode == 1 {
			refa.Chmode = refa.Cat.mode[0]
		}
	}
}

/* -------------------------------------------- */

/*  冷凍機／ヒ－トポンプのシステム方程式の係数  */

//
//             +------+
// [IN 1] ---> | REFA | --> [OUT 1]
//             +------+
//
func Refacfv(Refa []*REFA) {
	for _, refa := range Refa {
		if refa.Cmp.Control != OFF_SW {
			Eo1 := refa.Cmp.Elouts[0]

			cG := Spcheat(Eo1.Fluid) * Eo1.G
			refa.cG = cG
			Eo1.Coeffo = cG

			if Eo1.Control != OFF_SW {
				if Eo1.Sysld == 'y' {
					Eo1.Co = 0.0
					Eo1.Coeffin[0] = -cG
				} else {
					var err int
					refacoeff(refa, &err)
					if err == 0 {
						Eo1.Co = refa.Do
						Eo1.Coeffin[0] = refa.D1 - cG
					} else {
						s := fmt.Sprintf("xxxxx refacoeff xxx stop xx  %s chmode=%c  monitor=%s",
							refa.Name, refa.Chmode, refa.Cmp.Elouts[0].Emonitr.Cmp.Name)
						Eprint("<Refacfv>", s)
						os.Exit(EXIT_REFA)
					}
				}
			}
		}
	}
}

/* ------------------------------------------------------------- */

/*    冷凍機／ヒ－トポンプの能力特性一次式の係数  */

func refacoeff(Refa *REFA, err *int) {
	var E, EGex, Px float64
	*err = 0

	if Refa.Chmode == COOLING_SW {
		if Refa.Cat.cool != nil {
			EGex = Refa.Cat.cool.eex * Ca * Refa.Cat.cool.Gex
			Compca(&Refa.c_e, &Refa.c_d, EGex, Refa.Cat.rfc.Teo, *Refa.Ta, &Refa.Ho, &Refa.He)
			E = Refa.cG * Refa.Cat.cool.eo
		} else {
			*err = 1
		}
	} else if Refa.Chmode == HEATING_SW {
		if Refa.Cat.heat != nil {
			EGex = Refa.Cat.heat.eex * Ca * Refa.Cat.heat.Gex
			Compha(&Refa.h_e, &Refa.h_d, EGex, Refa.Cat.rfc.Tco, *Refa.Ta, &Refa.Ho, &Refa.He)
			E = Refa.cG * Refa.Cat.heat.eo
		} else {
			*err = 1
		}
	}

	if *err == 0 {
		Px = E / (E + Refa.He)
		Refa.Do = Refa.Ho * Px
		Refa.D1 = Refa.He * Px
	} else {
		Refa.Do = 0.0
		Refa.D1 = 0.0
	}
}

/* ------------------------------------------------------------- */

/*   冷却熱量/加熱量、エネルギーの計算 */

func Refaene(Refa []*REFA, LDreset *int) {
	var err, reset int
	var Emax float64
	var Eo *ELOUT

	for i, refa := range Refa {
		refa.Tin = refa.Cmp.Elins[0].Sysvin
		Eo = refa.Cmp.Elouts[0]
		refa.E = 0.0

		if Eo.Control != OFF_SW {
			refa.Ph = refa.Cat.Ph
			refa.Q = refa.cG * (Eo.Sysv - refa.Tin)

			if Eo.Sysld == 'n' {
				refa.Qmax = refa.Q

				if refa.Cat.Nmode > 0 {
					refa.E = Refpow(refa, refa.Q) / refa.Cat.rfc.Meff
				}

			} else {
				reset = chswreset(refa.Q, refa.Chmode, Eo)

				if reset != 0 {
					(*LDreset)++
					refa.Cmp.Control = OFF_SW
				} else {
					if refa.Cat.Nmode > 0 {
						refacoeff(refa, &err)

						if err == 0 {
							refa.Qmax = refa.Do - refa.D1*refa.Tin
							Emax = Refpow(refa, refa.Qmax) / refa.Cat.rfc.Meff
							refa.E = (refa.Q / refa.Qmax) * Emax

							if refa.Cat.unlimcap == 'n' {
								reset = maxcapreset(refa.Q, refa.Qmax, refa.Chmode, Eo)
							}
							if reset != 0 {
								Refacfv(Refa[i : i+1])
								(*LDreset)++
							}
						}
					} else {
						refa.Qmax = refa.Q
					}

				}
			}
		} else {
			refa.Q = 0.0
			refa.Ph = 0.0
		}
	}
}

func Refaene2(Refa []*REFA) {
	for _, refa := range Refa {
		if refa.Room != nil {
			refa.Room.Qeqp += (refa.Q * refa.Cmp.Eqpeff)
		}
	}
}

/* ------------------------------------------------------------- */

/* 負荷計算指定時の設定値のポインター */

func refaldptr(load *ControlSWType, key []string, Refa *REFA) (VPTR, error) {
	var err error
	var vptr VPTR

	if key[1] == "Tout" {
		vptr.Ptr = &Refa.Toset
		vptr.Type = VAL_CTYPE
		Refa.Load = load
	} else {
		err = errors.New("Tout expected")
	}
	return vptr, err
}

/* ------------------------------------------------------------- */

/* 冷暖運転切換のポインター */

func refaswptr(key []string, Refa *REFA) (VPTR, error) {
	if key[1] == "chmode" {
		return VPTR{
			Ptr:  &Refa.Chmode,
			Type: SW_CTYPE,
		}, nil
	}

	return VPTR{}, errors.New("refaswptr error")
}

/* --------------------------- */

/* 負荷計算指定時のスケジュール設定 */

func refaldschd(Refa *REFA) {
	Eo := Refa.Cmp.Elouts[0]

	if Refa.Load != nil {
		if Eo.Control != OFF_SW {
			if Refa.Toset > TEMPLIMIT {
				Eo.Control = LOAD_SW
				Eo.Sysv = Refa.Toset
			} else {
				Eo.Control = OFF_SW
			}
		}
	}
}

/* --------------------------- */

func refaprint(fo io.Writer, id int, Refa []*REFA) {
	switch id {
	case 0:
		if len(Refa) > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, len(Refa))
		}
		for _, refa := range Refa {
			fmt.Fprintf(fo, " %s 1 7\n", refa.Name)
		}
	case 1:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f ",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_P e f\n",
				refa.Name, refa.Name, refa.Name)
		}
	default:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %3.0f %3.0f %2.0f\n",
				refa.Cmp.Elouts[0].Control, refa.Cmp.Elouts[0].G, refa.Tin,
				refa.Cmp.Elouts[0].Sysv, refa.Q, refa.E, refa.Ph)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func refadyint(Refa []*REFA) {
	for _, refa := range Refa {
		svdyint(&refa.Tidy)
		qdyint(&refa.Qdy)
		edyint(&refa.Edy)
		edyint(&refa.Phdy)
	}
}

func refamonint(Refa []*REFA) {
	for _, refa := range Refa {
		svdyint(&refa.mTidy)
		qdyint(&refa.mQdy)
		edyint(&refa.mEdy)
		edyint(&refa.mPhdy)
	}
}

func refaday(Mon int, Day int, ttmm int, Refa []*REFA, Nday int, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)
	for _, refa := range Refa {

		// 日集計
		svdaysum(int64(ttmm), refa.Cmp.Control, refa.Tin, &refa.Tidy)
		qdaysum(int64(ttmm), refa.Cmp.Control, refa.Q, &refa.Qdy)
		edaysum(ttmm, refa.Cmp.Control, refa.E, &refa.Edy)
		edaysum(ttmm, refa.Cmp.Control, refa.Ph, &refa.Phdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, refa.Cmp.Control, refa.Tin, &refa.mTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, refa.Cmp.Control, refa.Q, &refa.mQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, refa.Cmp.Control, refa.E, &refa.mEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, refa.Cmp.Control, refa.Ph, &refa.mPhdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, refa.Cmp.Control, refa.E, &refa.mtEdy[Mo][tt])
		emtsum(Mon, Day, ttmm, refa.Cmp.Control, refa.E, &refa.mtPhdy[Mo][tt])
	}
}

func refadyprt(fo io.Writer, id int, Refa []*REFA) {
	switch id {
	case 0:
		if len(Refa) > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, len(Refa))
		}
		for _, refa := range Refa {
			fmt.Fprintf(fo, " %s 1 22\n", refa.Name)
		}
	case 1:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
		}
	default:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				refa.Tidy.Hrs, refa.Tidy.M,
				refa.Tidy.Mntime, refa.Tidy.Mn,
				refa.Tidy.Mxtime, refa.Tidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.Qdy.Hhr, refa.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", refa.Qdy.Chr, refa.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.Qdy.Hmxtime, refa.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.Qdy.Cmxtime, refa.Qdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.Edy.Hrs, refa.Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.Edy.Mxtime, refa.Edy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.Phdy.Hrs, refa.Phdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", refa.Phdy.Mxtime, refa.Phdy.Mx)
		}
	}
}

func refamonprt(fo io.Writer, id int, Refa []*REFA) {
	switch id {
	case 0:
		if len(Refa) > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, len(Refa))
		}
		for _, refa := range Refa {
			fmt.Fprintf(fo, " %s 1 22\n", refa.Name)
		}
	case 1:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				refa.Name, refa.Name, refa.Name, refa.Name)
		}
	default:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				refa.mTidy.Hrs, refa.mTidy.M,
				refa.mTidy.Mntime, refa.mTidy.Mn,
				refa.mTidy.Mxtime, refa.mTidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.mQdy.Hhr, refa.mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", refa.mQdy.Chr, refa.mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.mQdy.Hmxtime, refa.mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.mQdy.Cmxtime, refa.mQdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.mEdy.Hrs, refa.mEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", refa.mEdy.Mxtime, refa.mEdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", refa.mPhdy.Hrs, refa.mPhdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", refa.mPhdy.Mxtime, refa.mPhdy.Mx)
		}
	}
}

func refamtprt(fo io.Writer, id int, Refa []*REFA, Mo int, tt int) {
	switch id {
	case 0:
		if len(Refa) > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, len(Refa))
		}
		for _, refa := range Refa {
			fmt.Fprintf(fo, " %s 1 2\n", refa.Name)
		}
	case 1:
		for _, refa := range Refa {
			fmt.Fprintf(fo, "%s_E E f %s_Ph E f \n", refa.Name, refa.Name)
		}
	default:
		for _, refa := range Refa {
			fmt.Fprintf(fo, " %.2f %.2f\n", refa.mtEdy[Mo-1][tt-1].D*Cff_kWh, refa.mtPhdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}
