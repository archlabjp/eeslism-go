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

/*   rmschd.c      */

package eeslism

import (
	"fmt"
	"math"
)

/* --------------------------------------------- */

/*  窓の熱抵抗の変更 */

func Windowschdlr(isw []ControlSWType, windows []*WINDOW, N int, ds []*RMSRF) {
	for i := 0; i < N; i++ {
		sd := ds[i]
		if sd.ble == BLE_Window {
			nsw := sd.Nfn

			// デフォルトの窓
			j := 0

			// 窓の選択は条件による切り替えが優先される
			// スケジュールによる窓の変更
			if nsw > 1 {
				j = iswmode(rune(isw[sd.fnsw]), nsw, sd.fnmrk[:])
			}

			w := windows[sd.fnd[j]]

			// 動的な窓の変更
			if sd.ifwin != nil {
				// 条件の確認
				if contrlif(sd.Ctlif) != 0 {
					w = sd.ifwin
				}
			}

			sd.Eo = w.Eo
			sd.Ei = w.Ei
			sd.tgtn = w.tgtn
			sd.Bn = w.Bn
			sd.Rwall = w.Rwall
			sd.CAPwall = 0.0

			sd.fn = j
		}
	}
}

/* --------------------------------------------- */

/*  室内発熱の計算    */

func Qischdlr(_Room []*ROOM) {
	Ht := [9]float64{92, 106, 119, 131, 145, 198, 226, 264, 383}
	Hs24 := [9]float64{58, 62, 63, 64, 69, 76, 83, 99, 137}
	d := [9]float64{3.5, 3.6, 4.0, 4.2, 4.4, 6.5, 7.0, 7.3, 6.3}

	for i := range _Room {
		Room := _Room[i]

		Room.Hc = 0.0
		Room.Hr = 0.0
		Room.HL = 0.0
		Room.Lc = 0.0
		Room.Lr = 0.0
		Room.Ac = 0.0
		Room.Ar = 0.0
		Room.AL = 0.0

		if Room.Hmsch != nil && *Room.Hmsch > 0.0 {
			N := Room.Nhm * *Room.Hmsch

			if N > 0 && Room.Hmwksch != nil && *Room.Hmwksch > 0.0 {
				wk := int(*Room.Hmwksch - 1)

				if wk < 0 || wk > 8 {
					s := fmt.Sprintf("Room=%s wk=%d", Room.Name, wk)
					Eprint("<Qischdlr>", s)
				}

				Eo := Room.cmp.Elouts
				Tr := Room.Tr // 自然室温計算時は前時刻の室温で顕熱・潜熱分離する
				// 室温が高温となるときに、顕熱が負になるのを回避
				if Eo[0].Control == LOAD_SW {
					Tr = Room.rmld.Tset
				}

				Q := math.Max((Hs24[wk]-d[wk]*(Tr-24.0)), 0.) * N
				Room.Hc = Q * 0.5
				Room.Hr = Q * 0.5
				Room.HL = Ht[wk]*N - Q
			}
		}

		if Room.Lightsch != nil && *Room.Lightsch > 0.0 {
			Q := Room.Light * *Room.Lightsch
			Room.Lc = Q * 0.5
			Room.Lr = Q * 0.5
		} else {
			Room.Lc = 0.0
			Room.Lr = 0.0
		}

		if Room.Assch != nil && *Room.Assch > 0.0 {
			Room.Ac = Room.Apsc * *Room.Assch
			Room.Ar = Room.Apsr * *Room.Assch
		} else {
			Room.Ac = 0.0
			Room.Ar = 0.0
		}

		if Room.Alsch != nil && *Room.Alsch > 0.0 {
			Room.AL = Room.Apl * *Room.Alsch
		} else {
			Room.AL = 0.0
		}
	}
}

/* -------------------------------------------------------- */

/*  換気量の設定     */

func Vtschdlr(rooms []*ROOM) {
	for i := range rooms {
		Gvi := 0.0
		Gve := 0.0

		if rooms[i].Visc != nil && *rooms[i].Visc > 0.0 {
			Gvi = rooms[i].Gvi * *rooms[i].Visc
		} else {
			Gvi = 0.0
		}

		if rooms[i].Vesc != nil && *rooms[i].Vesc > 0.0 {
			Gve = rooms[i].Gve * *rooms[i].Vesc
		} else {
			Gve = 0.0
		}

		rooms[i].Gvent = Gvi + Gve
	}
}
