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

/*      schdlr.c     */

package main

/*  １日スケジュ－ルから設定値の選択   */

func schval(nday, ttmm int, Sch *SCH, Dsch []DSCH) float64 {
	sc := Sch.day[nday]

	if sc < 0 {
		return FNOP
	}

	Ds := &Dsch[sc]
	N := Ds.N

	for k := 0; k < N; k++ {
		stime := Ds.stime[k]
		etime := Ds.etime[k]
		val := Ds.val[k]
		if stime <= ttmm && ttmm <= etime {
			return val
		}
	}
	return FNOP
}

/* ----------------------------------------------------- */

/*  １日スケジュ－ルから設定モ－ドの選択   */

func scwmode(nday, ttmm int, Scw *SCH, Dscw []DSCW) rune {
	sw := Scw.day[nday]
	Ds := &Dscw[sw]
	N := Ds.N

	for k := 0; k < N; k++ {
		stime := Ds.stime[k]
		etime := Ds.etime[k]
		mode := Ds.mode[k]
		if stime <= ttmm && ttmm <= etime {
			return mode
		}
	}
	return 'x'
}

/* ----------------------------------------------------- */

/*  スケジュ－ルモ－ドから設定番号の検索   */

func iswmode(c rune, N int, mode []rune) int {
	if N == 1 {
		return 0
	} else {
		for i := 0; i < N; i++ {
			if c == mode[i] {
				return i
			}
		}
		return -1
	}
}

/* ----------------------------------------------------- */
