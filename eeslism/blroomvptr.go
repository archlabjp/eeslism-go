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

package eeslism

import "errors"

/* 室及び関連システム変数、内部変数のポインター  */

func roomvptr(Nk int, key []string, Room *ROOM) (VPTR, error) {
	var vptr VPTR
	vptr.Ptr = nil

	if Nk == 2 {
		switch string(key[1]) {
		case "Tr":
			vptr.Ptr = &Room.Tr
			vptr.Type = VAL_CTYPE
		case "xr":
			vptr.Ptr = &Room.xr
			vptr.Type = VAL_CTYPE
		case "RH":
			vptr.Ptr = &Room.RH
			vptr.Type = VAL_CTYPE
		case "PMV":
			vptr.Ptr = &Room.PMV
			vptr.Type = VAL_CTYPE
		case "Tsav":
			vptr.Ptr = &Room.Tsav
			vptr.Type = VAL_CTYPE
		case "Tot":
			vptr.Ptr = &Room.Tot
			vptr.Type = VAL_CTYPE
		case "hr":
			vptr.Ptr = &Room.hr
			vptr.Type = VAL_CTYPE
		}
	} else if Nk == 3 {
		for i := 0; i < Room.N; i++ {
			Sd := &Room.rsrf[i]
			if string(key[1]) == Sd.Name {
				switch string(key[2]) {
				case "Ts":
					vptr.Ptr = &Sd.Ts
					vptr.Type = VAL_CTYPE
				case "Tmrt":
					vptr.Ptr = &Sd.Tmrt
					vptr.Type = VAL_CTYPE
				case "Te":
					vptr.Ptr = &Sd.Tcole
					vptr.Type = VAL_CTYPE
				}
			}
		}
	}

	if vptr.Ptr == nil {
		return vptr, errors.New("roomvptr error")
	}

	return vptr, nil
}

/* ------------------------------------------- */

/* 室負荷計算時の設定値ポインター */

func roomldptr(load *rune, key []string, Room *ROOM, idmrk *byte) (VPTR, error) {
	var err error
	var i int
	var Sd *RMSRF
	var vptr VPTR

	if key[1] == "Tr" {
		vptr.Ptr = &Room.rmld.Tset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadt = load
		Room.rmld.tropt = 'a'
		*idmrk = 't'
	} else if key[1] == "Tot" {
		vptr.Ptr = &Room.rmld.Tset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadt = load
		Room.rmld.tropt = 'o'
		*idmrk = 't'
	} else if key[1] == "RH" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'r'
		*idmrk = 'x'
	} else if key[1] == "Tdp" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'd'
		*idmrk = 'x'
	} else if key[1] == "xr" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'x'
		*idmrk = 'x'
	} else if len(key) > 2 && key[2] == "Ts" {

		for i = 0; i < Room.N; i++ {
			Sd = &Room.rsrf[i]

			if Sd.Name == key[1] {
				vptr.Ptr = &Sd.Ts
				vptr.Type = VAL_CTYPE
				Room.rmld.loadt = load
				*idmrk = 't'
				err = nil
				break
			}
			err = errors.New("Surface not found: " + Sd.Name)
		}
	} else {
		err = errors.New("'Tr', 'Tot', 'RH', 'Tdp', 'xr' or '<roomname> Ts' are expected")
	}

	return vptr, err
}

/* ------------------------------------------- */

/* 室負荷計算時のスケジュール設定 */

func roomldschd(Room *ROOM) {
	var Eo *ELOUT
	var rmld *RMLOAD

	if rmld = Room.rmld; rmld != nil {
		Eo = Room.cmp.Elouts[0]
		if rmld.loadt != nil {
			if Eo == Eo.Eldobj || Eo.Eldobj.Control != OFF_SW {
				if rmld.Tset > TEMPLIMIT {
					Eo.Sysv = rmld.Tset
					Room.Tr = rmld.Tset
					Eo.Control = LOAD_SW
				} else {
					if Room.VAVcontrl != nil {
						Room.VAVcontrl.Cmp.Control = OFF_SW
						Room.VAVcontrl.Cmp.Elouts[0].Control = OFF_SW
					}
				}
			}
		}

		Eo = Room.cmp.Elouts[1]
		if rmld.loadx != nil {
			if Eo == Eo.Eldobj || Eo.Eldobj.Control != OFF_SW {
				switch rmld.hmopt {
				case 'r':
					if rmld.Xset > 0.0 {
						Eo.Sysv = FNXtr(Room.Tr, rmld.Xset)
						Eo.Control = LOAD_SW
					}
				case 'd':
					if rmld.Xset > TEMPLIMIT {
						Eo.Sysv = FNXp(FNPws(rmld.Xset))
						Eo.Control = LOAD_SW
					}
				case 'x':
					if rmld.Xset > 0.0 {
						Eo.Sysv = rmld.Xset
						Eo.Control = LOAD_SW
					}
				}
			}
		}
	}
}
