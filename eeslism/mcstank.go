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

/*  mcstank.c */

/*  95/11/17 rev  */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
)

/*　蓄熱槽仕様入力　　　　*/

func Stankdata(f *EeTokens, s string, Stankca *STANKCA) int {
	id := 0
	st := ""
	Stankca.gxr = 0.0

	var err error

	if stIdx := strings.IndexByte(s, '='); stIdx != -1 {
		s = strings.TrimSpace(s)
		st = s[stIdx+1:]

		switch {
		case strings.HasPrefix(s, "Vol"):
			Stankca.Vol, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAside"):
			Stankca.KAside, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAtop"):
			Stankca.KAtop, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAbtm"):
			Stankca.KAbtm, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "gxr"):
			Stankca.gxr, err = strconv.ParseFloat(st, 64)
		default:
			id = 1
		}

		if err != nil {
			fmt.Println(err)
		}

	} else if s == "-S" {
		st = ""
		s = f.GetToken()
		s += " *"
		Stankca.tparm = s
	} else {
		Stankca.name = s
		Stankca.Type = 'C'
		Stankca.tparm = ""
		Stankca.Vol = -999.0
		Stankca.KAside = -999.0
		Stankca.KAtop = -999.0
		Stankca.KAbtm = -999.0
		Stankca.gxr = 0.0
	}

	return id
}

/* ------------------------------------------------------- */

/* 蓄熱槽記憶域確保 */

func Stankmemloc(errkey string, Stank *STANK) {
	var np, Ndiv, Nin int
	var st, stt, ss string
	var parm []string = make([]string, 0)

	st = Stank.Cat.tparm[:]

	// 読み飛ばし処理
	np = 0
	for {
		_, err := fmt.Sscanf(st, "%s", &ss)
		if err != nil || ss == "*" {
			break
		}

		parm = append(parm, st)
		np++
		st = st[len(ss):]
		for st[0] == ' ' || st[0] == '\t' {
			st = st[1:]
		}
	}

	Stank.Pthcon = make([]ELIOType, np)
	Stank.Batchcon = make([]rune, np)
	Stank.Ihex = make([]rune, np)
	Stank.Jin = make([]int, np)
	Stank.Jout = make([]int, np)
	Stank.Ihxeff = make([]float64, np)
	Stank.KA = make([]float64, np)
	Stank.KAinput = make([]rune, np)

	i := 0

	for j := 0; j < np; j++ {
		_, err := fmt.Sscanf(parm[j], "%s", &ss)
		if err != nil {
			panic(err)
		}

		if strings.HasPrefix(ss, "N=") {
			Stank.Ndiv, err = strconv.Atoi(ss[2:])
			if err != nil {
				panic(err)
			}
		} else if stIdx := strings.IndexRune(ss, ':'); stIdx != -1 {
			Stank.Pthcon[i] = ELIOType(ss[0])
			tmp, err := strconv.Atoi(ss[stIdx+1:])
			if err != nil {
				panic(err)
			} else {
				Stank.Jin[i] = tmp - 1
			}

			if sttIdx := strings.IndexRune(ss[stIdx+1:], '-'); sttIdx != -1 {
				stt = ss[stIdx+1:]
				Stank.Ihex[i] = 'n'
				Stank.Ihxeff[i] = 1.0
				tmp, err := strconv.Atoi(stt)
				if err != nil {
					panic(err)
				} else {
					Stank.Jout[i] = tmp - 1
				}
			} else if sttIdx := strings.IndexRune(ss[stIdx+1:], '_'); sttIdx != -1 {
				stt = ss[stIdx+1 : sttIdx]
				Stank.Ihex[i] = 'y'

				if stt[1] == 'e' { // 温度効率が入力されている場合
					Stank.Ihxeff[i], err = strconv.ParseFloat(stt[5:], 64)
					if err != nil {
						panic(err)
					}
				} else if stt[1] == 'K' { // 内蔵熱交のKAが入力されている場合
					Stank.KAinput[i] = 'Y'
					Stank.KA[i], err = strconv.ParseFloat(stt[4:], 64)
					if err != nil {
						panic(err)
					}
				} else if stt[1] == 'd' {
					Stank.KAinput[i] = 'C' // 内蔵熱交換器の内径と長さが入力されている場合
					stpIdx := strings.IndexRune(stt[4:], '_')
					Stank.Dbld0, err = strconv.ParseFloat(stt[4:], 64)
					if err != nil {
						panic(err)
					}
					Stank.DblL, err = strconv.ParseFloat(stt[stpIdx+1:], 64)
					if err != nil {
						panic(err)
					}
					Stank.Ncalcihex++
				}

				Stank.Jout[i] = Stank.Jin[i]

				i++
			}
		}
	}

	Stank.Nin = i
	Nin = i

	Ndiv = Stank.Ndiv
	Stank.DtankF = make([]rune, Ndiv)

	Stank.B = make([]float64, Ndiv*Ndiv)
	Stank.R = make([]float64, Ndiv)
	Stank.D = make([]float64, Ndiv)
	Stank.Fg = make([]float64, Ndiv*Nin)
	Stank.Tss = make([]float64, Ndiv)

	Stank.Tssold = make([]float64, Ndiv)
	Stank.Dvol = make([]float64, Ndiv)
	Stank.Mdt = make([]float64, Ndiv)
	Stank.KS = make([]float64, Ndiv)
	Stank.CGwin = make([]float64, Nin)
	Stank.EGwin = make([]float64, Nin)
	Stank.Twin = make([]float64, Nin)
	Stank.Q = make([]float64, Nin)
	if Nin > 0 {
		Stank.Stkdy = make([]STKDAY, Nin)
	}
	if Nin > 0 {
		Stank.Mstkdy = make([]STKDAY, Nin)
	}
}

/* ------------------------------------------------------- */

/* 蓄熱槽初期設定 */

func Stankint(Nstank int, Stank []STANK, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Wd *WDAT) {
	var s, ss, Err, E string
	var mrk rune
	var Tso float64

	E = "Stankint"

	for i := 0; i < Nstank; i++ {
		Stank := &Stank[i]

		// 内蔵熱交換器の熱伝達率計算用温度の初期化
		Stank.DblTa = 20.0
		Stank.DblTw = 20.0

		s = Stank.Cmp.Tparm
		if s != "" {
			if s[0] == '(' {
				s = s[1:]
				for j := 0; j < Stank.Ndiv; j++ {
					_, err := fmt.Sscanf(s, " %s ", &ss)
					if err != nil {
						panic(err)
					}

					if ss[0] == TANK_EMPTY {
						Stank.DtankF[j] = TANK_EMPTY
						Stank.Tssold[j] = TANK_EMPTMP
					} else {
						Stank.DtankF[j] = TANK_FULL
						Stank.Tssold[j], err = strconv.ParseFloat(ss, 64)
						if err != nil {
							panic(err)
						}
					}
					s = s[len(ss):]
					for s[0] == ' ' {
						s = s[1:]
					}
				}
			} else {
				if s[0] == TANK_EMPTY {
					mrk = TANK_EMPTY
					Tso = TANK_EMPTMP
				} else {
					var err error
					mrk = TANK_FULL
					Tso, err = strconv.ParseFloat(s, 64)
					if err != nil {
						panic(err)
					}
				}
				for j := 0; j < Stank.Ndiv; j++ {
					Stank.DtankF[j] = mrk
					Stank.Tssold[j] = Tso
				}
			}
		}

		Stank.Tenv = envptr(Stank.Cmp.Envname, Simc, Ncompnt, Compnt, Wd, nil)
		stoint(Stank.Ndiv, Stank.Cat.Vol, Stank.Cat.KAside, Stank.Cat.KAtop, Stank.Cat.KAbtm,
			Stank.Dvol, Stank.Mdt, Stank.KS, Stank.Tss, Stank.Tssold, &Stank.Jva, &Stank.Jvb)

		if Stank.Cat.Vol < 0.0 {
			Err = fmt.Sprintf("Name=%s  Vol=%.4g", Stank.Cmp.Name, Stank.Cat.Vol)
			Eprint(E, Err)
		}
		if Stank.Cat.KAside < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAside=%.4g", Stank.Cmp.Name, Stank.Cat.KAside)
			Eprint(E, Err)
		}
		if Stank.Cat.KAtop < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAtop=%.4g", Stank.Cmp.Name, Stank.Cat.KAtop)
			Eprint(E, Err)
		}
		if Stank.Cat.KAbtm < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAbtm=%.4g", Stank.Cmp.Name, Stank.Cat.KAbtm)
			Eprint(E, Err)
		}
	}
}

/* ------------------------------------------------------- */

/* 蓄熱槽特性式係数 */

func Stankcfv(Nstank int, Stank []STANK) {
	for i := 0; i < Nstank; i++ {
		Stank := &Stank[i]

		for j := 0; j < Stank.Nin; j++ {
			elin := Stank.Cmp.Elins[j]
			cGwin := &Stank.CGwin[j]
			EGwin := &Stank.EGwin[j]
			ihxeff := &Stank.Ihxeff[j]
			ihex := &Stank.Ihex[j]

			if elin.Lpath.Batch == 'y' {
				*cGwin = 0.0
			} else {
				*cGwin = Spcheat('W') * elin.Lpath.G
			}

			// 内蔵熱交のKAが入力されている場合
			if *ihex == 'y' && *cGwin > 0.0 {
				// 内蔵熱交換器の内径、管長が入力されている場合
				if Stank.KAinput[j] == 'C' {
					dblT := (Stank.DblTa + Stank.DblTw) / 2.0
					// 内蔵熱交換器の表面温度は内外流体の平均温度で代用
					ho := FNhoutpipe(Stank.Dbld0, dblT, Stank.DblTw)
					// 流速の計算
					dblv := elin.Lpath.G / Row / (math.Pi * math.Pow(Stank.Dbld0/2.0, 2.0))
					hi := FNhinpipe(Stank.Dbld0, Stank.DblL, dblv, dblT)
					Stank.KA[j] = 1.0 / (1.0/ho + 1.0/hi) * math.Pi * Stank.Dbld0 * Stank.DblL
				}
				if Stank.KAinput[j] == 'Y' || Stank.KAinput[j] == 'C' {
					NTU := Stank.KA[j] / *cGwin
					*ihxeff = 1.0 - math.Exp(-NTU)
				}
			}
			*EGwin = *cGwin * *ihxeff
		}

		stofc(Stank.Ndiv, Stank.Nin, Stank.Jin,
			Stank.Jout, Stank.Ihex, Stank.Ihxeff, Stank.Jva, Stank.Jvb,
			Stank.Mdt, Stank.KS, Stank.Cat.gxr, Stank.Tenv,
			Stank.Tssold, Stank.CGwin, Stank.EGwin, Stank.B, Stank.R, Stank.D, Stank.Fg)

		fgIdx := 0
		cfinIdx := 0
		for j := 0; j < Stank.Nin; j++ {
			Eo := Stank.Cmp.Elouts[j]
			Eo.Coeffo = 1.0
			Eo.Co = Stank.D[Stank.Jout[j]]

			for k := 0; k < Stank.Nin; k++ {
				Eo.Coeffin[cfinIdx] = -Stank.Fg[fgIdx]
				cfinIdx++
				fgIdx++
			}
		}
	}
}

/* ------------------------------------------------------- */

/*  蓄熱槽内部水温のポインター  */

func stankvptr(key []string, Stank *STANK, vptr *VPTR) int {
	err := 0
	var s string
	if key[1] == "Ts" {
		s = key[2]
		if unicode.IsLetter(rune(s[0])) {
			if s[0] == 't' {
				vptr.Ptr = &Stank.Tssold[0]
				vptr.Type = VAL_CTYPE
			} else if s[0] == 'b' {
				vptr.Ptr = &Stank.Tssold[Stank.Ndiv-1]
				vptr.Type = VAL_CTYPE
			} else {
				err = 1
			}
		} else {
			i, _ := strconv.Atoi(s)
			if i >= 0 && i < Stank.Ndiv {
				vptr.Ptr = &Stank.Tssold[i]
				vptr.Type = VAL_CTYPE
			} else {
				err = 1
			}
		}
	} else {
		err = 1
	}

	return err
}

/* ------------------------------------------------------- */

/* 槽内水温、水温分布逆転の検討 */

func Stanktss(Nstank int, Stank []STANK, TKreset *int) {
	for i := 0; i < Nstank; i++ {

		for j := 0; j < Stank[i].Nin; j++ {
			eli := Stank[i].Cmp.Elins[j]
			Stank[i].Twin[j] = eli.Sysvin
		}

		stotss(Stank[i].Ndiv, Stank[i].Nin, Stank[i].Jin, Stank[i].B, Stank[i].R, Stank[i].EGwin, Stank[i].Twin,
			Stank[i].Tss)

		stotsexm(Stank[i].Ndiv, Stank[i].Tss, &Stank[i].Jva, &Stank[i].Jvb,
			Stank[i].DtankF, &Stank[i].Cfcalc)

		if Stank[i].Cfcalc == 'y' {
			*TKreset = 1
		}
	}
}

/* ------------------------------------------------------- */

/* 供給熱量、損失熱量計算、水温前時間値の置換 */

func Stankene(Nstank int, Stank []STANK) {
	for i := 0; i < Nstank; i++ {
		// バッチモードチェック（各層が空かどうかをチェック）
		for k := 0; k < Stank[i].Ndiv; k++ {
			if Stank[i].DtankF[k] == TANK_EMPTY {
				Stank[i].Tss[k] = TANK_EMPTMP
			}
		}

		// バッチモードの水供給
		if Stank[i].Batchop == BTFILL {
			Tsm := 0.0
			for k := 0; k < Stank[i].Ndiv; k++ {
				if Stank[i].DtankF[k] == TANK_EMPTY {
					Stank[i].DtankF[k] = TANK_FULL
					for j := 0; j < Stank[i].Nin; j++ {
						if Stank[i].Batchcon[j] == BTFILL {
							Stank[i].Tss[k] = Stank[i].Twin[j]
						}
					}
				}
				Tsm += Stank[i].Tss[k]
			}
			Tsm /= float64(Stank[i].Ndiv)
			for k := 0; k < Stank[i].Ndiv; k++ {
				Stank[i].Tss[k] = Tsm
			}
		}

		for j := 0; j < Stank[i].Nin; j++ {
			Jo := Stank[i].Jout[j]
			Q := &Stank[i].Q[j]
			EGwin := Stank[i].EGwin[j]
			Twin := Stank[i].Twin[j]
			// ihex := Stank[i].Ihex[j]

			*Q = EGwin * (Stank[i].Tss[Jo] - Twin)

			// // 内蔵熱交換器の場合
			if Stank[i].KAinput[j] == 'C' {
				Stank[i].DblTa = Stank[i].Tss[Jo]
				if EGwin > 0.0 {
					Stank[i].DblTw = Twin
				}
			}
		}

		Stank[i].Qloss = 0.0
		Stank[i].Qsto = 0.0
		for j := 0; j < Stank[i].Ndiv; j++ {
			if Stank[i].DtankF[j] == TANK_FULL {
				Stank[i].Qloss += Stank[i].KS[j] * (Stank[i].Tss[j] - *Stank[i].Tenv)
				if Stank[i].Tssold[j] > -273.0 {
					Stank[i].Qsto += Stank[i].Mdt[j] * (Stank[i].Tss[j] - Stank[i].Tssold[j])
				}
			}
			Stank[i].Tssold[j] = Stank[i].Tss[j]
		}
	}
}

/* ------------------------------------------------------- */

// 代表日の出力
func stankcmpprt(fo io.Writer, id int, Nstank int, Stank []STANK) {
	switch id {
	case 0:
		if Nstank > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, Nstank)
		}

		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)
			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*5+2+stank.Ndiv+stank.Ncalcihex)
		}
	case 1:
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_c c c %s:%c_G m f %s:%c_Ti t f %s:%c_To t f %s:%c_Q q f  ",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				if stank.KAinput[i] == 'C' {
					fmt.Fprintf(fo, "%s:%c_KA q f  ", stank.Name, c)
				}
				fmt.Fprintln(fo)
			}
			fmt.Fprintf(fo, "%s_Qls q f %s_Qst q f\n ", stank.Name, stank.Name)
			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, "%s_Ts[%d] t f ", stank.Name, i+1)
			}
			fmt.Fprintln(fo)
		}
	default:
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			Tss := &stank.Tss[0]
			for i := 0; i < stank.Nin; i++ {
				Ei := stank.Cmp.Elins[i]
				Twin := &stank.Twin[i]
				Q := &stank.Q[i]
				Eo := stank.Cmp.Elouts[i]
				fmt.Fprintf(fo, "%c %.5g %4.1f %4.1f %3.0f  ", Ei.Lpath.Control,
					Eo.G, *Twin, Eo.Sysv, *Q)

				if stank.KAinput[i] == 'C' {
					if Eo.G > 0.0 {
						fmt.Fprintf(fo, "%.2f  ", stank.KA[i])
					} else {
						fmt.Fprintf(fo, "%.2f  ", 0.0)
					}
				}
			}
			fmt.Fprintf(fo, "%2.0f %3.0f\n", stank.Qloss, stank.Qsto)

			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, " %4.1f", *Tss)
				Tss = &stank.Tss[i+1]
			}
			fmt.Fprintln(fo)
		}
	}
}

/* ------------------------------------------------------- */
func stankivprt(fo io.Writer, id int, Nstank int, Stank []STANK) {
	if id == 0 && Nstank > 0 {
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "m=%d  %s  %d\n", m, stank.Name, stank.Ndiv)
		}
	} else {
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "m=%d  ", m)

			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, " %5.1f", stank.Tss[i])
			}
			fmt.Fprintln(fo)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func stankdyint(Nstank int, Stank []STANK) {
	for i := 0; i < Nstank; i++ {
		stank := &Stank[i]
		stank.Qlossdy = 0.0
		stank.Qstody = 0.0

		for j := 0; j < stank.Nin; j++ {
			s := &stank.Stkdy[j]
			svdyint(&s.Tidy)
			svdyint(&s.Tsdy)
			qdyint(&s.Qdy)
		}
	}
}

func stankmonint(Nstank int, Stank []STANK) {
	for i := 0; i < Nstank; i++ {
		stank := &Stank[i]
		stank.MQlossdy = 0.0
		stank.MQstody = 0.0

		for j := 0; j < stank.Nin; j++ {
			s := &stank.Mstkdy[j]
			svdyint(&s.Tidy)
			svdyint(&s.Tsdy)
			qdyint(&s.Qdy)
		}
	}
}

// 日集計、月集計
func stankday(Mon, Day, ttmm, Nstank int, Stank []STANK, Nday, SimDayend int) {
	for i := 0; i < Nstank; i++ {
		stank := &Stank[i]

		// 日集計
		Ts := 0.0

		S := &stank.Stkdy[0]
		for j := 0; j < stank.Ndiv; j++ {
			Ts += stank.Tss[j] / float64(stank.Ndiv)
		}
		svdaysum(int64(ttmm), ON_SW, Ts, &S.Tsdy)

		stank.Qlossdy += stank.Qloss
		stank.Qstody += stank.Qsto

		for j := 0; j < stank.Nin; j++ {
			Ei := stank.Cmp.Elins[j]
			S := &stank.Stkdy[j]
			svdaysum(int64(ttmm), Ei.Lpath.Control, stank.Twin[j], &S.Tidy)
			qdaysum(int64(ttmm), Ei.Lpath.Control, stank.Q[j], &S.Qdy)
		}

		// 月集計
		S = &stank.Mstkdy[0]
		svmonsum(Mon, Day, ttmm, ON_SW, Ts, &S.Tsdy, Nday, SimDayend)

		stank.MQlossdy += stank.Qloss
		stank.MQstody += stank.Qsto

		for j := 0; j < stank.Nin; j++ {
			Ei := stank.Cmp.Elins[j]
			S := &stank.Mstkdy[j]
			svmonsum(Mon, Day, ttmm, Ei.Lpath.Control, stank.Twin[j], &S.Tidy, Nday, SimDayend)
			qmonsum(Mon, Day, ttmm, Ei.Lpath.Control, stank.Q[j], &S.Qdy, Nday, SimDayend)
		}
	}
}

// 日集計の出力
func stankdyprt(fo io.Writer, id, Nstank int, Stank []STANK) {
	switch id {
	case 0:
		if Nstank > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, Nstank)
		}

		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)

			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*14+2+1)
		}

	case 1:
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "%s_Ts t f \n", stank.Name)

			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
			}
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n", stank.Name, stank.Name)
		}

	default:
		for _, stank := range Stank {
			S := &stank.Stkdy[0]

			fmt.Fprintf(fo, "%.1f\n", S.Tsdy.M)
			for j := 0; j < stank.Nin; j++ {
				S := &stank.Stkdy[j]

				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					S.Tidy.Hrs, S.Tidy.M,
					S.Tidy.Mntime, S.Tidy.Mn,
					S.Tidy.Mxtime, S.Tidy.Mx)

				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Hhr, S.Qdy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Chr, S.Qdy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Hmxtime, S.Qdy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Cmxtime, S.Qdy.Cmx)
			}
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stank.Qlossdy*Cff_kWh, stank.Qstody*Cff_kWh)
		}
	}
}

// 月集計の出力
func stankmonprt(fo io.Writer, id, Nstank int, Stank []STANK) {
	switch id {
	case 0:
		if Nstank > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, Nstank)
		}

		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)

			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*14+2+1)
		}

	case 1:
		for m := 0; m < Nstank; m++ {
			stank := &Stank[m]
			fmt.Fprintf(fo, "%s_Ts t f \n", stank.Name)

			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
			}
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n", stank.Name, stank.Name)
		}

	default:
		for _, stank := range Stank {
			S := &stank.Mstkdy[0]

			fmt.Fprintf(fo, "%.1f\n", S.Tsdy.M)
			for j := 0; j < stank.Nin; j++ {
				S := &stank.Mstkdy[j]

				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					S.Tidy.Hrs, S.Tidy.M,
					S.Tidy.Mntime, S.Tidy.Mn,
					S.Tidy.Mxtime, S.Tidy.Mx)

				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Hhr, S.Qdy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Chr, S.Qdy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Hmxtime, S.Qdy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Cmxtime, S.Qdy.Cmx)
			}
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stank.MQlossdy*Cff_kWh, stank.MQstody*Cff_kWh)
		}
	}
}
