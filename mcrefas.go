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

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*  圧縮式冷凍機

機器仕様入力          */

func Refadata(s string, Refaca *REFACA, Nrfcmp int, Rfcmp []RFCMP) int {
	var hpch *HPCH
	var c byte
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
			for i = 0; i < Nrfcmp; i++ {
				rfc := &Rfcmp[i]
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
			c = s[1]
			Refaca.mode[Refaca.Nmode] = rune(c)
			if c == COOLING_SW {
				Refaca.cool = new(HPCH)
				hpch = Refaca.cool
			} else if c == HEATING_SW {
				Refaca.heat = new(HPCH)
				hpch = Refaca.heat
			}
			Refaca.Nmode++
		} else {
			dt, _ = strconv.ParseFloat(s[1:], 64)
			switch {
			case strings.HasPrefix(s, "Qo"):
				hpch.Qo = dt
			case strings.HasPrefix(s, "Go"):
				hpch.Go = dt
			case strings.HasPrefix(s, "Two"):
				hpch.Two = dt
			case strings.HasPrefix(s, "eo"):
				hpch.eo = dt
			case strings.HasPrefix(s, "Qex"):
				hpch.Qex = dt
			case strings.HasPrefix(s, "Gex"):
				hpch.Gex = dt
			case strings.HasPrefix(s, "Tex"):
				hpch.Tex = dt
			case strings.HasPrefix(s, "eex"):
				hpch.eex = dt
			case s[0] == 'W':
				hpch.Wo = dt
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

func Refaint(Nrefa int, Refa []REFA, Wd *WDAT, Ncompnt int, Compnt []COMPNT) {
	var Cmp *RFCMP
	var Teo, Tco, cGex, Qeo, Qco float64
	var Qes, Qcs, Ws, ke, kc, kw, E float64
	var i int

	for k := 0; k < Nrefa; k++ {
		Refa[k].Ta = &Wd.T

		if Refa[k].Cat.awtyp != 'a' {
			fmt.Printf("Refcfi   awtyp=%c\n", Refa[k].Cat.awtyp)
		}
		if Refa[k].Cmp.Roomname != "" {
			Refa[k].Room = roomptr(Refa[k].Cmp.Roomname, Ncompnt, Compnt)
			fmt.Printf("RefaRoom=%s\n", Refa[k].Cmp.Roomname)
		}

		for m := 0; m < Refa[k].Cat.Nmode; m++ {
			if Refa[k].Cat.mode[m] == COOLING_SW {
				cGex = Ca * Refa[k].Cat.cool.Gex
				E = (1.0 - Refa[k].Cat.cool.eo) / Refa[k].Cat.cool.eo
				Qeo = Refa[k].Cat.cool.Qo
				Qco = Refa[k].Cat.cool.Qex
				Teo = Qeo*E/(Cw*Refa[k].Cat.cool.Go) + Refa[k].Cat.cool.Two
				Tco = Qco/(Refa[k].Cat.cool.eex*cGex) + Refa[k].Cat.cool.Tex
			} else if Refa[k].Cat.mode[m] == HEATING_SW {
				cGex = Ca * Refa[k].Cat.heat.Gex
				E = (1.0 - Refa[k].Cat.heat.eo) / Refa[k].Cat.heat.eo
				Qeo = Refa[k].Cat.heat.Qex
				Qco = Refa[k].Cat.heat.Qo
				Tco = Qco*E/(Cw*Refa[k].Cat.heat.Go) + Refa[k].Cat.heat.Two
				Teo = Qeo/(Refa[k].Cat.heat.eex*cGex) + Refa[k].Cat.heat.Tex
			}

			Cmp = Refa[k].Cat.rfc
			Qes = Cmp.e[0] + Cmp.e[1]*Teo + Cmp.e[2]*Tco + Cmp.e[3]*Teo*Tco
			Qcs = Cmp.d[0] + Cmp.d[1]*Teo + Cmp.d[2]*Tco + Cmp.d[3]*Teo*Tco
			Ws = Cmp.w[0] + Cmp.w[1]*Teo + Cmp.w[2]*Tco + Cmp.w[3]*Teo*Tco
			ke = Qeo / Qes
			kc = Qco / Qcs
			if Refa[k].Cat.mode[m] == COOLING_SW {
				kw = Refa[k].Cat.cool.Wo / Ws
			} else if Refa[k].Cat.mode[m] == HEATING_SW {
				kw = Refa[k].Cat.heat.Wo / Ws
			}

			if Refa[k].Cat.mode[m] == COOLING_SW {
				for i = 0; i < 4; i++ {
					Refa[k].c_e[i] = ke * Cmp.e[i]
					Refa[k].c_d[i] = kc * Cmp.d[i]
					Refa[k].c_w[i] = kw * Cmp.w[i]
				}
			} else if Refa[k].Cat.mode[m] == HEATING_SW {
				for i = 0; i < 4; i++ {
					Refa[k].h_e[i] = ke * Cmp.e[i]
					Refa[k].h_d[i] = kc * Cmp.d[i]
					Refa[k].h_w[i] = kw * Cmp.w[i]
				}
			}
		}
		if Refa[k].Cat.Nmode == 1 {
			Refa[k].Chmode = Refa[k].Cat.mode[0]
		}
	}
}

/* -------------------------------------------- */

/*  冷凍機／ヒ－トポンプのシステム方程式の係数  */

func Refacfv(Nrefa int, Refa []REFA) {
	var Eo *ELOUT
	var cG float64
	var err int
	var s string

	for i := 0; i < Nrefa; i++ {
		if Refa[i].Cmp.Control != OFF_SW {
			Eo = Refa[i].Cmp.Elouts[0]

			cG = Spcheat(Eo.Fluid) * Eo.G
			Refa[i].cG = cG
			Eo.Coeffo = cG

			if Eo.Control != OFF_SW {
				if Eo.Sysld == 'y' {
					Eo.Co = 0.0
					Eo.Coeffin[0] = -cG
				} else {
					refacoeff(&Refa[i], &err)
					if err == 0 {
						Eo.Co = Refa[i].Do
						Eo.Coeffin[0] = Refa[i].D1 - cG
					} else {
						s = fmt.Sprintf("xxxxx refacoeff xxx stop xx  %s chmode=%c  monitor=%s",
							Refa[i].Name, Refa[i].Chmode, Refa[i].Cmp.Elouts[0].Emonitr.Cmp.Name)
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

func Refaene(Nrefa int, Refa []REFA, LDreset *int) {
	var err, reset int
	var Emax float64
	var Eo *ELOUT

	for i := 0; i < Nrefa; i++ {
		Refa[i].Tin = Refa[i].Cmp.Elins[0].Sysvin
		Eo = Refa[i].Cmp.Elouts[0]
		Refa[i].E = 0.0

		if Eo.Control != OFF_SW {
			Refa[i].Ph = Refa[i].Cat.Ph
			Refa[i].Q = Refa[i].cG * (Eo.Sysv - Refa[i].Tin)

			if Eo.Sysld == 'n' {
				Refa[i].Qmax = Refa[i].Q

				if Refa[i].Cat.Nmode > 0 {
					Refa[i].E = Refpow(&Refa[i], Refa[i].Q) / Refa[i].Cat.rfc.Meff
				}

			} else {
				reset = chswreset(Refa[i].Q, Refa[i].Chmode, Eo)

				if reset != 0 {
					(*LDreset)++
					Refa[i].Cmp.Control = OFF_SW
				} else {
					if Refa[i].Cat.Nmode > 0 {
						refacoeff(&Refa[i], &err)

						if err == 0 {
							Refa[i].Qmax = Refa[i].Do - Refa[i].D1*Refa[i].Tin
							Emax = Refpow(&Refa[i], Refa[i].Qmax) / Refa[i].Cat.rfc.Meff
							Refa[i].E = (Refa[i].Q / Refa[i].Qmax) * Emax

							if Refa[i].Cat.unlimcap == 'n' {
								reset = maxcapreset(Refa[i].Q, Refa[i].Qmax, Refa[i].Chmode, Eo)
							}
							if reset != 0 {
								Refacfv(1, Refa[i:i+1])
								(*LDreset)++
							}
						}
					} else {
						Refa[i].Qmax = Refa[i].Q
					}

				}
			}
		} else {
			Refa[i].Q = 0.0
			Refa[i].Ph = 0.0
		}
	}
}

func Refaene2(Nrefa int, Refa []REFA) {
	for i := 0; i < Nrefa; i++ {
		if Refa[i].Room != nil {
			Refa[i].Room.Qeqp += (Refa[i].Q * Refa[i].Cmp.Eqpeff)
		}
	}
}

/* ------------------------------------------------------------- */

/* 負荷計算指定時の設定値のポインター */

func refaldptr(load *rune, key []string, Refa *REFA, vptr *VPTR) int {
	err := 0

	if key[1] == "Tout" {
		vptr.Ptr = &Refa.Toset
		vptr.Type = VAL_CTYPE
		Refa.Load = load
	} else {
		err = 1
	}
	return err
}

/* ------------------------------------------------------------- */

/* 冷暖運転切換のポインター */

func refaswptr(key []string, Refa *REFA, vptr *VPTR) int {
	err := 0

	if key[1] == "chmode" {
		vptr.Ptr = &Refa.Chmode
		vptr.Type = SW_CTYPE
	} else {
		err = 1
	}
	return err
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

func refaprint(fo io.Writer, id, Nrefa int, Refa []REFA) {
	switch id {
	case 0:
		if Nrefa > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, Nrefa)
		}
		for i := 0; i < Nrefa; i++ {
			fmt.Fprintf(fo, " %s 1 7\n", Refa[i].Name)
		}
	case 1:
		for i := 0; i < Nrefa; i++ {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f ",
				Refa[i].Name, Refa[i].Name, Refa[i].Name, Refa[i].Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_P e f\n",
				Refa[i].Name, Refa[i].Name, Refa[i].Name)
		}
	default:
		for i := 0; i < Nrefa; i++ {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %3.0f %3.0f %2.0f\n",
				Refa[i].Cmp.Elouts[0].Control, Refa[i].Cmp.Elouts[0].G, Refa[i].Tin,
				Refa[i].Cmp.Elouts[0].Sysv, Refa[i].Q, Refa[i].E, Refa[i].Ph)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func refadyint(Nrefa int, Refa []REFA) {
	for i := 0; i < Nrefa; i++ {
		svdyint(&Refa[i].Tidy)
		qdyint(&Refa[i].Qdy)
		edyint(&Refa[i].Edy)
		edyint(&Refa[i].Phdy)
	}
}

func refamonint(Nrefa int, Refa []REFA) {
	for i := 0; i < Nrefa; i++ {
		svdyint(&Refa[i].mTidy)
		qdyint(&Refa[i].mQdy)
		edyint(&Refa[i].mEdy)
		edyint(&Refa[i].mPhdy)
	}
}

func refaday(Mon int, Day int, ttmm int, Nrefa int, Refa []REFA, Nday int, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)
	for i := 0; i < Nrefa; i++ {
		Refa := &Refa[i]

		// 日集計
		svdaysum(int64(ttmm), Refa.Cmp.Control, Refa.Tin, &Refa.Tidy)
		qdaysum(int64(ttmm), Refa.Cmp.Control, Refa.Q, &Refa.Qdy)
		edaysum(ttmm, Refa.Cmp.Control, Refa.E, &Refa.Edy)
		edaysum(ttmm, Refa.Cmp.Control, Refa.Ph, &Refa.Phdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.Tin, &Refa.mTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.Q, &Refa.mQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.E, &Refa.mEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.Ph, &Refa.mPhdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.E, &Refa.mtEdy[Mo][tt])
		emtsum(Mon, Day, ttmm, Refa.Cmp.Control, Refa.E, &Refa.mtPhdy[Mo][tt])
	}
}

func refadyprt(fo io.Writer, id int, Nrefa int, Refa []REFA) {
	switch id {
	case 0:
		if Nrefa > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, Nrefa)
		}
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, " %s 1 22\n", refa.Name)
		}
	case 1:
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
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
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, "%1ld %3.1f %1ld %3.1f %1ld %3.1f ",
				refa.Tidy.Hrs, refa.Tidy.M,
				refa.Tidy.Mntime, refa.Tidy.Mn,
				refa.Tidy.Mxtime, refa.Tidy.Mx)

			fmt.Fprintf(fo, "%1ld %3.1f ", refa.Qdy.Hhr, refa.Qdy.H)
			fmt.Fprintf(fo, "%1ld %3.1f ", refa.Qdy.Chr, refa.Qdy.C)
			fmt.Fprintf(fo, "%1ld %2.0f ", refa.Qdy.Hmxtime, refa.Qdy.Hmx)
			fmt.Fprintf(fo, "%1ld %2.0f ", refa.Qdy.Cmxtime, refa.Qdy.Cmx)

			fmt.Fprintf(fo, "%1ld %3.1f ", refa.Edy.Hrs, refa.Edy.D)
			fmt.Fprintf(fo, "%1ld %2.0f ", refa.Edy.Mxtime, refa.Edy.Mx)

			fmt.Fprintf(fo, "%1ld %3.1f ", refa.Phdy.Hrs, refa.Phdy.D)
			fmt.Fprintf(fo, "%1ld %2.0f\n", refa.Phdy.Mxtime, refa.Phdy.Mx)
		}
	}
}

func refamonprt(fo io.Writer, id int, Nrefa int, Refa []REFA) {
	switch id {
	case 0:
		if Nrefa > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, Nrefa)
		}
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, " %s 1 22\n", refa.Name)
		}
	case 1:
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
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
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
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

func refamtprt(fo io.Writer, id int, Nrefa int, Refa []REFA, Mo int, tt int) {
	switch id {
	case 0:
		if Nrefa > 0 {
			fmt.Fprintf(fo, "%s %d\n", REFACOMP_TYPE, Nrefa)
		}
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, " %s 1 2\n", refa.Name)
		}
	case 1:
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, "%s_E E f %s_Ph E f \n", refa.Name, refa.Name)
		}
	default:
		for i := 0; i < Nrefa; i++ {
			refa := &Refa[i]
			fmt.Fprintf(fo, " %.2f %.2f\n", refa.mtEdy[Mo-1][tt-1].D*Cff_kWh, refa.mtPhdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}
