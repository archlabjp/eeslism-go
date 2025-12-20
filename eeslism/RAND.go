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
)

// glibcRand implements glibc's TYPE_3 random number generator.
// This is the default algorithm used by glibc rand() when not explicitly initialized.
// C版のEESLISMはsrand()を呼び出さないため、デフォルトシード(1)から開始し、
// glibc rand()と同じ乱数シーケンスを生成します。
// これにより、モンテカルロ法による形態係数計算の結果がC版と一致します。
// Reference: glibc source code random_r.c
type glibcRand struct {
	state [32]int32 // 32-word state array (31 used for TYPE_3)
	fptr  int       // front pointer
	rptr  int       // rear pointer
}

// RAND_MAX はC言語のRAND_MAXに相当する値（2^31-1）
const RAND_MAX = 2147483647

// グローバル乱数生成器（C版と同じデフォルトシード1で初期化）
var cRand = newGlibcRand()

// newGlibcRand creates a new glibc-compatible RNG with seed 1 (default)
func newGlibcRand() *glibcRand {
	g := &glibcRand{}
	g.seed(1)
	return g
}

// seed initializes the state array using glibc's srandom algorithm
func (g *glibcRand) seed(s uint32) {
	// Initialize state array using LCG formula from glibc
	state := make([]int64, 31)
	state[0] = int64(s)
	for i := 1; i < 31; i++ {
		state[i] = (16807 * state[i-1]) % 2147483647
		if state[i] < 0 {
			state[i] += 2147483647
		}
	}

	// Copy to int32 state array
	for i := 0; i < 31; i++ {
		g.state[i] = int32(state[i])
	}
	g.state[31] = 0 // unused in TYPE_3

	// Initialize pointers for TYPE_3 (degree 31)
	g.fptr = 3
	g.rptr = 0

	// Warm up the generator (glibc does 310 = 10*31 iterations)
	for i := 0; i < 310; i++ {
		g.next()
	}
}

// next generates the next random number using TYPE_3 algorithm.
// TYPE_3: x[n] = x[n-3] + x[n-31]
func (g *glibcRand) next() int32 {
	val := g.state[g.fptr] + g.state[g.rptr]
	g.state[g.fptr] = val

	result := (val >> 1) & 0x7fffffff // Ensure non-negative

	g.fptr++
	if g.fptr >= 31 {
		g.fptr = 0
	}
	g.rptr++
	if g.rptr >= 31 {
		g.rptr = 0
	}

	return result
}

// Float64 は0.0から1.0の間の浮動小数点乱数を返します。
// C言語の ((double)rand() / RAND_MAX) と同等です。
func (g *glibcRand) Float64() float64 {
	return float64(g.next()) / float64(RAND_MAX)
}

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
	// glibc互換乱数生成器を使用（C版と同じ乱数シーケンス）
	*a = 2.0 * math.Pi * cRand.Float64()
	*v = mathAcos(mathSqrt(1.0 - cRand.Float64()))

	//TEST
	// *a = math.Pi
	// *v = math.Acos(math.Sqrt(0.5))
}
