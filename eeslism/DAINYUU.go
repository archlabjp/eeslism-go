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

  構造体から別の構造体へ代入
  FILE=DAINYUU.c
  Create Date=1999.6.7
*/

package eeslism

import (
	"fmt"
	"os"
)

/*
DAINYUU_MP (Assignment for Main Plane)

この関数は、日射量計算や日影計算で用いられる「受光面（Opening Plane, OP）」のデータ構造を、
「主面（Main Plane, MP）」のデータ構造に変換し、コピーします。
これにより、異なる計算モジュール間でデータの受け渡しを効率的に行い、
建物の熱的挙動を統合的にモデル化できます。

建築環境工学的な観点:
- **データ構造の変換**: 建物のシミュレーションでは、
  様々な目的のために異なるデータ構造が用いられます。
  この関数は、`OP_COORDNT`関数などで生成された受光面データ（`op`）を、
  日影計算や日射量計算のメインルーチンで扱いやすい主面データ（`mp`）に変換します。
- **窓の分離**: `_op.wd`（窓の数）が`0`より大きい場合、
  受光面内に含まれる窓（`opw`）を個別の主面として扱います。
  これにより、窓からの日射熱取得や日影の影響を、
  壁面とは別に詳細に評価できます。
  `_mpw.wlflg = 1` は、この主面が窓であることを示します。
- **反射率の考慮**: `_mpw.ref`（反射率）や`_mpw.refg`（前面地面の反射率）をコピーすることで、
  日射の反射による熱取得や、周囲の環境からの反射光の影響をモデル化できます。
- **幾何学的情報の保持**: 頂点座標（`_mpw.P`）、方位角（`_mpw.wa`）、傾斜角（`_mpw.wb`）、
  法線ベクトル（`_mpw.e`）などの幾何学的情報を保持することで、
  日射の入射角や日影の形状を正確に計算できます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ変換機能を提供します。
*/
func DAINYUU_MP(op []*P_MENN) []*P_MENN {
	mp := make([]*P_MENN, 0)

	for _, _op := range op {
		// op -> mp
		_mp := new(P_MENN)
		_mp.P = make([]XYZ, _op.polyd)
		*_mp = *_op
		_mp.wd = 0
		_mp.sbflg = 0
		_mp.wlflg = 0
		_mp.opname = _op.opname
		mp = append(mp, _mp)

		for j := 0; j < _op.wd; j++ {
			// opw -> mp
			_mpw := new(P_MENN)
			_mpw.wd = 0
			_mpw.sbflg = 0 // 0=その他
			_mpw.wlflg = 1 // 1=窓

			// 反射率、前面地面の反射率
			_mpw.refg = _op.refg
			_mpw.ref = _op.opw[j].ref

			// 色
			_mpw.rgb[0] = _op.opw[j].rgb[0]
			_mpw.rgb[1] = _op.opw[j].rgb[1]
			_mpw.rgb[2] = _op.opw[j].rgb[2]

			// 頂点
			_mpw.polyd = len(_op.opw[j].P)
			_mpw.P = make([]XYZ, len(_op.opw[j].P))
			for l := 0; l < _mpw.polyd; l++ {
				_mpw.P[l] = _op.opw[j].P[l]
			}

			// 前面地面の代表点までの距離
			_mpw.grpx = _op.opw[j].grpx

			// 方位角、傾斜角
			_mpw.wb = _op.wb
			_mpw.wa = _op.wa

			// 法線ベクトル
			_mpw.e = _op.e

			// 名前
			_mpw.opname = _op.opw[j].opwname

			mp = append(mp, _mpw)
		}

	}

	return mp
}

/*-------------------------------------------------------------------------*/

// p: 代入先
// O: 代入元
// E: 代入元
// ls: 代入元
// ms: 代入元
// ns: 代入元
/*
DAINYUU_GP (Assignment for Geometric Point)

この関数は、3次元空間における点`p`の座標を、
ある平面（法線ベクトル`E`）と、その平面上の点`O`、
そして光線ベクトル（`ls`, `ms`, `ns`）の関係に基づいて計算します。
これは、日影計算において、太陽光線が障害物のどの位置に当たるかを特定する際に用いられます。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、太陽光線（`ls`, `ms`, `ns`）が障害物表面（法線ベクトル`E`、平面上の点`O`）に
  どのように当たるかを幾何学的に計算し、
  その交点（`p`）の座標を求めます。
- **光線追跡の基礎**: この計算は、
  日影の形状や範囲を正確に特定するための光線追跡（Ray Tracing）の基礎となります。
  これにより、建物の窓面や壁面への日射入射量を正確に予測し、
  日射遮蔽効果を評価できます。
- **エラーハンドリング**: `if u != 0.0` の条件は、
  光線ベクトルと法線ベクトルが垂直でないことを確認しています。
  もし垂直であれば、光線が平面に平行であるため交点が存在しないか、
  計算が不安定になる可能性があります。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な幾何学的計算機能を提供します。
*/
func DAINYUU_GP(p *XYZ, O XYZ, E XYZ, ls float64, ms float64, ns float64) {
	var t, u float64

	// u : 法線ベクトルと光線ベクトルの内積
	u = E.X*ls + E.Y*ms + E.Z*ns

	if u != 0.0 {
		t = (E.X*O.X + E.Y*O.Y + E.Z*O.Z) / u
		p.X = ls * t
		p.Y = ms * t
		p.Z = ns * t
	} else {
		fmt.Println("error DAINYUU_GP")
		os.Exit(1)
	}
}

/*-----------------------------------------------------------------------*/
/*
DAINYUU_SMO2 (Assignment for Shaded Area from Previous Time Step)

この関数は、日影計算において、
現在の時刻の影面積を、前時刻の影面積データから取得または更新します。
これは、日影の動的な変化を効率的にモデル化するために用いられます。

建築環境工学的な観点:
- **日影の動的変化**: 日影の形状や大きさは、
  太陽の動き（時刻、日付）によって刻々と変化します。
  シミュレーションでは、この動的な変化を追跡し、
  各時刻における正確な影面積を把握する必要があります。
- **効率的なデータ更新**: `dcnt == 1` の場合（現在の時刻の影を計算する場合）、
  `mp[k].sum = _op.sum` のように、
  現在の計算結果を`mp`（主面データ）と`Sdstr`（日影データ）に格納します。
  `dcnt != 1` の場合（過去の影で近似する場合）、
  `mp[k].sum = Sdstr[k].sdsum[tm]` のように、
  過去の計算結果を再利用します。
  これにより、計算負荷を軽減し、シミュレーションの効率化を図ります。
- **日影率の計算**: `mp[k].sum`は、
  各主面（窓など）の影面積を示しており、
  これを総面積で割ることで日影率を計算できます。
  日影率は、日射熱取得の抑制効果を定量的に評価するために用いられます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ更新機能を提供します。
*/
func DAINYUU_SMO2(opn int, mpn int, op []*P_MENN, mp []*P_MENN, Sdstr []*SHADSTR, dcnt int, tm int) {
	k := 0

	if dcnt == 1 {
		// 過去の影で近似しない時間
		for _, _op := range op {
			mp[k].sum = _op.sum
			Sdstr[k].sdsum[tm] = mp[k].sum
			for j := 0; j < _op.wd; j++ {
				k++
				mp[k].sum = _op.opw[j].sumw
				Sdstr[k].sdsum[tm] = mp[k].sum
			}
			k++
		}
	} else {
		//過去の影で近似する時間
		for k := range mp {
			mp[k].sum = Sdstr[k].sdsum[tm]
		}
	}
}
