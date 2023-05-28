package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Exsfdata(section *EeTokens, dsn string, Exsf *EXSFS, Schdl *SCHDL, Simc *SIMCONTL) {
	var s, ename string
	//var st *string
	var dt, dfrg, rad, wa, wb, swa, cwa, swb, cwb float64
	var vall []float64
	var i, j, k int
	var ex *EXSF

	s = dsn
	Nd := ExsfCount(section)
	if Nd == 0 {
		Nd = 1
	}

	vall = Schdl.Val

	var exs []EXSF
	if Nd > 0 {
		exs = make([]EXSF, 0, Nd+1)

		s = fmt.Sprintf("%f", ALO)
		Exsf.Alosch = envptr(s, Simc, 0, nil, nil, nil)
		Exsf.Alotype = 'F'

		// for i := range exs {
		// 	Exsfinit(&exs[i])
		// }
	}

	line := section.GetLogicalLine()

	for _, s := range line[1:] {
		if strings.HasPrefix(s, "alo=") {
			if s[4:] == "Calc" {
				Exsf.Alotype = 'V'
			} else if k = idsch(s[4:], Schdl.Sch, ""); k >= 0 {
				Exsf.Alosch = &vall[k]
				Exsf.Alotype = 'S'
			} else {
				Exsf.Alosch = envptr(s[4:], Simc, 0, nil, nil, nil)
				if Exsf.Alosch != nil {
					Exsf.Alotype = 'S'
				}
			}
		} else if strings.HasPrefix(s, "r=") {
			dfrg, _ = strconv.ParseFloat(s[2:], 64)
			if dfrg < 0.0 || dfrg > 1.0 {
				fmt.Fprintf(os.Stderr, "%s の設置値が不適切です", s)
				os.Exit(1)
			}
		}
	}

	for section.IsEnd() == false {

		line = section.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		ex := new(EXSF)
		Exsfinit(ex)

		i++
		ex.Name = line[0]
		ex.Alotype = Exsf.Alotype
		ex.Alo = Exsf.Alosch
		if s == "Hor" {
			ex.Wb = 0.0
		} else if s == "EarthSf" {
			Exsf.EarthSrfFlg = 'Y'
			ex.Typ = 'e'
		} else {
			ex.Wb = 90.0
			ex.Rg = dfrg
		}

		for _, s := range line[1:] {
			if strings.HasPrefix(s, "a=") {
				var err error
				if dt, err = strconv.ParseFloat(s[2:], 64); err == nil {
					ex.Wa = dt
				} else {
					var dir rune = ' '
					if strings.Contains(s, "+") {
						st := strings.IndexRune(s, '+')
						dir = '+'
						ename = s[2:st]
						dt, err = strconv.ParseFloat(s[st+1:], 64)
						if err != nil {
							panic(err)
						}
					} else if strings.Contains(s, "-") {
						st := strings.IndexRune(s, '-')
						dir = '-'
						ename = s[2:st]
						dt, err = strconv.ParseFloat(s[st+1:], 64)
						if err != nil {
							panic(err)
						}
					} else {
						ename = s[2:]
					}

					for j := range exs {
						exj := &exs[j]
						if exj.Name == ename {
							if dir == '+' {
								ex.Wa = exj.Wa + dt
							} else if dir == '-' {
								ex.Wa = exj.Wa - dt
							} else {
								ex.Wa = exj.Wa
							}
							break
						}
					}
					if j == i+1 {
						Eprint("<Exsfdata>", s)
					}
				}
			} else {
				st := strings.IndexRune(s, '=')
				if strings.HasPrefix(s, "alo") {
					if s[st+1:] == "Calc" {
						ex.Alotype = 'V'
					} else if k = idsch(s[st+1:], Schdl.Sch, ""); k >= 0 {
						ex.Alo = &vall[k]
						ex.Alotype = 'S'
					} else {
						ex.Alo = envptr(s[st+1:], Simc, 0, nil, nil, nil)
						ex.Alotype = 'S'
					}
				} else {
					dt, _ = strconv.ParseFloat(s[st+1:], 64)
					switch s[0] {
					case 't':
						ex.Wb = dt
					case 'r':
						ex.Rg = dt
					case 'Z':
						ex.Z = dt
						ex.Typ = 'E'
					case 'd':
						ex.Erdff = dt
					default:
						Eprint("<Exsfdata>", s)
					}
				}
			}
		}

		exs = append(exs, *ex)
		ex.End = i
	}

	//Nd = i
	if Nd > 0 {
		s = strconv.FormatFloat(ALO, 'f', -1, 64)
		Exsf.Alosch = envptr(s, Simc, 0, nil, nil, nil)
		Exsf.Alotype = 'F'
		Exsf.Exs = exs
	}

	Exsf.Nexs = i
	Exsf.Exs[0].End = i

	for i = 0; i < Exsf.Nexs; i++ {
		ex = &Exsf.Exs[i]
		if ex.Typ == 'S' {
			wa = ex.Wa * rad
			wb = ex.Wb * rad
			cwa = math.Cos(wa)
			swa = math.Sin(wa)
			cwb = math.Cos(wb)
			swb = math.Sin(wb)
			ex.Wz = cwb
			ex.Ww = swb * swa
			ex.Ws = swb * cwa
			ex.Cbsa = cwb * swa
			ex.Cbca = cwb * cwa
			ex.Fs = 0.5 * (1.0 + cwb)
		}
	}
}

/*  外表面入射日射量の計算    */
func Exsfsol(Nexs int, Wd *WDAT, Exs []EXSF) {
	var cinc float64

	for i := 0; i < Nexs; i++ {
		ex := &Exs[i]
		if ex.Typ == 'S' {
			cinc = Wd.Sh*ex.Wz + Wd.Sw*ex.Ww + Wd.Ss*ex.Ws
			if cinc > 0.0 {
				ex.Tprof = (Wd.Sh*ex.Swb - Wd.Sw*ex.Cbsa - Wd.Ss*ex.Cbca) / cinc
				// プロファイル角の計算
				ex.Prof = math.Atan(ex.Tprof)
				ex.Tazm = (Wd.Sw*ex.Cwa - Wd.Ss*ex.Swa) / cinc
				// 見かけの方位角の計算
				ex.Gamma = math.Atan(ex.Tazm)
				ex.Cinc = cinc
			} else {
				ex.Prof = 0.0
				ex.Gamma = 0.0
				ex.Cinc = 0.0
			}
			ex.Idre = Wd.Idn * cinc // 外表面入射（直達）
			ex.Idf = Wd.Isky*ex.Fs + ex.Rg*Wd.Ihor*(1.0-ex.Fs)
			ex.Iw = ex.Idre + ex.Idf
			ex.Rn = Wd.RN * ex.Fs
		}
	}
}

/*  ガラス日射熱取得の計算         */
func Glasstga(Ag, tgtn, Bn, cinc, Fsdw, Idr, Idf float64, Qgt, Qga *float64, Cidtype string, Profile, Gamma float64) {
	var Cid, Cidf, Bid, Bidf, Qt, Qb float64

	Cid = 0.0
	Bid = 0.0
	Cidf = 0.01
	Bidf = 0.0
	Qt = 0.0
	Qb = 0.0

	if Cidtype == "N" {
		Cid = Glscid(cinc)
		Cidf = 0.91

		Bid = Cid
		Bidf = Cidf
	} else {
		fmt.Printf("xxxxx <eebslib.c  CidType=%s\n", Cidtype)
	}

	Qt = Ag * (Cid*Idr*(1.0-Fsdw) + Cidf*Idf)
	Qb = Ag * (Bid*Idr*(1.0-Fsdw) + Bidf*Idf)

	*Qgt = Qt * tgtn
	*Qga = Qb * Bn
}

/*  ガラスの直達日射透過率標準特性　　　　*/
func Glscid(cinc float64) float64 {
	return math.Max(0, cinc*(3.4167+cinc*(-4.389+cinc*(2.4948-0.5224*cinc))))
}

/*  ガラスの直達日射透過率標準特性　　　　*/
// 普通複層ガラス

func GlscidDG(cinc float64) float64 {
	return math.Max(0, cinc*(0.341819+cinc*(6.070709+cinc*(-9.899236+4.495774*cinc))))
}

func ExsfCount(section *EeTokens) int {
	var N int
	for section.IsEnd() == false {
		if section.GetToken() == ";" {
			N++
			break
		}
	}
	section.Reset()
	return N
}
