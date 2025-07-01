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

   乱数の発生
      FILE=RAND.c
      Create Date=1999.6.7

*/

package eeslism

import (
	"math"
	"math/rand"
)

/*
RAND (Random Number Generation for Monte Carlo Simulation)

この関数は、モンテカルロ法による日影計算や形態係数計算において、
ランダムな方向を持つ太陽光線（または放射線）を生成するための乱数を生成します。
aは方位角、vは仰角です。

建築環境工学的な観点:
- **モンテカルロ法の基礎**: モンテカルロ法は、
  乱数を用いて多数の試行を繰り返すことで、
  複雑な物理現象を統計的にシミュレーションする手法です。
  この関数は、そのモンテカルロ法における個々の光線の方向を決定するための乱数を生成します。
- **ランダムな方向の生成**: 
  - `*a = 2.0 * math.Pi * rand.Float64()`: 方位角（`a`）を0から2π（360度）の範囲でランダムに生成します。
  - `*v = math.Acos(math.Sqrt(1.0 - rand.Float64()))`: 仰角（`v`）をランダムに生成します。
    この式は、半球状の空間に均一に光線を分布させるためのものです。
- **日影計算と形態係数計算への応用**: 生成されたランダムな方向を持つ光線は、
  建物や周囲の障害物との交差判定に用いられ、
  日影の形状や面積、
  および形態係数を統計的に推定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func RAND(a, v *float64) {
	// a is azimuth, v is elevation
	*a = 2.0 * math.Pi * rand.Float64()
	*v = math.Acos(math.Sqrt(1.0 - rand.Float64()))

	//TEST
	// *a = math.Pi
	// *v = math.Acos(math.Sqrt(0.5))
}
