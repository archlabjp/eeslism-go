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

/*
Roomelm (Room Element Assignment)

この関数は、室間換気、放射パネル、およびシステム入力要素（空調機からの供給空気など）を、
各室の熱収支計算モデルに割り当てます。
これにより、建物全体の熱・空気の流れを統合的にモデル化できます。

建築環境工学的な観点:
- **室間相互換気の割り当て**: `room.Nachr`は、
  室間相互換気（隣接する室間での空気の移動）の数を表します。
  `elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]` のように、
  隣室の出口空気温度（`cmp.Elouts[0]`）と絶対湿度（`cmp.Elouts[1]`）を、
  現在の室の熱収支計算の入力として割り当てます。
  これにより、隣室の熱的状態が現在の室に与える影響をモデル化できます。
- **透過熱伝達の割り当て**: `room.Ntr`は、
  透過熱伝達（壁、床、天井などを介した熱の移動）の数を表します。
  `elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]` のように、
  隣室の出口空気温度を、現在の室の熱収支計算の入力として割り当てます。
  これにより、隣室の熱的状態が現在の室に与える影響をモデル化できます。
- **放射パネルの割り当て**: `room.Nrp`は、
  放射パネル（床暖房など）の数を表します。
  `elin.Upo = rmpnl.pnl.cmp.Elins[0]` のように、
  放射パネルの入口温度を、現在の室の熱収支計算の入力として割り当てます。
  これにより、放射パネルからの熱供給が室の熱負荷に与える影響をモデル化できます。
- **システム入力要素の割り当て**: `_Rdpnl`は、
  放射パネルや太陽電池パネルなどの熱的挙動をモデル化するためのデータ構造です。
  この関数は、これらのパネルの熱的挙動を、
  各室の熱収支計算モデルに組み込むための接続を行います。

この関数は、建物全体の熱・空気の流れを統合的にモデル化し、
各室の熱負荷、室内空気質、およびエネルギー消費量を評価するための重要な役割を果たします。
*/
func Roomelm(Room []*ROOM, _Rdpnl []*RDPNL) {

	for n := range Room {

		room := Room[n]
		compnt := room.cmp
		var elin_idx = 0

		for i := 0; i < room.Nachr; i++ {
			var elin *ELIN = compnt.Elins[elin_idx]
			var elinx *ELIN = compnt.Elins[elin_idx+compnt.Elouts[0].Ni]
			var achr *ACHIR = room.achr[i]

			cmp := achr.room.cmp

			elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]
			elinx.Upo, elinx.Upv = cmp.Elouts[1], cmp.Elouts[1]

			elin_idx++
		}

		for i := 0; i < room.Ntr; i++ {

			var elin *ELIN = compnt.Elins[elin_idx]
			trnx := room.trnx[i]

			cmp := trnx.nextroom.cmp
			elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]

			elin_idx++
		}

		for i := 0; i < room.Nrp; i++ {
			var elin *ELIN = compnt.Elins[elin_idx]
			rmpnl := room.rmpnl[i]

			elip := rmpnl.pnl.cmp.Elins[0]

			elin.Upo = elip.Upo

			elin_idx++
		}
	}

	for n := range _Rdpnl {
		Rdpnl := _Rdpnl[n]
		elin_idx := 1

		for m := 0; m < Rdpnl.MC; m++ {
			room := Rdpnl.rm[m]
			elin := Rdpnl.cmp.Elins[elin_idx]
			elin.Upo, elin.Upv = room.cmp.Elouts[0], room.cmp.Elouts[0]

			for i := 0; i < Rdpnl.Ntrm[m]; i++ {
				trnx := room.trnx[i]
				elin := Rdpnl.cmp.Elins[elin_idx]

				cmp := trnx.nextroom.cmp
				elin.Upo, elin.Upv = cmp.Elouts[0], cmp.Elouts[0]

				elin_idx++
			}

		}
	}
}

/*
Roomvar (Room System Equation Creation)

この関数は、各室の熱収支方程式を構築します。
これは、室温や絶対湿度を計算するための連立方程式の係数を設定するもので、
室内の熱的挙動を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
- **熱収支方程式の構築**: 室内の熱収支は、
  壁体からの熱伝達、内部発熱、換気、日射熱取得、
  そして空調システムからの熱供給など、様々な要素によって決まります。
  この関数は、これらの要素を線形方程式の係数として表現し、
  室温（`elout = compnt.Elouts[0]`）と絶対湿度（`elout = compnt.Elouts[1]`）に関する方程式を構築します。
- **係数の設定**: 
  - `elout.Coeffo`: 室の熱容量や、室内空気と表面間の熱伝達など、
    室温や絶対湿度の変化に影響を与える主要な係数です。
  - `elout.Co`: 内部発熱、日射熱取得、外気からの熱供給など、
    室温や絶対湿度を上昇させる定数項です。
  - `elout.Coeffin`: 隣室からの熱伝達、放射パネルからの熱供給、
    空調機からの供給空気など、他の要素からの影響を表す係数です。
- **室間相互換気**: `Room.Nachr`は、室間相互換気の数を表し、
  隣室の温度や湿度変化が現在の室に与える影響をモデル化します。
- **透過熱伝達**: `Room.Ntr`は、透過熱伝達の数を表し、
  隣室の温度変化が現在の室に与える影響をモデル化します。
- **放射パネル**: `Room.Nrp`は、放射パネルの数を表し、
  放射パネルからの熱供給が室温に与える影響をモデル化します。
- **空調機からの供給空気**: `Room.Nasup`は、空調機からの供給空気の数を表し、
  供給空気の温度や流量が室温に与える影響をモデル化します。
- **放射パネルのシステム方程式**: `_Rdpnl`は、
  放射パネルや太陽電池パネルなどの熱的挙動をモデル化するためのデータ構造です。
  この関数は、これらのパネルの熱収支方程式も構築し、
  建物全体の熱・空気の流れを統合的にモデル化します。

この関数は、室の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Roomvar(_Room []*ROOM, _Rdpnl []*RDPNL) {
	for i := range _Room {
		Room := _Room[i]

		compnt := Room.cmp
		elout := compnt.Elouts[0]

		elout.Coeffo = Room.RMt
		elout.Co = Room.RMC

		// 室間相互換気量
		for j := 0; j < Room.Nachr; j++ {
			Gvr := Ca * Room.achr[j].Gvr
			elout.Coeffin[j] = -Gvr
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
			elout.Coeffin[j+Room.Nachr+Room.Ntr+Room.Nrp] = -G
			elout.Coeffo += G
		}

		elout = compnt.Elouts[1]
		elout.Coeffo = Room.RMx
		elout.Co = Room.RMXC

		// 室間相互換気量
		for j := 0; j < Room.Nachr; j++ {
			Gvr := Room.achr[j].Gvr
			elout.Coeffin[j] = -Gvr
			elout.Coeffo += Gvr
		}

		// 流量
		for j := 0; j < Room.Nasup; j++ {
			G := compnt.Elins[j+Room.Nachr+Room.Ntr+Room.Nrp+Room.Nachr].Lpath.G
			elout.Coeffin[j+Room.Nachr] = -G
			elout.Coeffo += G
		}
	}

	for i := range _Rdpnl {
		Rdpnl := _Rdpnl[i]

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
