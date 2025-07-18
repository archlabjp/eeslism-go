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

// NOTE: すべてのポリゴンの頂点数が同じであることを前提としているようだ
//       POLYGON命令で指定した場合は、任意の頂点数になるので正しく動作しないように思う。

/*
URA (Under-Surface Relationship Analysis)

この関数は、日影計算において、
ある面（`OP`）から見た他の面（`LP`）の相対的な位置関係を計算します。
具体的には、`OP`の各頂点から`LP`の各頂点へのベクトルを`OP`の法線ベクトルに投影した値を計算します。

建築環境工学的な観点:
  - **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
    この関数は、ある面が別の面によって影になるかどうかを判断する際に、
    両面の相対的な位置関係を評価します。
    `t[j].ps[i][k]`に格納される値は、
    `OP`の頂点から`LP`の頂点へのベクトルが`OP`の法線ベクトルに沿ってどれだけ離れているかを示します。
    この値が正であれば`LP`が`OP`の「表側」にあることを示唆し、
    負であれば「裏側」にあることを示唆します。
  - **光線追跡の基礎**: この計算は、
    日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
    光線が障害物表面と交差する点を求める際に、
    この相対位置関係が用いられます。
  - **ポリゴン頂点数の前提**: コメントに記載されているように、
    この関数は全てのポリゴンの頂点数が同じであることを前提としています。
    これは、モデルの柔軟性を制限する可能性があります。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func URA(u, w int, LP []*P_MENN, t []*bekt, OP []*P_MENN) {
	for j := 0; j < u; j++ {
		for i := 0; i < w; i++ {
			// NOTE: 一旦、polydを4固定にする
			//for k := 0; k < LP[i].polyd; k++ {
			for k := 0; k < 4; k++ {

				t[j].ps[i][k] = -(OP[j].e.X*OP[j].P[k].X + OP[j].e.Y*OP[j].P[k].Y +
					OP[j].e.Z*OP[j].P[k].Z - OP[j].e.X*LP[i].P[k].X -
					OP[j].e.Y*LP[i].P[k].Y - OP[j].e.Z*LP[i].P[k].Z) /
					((OP[j].e.X)*(OP[j].e.X) + (OP[j].e.Y)*(OP[j].e.Y) +
						(OP[j].e.Z)*(OP[j].e.Z))
			}
		}
	}
}

/*-------------------------------------------------------------*/

/*
URA_M (Under-Surface Relationship Analysis for Main Plane)

この関数は、太陽光線ベクトル（`ls`, `ms`, `ns`）が、
傾斜角`wb`を持つ面に対して、
同じ方向を向いているか、逆方向を向いているか、
あるいは垂直であるかを判定します。
これは、日影計算において、太陽光線が障害物のどの面に当たるかを特定する際に用いられます。

建築環境工学的な観点:
  - **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
    この関数は、太陽光線が障害物表面に当たるかどうかを判断する際に、
    光線と表面の相対的な向きを評価します。
    `s`は、光線が表面を通過する距離を示し、
    正の値であれば光線が表面を通過し、
    負の値であれば光線が表面から遠ざかる方向であることを示唆します。
  - **光線追跡の基礎**: この計算は、
    日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
    光線が障害物表面と交差する点を求める際に、
    この方向判定が用いられます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func URA_M(ls, ms, ns float64, wb float64) float64 {
	ex := 0.0
	ey := -math.Sin((-wb) * math.Pi / 180)
	ez := math.Cos((-wb) * math.Pi / 180)

	s := (ex*ls + ey*ms + ez*ns) / (ex*ex + ey*ey + ez*ez)
	return s
}
