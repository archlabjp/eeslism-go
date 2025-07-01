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

///*
//
//                         法線ベクトルを求める
//                                        FILE=HOUSEN.c
//                                        Create Date=1998.10.26
//                                        Update 2007.10.11 higuchi
//
//*/

/*
HOUSEN (Normal Vector Calculation for Polygons)

この関数は、与えられた多角形（`LP`）の法線ベクトルを計算し、正規化します。
法線ベクトルは、その面が太陽光に対してどの方向を向いているかを示し、
日射入射角の計算や、日影の有無の判定に用いられます。

建築環境工学的な観点:
- **日射量計算の基礎**: 建物の壁面や窓面への日射入射量は、
  その面の向き（法線ベクトル）と太陽の方向との相対関係によって決まります。
  この関数は、多角形の頂点座標から法線ベクトルを計算することで、
  日射入射角を正確に計算し、日射熱取得量を予測するための基礎を提供します。
- **日影計算の基礎**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  法線ベクトルは、太陽光線が障害物のどの面に当たるかを判断し、
  日影の有無を判定するために用いられます。
- **幾何学的モデルの正確性**: 法線ベクトルを正確に計算し、正規化することで、
  建物の幾何学的モデルの正確性を確保します。
  これにより、シミュレーション結果の信頼性が向上します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func HOUSEN(LP []*P_MENN) {
	for _, _lp := range LP {
		// 多角形のうち2辺のベクトルを求める
		x := _lp.P[1].X - _lp.P[0].X
		y := _lp.P[1].Y - _lp.P[0].Y
		z := _lp.P[1].Z - _lp.P[0].Z
		x1 := _lp.P[2].X - _lp.P[0].X
		y1 := _lp.P[2].Y - _lp.P[0].Y
		z1 := _lp.P[2].Z - _lp.P[0].Z

		// 法線ベクトルを求める
		_lp.e.X = y*z1 - z*y1
		_lp.e.Y = z*x1 - x*z1
		_lp.e.Z = x*y1 - y*x1

		// 法線ベクトルの正規化
		el := math.Sqrt(_lp.e.X*_lp.e.X + _lp.e.Y*_lp.e.Y + _lp.e.Z*_lp.e.Z)
		_lp.e.X = _lp.e.X / el
		_lp.e.Y = _lp.e.Y / el
		_lp.e.Z = _lp.e.Z / el
	}
}

/*
HOUSEN2 (Normal Vector Calculation for Three Points)

この関数は、3つの点（`p0`, `p1`, `p2`）で定義される平面の法線ベクトルを計算し、正規化します。
これは、多角形が3つの頂点を持つ場合や、
特定の平面の向きを定義する際に用いられます。

建築環境工学的な観点:
- **日射量計算の基礎**: 建物の壁面や窓面への日射入射量は、
  その面の向き（法線ベクトル）と太陽の方向との相対関係によって決まります。
  この関数は、3つの点から法線ベクトルを計算することで、
  日射入射角を正確に計算し、日射熱取得量を予測するための基礎を提供します。
- **日影計算の基礎**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  法線ベクトルは、太陽光線が障害物のどの面に当たるかを判断し、
  日影の有無を判定するために用いられます。
- **幾何学的モデルの正確性**: 法線ベクトルを正確に計算し、正規化することで、
  建物の幾何学的モデルの正確性を確保します。
  これにより、シミュレーション結果の信頼性が向上します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func HOUSEN2(p0, p1, p2, e *XYZ) {
	x := p1.X - p0.X
	y := p1.Y - p0.Y
	z := p1.Z - p0.Z
	x1 := p2.X - p0.X
	y1 := p2.Y - p0.Y
	z1 := p2.Z - p0.Z

	e.X = y*z1 - z*y1
	e.Y = z*x1 - x*z1
	e.Z = x*y1 - y*x1

	el := math.Sqrt(e.X*e.X + e.Y*e.Y + e.Z*e.Z)
	e.X = e.X / el
	e.Y = e.Y / el
	e.Z = e.Z / el
}
