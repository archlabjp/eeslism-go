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
	"strings"
)

/* ------------------------------------------------------------------ */

/*  外気導入量および室間相互換気量の設定スケジュ－ル入力   */

func Ventdata(fi *EeTokens, dsn string, Schdl *SCHDL, Room []ROOM, Simc *SIMCONTL) {
	var achr *ACHIR
	var room, Rm *ROOM
	var name1, name2, s, ss, E string
	var val float64
	var i, j, v, k int

	E = fmt.Sprintf(ERRFMT, dsn)
	for fi.IsEnd() == false {
		name1 = fi.GetToken()

		i = idroom(name1, Room, E+name1)
		Rm = &Room[i]

		s = fi.GetToken()
		st := strings.IndexByte(s, '=')
		if st != -1 {
			for s != "" {
				_ss := strings.SplitN(s, "=", 2)
				key := _ss[0]
				valstr := _ss[1]
				switch key {
				case "Vent":
					var err error
					_, err = fmt.Sscanf(valstr, "(%f,%[^)])", &val, &ss)
					if err != nil {
						panic(err)
					}
					Rm.Gve = val

					if k = idsch(ss, Schdl.Sch, ""); k >= 0 {
						Rm.Vesc = &Schdl.Val[k]
					} else {
						Rm.Vesc = envptr(ss, Simc, 0, nil, nil, nil)
					}
				case "Inf":
					var err error
					_, err = fmt.Sscanf(valstr, "(%f,%[^)])", &val, &ss)
					if err != nil {
						panic(err)
					}
					Rm.Gvi = val

					if k = idsch(ss, Schdl.Sch, ""); k >= 0 {
						Rm.Visc = &Schdl.Val[k]
					} else {
						Rm.Visc = envptr(ss, Simc, 0, nil, nil, nil)
					}
				default:
					err := fmt.Sprintf("Room=%s  %s", Rm.Name, key)
					Eprint("<Ventedata>", err)
				}

				if st = strings.IndexByte(valstr, ';'); st != -1 {
					break
				}
				s = fi.GetToken()
			}
		} else {
			c := s[0]
			name2 = fi.GetToken()
			if fi.GetToken() != "v=" {
				panic("Invalid format of ventdata")
			}
			s = fi.GetToken()
			if c == ';' {
				break
			}
			j = idroom(name2, Room, E+name2)

			if ce := strings.IndexByte(s, ';'); ce != -1 {
				s = s[:ce]
			}
			v = idsch(s, Schdl.Sch, E+s)

			switch c {
			case '-':
				room = &Room[j]
				achr = &room.achr[room.Nachr]
				achr.rm = i
				achr.room = &Room[i]
				achr.sch = v
				room.Nachr++

				room = &Room[i]
				achr = &room.achr[room.Nachr]
				achr.rm = j
				achr.room = &Room[j]
				achr.sch = v
				room.Nachr++
			default:
				panic(c)
			}
		}
	}
}

/* ------------------------------------------------------------------ */

/*  室間相互換気量の設定   */

func Aichschdlr(val []float64, Nroom int, rooms []ROOM) {
	for i := 0; i < Nroom; i++ {
		room := &rooms[i]

		for j := 0; j < room.Nachr; j++ {
			achr := &room.achr[j]
			v := val[achr.sch]
			if v > 0.0 {
				achr.Gvr = v
			} else {
				achr.Gvr = 0.0
			}
		}
	}
}
