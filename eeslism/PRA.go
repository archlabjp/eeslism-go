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

   ベクトルの向きを判定する
              FILE=PRA.c
              Create Date=1998.8.15

*/

package eeslism

import (
	"fmt"
	"math"
	"os"
)

/*
PRA (Vector Direction Determination)

この関数は、与えられたベクトル（`x`, `y`, `z`）が、
基準となる太陽方位ベクトル（`ls`, `ms`, `ns`）に対して、
同じ方向を向いているか、逆方向を向いているか、
あるいは垂直であるかを判定します。
これは、日影計算において、太陽光線が障害物のどの面に当たるかを特定する際に用いられます。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、太陽光線が障害物表面に当たるかどうかを判断する際に、
  光線と表面の相対的な向きを評価します。
  `U`は、光線が表面を通過する距離を示し、
  正の値であれば光線が表面を通過し、
  負の値であれば光線が表面から遠ざかる方向であることを示唆します。
- **光線追跡の基礎**: この計算は、
  日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
  光線が障害物表面と交差する点を求める際に、
  この方向判定が用いられます。
- **エラーハンドリング**: `math.Abs(ls) > epsilon` などの条件は、
  分母がゼロになることを防ぎ、計算の安定性を確保するためのものです。
  もし、太陽方位ベクトルがゼロベクトルに近い場合、
  エラーメッセージを出力し、プログラムを終了します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func PRA(U *float64, ls, ms, ns, x, y, z float64) {
	epsilon := 1.0e-6

	if math.Abs(ls) > epsilon {
		*U = x / ls
	} else if math.Abs(ms) > epsilon {
		*U = y / ms
	} else if math.Abs(ns) > epsilon {
		*U = z / ns
	} else {
		fmt.Printf("ls=%f ms=%f ns=%f\n", ls, ms, ns)
		fmt.Println("errorPRA")
		os.Exit(1)
	}
}
