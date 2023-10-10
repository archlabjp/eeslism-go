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
	"strings"
)

/* --------------------------------------------------------- */
/*
居住者スケジュ－ルの入力              */

func Residata(fi *EeTokens, dsn string, schdl *SCHDL, rooms []ROOM, pmvpri *int, simc *SIMCONTL) {
	errFmt := fmt.Sprintf(ERRFMT, dsn)
	vall := schdl.Val

	for fi.IsEnd() == false {
		var s, ss, sss, s4 string
		s = fi.GetToken()
		if s == "*" {
			break
		}

		errMsg := errFmt + s
		i := idroom(s, rooms, errMsg)
		rm := rooms[i]

		for {
			s = fi.GetToken()
			if s == ";" {
				break
			}

			errMsg := errFmt + s
			ce := strings.Index(s, ";")
			st := strings.Index(s, "=")
			stVal := s[st+1 : ce]

			switch s[:st] {
			case "H":
				fmt.Sscanf(stVal, "(%f,%[^,],%[^)])", &rm.Nhm, &ss, &sss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.Hmsch = &vall[k]
				} else {
					rm.Hmsch = envptr(ss, simc, 0, nil, nil, nil)
				}

				if k := idsch(sss, schdl.Sch, ""); k >= 0 {
					rm.Hmwksch = &vall[k]
				} else {
					rm.Hmwksch = envptr(sss, simc, 0, nil, nil, nil)
				}
			case "comfrt":
				fmt.Sscanf(stVal, "(%[^,],%[^,],%[^)])", &ss, &sss, &s4)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.Metsch = &vall[k]
				} else {
					rm.Metsch = envptr(ss, simc, 0, nil, nil, nil)
				}

				if k := idsch(sss, schdl.Sch, ""); k >= 0 {
					rm.Closch = &vall[k]
				} else {
					rm.Closch = envptr(sss, simc, 0, nil, nil, nil)
				}

				if k := idsch(s4, schdl.Sch, ""); k >= 0 {
					rm.Wvsch = &vall[k]
				} else {
					rm.Wvsch = envptr(s4, simc, 0, nil, nil, nil)
				}

				*pmvpri = 1
				if SETprint == 1 {
					rm.setpri = 1
				}
			default:
				Eprint("<Residata>", errMsg)
			}

			if ce != -1 {
				break
			}
		}
	}
}

/* --------------------------------------------------------- */
/*
照明・機器利用スケジュ－ルの入力              */

func Appldata(fi *EeTokens, dsn string, schdl *SCHDL, rooms []ROOM, simc *SIMCONTL) {
	errFmt := fmt.Sprintf(ERRFMT, dsn)
	vall := schdl.Val

	for fi.IsEnd() == false {
		var s, ss string
		s = fi.GetToken()
		if s == "*" {
			break
		}

		errMsg := errFmt + s
		i := idroom(s, rooms, errMsg)
		rm := rooms[i]

		for fi.IsEnd() == false {
			s = fi.GetToken()
			if s == ";" {
				break
			}

			errMsg := errFmt + s
			ce := strings.Index(s, ";")
			st := strings.Index(s, "=")
			stVal := s[st+1 : ce]

			switch s[:st] {
			case "L":
				fmt.Sscanf(stVal, "(%f,%c,%[^)])", &rm.Light, &rm.Ltyp, &ss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.Lightsch = &vall[k]
				} else {
					rm.Lightsch = envptr(ss, simc, 0, nil, nil, nil)
				}
			case "As":
				fmt.Sscanf(stVal, "(%f,%f,%[^)])", &rm.Apsc, &rm.Apsr, &ss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.Assch = &vall[k]
				} else {
					rm.Assch = envptr(ss, simc, 0, nil, nil, nil)
				}
			case "Al":
				fmt.Sscanf(stVal, "(%f,%[^)])", &rm.Apl, &ss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.Alsch = &vall[k]
				} else {
					rm.Alsch = envptr(ss, simc, 0, nil, nil, nil)
				}
			case "AE":
				fmt.Sscanf(stVal, "(%f,%[^)])", &rm.AE, &ss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.AEsch = &vall[k]
				} else {
					rm.AEsch = envptr(ss, simc, 0, nil, nil, nil)
				}
			case "AG":
				fmt.Sscanf(stVal, "(%f,%[^)])", &rm.AG, &ss)

				if k := idsch(ss, schdl.Sch, ""); k >= 0 {
					rm.AGsch = &vall[k]
				} else {
					rm.AGsch = envptr(ss, simc, 0, nil, nil, nil)
				}
			default:
				Eprint("<Appldata>", errMsg)
			}

			if ce != -1 {
				break
			}
		}
	}
}
