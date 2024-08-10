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

// KOUTEN は直線と平面の交点を計算します。
// Qx, Qy, Qz: 直線の座標
// ls, ms, ns: 直線の方向
// Px, Py, Pz: 交点の座標 [出力]
// lp: 交点の座標
// E: 平面の座標
func KOUTEN(Qx, Qy, Qz, ls, ms, ns float64, Px, Py, Pz *float64, lp, E XYZ) {
	t := (E.X*lp.X + E.Y*lp.Y + E.Z*lp.Z - E.X*Qx - E.Y*Qy - E.Z*Qz) / (E.X*ls + E.Y*ms + E.Z*ns)
	*Px = t*ls + Qx
	*Py = t*ms + Qy
	*Pz = t*ns + Qz
}
