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
MATINIT (Matrix Initialization for P_MENN)

この関数は、`P_MENN`構造体の配列（`q`）を初期化します。
`P_MENN`構造体は、日影計算や日射量計算で用いられる「受光面（Opening Plane, OP）」や
「被受照面（Light-Receiving Plane, LP）」の幾何学的情報を格納するために用いられます。

建築環境工学的な観点:
- **幾何学的モデルの準備**: 建物の日射環境をシミュレーションする前に、
  窓、壁、日よけ、障害物などの幾何学的情報を格納するデータ構造を準備します。
  この関数は、各`P_MENN`構造体の法線ベクトル（`e`）と頂点座標（`P`）をゼロで初期化します。
- **計算の正確性**: 適切な初期化は、
  その後の日影計算や日射量計算の正確性を確保するために重要です。
  未初期化のデータを使用すると、
  計算結果に予期せぬ誤差が生じる可能性があります。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func MATINIT(q []*P_MENN, n int) {
	for i := 0; i < n; i++ {
		q[i].e.X = 0.0
		q[i].e.Y = 0.0
		q[i].e.Z = 0.0
		for j := 0; j < q[i].polyd; j++ {
			q[i].P[j].X = 0.0
			q[i].P[j].Y = 0.0
			q[i].P[j].Z = 0.0
		}
	}
}

/*
MATINIT_sum (Summation Initialization for P_MENN)

この関数は、`P_MENN`構造体の配列（`op`）内の影面積の合計（`sum`）と、
窓の影面積の合計（`opw[i].sumw`）をゼロに初期化します。

建築環境工学的な観点:
- **影面積の集計の準備**: 日影計算では、
  各面や窓の影面積を合計して、
  日影率や日射熱取得量を計算します。
  この関数は、新しい計算期間の開始前に、
  これらの合計値をゼロにリセットすることで、
  正確な集計を可能にします。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func MATINIT_sum(opn int, op []*P_MENN) {
	for j := 0; j < opn; j++ {
		op[j].sum = 0.0
		for i := 0; i < op[j].wd; i++ {
			op[j].opw[i].sumw = 0.0
		}
	}
}

/*
MATINIT_sdstr (Shaded Data Structure Initialization)

この関数は、日影データ構造（`SHADSTR`）内の影面積の合計（`sdsum`）をゼロに初期化します。

建築環境工学的な観点:
- **日影データの集計の準備**: 日影計算では、
  各時刻における影面積を合計して、
  日影率や日射熱取得量を計算します。
  この関数は、新しい計算期間の開始前に、
  これらの合計値をゼロにリセットすることで、
  正確な集計を可能にします。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func MATINIT_sdstr(mpn, mtb int, Sdstr []*SHADSTR) {
	for j := 0; j < mpn; j++ {
		for i := 0; i < mtb; i++ {
			Sdstr[j].sdsum[i] = 0.0
		}
	}
}
