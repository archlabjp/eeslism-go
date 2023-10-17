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

package eeslism

import (
	"errors"
	"fmt"
)

/*  スケジュ－ル名からスケジュ－ル番号の検索   */
/* --------------------------------------------*/

func idssn(code string, _Seasn []SEASN) (int, error) {
	N := len(_Seasn)

	for j := 0; j < N; j++ {
		Seasn := &_Seasn[j]
		if code == Seasn.name {
			return j, nil
		}
	}

	return -1, errors.New("SEASN Not Found")
}

/* ---------------------------------------- */

func idwkd(code string, Wkdy []WKDY) (int, error) {
	N := len(Wkdy)

	if N != len(Wkdy) {
		panic("N != len(Wkdy)")
	}

	for j := 0; j < N; j++ {
		_Wkdy := &Wkdy[j]
		if code == _Wkdy.name {
			return j, nil
		}
	}

	return -1, errors.New("WKDY Not Found")
}

/* ---------------------------------------- */

func iddsc(code string, Dsch []DSCH) (int, error) {
	N := len(Dsch)

	for j := 0; j < N; j++ {
		_Dsch := &Dsch[j]
		if code == _Dsch.name {
			return j, nil
		}
	}

	return -1, errors.New("DSCH Not Found")
}

/* ---------------------------------------- */

func iddsw(code string, Dscw []DSCW) (int, error) {
	N := len(Dscw)

	for j := 0; j < N; j++ {
		_Dscw := &Dscw[j]
		if code == _Dscw.name {
			return j, nil
		}
	}

	return -1, errors.New("DSCW Not Found")
}

/* ---------------------------------------- */

/* ---------------------------------------- */

//
// スケジュールcodeを Sch から検索し、インデックス番号を返す
// ただし、検索しても見つからない場合は -1 を返す
func idsch(code string, Sch []SCH, err string) (int, error) {
	N := len(Sch)

	if N != len(Sch) {
		panic("N != len(Sch)")
	}

	for j := 0; j < N; j++ {
		_Sch := &Sch[j]
		if code == _Sch.name {
			return j, nil
		}
	}

	if err != "" {
		Eprint("<idsch>", err)
	}
	return -1, errors.New("Schedule Not Found")
}

/* ---------------------------------------- */

// スケジュールcodeを Scw から検索し、インデックス番号を返す
// ただし、検索しても見つからない場合は -1 を返す
func idscw(code string, Scw []SCH, err string) (int, error) {
	N := len(Scw)

	if N != len(Scw) {
		panic("N != len(Scw)")
	}

	for j := 0; j < N; j++ {
		_Scw := &Scw[j]
		if code == _Scw.name {
			return j, nil
		}
	}

	if err != "" {
		Eprint("<idscw>", err)
	}
	return -1, errors.New("Schedule Not Found")
}

/* ---------------------------------------- */

// 室名 `code` に一致する部屋を 部屋の一覧 `Room` から検索し、その番号を返す
// ただし、検索しても見つからない場合はエラーを返す
func idroom(code string, rooms []ROOM, err string) (int, error) {
	for j := range rooms {
		_Room := &rooms[j]
		if code == _Room.Name {
			return j, nil
		}
	}

	E := fmt.Sprintf("Room=%s %s", code, err)
	Eprint("<idroom>", E)

	return -1, errors.New("Room Not Found")
}

/* ---------------------------------------- */

/* ---------------------------------------- */
