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

/*  bl_roomvar.c  */

package eeslism

/* ------------------------------------- */

/* 室間換気、放射パネルとシステム入力要素への割り付け */

func Roomelm(Room []ROOM, _Rdpnl []RDPNL) {

	for n := range Room {

		room := &Room[n]
		compnt := room.cmp
		var elin_idx = 0

		for i := 0; i < room.Nachr; i++ {
			var elin *ELIN = compnt.Elins[elin_idx]
			var elinx *ELIN = compnt.Elins[elin_idx+compnt.Elouts[0].Ni]
			var achr *ACHIR = &room.achr[i]

			cmp := achr.room.cmp

			elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]
			elinx.Upo, elinx.Upv = cmp.Elouts[1], cmp.Elouts[1]

			elin_idx++
		}

		for i := 0; i < room.Ntr; i++ {

			var elin *ELIN = compnt.Elins[elin_idx]
			trnx := &room.trnx[i]

			cmp := trnx.nextroom.cmp
			elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]

			elin_idx++
		}

		for i := 0; i < room.Nrp; i++ {
			var elin *ELIN = compnt.Elins[elin_idx]
			rmpnl := &room.rmpnl[i]

			elip := rmpnl.pnl.cmp.Elins[0]

			elin.Upo = elip.Upo

			elin_idx++
		}
	}

	for n := range _Rdpnl {
		Rdpnl := &_Rdpnl[n]
		elin_idx := 1

		for m := 0; m < Rdpnl.MC; m++ {
			room := Rdpnl.rm[m]
			elin := Rdpnl.cmp.Elins[elin_idx]
			elin.Upo, elin.Upv = room.cmp.Elouts[0], room.cmp.Elouts[0]

			for i := 0; i < Rdpnl.Ntrm[m]; i++ {
				trnx := &room.trnx[i]
				elin := Rdpnl.cmp.Elins[elin_idx]

				cmp := trnx.nextroom.cmp
				elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]

				elin_idx++
			}

		}
	}
}

/* -------------------------------  */

/* 室、放射パネルのシステム方程式作成 */

func Roomvar(_Room []ROOM, _Rdpnl []RDPNL) {
	for i := range _Room {
		Room := &_Room[i]

		compnt := Room.cmp
		elout := compnt.Elouts[0]

		elout.Coeffo = Room.RMt
		elout.Co = Room.RMC

		// 室間相互換気量
		for j := 0; j < Room.Nachr; j++ {
			Gvr := Ca * Room.achr[j].Gvr
			elout.Coeffin[j] -= Gvr
			elout.Coeffo += Gvr
		}

		// ARN
		for j := 0; j < Room.Ntr; j++ {
			elout.Coeffin[j+Room.Nachr] = -Room.ARN[j]
		}

		// RMP
		for j := 0; j < Room.Nrp; j++ {
			elout.Coeffin[j+Room.Nachr+Room.Ntr] = -Room.RMP[j]
		}

		// 流量
		for j := 0; j < Room.Nasup; j++ {
			G := Ca * compnt.Elins[j+Room.Nachr+Room.Ntr+Room.Nrp].Lpath.G
			elout.Coeffin[j+Room.Nachr+Room.Ntr+Room.Nrp] -= G
			elout.Coeffo += G
		}

		elout = compnt.Elouts[1]
		elout.Coeffo = Room.RMx
		elout.Co = Room.RMXC

		// 室間相互換気量
		for j := 0; j < Room.Nachr; j++ {
			Gvr := Room.achr[j].Gvr
			elout.Coeffin[j] -= Gvr
			elout.Coeffo += Gvr
		}

		// 流量
		for j := 0; j < Room.Nasup; j++ {
			G := compnt.Elins[j+Room.Nachr+Room.Ntr+Room.Nrp+Room.Nachr].Lpath.G
			elout.Coeffin[j+Room.Nachr] -= G
			elout.Coeffo += G
		}
	}

	for i := range _Rdpnl {
		Rdpnl := &_Rdpnl[i]

		compnt := Rdpnl.cmp
		G := compnt.Elouts[0].Lpath.G
		cG := Spcheat(compnt.Elouts[0].Fluid) * G
		compnt.Elouts[0].Coeffo = cG
		compnt.Elouts[0].Co = Rdpnl.EPC

		cfin := &compnt.Elouts[0].Coeffin[0]
		if Rdpnl.sd[0].mw.wall.WallType == WallType_P {
			// 通常の床暖房パネル
			*cfin = Rdpnl.Epw - cG
		} else {
			// 屋根一体型空気集熱器
			*cfin = -Rdpnl.Epw
		}

		off := 1
		for m := 0; m < Rdpnl.MC; m++ {
			cfin := &compnt.Elouts[0].Coeffin[off]
			off++

			*cfin = -Rdpnl.EPt[m]
			for j := 0; j < Rdpnl.Ntrm[m]; j++ {
				compnt.Elouts[0].Coeffin[off] = -Rdpnl.EPR[m][j]
				off++
			}

			for j := 0; j < Rdpnl.Nrp[m]; j++ {
				compnt.Elouts[0].Coeffin[off] = -Rdpnl.EPW[m][j]
				off++
			}
		}

		/* 空気系統湿度計算用ダミー */
		elout := compnt.Elouts[1]
		elout.Coeffo = G
		elout.Co = 0.0
		elout.Coeffin[0] = -G
	}
}
