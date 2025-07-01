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

import "math"

/*
YOGEN (Vector Direction and Dot Product Calculation)

この関数は、2つの点（`Qx, Qy, Qz`と`Px, Py, Pz`）を結ぶベクトルと、
与えられた法線ベクトル`e`の内積を計算し、
その結果からベクトルの向きを判定します。
これは、日影計算において、太陽光線が障害物のどの面に当たるかを特定する際に用いられます。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、太陽光線が障害物表面に当たるかどうかを判断する際に、
  光線と表面の相対的な向きを評価します。
  `S`は、光線が表面を通過する距離を示し、
  正の値であれば光線が表面を通過し、
  負の値であれば光線が表面から遠ざかる方向であることを示唆します。
- **光線追跡の基礎**: この計算は、
  日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
  光線が障害物表面と交差する点を求める際に、
  この方向判定が用いられます。
- **エラーハンドリング**: `PQ == 0.0` の条件は、
  2つの点が同じ位置にある場合（ベクトルがゼロベクトルになる）に、
  計算が不安定になることを防ぐためのものです。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func YOGEN(Qx, Qy, Qz, Px, Py, Pz float64, S *float64, e XYZ) {
	PQx := Px - Qx
	PQy := Py - Qy
	PQz := Pz - Qz

	CAT(&PQx, &PQy, &PQz) // //20170422 higuchi add

	PQ := math.Sqrt(PQx*PQx + PQy*PQy + PQz*PQz)
	E := math.Sqrt(e.X*e.X + e.Y*e.Y + e.Z*e.Z)

	// ↓条件文にした。　20170422 higuchi add
	if PQ == 0.0 {
		*S = -777
	} else {
		*S = (PQx*e.X + PQy*e.Y + PQz*e.Z) / (PQ * E)
	}
}
