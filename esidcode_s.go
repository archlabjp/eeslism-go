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

/*     esidcode_s.c       */

package main

import "fmt"

/*  スケジュ－ル名からスケジュ－ル番号の検索   */
/* --------------------------------------------*/

func idssn(code string, _Seasn []SEASN, err string) int {
	N := _Seasn[0].end
	var j int
	for j = 0; j < N; j++ {
		Seasn := &_Seasn[j]
		if code == Seasn.name {
			break
		}
	}
	if j == N {
		fmt.Println(err)
	}
	return j
}

/* ---------------------------------------- */

func idwkd(code string, Wkdy []WKDY, err string) int {
	N := Wkdy[0].end
	var j int
	for j = 0; j < N; j++ {
		_Wkdy := &Wkdy[j]
		if code == _Wkdy.name {
			break
		}
	}
	if j == N {
		fmt.Println(err)
	}
	return j
}

/* ---------------------------------------- */

func iddsc(code string, Dsch []DSCH, err string) int {
	N := Dsch[0].end
	var j int
	for j = 0; j < N; j++ {
		_Dsch := &Dsch[j]
		if code == _Dsch.name {
			break
		}
	}
	if j == N {
		fmt.Println(err)
	}
	return j
}

/* ---------------------------------------- */

func iddsw(code string, Dscw []DSCW, err string) int {
	N := Dscw[0].end
	var j int
	for j = 0; j < N; j++ {
		_Dscw := &Dscw[j]
		if code == _Dscw.name {
			break
		}
	}
	if j == N {
		fmt.Println(err)
	}
	return j
}

/* ---------------------------------------- */

/* ---------------------------------------- */

func idsch(code string, Sch []SCH, err string) int {
	N := Sch[0].end
	var j int
	for j = 0; j < N; j++ {
		_Sch := &Sch[j]
		if code == _Sch.name {
			break
		}
	}
	if j == N {
		j = -1
		if err != "" {
			Eprint("<idsch>", err)
		}
	}
	return j
}

/* ---------------------------------------- */

func idscw(code string, Scw []SCH, err string) int {
	N := Scw[0].end
	var j int
	for j = 0; j < N; j++ {
		_Scw := &Scw[j]
		if code == _Scw.name {
			break
		}
	}
	if j == N {
		j = -1
		if err != "" {
			Eprint("<idscw>", err)
		}
	}
	return j
}

/* ---------------------------------------- */

// 室名 `code` に一致する部屋を 部屋の一覧 `Room` から検索し、その番号を返す
func idroom(code string, Room []ROOM, err string) int {
	N := Room[0].end

	var j int
	for j = 0; j < N; j++ {
		_Room := &Room[j]
		if code == _Room.Name {
			break
		}
	}

	if j == N {
		E := fmt.Sprintf("Room=%s %s", code, err)
		Eprint("<idroom>", E)
	}

	return j
}

/* ---------------------------------------- */

/* ---------------------------------------- */
