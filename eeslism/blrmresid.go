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

/*  rmresid.c   */
package eeslism

import (
	"fmt"
	"regexp"
	"strings"
)

/* --------------------------------------------------------- */
/*
居住者スケジュ－ルの入力              */

func Residata(fi *EeTokens, dsn string, schdl *SCHDL, rooms []*ROOM, pmvpri *int, simc *SIMCONTL) {
	errFmt := fmt.Sprintf(ERRFMT, dsn)

	for fi.IsEnd() == false {
		var s, ss, sss, s4 string
		s = fi.GetToken()
		if s == "*" {
			break
		}
		if s == "\n" {
			continue
		}
		if s == "%s" || s == "%sn" {
			fi.GetLogicalLine()
			continue
		}

		errMsg := errFmt + s
		i, err := idroom(s, rooms, errMsg)
		if err != nil {
			panic(err)
		}
		rm := rooms[i]

		for {
			s = fi.GetToken()
			if s == ";" {
				break
			}

			errMsg := errFmt + s
			st := strings.Index(s, "=")
			stVal := s[st+1:]

			switch s[:st] {
			case "H":
				// 人体発熱
				// H=(基準人数(人),在室率設定値名,作業強度設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 4 {
					// 基準人数(人)
					rm.Nhm, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 在室率設定値名
					ss = match[2]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.Hmsch = &schdl.Val[k]
					} else {
						rm.Hmsch = envptr(ss, simc, nil, nil, nil)
					}

					// 作業強度設定値名
					sss = match[3]
					if k, err := idsch(sss, schdl.Sch, ""); err == nil {
						rm.Hmwksch = &schdl.Val[k]
					} else {
						rm.Hmwksch = envptr(sss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}
			case "comfrt":
				// 熱環境条件
				// comfrt=(代謝率設定値名,着衣量設定値名,室内風速設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 4 {
					// 代謝率(Met値)設定値名
					ss = match[1]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.Metsch = &schdl.Val[k]
					} else {
						rm.Metsch = envptr(ss, simc, nil, nil, nil)
					}

					// 着衣量(Clo値)設定値名
					sss = match[2]
					if k, err := idsch(sss, schdl.Sch, ""); err == nil {
						rm.Closch = &schdl.Val[k]
					} else {
						rm.Closch = envptr(sss, simc, nil, nil, nil)
					}

					// 室内風速設定値名
					s4 = match[3]
					if k, err := idsch(s4, schdl.Sch, ""); err == nil {
						rm.Wvsch = &schdl.Val[k]
					} else {
						rm.Wvsch = envptr(s4, simc, nil, nil, nil)
					}

					*pmvpri = 1
					if SETprint {
						rm.setpri = true
					}
				} else {
					fmt.Println("No match found.")
				}

			default:
				Eprint("<Residata>", errMsg)
			}
		}
	}
}

/* --------------------------------------------------------- */
/*
照明・機器利用スケジュ－ルの入力              */

func Appldata(fi *EeTokens, dsn string, schdl *SCHDL, rooms []*ROOM, simc *SIMCONTL) {
	errFmt := fmt.Sprintf(ERRFMT, dsn)

	for fi.IsEnd() == false {
		var s, ss string
		s = fi.GetToken()
		if s == "*" {
			break
		}
		if s == "\n" {
			continue
		}
		if s == "%s" || s == "%sn" {
			fi.GetLogicalLine()
			continue
		}

		errMsg := errFmt + s
		i, err := idroom(s, rooms, errMsg)
		if err != nil {
			panic(err)
		}
		rm := rooms[i]

		for fi.IsEnd() == false {
			s = fi.GetToken()
			if s == ";" {
				break
			}

			errMsg := errFmt + s
			st := strings.Index(s, "=")
			stVal := s[st+1:]

			switch s[:st] {
			case "L":
				// 照明
				// L=(基準値[W],器具タイプ,照明入力設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 4 {
					// 基準値[W]
					rm.Light, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 器具タイプ
					rm.Ltyp = rune(match[2][0])

					// 照明入力設定値名
					ss = match[3]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.Lightsch = &schdl.Val[k]
					} else {
						rm.Lightsch = envptr(ss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "As":
				// 機器顕熱
				// As=(対流成分基準値[W],放射成分基準値[W],設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 4 {
					// 対流成分基準値[W]
					rm.Apsc, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 放射成分基準値[W]
					rm.Apsr, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 設定値名
					ss = match[3]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.Assch = &schdl.Val[k]
					} else {
						rm.Assch = envptr(ss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "Al":
				// 機器潜熱
				// Al=(基準値[W],設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 3 {
					// 基準値[W]
					rm.Apl, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 設定値名
					ss = match[2]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.Alsch = &schdl.Val[k]
					} else {
						rm.Alsch = envptr(ss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "AE":
				// 電力に関する集計を行う
				// AE=(基準値[W],電力設定値名
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 3 {
					// 基準値[W]
					rm.AE, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 電力設定値名
					ss = match[2]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.AEsch = &schdl.Val[k]
					} else {
						rm.AEsch = envptr(ss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "AG":
				// ガスに関する集計を行う
				// AG=(基準値[W],ガス設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(stVal)
				if len(match) == 3 {
					// 基準値[W]
					rm.AG, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// ガス設定値名
					ss = match[2]
					if k, err := idsch(ss, schdl.Sch, ""); err == nil {
						rm.AGsch = &schdl.Val[k]
					} else {
						rm.AGsch = envptr(ss, simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			default:
				Eprint("<Appldata>", errMsg)
			}
		}
	}
}
