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

	入力データの計算用構造体への変換
	FILE=COORDNT.c
	Create Date=1999.6.7
*/

package eeslism

import (
	"fmt"
	"math"
	"os"
)

/*
LP_COORDNT (Light-Receiving Plane Coordinate Transformation)

この関数は、建物周辺の障害物（庇、バルコニー、袖壁、樹木など）や、
多角形データで定義された障害物の幾何学的情報を、
日影計算や日射量計算に用いられる「被受照面（Light-Receiving Plane）」の座標データに変換します。

建築環境工学的な観点:
- **日影計算の基礎**: 建物の窓面や壁面への日射入射量は、
  周囲の建物や地形、植栽、そして建物自体に付随する日よけ（庇、バルコニーなど）によって形成される日影によって大きく影響されます。
  この関数は、これらの障害物の形状と位置を正確にモデル化し、
  日影計算の基礎となる座標データ（`lp`）を生成します。
- **座標変換**: 障害物の相対的な位置（`x`, `y`, `z`）や、
  方位角（`Wa`）と傾斜角（`Wb`）を考慮して、
  各障害物表面の頂点座標（`lp_k.P`）を計算します。
  これにより、3次元空間における障害物の正確な形状を表現できます。
- **法線ベクトルの算出**: 各被受照面（障害物表面）の法線ベクトル（`lp_k.e`）を算出します。
  法線ベクトルは、その面が太陽光に対してどの方向を向いているかを示し、
  日射入射角の計算や、日影の有無の判定に用いられます。
- **日よけの種類に応じたモデル化**: `sblk.sbfname`によって、
  庇（`HISASI`）、バルコニー（`BARUKONI`）、袖壁（`SODEKABE`）、窓日よけ（`MADOHIYOKE`）など、
  様々な種類の日よけを識別し、それぞれの形状に応じた座標変換ロジックを適用します。
- **日影率の計算**: `lp_k.shad[p]`は、
  各被受照面が年間を通じてどれだけの期間日影になるかを示す配列です。
  これにより、日射遮蔽効果を定量的に評価できます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func LP_COORDNT(
	poly []*POLYGN,
	tree []*TREE,
	obs []*OBS,
	BDP []*BBDP,
) []*P_MENN {
	var name string
	var j, p int
	var sWa, cWa, sWb, cWb, a, b, c, chWA, shWA float64
	var csWA, ssWA, csWb, ssWb float64
	var caWb, saWa, caWa, saWb float64
	var cbWa, sbWa float64

	lp := make([]*P_MENN, 0)

	const M_rad = math.Pi / 180.0
	for i := range BDP {
		sWb = math.Sin(BDP[i].Wb * M_rad)
		cWb = math.Cos(BDP[i].Wb * M_rad)
		sWa = math.Sin(-BDP[i].Wa * M_rad)
		cWa = math.Cos(-BDP[i].Wa * M_rad)

		if len(BDP[i].SBLK) != 0 {
			for _, sblk := range BDP[i].SBLK {

				if sblk.sbfname == "HISASI" {
					a = 0.0
					b = 0.0
					c = 0.0
					chWA = 0.0
					shWA = 0.0

					a = BDP[i].x0 + sblk.x*cWa - sblk.y*cWb*sWa
					b = BDP[i].y0 + sblk.x*sWa + sblk.y*cWb*cWa
					c = BDP[i].z0 + sblk.y*sWb

					chWA = math.Cos((sblk.WA - BDP[i].Wb) * M_rad)
					shWA = math.Sin((sblk.WA - BDP[i].Wb) * M_rad)

					lp_k := NewP_MENN()
					lp_k.opname = sblk.snbname
					lp_k.wa = BDP[i].Wa
					lp_k.wb = BDP[i].Wb + (180.0 - sblk.WA)
					if lp_k.wb > 180.0 {
						lp_k.wb = 360.0 - lp_k.wb
					}

					lp_k.wd = 0
					lp_k.ref = 0.0 /*--付設障害物からの反射率を０としている-*/
					for p = 0; p < 366; p++ {
						lp_k.shad[p] = 1
					}

					lp_k.sbflg = 1

					lp_k.rgb[0] = sblk.rgb[0]
					lp_k.rgb[1] = sblk.rgb[1]
					lp_k.rgb[2] = sblk.rgb[2]
					lp_k.polyd = 4
					lp_k.P = make([]XYZ, 4)

					lp_k.P[0].X = a
					lp_k.P[0].Y = b
					lp_k.P[0].Z = c

					lp_k.P[1].X = a + sblk.D*chWA*sWa
					lp_k.P[1].Y = b - sblk.D*chWA*cWa
					lp_k.P[1].Z = c + sblk.D*shWA

					lp_k.P[2].X = lp_k.P[1].X + sblk.W*cWa
					lp_k.P[2].Y = lp_k.P[1].Y + sblk.W*sWa
					lp_k.P[2].Z = lp_k.P[1].Z

					lp_k.P[3].X = a + sblk.W*cWa
					lp_k.P[3].Y = b + sblk.W*sWa
					lp_k.P[3].Z = c

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
					lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
					lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
					CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)

					//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e)

					lp = append(lp, lp_k)
				} else if sblk.sbfname == "BARUKONI" {

					a = 0.0
					b = 0.0
					c = 0.0
					chWA = 0.0
					shWA = 0.0

					a = BDP[i].x0 + sblk.x*cWa - sblk.y*cWb*sWa
					b = BDP[i].y0 + sblk.x*sWa + sblk.y*cWb*cWa
					c = BDP[i].z0 + sblk.y*sWb

					chWA = math.Cos((90.0 - BDP[i].Wb) * M_rad)
					shWA = math.Sin((90.0 - BDP[i].Wb) * M_rad)

					lp_1 := NewP_MENN()
					lp_2 := NewP_MENN()
					lp_3 := NewP_MENN()
					lp_4 := NewP_MENN()
					lp_5 := NewP_MENN()

					lp_1.opname = sblk.snbname

					lp_1.wa = BDP[i].Wa
					lp_1.wb = BDP[i].Wb + 90.0
					if lp_1.wb > 180.0 {
						lp_1.wb = -(360.0 - lp_1.wb)
					}

					lp_1.ref = 0.0
					lp_1.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_1.shad[p] = 1
					}

					lp_1.rgb[0] = sblk.rgb[0]
					lp_1.rgb[1] = sblk.rgb[1]
					lp_1.rgb[2] = sblk.rgb[2]
					lp_1.polyd = 4
					lp_1.P = make([]XYZ, 4)

					lp_1.wd = 0

					lp_1.P[0].X = a
					lp_1.P[0].Y = b
					lp_1.P[0].Z = c

					lp_1.P[1].X = a + sblk.D*chWA*sWa
					lp_1.P[1].Y = b - sblk.D*chWA*cWa
					lp_1.P[1].Z = c + sblk.D*shWA

					lp_1.P[2].X = lp_1.P[1].X + sblk.W*cWa
					lp_1.P[2].Y = lp_1.P[1].Y + sblk.W*sWa
					lp_1.P[2].Z = lp_1.P[1].Z

					lp_1.P[3].X = a + sblk.W*cWa
					lp_1.P[3].Y = b + sblk.W*sWa
					lp_1.P[3].Z = c

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_1.e.Z = math.Cos(lp_1.wb * M_rad)
					lp_1.e.Y = -math.Sin(lp_1.wb*M_rad) * math.Cos(lp_1.wa*M_rad)
					lp_1.e.X = -math.Sin(lp_1.wb*M_rad) * math.Sin(lp_1.wa*M_rad)
					CAT(&lp_1.e.X, &lp_1.e.Y, &lp_1.e.Z)

					//HOUSEN2(&lp_1.P[0],&lp_1.P[1],&lp_1.P[2],&lp_1.e)

					name = "2" + lp_1.opname
					lp_2.opname = name

					lp_2.wa = BDP[i].Wa - 90.0
					if lp_2.wa <= -180.0 {
						lp_2.wa = 360.0 + lp_2.wa
					}
					lp_2.wb = 90.0
					lp_2.ref = 0.0
					lp_2.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_2.shad[p] = 1
					}

					lp_2.rgb[0] = sblk.rgb[0]
					lp_2.rgb[1] = sblk.rgb[1]
					lp_2.rgb[2] = sblk.rgb[2]
					lp_2.polyd = 4
					lp_2.P = make([]XYZ, 4)

					lp_2.wd = 0

					lp_2.P[0] = lp_1.P[0]
					lp_2.P[1] = lp_1.P[1]
					lp_2.P[2].X = lp_2.P[1].X + sblk.H*cWb*sWa
					lp_2.P[2].Y = lp_2.P[1].Y - sblk.H*cWb*cWa
					lp_2.P[2].Z = lp_2.P[1].Z - sblk.H*sWb
					lp_2.P[3].X = a + sblk.H*cWb*sWa
					lp_2.P[3].Y = b - sblk.H*cWb*cWa
					lp_2.P[3].Z = c - sblk.H*sWb

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_2.e.Z = math.Cos(lp_2.wb * M_rad)
					lp_2.e.Y = -math.Sin(lp_2.wb*M_rad) * math.Cos(lp_2.wa*M_rad)
					lp_2.e.X = -math.Sin(lp_2.wb*M_rad) * math.Sin(lp_2.wa*M_rad)
					CAT(&lp_2.e.X, &lp_2.e.Y, &lp_2.e.Z)

					//HOUSEN2(&lp_2.P[0],&lp_2.P[1],&lp_2.P[2],&lp_2.e)

					name = "3" + lp_1.opname
					lp_3.opname = name

					lp_3.wa = BDP[i].Wa
					lp_3.wb = BDP[i].Wb - 90.0
					lp_3.ref = sblk.ref
					lp_3.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_3.shad[p] = 1
					}

					lp_3.wd = 0

					lp_3.rgb[0] = sblk.rgb[0]
					lp_3.rgb[1] = sblk.rgb[1]
					lp_3.rgb[2] = sblk.rgb[2]
					lp_3.polyd = 4
					lp_3.P = make([]XYZ, 4)

					lp_3.P[0] = lp_2.P[3]
					lp_3.P[1] = lp_2.P[2]
					lp_3.P[2].X = lp_3.P[1].X + sblk.W*cWa
					lp_3.P[2].Y = lp_3.P[1].Y + sblk.W*sWa
					lp_3.P[2].Z = lp_3.P[1].Z
					lp_3.P[3].X = lp_3.P[0].X + sblk.W*cWa
					lp_3.P[3].Y = lp_3.P[0].Y + sblk.W*sWa
					lp_3.P[3].Z = lp_3.P[0].Z

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_3.e.Z = math.Cos(lp_3.wb * M_rad)
					lp_3.e.Y = -math.Sin(lp_3.wb*M_rad) * math.Cos(lp_3.wa*M_rad)
					lp_3.e.X = -math.Sin(lp_3.wb*M_rad) * math.Sin(lp_3.wa*M_rad)
					CAT(&lp_3.e.X, &lp_3.e.Y, &lp_3.e.Z)

					//HOUSEN2(&lp_3.P[0],&lp_3.P[1],&lp_3.P[2],&lp_3.e) ;

					name = "4" + lp_1.opname
					lp_4.opname = name

					lp_4.wa = BDP[i].Wa + 90.0
					if lp_4.wa >= 180.0 {
						lp_4.wa = lp_4.wa - 360.0
					}
					lp_4.wb = 90.0
					lp_4.ref = 0.0
					lp_4.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_4.shad[p] = 1
					}

					lp_4.rgb[0] = sblk.rgb[0]
					lp_4.rgb[1] = sblk.rgb[1]
					lp_4.rgb[2] = sblk.rgb[2]
					lp_4.polyd = 4
					lp_4.P = make([]XYZ, 4)

					lp_4.wd = 0

					lp_4.P[0] = lp_1.P[3]
					lp_4.P[1] = lp_1.P[2]
					lp_4.P[2] = lp_3.P[2]
					lp_4.P[3] = lp_3.P[3]

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_4.e.Z = math.Cos(lp_4.wb * M_rad)
					lp_4.e.Y = -math.Sin(lp_4.wb*M_rad) * math.Cos(lp_4.wa*M_rad)
					lp_4.e.X = -math.Sin(lp_4.wb*M_rad) * math.Sin(lp_4.wa*M_rad)
					CAT(&lp_4.e.X, &lp_4.e.Y, &lp_4.e.Z)

					//HOUSEN2(&lp_4.P[0],&lp_4.P[1],&lp_4.P[2],&lp_4.e)

					name = "5" + lp_1.opname
					lp_5.opname = name

					lp_5.wa = BDP[i].Wa - 180
					if lp_5.wa > 180 {
						lp_5.wa = lp_5.wa - 360
					} else if lp_5.wa < 180 {
						lp_5.wa = lp_5.wa + 360
					}

					lp_5.wb = BDP[i].Wb
					lp_5.ref = 0.0
					lp_5.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_5.shad[p] = 1
					}

					lp_5.rgb[0] = sblk.rgb[0]
					lp_5.rgb[1] = sblk.rgb[1]
					lp_5.rgb[2] = sblk.rgb[2]
					lp_5.polyd = 4
					lp_5.P = make([]XYZ, 4)

					lp_5.wd = 0

					lp_5.P[0] = lp_2.P[2]
					lp_5.P[1].X = lp_5.P[0].X - sblk.h*sWa*cWb
					lp_5.P[1].Y = lp_5.P[0].Y + sblk.h*cWa*cWb
					lp_5.P[1].Z = lp_5.P[0].Z + sblk.h*sWb
					lp_5.P[2].X = lp_5.P[1].X + sblk.W*cWa
					lp_5.P[2].Y = lp_5.P[1].Y + sblk.W*sWa
					lp_5.P[2].Z = lp_5.P[1].Z
					lp_5.P[3] = lp_4.P[2]

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_5.e.Z = math.Cos(lp_5.wb * M_rad)
					lp_5.e.Y = -math.Sin(lp_5.wb*M_rad) * math.Cos(lp_5.wa*M_rad)
					lp_5.e.X = -math.Sin(lp_5.wb*M_rad) * math.Sin(lp_5.wa*M_rad)
					CAT(&lp_5.e.X, &lp_5.e.Y, &lp_5.e.Z)

					//HOUSEN2(&lp_5.P[0],&lp_5.P[1],&lp_5.P[2],&lp_5.e)

					lp = append(lp, lp_1)
					lp = append(lp, lp_2)
					lp = append(lp, lp_3)
					lp = append(lp, lp_4)
					lp = append(lp, lp_5)
				} else if sblk.sbfname == "SODEKABE" {

					a = 0.0
					b = 0.0
					c = 0.0
					csWA = 0.0
					ssWA = 0.0
					csWb = 0.0
					ssWb = 0.0

					a = BDP[i].x0 + sblk.x*cWa - sblk.y*cWb*sWa
					b = BDP[i].y0 + sblk.x*sWa + sblk.y*cWb*cWa
					c = BDP[i].z0 + sblk.y*sWb
					csWA = math.Cos((-BDP[i].Wa - sblk.WA) * M_rad)
					ssWA = math.Sin((-BDP[i].Wa - sblk.WA) * M_rad)
					csWb = math.Cos((90.0 - BDP[i].Wb) * M_rad)
					ssWb = math.Sin((90.0 - BDP[i].Wb) * M_rad)

					lp_k := NewP_MENN()

					lp_k.opname = sblk.snbname

					lp_k.wa = BDP[i].Wa - 90.0
					if lp_k.wa <= -180.0 {
						lp_k.wa = 360.0 + lp_k.wa
					}
					lp_k.wb = 90.0
					lp_k.ref = 0.0
					lp_k.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_k.shad[p] = 1
					}

					lp_k.rgb[0] = sblk.rgb[0]
					lp_k.rgb[1] = sblk.rgb[1]
					lp_k.rgb[2] = sblk.rgb[2]
					lp_k.polyd = 4
					lp_k.P = make([]XYZ, 4)

					lp_k.wd = 0

					lp_k.P[0].X = a
					lp_k.P[0].Y = b
					lp_k.P[0].Z = c
					lp_k.P[1].X = a + sblk.D*csWb*csWA
					lp_k.P[1].Y = b + sblk.D*csWb*ssWA
					lp_k.P[1].Z = c + sblk.D*ssWb
					lp_k.P[2].X = lp_k.P[1].X + sblk.H*cWb*sWa
					lp_k.P[2].Y = lp_k.P[1].Y - sblk.H*cWb*cWa
					lp_k.P[2].Z = lp_k.P[1].Z - sblk.H*sWb
					lp_k.P[3].X = a + sblk.H*cWb*sWa
					lp_k.P[3].Y = b - sblk.H*cWb*cWa
					lp_k.P[3].Z = c - sblk.H*sWb

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
					lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
					lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
					CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)

					//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e) ;

					lp = append(lp, lp_k)

				} else if sblk.sbfname == "MADOHIYOKE" {

					a = 0.0
					b = 0.0
					c = 0.0
					chWA = 0.0
					shWA = 0.0

					a = BDP[i].x0 + sblk.x*cWa - sblk.y*cWb*sWa
					b = BDP[i].y0 + sblk.x*sWa + sblk.y*cWb*cWa
					c = BDP[i].z0 + sblk.y*sWb
					chWA = math.Cos((90.0 - BDP[i].Wb) * M_rad)
					shWA = math.Sin((90.0 - BDP[i].Wb) * M_rad)

					lp_k := NewP_MENN()

					lp_k.opname = sblk.snbname

					lp_k.wa = BDP[i].Wa - 180
					if lp_k.wa > 180 {
						lp_k.wa = lp_k.wa - 360
					} else if lp_k.wa < 180 {
						lp_k.wa = lp_k.wa + 360
					}

					lp_k.wb = BDP[i].Wb
					lp_k.ref = 0.0
					lp_k.sbflg = 1

					for p = 0; p < 366; p++ {
						lp_k.shad[p] = 1
					}

					lp_k.rgb[0] = sblk.rgb[0]
					lp_k.rgb[1] = sblk.rgb[1]
					lp_k.rgb[2] = sblk.rgb[2]
					lp_k.polyd = 4
					lp_k.P = make([]XYZ, 4)

					lp_k.wd = 0

					lp_k.P[0].X = a + sblk.D*chWA*sWa
					lp_k.P[0].Y = b - sblk.D*chWA*cWa
					lp_k.P[0].Z = c + sblk.H*shWA
					lp_k.P[1].X = lp_k.P[0].X + sblk.H*cWb*sWa
					lp_k.P[1].Y = lp_k.P[0].Y - sblk.H*cWb*cWa
					lp_k.P[1].Z = lp_k.P[0].Z - sblk.H*sWb
					lp_k.P[3].X = lp_k.P[0].X + sblk.W*cWa
					lp_k.P[3].Y = lp_k.P[0].Y + sblk.W*sWa
					lp_k.P[3].Z = lp_k.P[0].Z
					lp_k.P[2].X = lp_k.P[1].X + sblk.W*cWa
					lp_k.P[2].Y = lp_k.P[1].Y + sblk.W*sWa
					lp_k.P[2].Z = lp_k.P[1].Z

					/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
					lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
					lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
					lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
					CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)

					//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e) ;

					lp = append(lp, lp_k)
				}
			}
		}

	}

	/*-----------------------------------------------------*/

	for i := range obs {
		if obs[i].fname == "rect" {
			caWb = 0.0
			saWa = 0.0
			caWa = 0.0
			saWb = 0.0

			caWb = math.Cos(obs[i].Wb * M_rad)
			saWa = math.Sin(-obs[i].Wa * M_rad)
			saWb = math.Sin(obs[i].Wb * M_rad)
			caWa = math.Cos(-obs[i].Wa * M_rad)

			lp_k := NewP_MENN()

			lp_k.opname = obs[i].obsname

			lp_k.wa = obs[i].Wa
			lp_k.wb = obs[i].Wb
			lp_k.ref = obs[i].ref[0]
			for p = 0; p < 366; p++ {
				lp_k.shad[p] = 1
			}

			lp_k.rgb[0] = obs[i].rgb[0]
			lp_k.rgb[1] = obs[i].rgb[1]
			lp_k.rgb[2] = obs[i].rgb[2]
			lp_k.polyd = 4
			lp_k.P = make([]XYZ, 4)

			lp_k.wd = 0
			lp_k.sbflg = 0

			lp_k.P[0].X = obs[i].x
			lp_k.P[0].Y = obs[i].y
			lp_k.P[0].Z = obs[i].z
			lp_k.P[1].X = obs[i].x - obs[i].H*caWb*saWa
			lp_k.P[1].Y = obs[i].y + obs[i].H*caWb*caWa
			lp_k.P[1].Z = obs[i].z + obs[i].H*saWb
			lp_k.P[2].X = obs[i].x + obs[i].W*caWa - obs[i].H*caWb*saWa
			lp_k.P[2].Y = obs[i].y + obs[i].H*caWb*caWa + obs[i].W*saWa
			lp_k.P[2].Z = obs[i].z + obs[i].H*saWb
			lp_k.P[3].X = obs[i].x + obs[i].W*caWa
			lp_k.P[3].Y = obs[i].y + obs[i].W*saWa
			lp_k.P[3].Z = obs[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
			lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
			lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
			CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)
			//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e)

			lp = append(lp, lp_k)
		} else if obs[i].fname == "cube" {
			cbWa = 0.0
			sbWa = 0.0

			cbWa = math.Cos(-obs[i].Wa * M_rad)
			sbWa = math.Sin(-obs[i].Wa * M_rad)

			lp_1 := NewP_MENN()
			lp_2 := NewP_MENN()
			lp_3 := NewP_MENN()
			lp_4 := NewP_MENN()

			lp_1.opname = obs[i].obsname

			lp_1.wa = obs[i].Wa
			lp_1.wb = 90.0
			lp_1.ref = obs[i].ref[0]
			for p = 0; p < 366; p++ {
				lp_1.shad[p] = 1
			}

			lp_1.rgb[0] = obs[i].rgb[0]
			lp_1.rgb[1] = obs[i].rgb[1]
			lp_1.rgb[2] = obs[i].rgb[2]
			lp_1.polyd = 4
			lp_1.P = make([]XYZ, 4)

			lp_1.wd = 0
			lp_1.sbflg = 0

			lp_1.P[0].X = obs[i].x
			lp_1.P[0].Y = obs[i].y
			lp_1.P[0].Z = obs[i].z
			lp_1.P[1].X = obs[i].x
			lp_1.P[1].Y = obs[i].y
			lp_1.P[1].Z = obs[i].z + obs[i].H
			lp_1.P[2].X = obs[i].x + obs[i].W*cbWa
			lp_1.P[2].Y = obs[i].y + obs[i].W*sbWa
			lp_1.P[2].Z = lp_1.P[1].Z
			lp_1.P[3].X = lp_1.P[2].X
			lp_1.P[3].Y = lp_1.P[2].Y
			lp_1.P[3].Z = lp_1.P[0].Z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_1.e.Z = math.Cos(lp_1.wb * M_rad)
			lp_1.e.Y = -math.Sin(lp_1.wb*M_rad) * math.Cos(lp_1.wa*M_rad)
			lp_1.e.X = -math.Sin(lp_1.wb*M_rad) * math.Sin(lp_1.wa*M_rad)
			CAT(&lp_1.e.X, &lp_1.e.Y, &lp_1.e.Z)

			//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e)

			name = "2" + lp_1.opname
			lp_2.opname = name

			lp_2.wa = obs[i].Wa - 90.0
			if lp_2.wa <= -180.0 {
				lp_2.wa = 360.0 + lp_2.wa
			}
			lp_2.wb = 90.0
			lp_2.ref = obs[i].ref[1]
			for p = 0; p < 366; p++ {
				lp_2.shad[p] = 1
			}

			lp_2.rgb[0] = obs[i].rgb[0]
			lp_2.rgb[1] = obs[i].rgb[1]
			lp_2.rgb[2] = obs[i].rgb[2]
			lp_2.polyd = 4
			lp_2.P = make([]XYZ, 4)

			lp_2.wd = 0
			lp_2.sbflg = 0

			lp_2.P[0] = lp_1.P[3]
			lp_2.P[1] = lp_1.P[2]
			lp_2.P[2].X = lp_2.P[1].X - obs[i].D*sbWa
			lp_2.P[2].Y = lp_2.P[1].Y + obs[i].D*cbWa
			lp_2.P[2].Z = lp_1.P[2].Z
			lp_2.P[3].X = lp_2.P[2].X
			lp_2.P[3].Y = lp_2.P[2].Y
			lp_2.P[3].Z = lp_2.P[0].Z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_2.e.Z = math.Cos(lp_2.wb * M_rad)
			lp_2.e.Y = -math.Sin(lp_2.wb*M_rad) * math.Cos(lp_2.wa*M_rad)
			lp_2.e.X = -math.Sin(lp_2.wb*M_rad) * math.Sin(lp_2.wa*M_rad)
			CAT(&lp_2.e.X, &lp_2.e.Y, &lp_2.e.Z)
			//HOUSEN2(&lp_2.P[0],&lp_2.P[1],&lp_2.P[2],&lp_2.e)

			name = "3" + lp_1.opname
			lp_3.opname = name

			lp_3.wa = obs[i].Wa - 180.0
			if lp_3.wa <= -180.0 {
				lp_3.wa = 360.0 + lp_3.wa
			}
			lp_3.wb = 90.0
			lp_3.ref = obs[i].ref[2]
			for p = 0; p < 366; p++ {
				lp_3.shad[p] = 1
			}

			lp_3.rgb[0] = obs[i].rgb[0]
			lp_3.rgb[1] = obs[i].rgb[1]
			lp_3.rgb[2] = obs[i].rgb[2]
			lp_3.polyd = 4
			lp_3.P = make([]XYZ, 4)

			lp_3.wd = 0
			lp_3.sbflg = 0
			lp_3.P[0] = lp_2.P[3]
			lp_3.P[1] = lp_2.P[2]
			lp_3.P[2].X = obs[i].x - obs[i].D*sbWa
			lp_3.P[2].Y = obs[i].y + obs[i].D*cbWa
			lp_3.P[2].Z = lp_3.P[1].Z
			lp_3.P[3].X = lp_3.P[2].X
			lp_3.P[3].Y = lp_3.P[2].Y
			lp_3.P[3].Z = lp_3.P[0].Z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_3.e.Z = math.Cos(lp_3.wb * M_rad)
			lp_3.e.Y = -math.Sin(lp_3.wb*M_rad) * math.Cos(lp_3.wa*M_rad)
			lp_3.e.X = -math.Sin(lp_3.wb*M_rad) * math.Sin(lp_3.wa*M_rad)
			CAT(&lp_3.e.X, &lp_3.e.Y, &lp_3.e.Z)
			//HOUSEN2(&lp_3.P[0],&lp_3.P[1],&lp_3.P[2],&lp_3.e)

			name = "4" + lp_1.opname
			lp_4.opname = name

			lp_4.wa = obs[i].Wa + 90.0
			if lp_4.wa > 180.0 {
				lp_4.wa = 360.0 - lp_4.wa
			}
			lp_4.wb = 90.0
			lp_4.ref = obs[i].ref[3]
			for p = 0; p < 366; p++ {
				lp_4.shad[p] = 1
			}

			lp_4.rgb[0] = obs[i].rgb[0]
			lp_4.rgb[1] = obs[i].rgb[1]
			lp_4.rgb[2] = obs[i].rgb[2]
			lp_4.polyd = 4
			lp_4.P = make([]XYZ, 4)

			lp_4.wd = 0
			lp_4.sbflg = 0
			lp_4.P[0] = lp_3.P[3]
			lp_4.P[1] = lp_3.P[2]
			lp_4.P[2] = lp_1.P[1]
			lp_4.P[3] = lp_1.P[0]

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_4.e.Z = math.Cos(lp_4.wb * M_rad)
			lp_4.e.Y = -math.Sin(lp_4.wb*M_rad) * math.Cos(lp_4.wa*M_rad)
			lp_4.e.X = -math.Sin(lp_4.wb*M_rad) * math.Sin(lp_4.wa*M_rad)
			CAT(&lp_4.e.X, &lp_4.e.Y, &lp_4.e.Z)
			//HOUSEN2(&lp_4.P[0],&lp_4.P[1],&lp_4.P[2],&lp_4.e)

			lp = append(lp, lp_1)
			lp = append(lp, lp_2)
			lp = append(lp, lp_3)
			lp = append(lp, lp_4)
		} else if obs[i].fname == "r_tri" {

			caWb = 0.0
			saWa = 0.0
			caWa = 0.0
			saWb = 0.0

			caWb = math.Cos(obs[i].Wb * M_rad)
			saWa = math.Sin(obs[i].Wa * M_rad)
			saWb = math.Sin(obs[i].Wb * M_rad)
			caWa = math.Cos(obs[i].Wa * M_rad)

			lp_k := NewP_MENN()

			lp_k.opname = obs[i].obsname

			lp_k.wa = obs[i].Wa
			lp_k.wb = obs[i].Wb
			lp_k.ref = obs[i].ref[0]
			for p = 0; p < 366; p++ {
				lp_k.shad[p] = 1
			}

			lp_k.rgb[0] = obs[i].rgb[0]
			lp_k.rgb[1] = obs[i].rgb[1]
			lp_k.rgb[2] = obs[i].rgb[2]
			lp_k.polyd = 3
			lp_k.P = make([]XYZ, 3)

			lp_k.wd = 0
			lp_k.sbflg = 0

			lp_k.P[0].X = obs[i].x
			lp_k.P[0].Y = obs[i].y
			lp_k.P[0].Z = obs[i].z
			lp_k.P[1].X = obs[i].x - obs[i].H*caWb*saWa
			lp_k.P[1].Y = obs[i].y + obs[i].H*caWb*caWa
			lp_k.P[1].Z = obs[i].z + obs[i].H*saWb
			lp_k.P[2].X = obs[i].x + obs[i].W*caWa
			lp_k.P[2].Y = obs[i].y + obs[i].W*saWa
			lp_k.P[2].Z = obs[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
			lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
			lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
			CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)
			//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e)

			lp = append(lp, lp_k)
		} else if obs[i].fname == "i_tri" {

			caWb = 0.0
			saWa = 0.0
			caWa = 0.0
			saWb = 0.0

			caWb = math.Cos(obs[i].Wb * M_rad)
			saWa = math.Sin(-obs[i].Wa * M_rad)
			saWb = math.Sin(obs[i].Wb * M_rad)
			caWa = math.Cos(-obs[i].Wa * M_rad)

			lp_k := NewP_MENN()

			lp_k.opname = obs[i].obsname

			lp_k.wa = obs[i].Wa
			lp_k.wb = obs[i].Wb
			lp_k.ref = obs[i].ref[0]
			for p = 0; p < 366; p++ {
				lp_k.shad[p] = 1
			}

			lp_k.rgb[0] = obs[i].rgb[0]
			lp_k.rgb[1] = obs[i].rgb[1]
			lp_k.rgb[2] = obs[i].rgb[2]
			lp_k.polyd = 3
			lp_k.P = make([]XYZ, 3)

			lp_k.wd = 0
			lp_k.sbflg = 0

			lp_k.P[0].X = obs[i].x
			lp_k.P[0].Y = obs[i].y
			lp_k.P[0].Z = obs[i].z
			lp_k.P[1].X = obs[i].x + ((obs[i].W)/2)*caWa - obs[i].H*caWb*saWa
			lp_k.P[1].Y = obs[i].y + obs[i].H*caWb*caWa + ((obs[i].W)/2)*saWa
			lp_k.P[1].Z = obs[i].z + obs[i].H*saWb
			lp_k.P[2].X = obs[i].x + obs[i].W*caWa
			lp_k.P[2].Y = obs[i].y + obs[i].W*saWa
			lp_k.P[2].Z = obs[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
			lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
			lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
			CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)
			//HOUSEN2(&lp_k.P[0],&lp_k.P[1],&lp_k.P[2],&lp_k.e)

			lp = append(lp, lp_k)
		} else {
			fmt.Printf("error--**COORDNT-lp\n")
			os.Exit(1)
		}

	}

	/*--------------------------------------------------------*/
	for i := range tree {
		if tree[i].treetype == "treeA" {
			/*----1----*/
			lp_1 := NewP_MENN()

			name = tree[i].treename + "-m1"
			lp_1.opname = name

			for p = 0; p < 366; p++ {
				lp_1.shad[p] = 1
			}

			lp_1.rgb[0] = 0.4
			lp_1.rgb[1] = 0.3
			lp_1.rgb[2] = 0.01
			lp_1.polyd = 4
			lp_1.P = make([]XYZ, 4)

			lp_1.wa = 0
			lp_1.wb = 90.0
			lp_1.ref = 0.0
			lp_1.wd = 0
			lp_1.sbflg = 0

			lp_1.P[0].X = tree[i].x - (tree[i].W1 * 0.5)
			lp_1.P[0].Y = tree[i].y - (tree[i].W1 * 0.5)
			lp_1.P[0].Z = tree[i].z
			lp_1.P[1].X = lp_1.P[0].X
			lp_1.P[1].Y = lp_1.P[0].Y
			lp_1.P[1].Z = tree[i].z + tree[i].H1
			lp_1.P[2].X = tree[i].x + (tree[i].W1 * 0.5)
			lp_1.P[2].Y = lp_1.P[0].Y
			lp_1.P[2].Z = tree[i].z + tree[i].H1
			lp_1.P[3].X = lp_1.P[2].X
			lp_1.P[3].Y = lp_1.P[0].Y
			lp_1.P[3].Z = tree[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_1.e.Z = math.Cos(lp_1.wb * M_rad)
			lp_1.e.Y = -math.Sin(lp_1.wb*M_rad) * math.Cos(lp_1.wa*M_rad)
			lp_1.e.X = -math.Sin(lp_1.wb*M_rad) * math.Sin(lp_1.wa*M_rad)
			CAT(&lp_1.e.X, &lp_1.e.Y, &lp_1.e.Z)
			//HOUSEN2(&lp_1.P[0],&lp_1.P[1],&lp_1.P[2],&lp_1.e)

			lp = append(lp, lp_1)

			/*----2----*/
			lp_2 := NewP_MENN()

			name = tree[i].treename + "-m2"
			lp_2.opname = name
			lp_2.wa = -90.0
			lp_2.wb = 90.0
			lp_2.ref = 0.0
			lp_2.wd = 0
			lp_2.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_2.shad[p] = 1
			}

			lp_2.rgb[0] = 0.4
			lp_2.rgb[1] = 0.3
			lp_2.rgb[2] = 0.01
			lp_2.polyd = 4
			lp_2.P = make([]XYZ, 4)

			lp_2.P[0] = lp_1.P[3]
			lp_2.P[1] = lp_1.P[2]
			lp_2.P[2].X = tree[i].x + (tree[i].W1 * 0.5)
			lp_2.P[2].Y = tree[i].y + (tree[i].W1 * 0.5)
			lp_2.P[2].Z = tree[i].z + tree[i].H1
			lp_2.P[3].X = lp_2.P[2].X
			lp_2.P[3].Y = lp_2.P[2].Y
			lp_2.P[3].Z = tree[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_2.e.Z = math.Cos(lp_2.wb * M_rad)
			lp_2.e.Y = -math.Sin(lp_2.wb*M_rad) * math.Cos(lp_2.wa*M_rad)
			lp_2.e.X = -math.Sin(lp_2.wb*M_rad) * math.Sin(lp_2.wa*M_rad)
			CAT(&lp_2.e.X, &lp_2.e.Y, &lp_2.e.Z)
			//HOUSEN2(&lp_2.P[0],&lp_2.P[1],&lp_2.P[2],&lp_2.e)

			lp = append(lp, lp_2)

			/*----3----*/
			lp_3 := NewP_MENN()

			name = tree[i].treename + "-m3"
			lp_3.opname = name

			lp_3.wa = 180.0
			lp_3.wb = 90.0
			lp_3.ref = 0.0
			lp_3.wd = 0
			lp_3.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_3.shad[p] = 1
			}

			lp_3.rgb[0] = 0.4
			lp_3.rgb[1] = 0.3
			lp_3.rgb[2] = 0.01
			lp_3.polyd = 4
			lp_3.P = make([]XYZ, 4)

			lp_3.P[0] = lp_2.P[3]
			lp_3.P[1] = lp_2.P[2]
			lp_3.P[2].X = tree[i].x - (tree[i].W1 * 0.5)
			lp_3.P[2].Y = tree[i].y + (tree[i].W1 * 0.5)
			lp_3.P[2].Z = tree[i].z + tree[i].H1
			lp_3.P[3].X = lp_3.P[2].X
			lp_3.P[3].Y = lp_3.P[2].Y
			lp_3.P[3].Z = tree[i].z

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_3.e.Z = math.Cos(lp_3.wb * M_rad)
			lp_3.e.Y = -math.Sin(lp_3.wb*M_rad) * math.Cos(lp_3.wa*M_rad)
			lp_3.e.X = -math.Sin(lp_3.wb*M_rad) * math.Sin(lp_3.wa*M_rad)
			CAT(&lp_3.e.X, &lp_3.e.Y, &lp_3.e.Z)
			//HOUSEN2(&lp_3.P[0],&lp_3.P[1],&lp_3.P[2],&lp_3.e)

			lp = append(lp, lp_3)

			/*----4----*/
			lp_4 := NewP_MENN()

			name = tree[i].treename + "-m4"
			lp_4.opname = name

			lp_4.wa = 90.0
			lp_4.wb = 90.0
			lp_4.ref = 0.0
			lp_4.wd = 0
			lp_4.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_4.shad[p] = 1
			}

			lp_4.rgb[0] = 0.4
			lp_4.rgb[1] = 0.3
			lp_4.rgb[2] = 0.01
			lp_4.polyd = 4
			lp_4.P = make([]XYZ, 4)

			lp_4.P[0] = lp_3.P[3]
			lp_4.P[1] = lp_3.P[2]
			lp_4.P[2] = lp_1.P[1]
			lp_4.P[3] = lp_1.P[0]

			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_4.e.Z = math.Cos(lp_4.wb * M_rad)
			lp_4.e.Y = -math.Sin(lp_4.wb*M_rad) * math.Cos(lp_4.wa*M_rad)
			lp_4.e.X = -math.Sin(lp_4.wb*M_rad) * math.Sin(lp_4.wa*M_rad)
			CAT(&lp_4.e.X, &lp_4.e.Y, &lp_4.e.Z)
			//HOUSEN2(&lp_4.P[0],&lp_4.P[1],&lp_4.P[2],&lp_4.e)

			lp = append(lp, lp_4)

			/*----5----*/
			lp_5 := NewP_MENN()

			name = tree[i].treename
			lp_5.opname = name

			lp_5.wa = -22.5
			lp_5.ref = 0.0
			lp_5.wd = 0
			lp_5.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_5.shad[p] = 1
			}

			lp_5.rgb[0] = 0.0
			lp_5.rgb[1] = 1
			lp_5.rgb[2] = 0.0
			lp_5.polyd = 4
			lp_5.P = make([]XYZ, 4)

			lp_5.P[0].X = tree[i].x
			lp_5.P[0].Y = tree[i].y - (tree[i].W2 * 0.5)
			lp_5.P[0].Z = tree[i].z + tree[i].H1
			lp_5.P[1].X = tree[i].x
			lp_5.P[1].Y = tree[i].y - (tree[i].W3 * 0.5)
			lp_5.P[1].Z = tree[i].z + tree[i].H1 + tree[i].H2
			lp_5.P[2].X = tree[i].x + (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_5.P[2].Y = tree[i].y - (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_5.P[2].Z = lp_5.P[1].Z
			lp_5.P[3].X = tree[i].x + (tree[i].W2*0.5)*math.Cos(45*M_rad)
			lp_5.P[3].Y = tree[i].y - (tree[i].W2*0.5)*math.Sin(45*M_rad)
			lp_5.P[3].Z = tree[i].z + tree[i].H1

			HOUSEN2(&lp_5.P[0], &lp_5.P[1], &lp_5.P[2], &lp_5.e)
			lp_5.wb = math.Acos(lp_5.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_5.e.Z = math.Cos(lp_5.wb * M_rad)
			lp_5.e.Y = -math.Sin(lp_5.wb*M_rad) * math.Cos(lp_5.wa*M_rad)
			lp_5.e.X = -math.Sin(lp_5.wb*M_rad) * math.Sin(lp_5.wa*M_rad)
			CAT(&lp_5.e.X, &lp_5.e.Y, &lp_5.e.Z)

			lp = append(lp, lp_5)

			/*----6----*/
			lp_6 := NewP_MENN()

			name = tree[i].treename
			lp_6.opname = name

			lp_6.wa = lp_5.wa - 45
			lp_6.ref = 0.0
			lp_6.wd = 0
			lp_6.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_6.shad[p] = 1
			}

			lp_6.rgb[0] = 0.0
			lp_6.rgb[1] = 1
			lp_6.rgb[2] = 0.0
			lp_6.polyd = 4
			lp_6.P = make([]XYZ, 4)

			lp_6.P[0] = lp_5.P[3]
			lp_6.P[1] = lp_5.P[2]
			lp_6.P[2].X = tree[i].x + (tree[i].W3 * 0.5)
			lp_6.P[2].Y = tree[i].y
			lp_6.P[2].Z = lp_6.P[1].Z
			lp_6.P[3].X = tree[i].x + (tree[i].W2 * 0.5)
			lp_6.P[3].Y = tree[i].y
			lp_6.P[3].Z = lp_6.P[0].Z

			HOUSEN2(&lp_6.P[0], &lp_6.P[1], &lp_6.P[2], &lp_6.e)
			lp_6.wb = math.Acos(lp_6.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_6.e.Z = math.Cos(lp_6.wb * M_rad)
			lp_6.e.Y = -math.Sin(lp_6.wb*M_rad) * math.Cos(lp_6.wa*M_rad)
			lp_6.e.X = -math.Sin(lp_6.wb*M_rad) * math.Sin(lp_6.wa*M_rad)
			CAT(&lp_6.e.X, &lp_6.e.Y, &lp_6.e.Z)

			lp = append(lp, lp_6)

			/*----7----*/
			lp_7 := NewP_MENN()

			name = tree[i].treename
			lp_7.opname = name

			lp_7.wa = lp_6.wa - 45
			lp_7.ref = 0.0
			lp_7.wd = 0
			lp_7.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_7.shad[p] = 1
			}

			lp_7.rgb[0] = 0.0
			lp_7.rgb[1] = 1
			lp_7.rgb[2] = 0.0
			lp_7.polyd = 4
			lp_7.P = make([]XYZ, 4)

			lp_7.P[0] = lp_6.P[3]
			lp_7.P[1] = lp_6.P[2]
			lp_7.P[2].X = tree[i].x + (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_7.P[2].Y = tree[i].y + (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_7.P[2].Z = lp_7.P[1].Z
			lp_7.P[3].X = tree[i].x + (tree[i].W2*0.5)*math.Cos(45*M_rad)
			lp_7.P[3].Y = tree[i].y + (tree[i].W2*0.5)*math.Sin(45*M_rad)
			lp_7.P[3].Z = lp_7.P[0].Z

			HOUSEN2(&lp_7.P[0], &lp_7.P[1], &lp_7.P[2], &lp_7.e)
			lp_7.wb = math.Acos(lp_7.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_7.e.Z = math.Cos(lp_7.wb * M_rad)
			lp_7.e.Y = -math.Sin(lp_7.wb*M_rad) * math.Cos(lp_7.wa*M_rad)
			lp_7.e.X = -math.Sin(lp_7.wb*M_rad) * math.Sin(lp_7.wa*M_rad)
			CAT(&lp_7.e.X, &lp_7.e.Y, &lp_7.e.Z)

			lp = append(lp, lp_7)

			/*----8----*/
			lp_8 := NewP_MENN()

			name = tree[i].treename
			lp_8.opname = name

			lp_8.wa = lp_7.wa - 45
			lp_8.ref = 0.0
			lp_8.wd = 0
			lp_8.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_8.shad[p] = 1
			}

			lp_8.rgb[0] = 0.0
			lp_8.rgb[1] = 1
			lp_8.rgb[2] = 0.0
			lp_8.polyd = 4
			lp_8.P = make([]XYZ, 4)

			lp_8.P[0] = lp_7.P[3]
			lp_8.P[1] = lp_7.P[2]
			lp_8.P[2].X = tree[i].x
			lp_8.P[2].Y = tree[i].y + (tree[i].W3 * 0.5)
			lp_8.P[2].Z = lp_8.P[1].Z
			lp_8.P[3].X = tree[i].x
			lp_8.P[3].Y = tree[i].y + (tree[i].W2 * 0.5)
			lp_8.P[3].Z = lp_8.P[0].Z

			HOUSEN2(&lp_8.P[0], &lp_8.P[1], &lp_8.P[2], &lp_8.e)
			lp_8.wb = math.Acos(lp_8.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_8.e.Z = math.Cos(lp_8.wb * M_rad)
			lp_8.e.Y = -math.Sin(lp_8.wb*M_rad) * math.Cos(lp_8.wa*M_rad)
			lp_8.e.X = -math.Sin(lp_8.wb*M_rad) * math.Sin(lp_8.wa*M_rad)
			CAT(&lp_8.e.X, &lp_8.e.Y, &lp_8.e.Z)

			lp = append(lp, lp_8)

			/*----9----*/
			lp_9 := NewP_MENN()

			name = tree[i].treename
			lp_9.opname = name

			lp_9.wa = 360 + (lp_8.wa - 45)
			lp_9.ref = 0.0
			lp_9.wd = 0
			lp_9.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_9.shad[p] = 1
			}

			lp_9.rgb[0] = 0.0
			lp_9.rgb[1] = 1
			lp_9.rgb[2] = 0.0
			lp_9.polyd = 4
			lp_9.P = make([]XYZ, 4)

			lp_9.P[0] = lp_8.P[3]
			lp_9.P[1] = lp_8.P[2]
			lp_9.P[2].X = tree[i].x - (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_9.P[2].Y = tree[i].y + (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_9.P[2].Z = lp_9.P[1].Z
			lp_9.P[3].X = tree[i].x - (tree[i].W2*0.5)*math.Cos(45*M_rad)
			lp_9.P[3].Y = tree[i].y + (tree[i].W2*0.5)*math.Sin(45*M_rad)
			lp_9.P[3].Z = lp_9.P[0].Z

			HOUSEN2(&lp_9.P[0], &lp_9.P[1], &lp_9.P[2], &lp_9.e)
			lp_9.wb = math.Acos(lp_9.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_9.e.Z = math.Cos(lp_9.wb * M_rad)
			lp_9.e.Y = -math.Sin(lp_9.wb*M_rad) * math.Cos(lp_9.wa*M_rad)
			lp_9.e.X = -math.Sin(lp_9.wb*M_rad) * math.Sin(lp_9.wa*M_rad)
			CAT(&lp_9.e.X, &lp_9.e.Y, &lp_9.e.Z)

			lp = append(lp, lp_9)

			/*----10----*/
			lp_10 := NewP_MENN()

			name = tree[i].treename
			lp_10.opname = name

			lp_10.wa = lp_9.wa - 45
			lp_10.ref = 0.0
			lp_10.wd = 0
			lp_10.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_10.shad[p] = 1
			}

			lp_10.rgb[0] = 0.0
			lp_10.rgb[1] = 1
			lp_10.rgb[2] = 0.0
			lp_10.polyd = 4
			lp_10.P = make([]XYZ, 4)

			lp_10.P[0] = lp_9.P[3]
			lp_10.P[1] = lp_9.P[2]
			lp_10.P[2].X = tree[i].x - (tree[i].W3 * 0.5)
			lp_10.P[2].Y = tree[i].y
			lp_10.P[2].Z = lp_10.P[1].Z
			lp_10.P[3].X = tree[i].x - (tree[i].W2 * 0.5)
			lp_10.P[3].Y = tree[i].y
			lp_10.P[3].Z = lp_10.P[0].Z

			HOUSEN2(&lp_10.P[0], &lp_10.P[1], &lp_10.P[2], &lp_10.e)
			lp_10.wb = math.Acos(lp_10.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_10.e.Z = math.Cos(lp_10.wb * M_rad)
			lp_10.e.Y = -math.Sin(lp_10.wb*M_rad) * math.Cos(lp_10.wa*M_rad)
			lp_10.e.X = -math.Sin(lp_10.wb*M_rad) * math.Sin(lp_10.wa*M_rad)
			CAT(&lp_10.e.X, &lp_10.e.Y, &lp_10.e.Z)

			lp = append(lp, lp_10)

			/*----11----*/
			lp_11 := NewP_MENN()

			name = tree[i].treename
			lp_11.opname = name

			lp_11.wa = lp_10.wa - 45
			lp_11.ref = 0.0
			lp_11.wd = 0
			lp_11.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_11.shad[p] = 1
			}

			lp_11.rgb[0] = 0.0
			lp_11.rgb[1] = 1
			lp_11.rgb[2] = 0.0
			lp_11.polyd = 4
			lp_11.P = make([]XYZ, 4)

			lp_11.P[0] = lp_10.P[3]
			lp_11.P[1] = lp_10.P[2]
			lp_11.P[2].X = tree[i].x - (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_11.P[2].Y = tree[i].y - (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_11.P[2].Z = lp_11.P[1].Z
			lp_11.P[3].X = tree[i].x - (tree[i].W2*0.5)*math.Cos(45*M_rad)
			lp_11.P[3].Y = tree[i].y - (tree[i].W2*0.5)*math.Sin(45*M_rad)
			lp_11.P[3].Z = lp_11.P[0].Z

			HOUSEN2(&lp_11.P[0], &lp_11.P[1], &lp_11.P[2], &lp_11.e)
			lp_11.wb = math.Acos(lp_11.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_11.e.Z = math.Cos(lp_11.wb * M_rad)
			lp_11.e.Y = -math.Sin(lp_11.wb*M_rad) * math.Cos(lp_11.wa*M_rad)
			lp_11.e.X = -math.Sin(lp_11.wb*M_rad) * math.Sin(lp_11.wa*M_rad)
			CAT(&lp_11.e.X, &lp_11.e.Y, &lp_11.e.Z)

			lp = append(lp, lp_11)

			/*----12----*/
			lp_12 := NewP_MENN()

			name = tree[i].treename
			lp_12.opname = name

			lp_12.wa = lp_11.wa - 45
			lp_12.ref = 0.0
			lp_12.wd = 0
			lp_12.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_12.shad[p] = 1
			}

			lp_12.rgb[0] = 0.0
			lp_12.rgb[1] = 1
			lp_12.rgb[2] = 0.0
			lp_12.polyd = 4
			lp_12.P = make([]XYZ, 4)

			lp_12.P[0] = lp_11.P[3]
			lp_12.P[1] = lp_11.P[2]
			lp_12.P[2].X = tree[i].x
			lp_12.P[2].Y = tree[i].y - (tree[i].W3 * 0.5)
			lp_12.P[2].Z = lp_12.P[1].Z
			lp_12.P[3].X = tree[i].x
			lp_12.P[3].Y = tree[i].y - (tree[i].W2 * 0.5)
			lp_12.P[3].Z = lp_12.P[0].Z

			HOUSEN2(&lp_12.P[0], &lp_12.P[1], &lp_12.P[2], &lp_12.e)
			lp_12.wb = math.Acos(lp_12.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_12.e.Z = math.Cos(lp_12.wb * M_rad)
			lp_12.e.Y = -math.Sin(lp_12.wb*M_rad) * math.Cos(lp_12.wa*M_rad)
			lp_12.e.X = -math.Sin(lp_12.wb*M_rad) * math.Sin(lp_12.wa*M_rad)
			CAT(&lp_12.e.X, &lp_12.e.Y, &lp_12.e.Z)

			lp = append(lp, lp_12)

			/*----13----*/
			lp_13 := NewP_MENN()

			name = tree[i].treename
			lp_13.opname = name

			lp_13.wa = -22.5
			lp_13.ref = 0.0
			lp_13.wd = 0
			lp_13.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_13.shad[p] = 1
			}

			lp_13.rgb[0] = 0.0
			lp_13.rgb[1] = 1
			lp_13.rgb[2] = 0.0
			lp_13.polyd = 4
			lp_13.P = make([]XYZ, 4)

			lp_13.P[0].X = tree[i].x
			lp_13.P[0].Y = tree[i].y - (tree[i].W3 * 0.5)
			lp_13.P[0].Z = tree[i].z + tree[i].H1 + tree[i].H2
			lp_13.P[1].X = tree[i].x
			lp_13.P[1].Y = tree[i].y - (tree[i].W4 * 0.5)
			lp_13.P[1].Z = tree[i].z + tree[i].H1 + tree[i].H2 + tree[i].H3
			lp_13.P[2].X = tree[i].x + (tree[i].W4*0.5)*math.Cos(45*M_rad)
			lp_13.P[2].Y = tree[i].y - (tree[i].W4*0.5)*math.Sin(45*M_rad)
			lp_13.P[2].Z = lp_13.P[1].Z
			lp_13.P[3].X = tree[i].x + (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_13.P[3].Y = tree[i].y - (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_13.P[3].Z = lp_13.P[0].Z

			HOUSEN2(&lp_13.P[0], &lp_13.P[1], &lp_13.P[2], &lp_13.e)
			lp_13.wb = math.Acos(lp_13.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_13.e.Z = math.Cos(lp_13.wb * M_rad)
			lp_13.e.Y = -math.Sin(lp_13.wb*M_rad) * math.Cos(lp_13.wa*M_rad)
			lp_13.e.X = -math.Sin(lp_13.wb*M_rad) * math.Sin(lp_13.wa*M_rad)
			CAT(&lp_13.e.X, &lp_13.e.Y, &lp_13.e.Z)

			lp = append(lp, lp_13)

			/*----14----*/
			lp_14 := NewP_MENN()

			name = tree[i].treename
			lp_14.opname = name

			lp_14.wa = lp_13.wa - 45
			lp_14.ref = 0.0
			lp_14.wd = 0
			lp_14.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_14.shad[p] = 1
			}

			lp_14.rgb[0] = 0.0
			lp_14.rgb[1] = 1
			lp_14.rgb[2] = 0.0
			lp_14.polyd = 4
			lp_14.P = make([]XYZ, 4)

			lp_14.P[0] = lp_13.P[3]
			lp_14.P[1] = lp_13.P[2]
			lp_14.P[2].X = tree[i].x + (tree[i].W4 * 0.5)
			lp_14.P[2].Y = tree[i].y
			lp_14.P[2].Z = lp_14.P[1].Z
			lp_14.P[3].X = tree[i].x + (tree[i].W3 * 0.5)
			lp_14.P[3].Y = tree[i].y
			lp_14.P[3].Z = lp_14.P[0].Z

			HOUSEN2(&lp_14.P[0], &lp_14.P[1], &lp_14.P[2], &lp_14.e)
			lp_14.wb = math.Acos(lp_14.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_14.e.Z = math.Cos(lp_14.wb * M_rad)
			lp_14.e.Y = -math.Sin(lp_14.wb*M_rad) * math.Cos(lp_14.wa*M_rad)
			lp_14.e.X = -math.Sin(lp_14.wb*M_rad) * math.Sin(lp_14.wa*M_rad)
			CAT(&lp_14.e.X, &lp_14.e.Y, &lp_14.e.Z)

			lp = append(lp, lp_14)

			/*----15----*/
			lp_15 := NewP_MENN()

			name = tree[i].treename
			lp_15.opname = name

			lp_15.wa = lp_14.wa - 45
			lp_15.ref = 0.0
			lp_15.wd = 0
			lp_15.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_15.shad[p] = 1
			}

			lp_15.rgb[0] = 0.0
			lp_15.rgb[1] = 1
			lp_15.rgb[2] = 0.0
			lp_15.polyd = 4
			lp_15.P = make([]XYZ, 4)

			lp_15.P[0] = lp_14.P[3]
			lp_15.P[1] = lp_14.P[2]
			lp_15.P[2].X = tree[i].x + (tree[i].W4*0.5)*math.Cos(45*M_rad)
			lp_15.P[2].Y = tree[i].y + (tree[i].W4*0.5)*math.Sin(45*M_rad)
			lp_15.P[2].Z = lp_15.P[1].Z
			lp_15.P[3].X = tree[i].x + (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_15.P[3].Y = tree[i].y + (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_15.P[3].Z = lp_15.P[0].Z

			HOUSEN2(&lp_15.P[0], &lp_15.P[1], &lp_15.P[2], &lp_15.e)
			lp_15.wb = math.Acos(lp_15.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_15.e.Z = math.Cos(lp_15.wb * M_rad)
			lp_15.e.Y = -math.Sin(lp_15.wb*M_rad) * math.Cos(lp_15.wa*M_rad)
			lp_15.e.X = -math.Sin(lp_15.wb*M_rad) * math.Sin(lp_15.wa*M_rad)
			CAT(&lp_15.e.X, &lp_15.e.Y, &lp_15.e.Z)

			lp = append(lp, lp_15)

			/*----16----*/
			lp_16 := NewP_MENN()

			name = tree[i].treename
			lp_16.opname = name

			lp_16.wa = lp_15.wa - 45
			lp_16.ref = 0.0
			lp_16.wd = 0
			lp_16.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_16.shad[p] = 1
			}

			lp_16.rgb[0] = 0.0
			lp_16.rgb[1] = 1
			lp_16.rgb[2] = 0.0
			lp_16.polyd = 4
			lp_16.P = make([]XYZ, 4)

			lp_16.P[0] = lp_15.P[3]
			lp_16.P[1] = lp_15.P[2]
			lp_16.P[2].X = tree[i].x
			lp_16.P[2].Y = tree[i].y + (tree[i].W4 * 0.5)
			lp_16.P[2].Z = lp_16.P[1].Z
			lp_16.P[3].X = tree[i].x
			lp_16.P[3].Y = tree[i].y + (tree[i].W3 * 0.5)
			lp_16.P[3].Z = lp_16.P[0].Z

			HOUSEN2(&lp_16.P[0], &lp_16.P[1], &lp_16.P[2], &lp_16.e)
			lp_16.wb = math.Acos(lp_16.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_16.e.Z = math.Cos(lp_16.wb * M_rad)
			lp_16.e.Y = -math.Sin(lp_16.wb*M_rad) * math.Cos(lp_16.wa*M_rad)
			lp_16.e.X = -math.Sin(lp_16.wb*M_rad) * math.Sin(lp_16.wa*M_rad)
			CAT(&lp_16.e.X, &lp_16.e.Y, &lp_16.e.Z)

			lp = append(lp, lp_16)

			/*---17---*/
			lp_17 := NewP_MENN()

			name = tree[i].treename
			lp_17.opname = name

			lp_17.wa = 360 + (lp_16.wa - 45)
			lp_17.ref = 0.0
			lp_17.wd = 0
			lp_17.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_17.shad[p] = 1
			}

			lp_17.rgb[0] = 0.0
			lp_17.rgb[1] = 1
			lp_17.rgb[2] = 0.0
			lp_17.polyd = 4
			lp_17.P = make([]XYZ, 4)

			lp_17.P[0] = lp_16.P[3]
			lp_17.P[1] = lp_16.P[2]
			lp_17.P[2].X = tree[i].x - (tree[i].W4*0.5)*math.Cos(45*M_rad)
			lp_17.P[2].Y = tree[i].y + (tree[i].W4*0.5)*math.Sin(45*M_rad)
			lp_17.P[2].Z = lp_17.P[1].Z
			lp_17.P[3].X = tree[i].x - (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_17.P[3].Y = tree[i].y + (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_17.P[3].Z = lp_17.P[0].Z

			HOUSEN2(&lp_17.P[0], &lp_17.P[1], &lp_17.P[2], &lp_17.e)
			lp_17.wb = math.Acos(lp_17.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_17.e.Z = math.Cos(lp_17.wb * M_rad)
			lp_17.e.Y = -math.Sin(lp_17.wb*M_rad) * math.Cos(lp_17.wa*M_rad)
			lp_17.e.X = -math.Sin(lp_17.wb*M_rad) * math.Sin(lp_17.wa*M_rad)
			CAT(&lp_17.e.X, &lp_17.e.Y, &lp_17.e.Z)

			lp = append(lp, lp_17)

			/*----18----*/
			lp_18 := NewP_MENN()

			name = tree[i].treename
			lp_18.opname = name

			lp_18.wa = lp_17.wa - 45
			lp_18.ref = 0.0
			lp_18.wd = 0
			lp_18.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_18.shad[p] = 1
			}

			lp_18.rgb[0] = 0.0
			lp_18.rgb[1] = 1
			lp_18.rgb[2] = 0.0
			lp_18.polyd = 4
			lp_18.P = make([]XYZ, 4)

			lp_18.P[0] = lp_17.P[3]
			lp_18.P[1] = lp_17.P[2]
			lp_18.P[2].X = tree[i].x - (tree[i].W4 * 0.5)
			lp_18.P[2].Y = tree[i].y
			lp_18.P[2].Z = lp_18.P[1].Z
			lp_18.P[3].X = tree[i].x - (tree[i].W3 * 0.5)
			lp_18.P[3].Y = tree[i].y
			lp_18.P[3].Z = lp_18.P[0].Z

			HOUSEN2(&lp_18.P[0], &lp_18.P[1], &lp_18.P[2], &lp_18.e)
			lp_18.wb = math.Acos(lp_18.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_18.e.Z = math.Cos(lp_18.wb * M_rad)
			lp_18.e.Y = -math.Sin(lp_18.wb*M_rad) * math.Cos(lp_18.wa*M_rad)
			lp_18.e.X = -math.Sin(lp_18.wb*M_rad) * math.Sin(lp_18.wa*M_rad)
			CAT(&lp_18.e.X, &lp_18.e.Y, &lp_18.e.Z)

			lp = append(lp, lp_18)

			/*---19----*/
			lp_19 := NewP_MENN()

			name = tree[i].treename
			lp_19.opname = name

			lp_19.wa = lp_18.wa - 45
			lp_19.ref = 0.0
			lp_19.wd = 0
			lp_19.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_19.shad[p] = 1
			}

			lp_19.rgb[0] = 0.0
			lp_19.rgb[1] = 1
			lp_19.rgb[2] = 0.0
			lp_19.polyd = 4
			lp_19.P = make([]XYZ, 4)

			lp_19.P[0] = lp_18.P[3]
			lp_19.P[1] = lp_18.P[2]
			lp_19.P[2].X = tree[i].x - (tree[i].W4*0.5)*math.Cos(45*M_rad)
			lp_19.P[2].Y = tree[i].y - (tree[i].W4*0.5)*math.Sin(45*M_rad)
			lp_19.P[2].Z = lp_19.P[1].Z
			lp_19.P[3].X = tree[i].x - (tree[i].W3*0.5)*math.Cos(45*M_rad)
			lp_19.P[3].Y = tree[i].y - (tree[i].W3*0.5)*math.Sin(45*M_rad)
			lp_19.P[3].Z = lp_19.P[0].Z

			HOUSEN2(&lp_19.P[0], &lp_19.P[1], &lp_19.P[2], &lp_19.e)
			lp_19.wb = math.Acos(lp_19.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_19.e.Z = math.Cos(lp_19.wb * M_rad)
			lp_19.e.Y = -math.Sin(lp_19.wb*M_rad) * math.Cos(lp_19.wa*M_rad)
			lp_19.e.X = -math.Sin(lp_19.wb*M_rad) * math.Sin(lp_19.wa*M_rad)
			CAT(&lp_19.e.X, &lp_19.e.Y, &lp_19.e.Z)

			lp = append(lp, lp_19)

			/*---20----*/
			lp_20 := NewP_MENN()

			name = tree[i].treename
			lp_20.opname = name

			lp_20.wa = lp_19.wa - 45
			lp_20.ref = 0.0
			lp_20.wd = 0
			lp_20.sbflg = 0

			for p = 0; p < 366; p++ {
				lp_20.shad[p] = 1
			}

			lp_20.rgb[0] = 0.0
			lp_20.rgb[1] = 1
			lp_20.rgb[2] = 0.0
			lp_20.polyd = 4
			lp_20.P = make([]XYZ, 4)

			lp_20.P[0] = lp_19.P[3]
			lp_20.P[1] = lp_19.P[2]
			lp_20.P[2].X = tree[i].x
			lp_20.P[2].Y = tree[i].y - (tree[i].W4 * 0.5)
			lp_20.P[2].Z = lp_20.P[1].Z
			lp_20.P[3].X = tree[i].x
			lp_20.P[3].Y = tree[i].y - (tree[i].W3 * 0.5)
			lp_20.P[3].Z = lp_20.P[0].Z

			HOUSEN2(&lp_20.P[0], &lp_20.P[1], &lp_20.P[2], &lp_20.e)
			lp_20.wb = math.Acos(lp_20.e.Z) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_20.e.Z = math.Cos(lp_20.wb * M_rad)
			lp_20.e.Y = -math.Sin(lp_20.wb*M_rad) * math.Cos(lp_20.wa*M_rad)
			lp_20.e.X = -math.Sin(lp_20.wb*M_rad) * math.Sin(lp_20.wa*M_rad)
			CAT(&lp_20.e.X, &lp_20.e.Y, &lp_20.e.Z)

			lp = append(lp, lp_20)
		} else {
			fmt.Printf("error--**COORDNT-lp TREE \n")
			os.Exit(1)
		}
	}

	/*-------多角形の障害物の直接入力----------------------*/
	for i := range poly {

		if poly[i].polyknd == "OBS" {
			lp_k := NewP_MENN()
			lp_k.opname = poly[i].polyname

			lp_k.polyd = poly[i].polyd
			lp_k.P = make([]XYZ, lp_k.polyd)

			for j = 0; j < lp_k.polyd; j++ {
				lp_k.P[j] = poly[i].P[j]
			}

			HOUSEN2(&lp_k.P[0], &lp_k.P[1], &lp_k.P[2], &lp_k.e)
			lp_k.wb = math.Acos(lp_k.e.Z) * (180 / math.Pi)
			lp_k.wa = math.Asin(lp_k.e.X/math.Sin(lp_k.wb*M_rad)) * (180 / math.Pi)
			/*--法線ベクトルの算出　HOUSEN2関数を使うと、向きが逆になるので、変更 091128 higuchi --*/
			lp_k.e.Z = math.Cos(lp_k.wb * M_rad)
			lp_k.e.Y = -math.Sin(lp_k.wb*M_rad) * math.Cos(lp_k.wa*M_rad)
			lp_k.e.X = -math.Sin(lp_k.wb*M_rad) * math.Sin(lp_k.wa*M_rad)
			CAT(&lp_k.e.X, &lp_k.e.Y, &lp_k.e.Z)

			lp_k.ref = 0.0
			lp_k.wd = 0
			lp_k.sbflg = 0

			lp_k.rgb[0] = poly[i].rgb[0]
			lp_k.rgb[1] = poly[i].rgb[1]
			lp_k.rgb[2] = poly[i].rgb[2]

			for p = 0; p < 366; p++ {
				lp_k.shad[p] = 1
			}

			lp = append(lp, lp_k)
		}
	}

	return lp
}

/*
OP_COORDNT (Opening Plane Coordinate Transformation)

この関数は、建物の開口部（窓など）や、
多角形データで定義された受光面（窓、壁面など）の幾何学的情報を、
日射量計算や昼光利用計算に用いられる「受光面（Opening Plane）」の座標データに変換します。

建築環境工学的な観点:
- **受光面のモデル化**: 建物の窓面や壁面は、太陽光を受け入れる主要な面です。
  この関数は、これらの受光面の形状と位置を正確にモデル化し、
  日射量計算や昼光利用計算の基礎となる座標データ（`op`）を生成します。
- **座標変換**: 建物の相対的な位置（`bp.x0`, `bp.y0`, `bp.z0`）や、
  方位角（`bp.Wa`）と傾斜角（`bp.Wb`）を考慮して、
  各受光面（`rmp`）の頂点座標（`opj.P`）を計算します。
  これにより、3次元空間における受光面の正確な形状を表現できます。
- **法線ベクトルの算出**: 各受光面の法線ベクトル（`opj.e`）を算出します。
  法線ベクトルは、その面が太陽光に対してどの方向を向いているかを示し、
  日射入射角の計算や、日射熱取得量の計算に用いられます。
- **窓のモデル化**: `rmp.WD`は、
  受光面内に含まれる窓（`winname`）の情報を格納しており、
  窓の形状や位置を詳細にモデル化できます。
- **日射熱取得と昼光利用の評価**: この関数で生成される受光面の座標データは、
  - **日射熱取得の予測**: 窓を透過して室内に侵入する日射熱量を正確に予測し、
    冷房負荷を評価します。
  - **昼光利用の検討**: 窓からの昼光の取り込み量を評価し、
    照明エネルギーの削減効果を予測します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func OP_COORDNT(BDP []*BBDP, poly []*POLYGN) []*P_MENN {
	const M_rad = math.Pi / 180.0

	op := make([]*P_MENN, 0)

	for _, bp := range BDP {
		// 方位角と傾斜角の三角関数の計算
		cWa := math.Cos(-bp.Wa * M_rad)
		sWa := math.Sin(-bp.Wa * M_rad)
		cWb := math.Cos(bp.Wb * M_rad)
		sWb := math.Sin(bp.Wb * M_rad)

		for _, rmp := range bp.RMP {
			a := 0.0
			b := 0.0
			c := 0.0

			opj := NewP_MENN()
			opj.wd = len(rmp.WD)
			opj.wlflg = 0

			opj.opname = rmp.rmpname

			opj.rgb[0] = rmp.rgb[0]
			opj.rgb[1] = rmp.rgb[1]
			opj.rgb[2] = rmp.rgb[2]
			opj.polyd = 4
			opj.P = make([]XYZ, 4)

			a = bp.x0 + rmp.xb0*cWa - rmp.yb0*cWb*sWa
			b = bp.y0 + rmp.xb0*sWa + rmp.yb0*cWb*cWa
			c = bp.z0 + rmp.yb0*sWb

			opj.P[0].X = a
			opj.P[0].Y = b
			opj.P[0].Z = c

			opj.P[1].X = a - rmp.Rh*cWb*sWa
			opj.P[1].Y = b + rmp.Rh*cWb*cWa
			opj.P[1].Z = c + rmp.Rh*sWb

			opj.P[2].X = rmp.Rw*cWa + opj.P[1].X
			opj.P[2].Y = rmp.Rw*sWa + opj.P[1].Y
			opj.P[2].Z = opj.P[1].Z

			opj.P[3].X = a + rmp.Rw*cWa
			opj.P[3].Y = b + rmp.Rw*sWa
			opj.P[3].Z = c

			opj.grpx = rmp.grpx
			opj.wa = bp.Wa
			opj.wb = bp.Wb
			opj.ref = rmp.ref
			opj.e.Z = math.Cos(opj.wb * M_rad)
			opj.e.Y = -math.Sin(opj.wb*M_rad) * math.Cos(opj.wa*M_rad)
			opj.e.X = -math.Sin(opj.wb*M_rad) * math.Sin(opj.wa*M_rad)
			CAT(&opj.e.X, &opj.e.Y, &opj.e.Z)

			opj.Nopw = len(rmp.WD)

			// opw構造体のメモリ確保
			if opj.Nopw > 0 {
				opj.opw = make([]WD_MENN, opj.Nopw)
				for m := 0; m < opj.Nopw; m++ {
					opw := opj.opw[m]

					opw.opwname = ""
					opw.P = make([]XYZ, 0)
					opw.ref = 0.0
					opw.grpx = 0.0
					opw.sumw = 0.0
					matinit(opw.rgb[:], 3)
				}
			}
			if len(rmp.WD) != 0 {
				for i := range rmp.WD {
					ax := 0.0
					by := 0.0
					cz := 0.0

					opj.opw[i].opwname = rmp.WD[i].winname

					ax = a + rmp.WD[i].xr*cWa - rmp.WD[i].yr*cWb*sWa
					by = b + rmp.WD[i].xr*sWa + rmp.WD[i].yr*cWb*cWa
					cz = c + rmp.WD[i].yr*sWb

					opj.opw[i].grpx = rmp.WD[i].grpx
					opj.opw[i].ref = rmp.WD[i].ref

					opj.opw[i].rgb[0] = rmp.WD[i].rgb[0]
					opj.opw[i].rgb[1] = rmp.WD[i].rgb[1]
					opj.opw[i].rgb[2] = rmp.WD[i].rgb[2]
					opj.opw[i].P = make([]XYZ, 4)

					opj.opw[i].P[0].X = ax
					opj.opw[i].P[0].Y = by
					opj.opw[i].P[0].Z = cz

					opj.opw[i].P[1].X = ax - rmp.WD[i].Wh*cWb*sWa
					opj.opw[i].P[1].Y = by + rmp.WD[i].Wh*cWb*cWa
					opj.opw[i].P[1].Z = cz + rmp.WD[i].Wh*sWb
					opj.opw[i].P[2].X = ax + rmp.WD[i].Ww*cWa - rmp.WD[i].Wh*cWb*sWa
					opj.opw[i].P[2].Y = by + rmp.WD[i].Ww*sWa + rmp.WD[i].Wh*cWb*cWa
					opj.opw[i].P[2].Z = cz + rmp.WD[i].Wh*sWb
					opj.opw[i].P[3].X = ax + rmp.WD[i].Ww*cWa
					opj.opw[i].P[3].Y = by + rmp.WD[i].Ww*sWa
					opj.opw[i].P[3].Z = cz
				}
			}

			op = append(op, opj)
		}
	}

	for _, p := range poly {
		if p.polyknd == "RMP" {

			op_sum := NewP_MENN()
			op_sum.opname = p.polyname
			op_sum.ref = p.ref
			op_sum.wd = 0
			op_sum.grpx = p.grpx

			op_sum.rgb[0] = p.rgb[0]
			op_sum.rgb[1] = p.rgb[1]
			op_sum.rgb[2] = p.rgb[2]
			op_sum.polyd = p.polyd
			op_sum.P = make([]XYZ, op_sum.polyd)
			for j := 0; j < op_sum.polyd; j++ {
				op_sum.P[j] = p.P[j]
			}

			HOUSEN2(&op_sum.P[0], &op_sum.P[1], &op_sum.P[2], &op_sum.e)
			op_sum.wb = math.Acos(op_sum.e.Z) * (180 / math.Pi)
			op_sum.wa = math.Asin(op_sum.e.X/math.Sin(op_sum.wb*M_rad)) * (180 / math.Pi)

			/*--ポリゴンの法線ベクトルが逆で出てしまうのを修正 091128 higuchi--*/
			op_sum.e.Z = math.Cos(op_sum.wb * M_rad)
			op_sum.e.Y = -math.Sin(op_sum.wb*M_rad) * math.Cos(op_sum.wa*M_rad)
			op_sum.e.X = -math.Sin(op_sum.wb*M_rad) * math.Sin(op_sum.wa*M_rad)

			op = append(op, op_sum)
		}

	}

	return op
}
