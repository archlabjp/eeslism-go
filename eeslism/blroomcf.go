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

/*   bl_roomcf.c  */

package eeslism

/*
eeroomcf (Room Coefficient Calculation for Energy Simulation)

この関数は、建物のエネルギーシミュレーションにおいて、
室内の熱収支計算に必要な各種熱的パラメータを計算します。
これには、表面熱伝達率、熱貫流率、透過日射量、相当外気温度などが含まれます。

建築環境工学的な観点:
- **熱伝達率の計算**: 建物の熱負荷は、
  室内外の温度差や日射などによって生じる熱の移動によって決まります。
  この関数は、以下の熱伝達率を計算します。
  - `Rmhtrcf`: 室内表面における対流熱伝達率と放射熱伝達率を計算します。
    これらの熱伝達率は、室内空気と表面、および表面間の放射熱交換を介した熱移動をモデル化するために重要です。
  - `Rmhtrsmcf`: 壁体や窓などの熱貫流率を計算します。
    熱貫流率は、建物の断熱性能を示す指標であり、
    熱損失・熱取得量を評価する上で不可欠です。
- **透過日射と相当外気温度の計算 (Rmexct)**:
  - 透過日射量: 窓などを透過して室内に侵入する日射熱量を計算します。
    これは、夏季の冷房負荷や冬季の暖房負荷軽減に大きく影響します。
  - 相当外気温度: 日射や夜間放射などの影響を、
    あたかも外気温度が変化したかのように見なして、
    熱伝達計算を簡略化するために導入される仮想的な温度です。
    これにより、外皮からの熱損失・熱取得をより簡単に計算できるようになります。
- **室の係数と定数項の計算 (Roomcf)**:
  `Roomcf`関数は、壁体内部の熱伝導、家具の熱容量、内部発熱、換気など、
  室内の熱収支を構成する様々な要素を考慮した係数と定数項を計算します。
  これらの係数は、室温や熱負荷の予測に用いられます。
- **シミュレーションの統合**: この関数は、
  建物の熱的挙動をモデル化するための様々な計算モジュールを統合し、
  一連の熱収支計算を実行します。
  これにより、建物全体のエネルギー性能を総合的に評価できます。

この関数は、建物の熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要な役割を果たします。
*/
func eeroomcf(Wd *WDAT, Exs *EXSFS, Rmvls *RMVLS, nday int, mt int) {
	// 熱伝達率の計算

	// 表面熱伝達率（対流・放射））の計算
	Rmhtrcf(Exs, Rmvls.Emrk, Rmvls.Room, Rmvls.Sd, Wd)

	if DEBUG {
		// 表面熱伝達率の表示
		xpralph(Rmvls.Room, Rmvls.Sd)
	}

	// 熱貫流率の計算
	Rmhtrsmcf(Rmvls.Sd)

	// 透過日射、相当外気温度の計算
	Rmexct(Rmvls.Room, Rmvls.Sd, Wd, Exs.Exs, Rmvls.Snbk, Rmvls.Qrm, nday, mt)

	// 室の係数（壁体熱伝導等））、定数項の計算
	Roomcf(Rmvls.Mw, Rmvls.Room, Rmvls.Rdpnl, Wd, Exs)

	xprroom(Rmvls.Room)
	xprxas(Rmvls.Room, Rmvls.Sd)
}
