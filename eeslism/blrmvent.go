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

/* rmvent.c  */
package eeslism

import (
	"fmt"
	"regexp"
	"strings"
)

/* ------------------------------------------------------------------ */

/*  外気導入量および室間相互換気量の設定スケジュ－ル入力   */

// VENTデータセット
func Ventdata(fi *EeTokens, dsn string, Schdl *SCHDL, Room []*ROOM, Simc *SIMCONTL) {
	var Rm *ROOM
	var name1, ss, E string
	var k int

	E = fmt.Sprintf(ERRFMT, dsn)
	for fi.IsEnd() == false {
		line := fi.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		// 室名
		name1 = line[0]

		// 室検索
		i, err := idroom(name1, Room, E+name1)
		if err != nil {
			panic(err)
		}
		Rm = Room[i] //室の参照

		for _, s := range line[1:] {
			_ss := strings.SplitN(s, "=", 2)
			key := _ss[0]
			valstr := _ss[1]

			switch key {
			case "Vent":
				// 換気量
				// Vent=(基準値[kg/s],換気量設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(valstr)
				if len(match) == 3 {
					// 基準値[kg/s]
					Rm.Gve, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 換気量設定値名
					ss = match[2]
					if k, err := idsch(ss, Schdl.Sch, ""); err == nil {
						Rm.Vesc = &Schdl.Val[k]
					} else {
						Rm.Vesc = envptr(ss, Simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "Inf":
				// すきま風
				// Inf=(基準値[kg/s],隙間風量設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(valstr)
				if len(match) == 3 {
					// 基準値[kg/s]
					Rm.Gvi, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 隙間風量設定値名
					ss = match[2]
					if k, err = idsch(ss, Schdl.Sch, ""); err == nil {
						Rm.Visc = &Schdl.Val[k]
					} else {
						Rm.Visc = envptr(ss, Simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			default:
				err := fmt.Sprintf("Room=%s  %s", Rm.Name, key)
				Eprint("<Ventedata>", err)
			}
		}
	}
}

/* ------------------------------------------------------------------ */

/*  室間相互換気量の設定   */

func Aichschdlr(val []float64, rooms []*ROOM) {
	for i := range rooms {
		room := rooms[i]

		for j := 0; j < room.Nachr; j++ {
			achr := room.achr[j]
			v := val[achr.sch]
			if v > 0.0 {
				achr.Gvr = v
			} else {
				achr.Gvr = 0.0
			}
		}
	}
}
