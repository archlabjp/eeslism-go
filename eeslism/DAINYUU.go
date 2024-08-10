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

/*

  構造体から別の構造体へ代入
  FILE=DAINYUU.c
  Create Date=1999.6.7
*/

package eeslism

import (
	"fmt"
	"os"
)

/*---------------------------------------------------------------------*/
func DAINYUU_MP(mp *[]*P_MENN, op []*P_MENN, opn int, mpn int) {
	k := 0
	for i := 0; i < opn; i++ {
		(*mp)[k].P = make([]XYZ, op[i].polyd)
		*(*mp)[k] = *op[i]
		(*mp)[k].wd = 0
		(*mp)[k].sbflg = 0
		(*mp)[k].wlflg = 0
		(*mp)[k].opname = op[i].opname

		for j := 0; j < op[i].wd; j++ {
			k++
			(*mp)[k].wd = 0
			(*mp)[k].sbflg = 0
			(*mp)[k].wlflg = 1
			(*mp)[k].refg = op[i].refg
			(*mp)[k].ref = op[i].opw[j].ref
			(*mp)[k].rgb[0] = op[i].opw[j].rgb[0]
			(*mp)[k].rgb[1] = op[i].opw[j].rgb[1]
			(*mp)[k].rgb[2] = op[i].opw[j].rgb[2]
			(*mp)[k].polyd = op[i].opw[j].polyd
			(*mp)[k].P = make([]XYZ, op[i].opw[j].polyd)
			(*mp)[k].grpx = op[i].opw[j].grpx
			(*mp)[k].wb = op[i].wb
			(*mp)[k].wa = op[i].wa
			(*mp)[k].e = op[i].e
			(*mp)[k].opname = op[i].opw[j].opwname

			for l := 0; l < (*mp)[k].polyd; l++ {
				(*mp)[k].P[l] = op[i].opw[j].P[l]
			}
		}
		k++
	}
}

/*-------------------------------------------------------------------------*/
func DAINYUU_GP(p *XYZ, O XYZ, E XYZ, ls float64, ms float64, ns float64) {
	var t, u float64

	u = E.X*ls + E.Y*ms + E.Z*ns

	if u != 0.0 {
		t = (E.X*O.X + E.Y*O.Y + E.Z*O.Z) / u
		p.X = ls * t
		p.Y = ms * t
		p.Z = ns * t
	} else {
		fmt.Println("error DAINYUU_GP")
		os.Exit(1)
	}
}

/*-----------------------------------------------------------------------*/
func DAINYUU_SMO(opn int, mpn int, op []P_MENN, mp []P_MENN) {
	k := 0
	for i := 0; i < opn; i++ {
		mp[k].sum = op[i].sum
		for j := 0; j < op[i].wd; j++ {
			k++
			mp[k].sum = op[i].opw[j].sumw
		}
		k++
	}
}

/*-----------------------------------------------------------------------*/
func DAINYUU_SMO2(opn int, mpn int, op []*P_MENN, mp []*P_MENN, Sdstr []*SHADSTR, dcnt int, tm int) {
	k := 0

	if dcnt == 1 {
		// 過去の影で近似しない時間
		for i := 0; i < opn; i++ {
			mp[k].sum = op[i].sum
			Sdstr[k].sdsum[tm] = mp[k].sum
			for j := 0; j < op[i].wd; j++ {
				k++
				mp[k].sum = op[i].opw[j].sumw
				Sdstr[k].sdsum[tm] = mp[k].sum
			}
			k++
		}
	} else {
		//過去の影で近似する時間
		for k = 0; k < mpn; k++ {
			mp[k].sum = Sdstr[k].sdsum[tm]
		}
	}
}
