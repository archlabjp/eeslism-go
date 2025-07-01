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

/*
直線と平面の交点を求める
FILE=KOUTEN.c
Create Date=1998.10.26
*/

/*
KOUTEN (Intersection Point Calculation for Line and Plane)

この関数は、3次元空間における直線と平面の交点座標を計算します。
これは、日影計算において、太陽光線が障害物のどの面に当たるかを特定する際に用いられます。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、太陽光線（直線）が障害物表面（平面）と交差する点を計算します。
  - `Qx, Qy, Qz`: 直線上の点の座標（通常は太陽光線の出発点）。
  - `ls, ms, ns`: 直線の方向ベクトル（太陽光線の方向）。
  - `lp`: 平面上の点の座標。
  - `E`: 平面の法線ベクトル。
  - `Px, Py, Pz`: 計算される交点の座標。
- **光線追跡の基礎**: この計算は、
  日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
  太陽光線が障害物表面と交差する点を求めることで、
  その光線が障害物によって遮られるかどうかを判断できます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func KOUTEN(Qx, Qy, Qz, ls, ms, ns float64, Px, Py, Pz *float64, lp, E XYZ) {
	t := (E.X*lp.X + E.Y*lp.Y + E.Z*lp.Z - E.X*Qx - E.Y*Qy - E.Z*Qz) / (E.X*ls + E.Y*ms + E.Z*ns)
	*Px = t*ls + Qx
	*Py = t*ms + Qy
	*Pz = t*ns + Qz
}
