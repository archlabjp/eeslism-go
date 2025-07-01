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

	 壁面の中心点の座標を求める
					FILE=GDATA.c
					Create Date=1999.10.26

*/

package eeslism

/*
GDATA (Geometric Data for Center of Gravity)

この関数は、与えられた多角形（`OP`）の重心座標を計算します。
これは、建物の幾何学的モデルにおいて、
各面の代表点や、熱伝達計算における基準点を特定する際に用いられます。

建築環境工学的な観点:
- **幾何学的モデルの簡略化**: 建物の壁面や窓面は、
  熱計算において一つの代表点として扱われることがあります。
  この関数は、多角形の頂点座標から重心を計算することで、
  その多角形を代表する点を特定します。
- **熱伝達計算の基準点**: 熱伝達計算では、
  熱流の方向や大きさを決定するために、
  熱が伝達する面の基準点が必要となります。
  重心は、その基準点として用いられることがあります。
- **日射量計算の補助**: 日射量計算において、
  多角形表面への日射入射角を計算する際に、
  その面の代表点（重心）の座標が用いられることがあります。

この関数は、建物の幾何学的モデルを構築し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための基礎的な役割を果たします。
*/
func GDATA(OP *P_MENN) XYZ {
	var x, y, z float64

	for i := range OP.P {
		x += OP.P[i].X
		y += OP.P[i].Y
		z += OP.P[i].Z
	}

	// the center of gravity of the polygon.
	d := float64(len(OP.P))
	return XYZ{
		X: x / d,
		Y: y / d,
		Z: z / d,
	}
}
