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
func DAINYUU_MP(op []*P_MENN) []*P_MENN {
	mp := make([]*P_MENN, 0)

	for _, _op := range op {
		// op -> mp
		_mp := new(P_MENN)
		_mp.P = make([]XYZ, _op.polyd)
		*_mp = *_op
		_mp.wd = 0
		_mp.sbflg = 0
		_mp.wlflg = 0
		_mp.opname = _op.opname
		mp = append(mp, _mp)

		for j := 0; j < _op.wd; j++ {
			// opw -> mp
			_mpw := new(P_MENN)
			_mpw.wd = 0
			_mpw.sbflg = 0 // 0=その他
			_mpw.wlflg = 1 // 1=窓

			// 反射率、前面地面の反射率
			_mpw.refg = _op.refg
			_mpw.ref = _op.opw[j].ref

			// 色
			_mpw.rgb[0] = _op.opw[j].rgb[0]
			_mpw.rgb[1] = _op.opw[j].rgb[1]
			_mpw.rgb[2] = _op.opw[j].rgb[2]

			// 頂点
			_mpw.polyd = len(_op.opw[j].P)
			_mpw.P = make([]XYZ, len(_op.opw[j].P))
			for l := 0; l < _mpw.polyd; l++ {
				_mpw.P[l] = _op.opw[j].P[l]
			}

			// 前面地面の代表点までの距離
			_mpw.grpx = _op.opw[j].grpx

			// 方位角、傾斜角
			_mpw.wb = _op.wb
			_mpw.wa = _op.wa

			// 法線ベクトル
			_mpw.e = _op.e

			// 名前
			_mpw.opname = _op.opw[j].opwname

			mp = append(mp, _mpw)
		}

	}

	return mp
}

/*-------------------------------------------------------------------------*/

// p: 代入先
// O: 代入元
// E: 代入元
// ls: 代入元
// ms: 代入元
// ns: 代入元
func DAINYUU_GP(p *XYZ, O XYZ, E XYZ, ls float64, ms float64, ns float64) {
	var t, u float64

	// u : 法線ベクトルと光線ベクトルの内積
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
func DAINYUU_SMO2(opn int, mpn int, op []*P_MENN, mp []*P_MENN, Sdstr []*SHADSTR, dcnt int, tm int) {
	k := 0

	if dcnt == 1 {
		// 過去の影で近似しない時間
		for _, _op := range op {
			mp[k].sum = _op.sum
			Sdstr[k].sdsum[tm] = mp[k].sum
			for j := 0; j < _op.wd; j++ {
				k++
				mp[k].sum = _op.opw[j].sumw
				Sdstr[k].sdsum[tm] = mp[k].sum
			}
			k++
		}
	} else {
		//過去の影で近似する時間
		for k := range mp {
			mp[k].sum = Sdstr[k].sdsum[tm]
		}
	}
}
