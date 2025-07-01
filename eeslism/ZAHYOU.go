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

   モンテカルロ法を用いる際の座標変換
              FILE=ZAHYOU.c
              Create Date=1999.6.7

*/

package eeslism

import "math"

/*
ZAHYOU (Coordinate Transformation for Monte Carlo Method)

この関数は、モンテカルロ法による日影計算や形態係数計算において、
座標変換を行います。
具体的には、元の座標（`Op`）を、
重心（`G`）を原点とし、
方位角（`wa`）と傾斜角（`wb`）で回転させた新しい座標系に変換します。

建築環境工学的な観点:
- **幾何学的モデルの簡略化**: 複雑な形状の建物や周囲の障害物による日影を正確に計算するために、
  座標変換は不可欠です。
  この関数は、各面をその重心を原点とする局所座標系に変換することで、
  計算を簡略化し、効率化を図ります。
- **日影計算の基礎**: 変換された座標は、
  太陽光線と障害物表面の交差判定や、
  影の形状と面積の計算に用いられます。
- **回転変換**: 
  - `wa`: Z軸周りの回転（方位角）。
  - `wb`: X軸周りの回転（傾斜角）。
  これらの回転変換を適用することで、
  様々な向きを持つ面を統一的に扱うことができます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func ZAHYOU(Op, G XYZ, op *XYZ, wa, wb float64) {
	Cwa := math.Cos(wa * math.Pi / 180)
	Swa := math.Sin(wa * math.Pi / 180)
	Cwb := math.Cos((-wb) * math.Pi / 180)
	Swb := math.Sin((-wb) * math.Pi / 180)

	p := XYZ{
		X: Op.X - G.X,
		Y: Op.Y - G.Y,
		Z: Op.Z - G.Z,
	}

	q := XYZ{
		X: p.X*Cwa - p.Y*Swa,
		Y: p.X*Swa + p.Y*Cwa,
		Z: p.Z,
	}

	op.X = q.X
	op.Y = q.Y*Cwb - q.Z*Swb
	op.Z = q.Y*Swb + q.Z*Cwb

	CAT(&op.X, &op.Y, &op.Z)
}

/*------------------------------------------------------------------*/

/*
R_ZAHYOU (Reverse Coordinate Transformation)

この関数は、`ZAHYOU`関数によって変換された座標を、
元の座標系に戻す逆変換を行います。

建築環境工学的な観点:
- **座標変換の逆操作**: シミュレーションの途中で座標変換を行った後、
  最終的な結果を元の座標系で表現する必要がある場合に用いられます。
  例えば、日影計算で局所座標系で影の形状を計算した後、
  それを建物全体の座標系に戻して可視化する際に利用されます。
- **幾何学的モデルの整合性**: 座標変換とその逆変換を正確に行うことで、
  シミュレーションモデルの幾何学的整合性を確保します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func R_ZAHYOU(Op, G XYZ, op *XYZ, wa, wb float64) {
	Cwa := math.Cos((-wa) * math.Pi / 180)
	Swa := math.Sin((-wa) * math.Pi / 180)
	Cwb := math.Cos(wb * math.Pi / 180)
	Swb := math.Sin(wb * math.Pi / 180)

	p := XYZ{
		X: Op.X,
		Y: Op.Y*Cwb - Op.Z*Swb,
		Z: Op.Y*Swb + Op.Z*Cwb,
	}

	q := XYZ{
		X: p.X*Cwa - p.Y*Swa,
		Y: p.X*Swa + p.Y*Cwa,
		Z: p.Z,
	}

	op.X = q.X + G.X
	op.Y = q.Y + G.Y
	op.Z = q.Z + G.Z

	CAT(&op.X, &op.Y, &op.Z)
}
