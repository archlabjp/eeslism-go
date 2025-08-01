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

/*
SHADOW (Shading Calculation for Opening Planes)

この関数は、建物の開口部（窓など）の主面（`op`）に対する影の面積を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、主面を微小なグリッドに分割し、
  各グリッド点から太陽方向へ光線を射出します。
  その光線が他の障害物（建物自身や外部障害物）によって遮られるかどうかを判定することで、
  影の面積を計算します。
  - `DEM`: グリッドの細かさ（微小四角形の辺の長さ）。
  - `AMAX`, `BMAX`: グリッドの分割数。
  - `Qx, Qy, Qz`: グリッド点の座標。
  - `ls, ms, ns`: 太陽光線ベクトル。
  - `KOUTEN`: 光線と障害物表面の交点を計算します。
  - `INOROUT`: 交点が障害物表面の多角形の内部にあるかどうかを判断します。
  - `YOGEN`: 光線が障害物表面に当たる角度を計算します。
- **影面積の計算**: 光線が障害物に当たった場合、
  そのグリッド点が影になったと判断し、
  影面積（`op.sum`）に寄与させます。
  窓面（`op.opw`）の影も個別に計算されます。
- **日射熱取得の予測**: この関数で計算される影面積は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func SHADOW(
	g int,
	DE float64,
	opn int,
	lpn int,
	ls float64,
	ms float64,
	ns float64,
	s *bekt,
	t *bekt,
	op *P_MENN,
	OP []*P_MENN,
	LP []*P_MENN,
	wap *float64,
	wip []float64,
	nday int,
) {
	var Am, Bm, Amx, Amy, Amz, Bmx, Bmy, Bmz, A, B float64
	var i, j, h, k int
	var Px, Py, Pz, Qx, Qy, Qz, OX, OY, OZ, OX1, OY1, OZ1, SHITA float64
	var rp, rp2 int /*--多角形ループ-*/
	var naibu int   /*--壁面の内部かフラグ*/
	var naibuw int  /*--窓面の内部かフラグ*/
	var omote int   /*--面から見えるかフラグ　0:見えない 1:見える--*/
	var DEM float64
	var S, T float64
	var flg int
	var tau float64 /*--透過率（合計）--*/

	var loopA, loopB int /* ループカウンター　横・縦  20170821 higuchi*/
	var AMAX, BMAX int   /* 分割数　幅・高さ*/

	DEM = DE / 1000. /*--mmをmに変換--*/
	rp = op.polyd - 2

	for i = 0; i < rp; i++ {
		Amx = op.P[i+1].X - op.P[0].X
		Amy = op.P[i+1].Y - op.P[0].Y
		Amz = op.P[i+1].Z - op.P[0].Z
		Am = math.Sqrt(Amx*Amx + Amy*Amy + Amz*Amz)
		Bmx = op.P[i+2].X - op.P[0].X
		Bmy = op.P[i+2].Y - op.P[0].Y
		Bmz = op.P[i+2].Z - op.P[0].Z
		Bm = math.Sqrt(Bmx*Bmx + Bmy*Bmy + Bmz*Bmz)

		/* 20170821 higuchi add*/
		AMAX = int(math.Ceil(Am / DEM))
		BMAX = int(math.Ceil(Bm / DEM))

		for loopA = 0; loopA < AMAX-1; loopA++ {
			for loopB = 0; loopB < BMAX-1; loopB++ {

				A = (DEM / 2.0) + DEM*float64(loopA)
				B = (DEM / 2.0) + DEM*float64(loopB)

				//for (A = DEM / 2.; A < Am; A = A + DEM){
				//for (B = DEM / 2.; B < Bm; B = B + DEM){
				Px = 0.0
				Py = 0.0
				Pz = 0.0
				Qx = 0.0
				Qy = 0.0
				Qz = 0.0
				S = 0.0
				T = 0.0
				OX = 0.0
				OY = 0.0
				OZ = 0.0
				OX1 = 0.0
				OY1 = 0.0
				OZ1 = 0.0
				SHITA = 0.0

				/*----
				printf("DEM=%f,A=%f,B=%f\n",DEM,A,B) ;
				----*/
				/*--グリッドを移動するＱ点の座標を求める---*/
				OX = (A / Am) * Amx
				OY = (A / Am) * Amy
				OZ = (A / Am) * Amz
				OX1 = (B / Bm) * Bmx
				OY1 = (B / Bm) * Bmy
				OZ1 = (B / Bm) * Bmz

				Qx = op.P[0].X + OX + OX1
				Qy = op.P[0].Y + OY + OY1
				Qz = op.P[0].Z + OZ + OZ1

				naibu = 0
				/*------QがＯＰの内部にあるか------*/
				rp2 = op.polyd - 2

				for j = 0; j < rp2; j++ {
					INOROUT(Qx, Qy, Qz, op.P[0], op.P[j+1], op.P[j+2], &S, &T)

					if (S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0) {
						(*wap) = (*wap) + 1.0
						naibu = 1   /*--壁内部--*/
						naibuw = -1 /*--窓番号初期化--*/
						/*---Qが窓面にあるか----*/
						for k = 0; k < op.wd; k++ {
							for h = 0; h < 2; h++ {
								INOROUT(Qx, Qy, Qz, op.opw[k].P[0], op.opw[k].P[h+1], op.opw[k].P[h+2], &S, &T)
								if (S >= 0.0 && T >= 0.0) && (S+T) <= 1.0 {
									wip[k] = wip[k] + 1.0
									naibuw = k
									break
								}
							}
							if naibuw != -1 {
								break
							}
						}
						break
					}
				}
				if naibu == 0 {
					continue
				}

				/*---建物自身による影を考慮----*/
				for j = 0; j < opn; j++ {
					omote = 0
					for k = 0; k < OP[j].polyd; k++ {
						if s.ps[j][k] > 0.0 {
							omote = 1
							break
						}
					}
					if (g != j) && (omote > 0) {
						KOUTEN(Qx, Qy, Qz, ls, ms, ns, &Px, &Py, &Pz, OP[j].P[0], OP[j].e)
						YOGEN(Qx, Qy, Qz, Px, Py, Pz, &SHITA, op.e)
						rp2 = OP[j].polyd - 2
						for k = 0; k < rp2; k++ { /*--多角形ループ---*/

							INOROUT(Px, Py, Pz, OP[j].P[0], OP[j].P[k+1], OP[j].P[k+2], &S, &T)
							if ((S > 0.0 && T > 0.0) && ((S + T) < 1.0)) && (SHITA > 0) {

								if naibu == 1 {
									op.sum = op.sum + 1.0
								}
								if naibuw >= 0 {
									op.opw[naibuw].sumw = op.opw[naibuw].sumw + 1.0
								}
								goto koko
							}
						}
					}
				}

				/*-------------障害物による影を考慮----------------------*/
				flg = 0 // 100703 higuchi add
				for h = 0; h < lpn; h++ {
					omote = 0
					for k = 0; k < LP[h].polyd; k++ {
						if t.ps[h][k] > 0.0 {
							omote = 1
							break
						}
					}
					if omote > 0 {

						//	    flg = 0 ;  100703 higuchi dell
						KOUTEN(Qx, Qy, Qz, ls, ms, ns, &Px, &Py, &Pz, LP[h].P[0], LP[h].e)
						YOGEN(Qx, Qy, Qz, Px, Py, Pz, &SHITA, op.e)

						rp2 = LP[h].polyd - 2
						for k = 0; k < rp2; k++ {

							INOROUT(Px, Py, Pz, LP[h].P[0], LP[h].P[k+1], LP[h].P[k+2], &S, &T)
							if ((S > 0.0 && T > 0.0) && ((S + T) < 1.0)) && (SHITA > 0) {

								if flg == 0 {
									tau = 1. - LP[h].shad[nday]
									flg = 1
								} else {
									tau = tau * (1. - LP[h].shad[nday])
								}

								break //100703 higuchi add
								//		    op.sum = op.sum + LP[h].shad[nday] ;
								//            if(naibuw >= 0)
								//		       op.opw[naibuw].sumw=op.opw[naibuw].sumw+LP[h].shad[nday] ;
								//		    goto koko ;
							}
						} // for

					}
				} // for

				op.sum = op.sum + (1. - tau)

				if naibuw >= 0 {
					op.opw[naibuw].sumw = op.opw[naibuw].sumw + (1 - tau)
				}
				tau = 1
			koko:
			}
		}
	}

	for i = 0; i < op.wd; i++ {
		(*wap) = (*wap) - wip[i]
		op.sum = op.sum - op.opw[i].sumw
		if wip[i] == 0 {
			op.opw[i].sumw = 0
		} else {
			op.opw[i].sumw = (op.opw[i].sumw / wip[i])
		}
	}

	if *wap == 0 {
		op.sum = 0
	} else {
		op.sum = op.sum / (*wap)
	}
}

/*----------------------------------------------------------------------*/
/*
SHADOWlp (Shading Calculation for Light-Receiving Planes)

この関数は、建物の被受照面（`lp`）に対する影の面積を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **日影計算の幾何学**: 日影は、太陽光線が障害物によって遮られることで形成されます。
  この関数は、被受照面を微小なグリッドに分割し、
  各グリッド点から太陽方向へ光線を射出します。
  その光線が他の障害物（建物自身や外部障害物）によって遮られるかどうかを判定することで、
  影の面積を計算します。
  - `DEM`: グリッドの細かさ（微小四角形の辺の長さ）。
  - `AMAX`, `BMAX`: グリッドの分割数。
  - `Qx, Qy, Qz`: グリッド点の座標。
  - `ls, ms, ns`: 太陽光線ベクトル。
  - `KOUTEN`: 光線と障害物表面の交点を計算します。
  - `INOROUT`: 交点が障害物表面の多角形の内部にあるかどうかを判断します。
  - `YOGEN`: 光線が障害物表面に当たる角度を計算します。
- **影面積の計算**: 光線が障害物に当たった場合、
  そのグリッド点が影になったと判断し、
  影面積（`lp.sum`）に寄与させます。
- **日射熱取得の予測**: この関数で計算される影面積は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func SHADOWlp(
	g int,
	DE float64,
	lpn int,
	mpn int,
	ls float64,
	ms float64,
	ns float64,
	s *bekt,
	t *bekt,
	lp *P_MENN,
	LP []P_MENN,
	MP []P_MENN,
) {
	var Am, Bm, Amx, Amy, Amz, Bmx, Bmy, Bmz, A, B, wap float64
	var i, j, k, h int
	var Px, Py, Pz, Qx, Qy, Qz float64
	var S, T float64
	var OX, OY, OZ, OX1, OY1, OZ1, SHITA float64
	var omote int   /*--面から見えるかフラグ-*/
	var rp, rp2 int /*--多角形ループ回数--*/

	var loopA, loopB int /* ループカウンター　横・縦  20170821 higuchi*/
	var AMAX, BMAX int   /* 分割数　幅・高さ*/

	var DEM float64         /*higuchi add 20180125*/
	DEM = (DE * 5) / 1000.0 /*--mmをmに変換-- 粗目のグリッドにする　higuchi add 20180125*/
	rp = lp.polyd - 2
	for i = 0; i < rp; i++ {

		Amx = lp.P[i+1].X - lp.P[0].X
		Amy = lp.P[i+1].Y - lp.P[0].Y
		Amz = lp.P[i+1].Z - lp.P[0].Z
		Am = math.Sqrt(Amx*Amx + Amy*Amy + Amz*Amz)
		Bmx = lp.P[i+2].X - lp.P[0].X
		Bmy = lp.P[i+2].Y - lp.P[0].Y
		Bmz = lp.P[i+2].Z - lp.P[0].Z
		Bm = math.Sqrt(Bmx*Bmx + Bmy*Bmy + Bmz*Bmz)
		//AM = Am * 100.0;
		//BM = Bm * 100.0;

		//dea = AM / DE;
		//deb = BM / DE;

		/* 20180125 higuchi add*/
		AMAX = int(math.Ceil(Am / DEM))
		BMAX = int(math.Ceil(Bm / DEM))

		/* 20170821 higuchi add*/
		//AMAX = (int)ceil(Am / dea);
		//BMAX = (int)ceil(Bm / deb);

		for loopA = 0; loopA < AMAX-1; loopA++ {
			for loopB = 0; loopB < BMAX-1; loopB++ {

				/*--higuchi upd 20180125--*/
				A = (DEM / 2.0) + DEM*float64(loopA)
				B = (DEM / 2.0) + DEM*float64(loopB)

				//for (A = (dea / 2); A < AM; A = A + dea){
				//for (B = (deb / 2); B < BM; B = B + deb){
				Px = 0.0
				Py = 0.0
				Pz = 0.0
				Qx = 0.0
				Qy = 0.0
				Qz = 0.0
				S = 0.0
				T = 0.0
				OX = 0.0
				OY = 0.0
				OZ = 0.0
				OX1 = 0.0
				OY1 = 0.0
				OZ1 = 0.0
				SHITA = 0.0

				wap = wap + 1.0

				OX = (A / Am) * Amx
				OY = (A / Am) * Amy
				OZ = (A / Am) * Amz
				OX1 = (B / Bm) * Bmx
				OY1 = (B / Bm) * Bmy
				OZ1 = (B / Bm) * Bmz

				//OX = (A / AM)*Amx;
				//OY = (A / AM)*Amy;
				//OZ = (A / AM)*Amz;
				//OX1 = (B / BM)*Bmx;
				//OY1 = (B / BM)*Bmy;
				//OZ1 = (B / BM)*Bmz;

				Qx = lp.P[0].X + OX + OX1
				Qy = lp.P[0].Y + OY + OY1
				Qz = lp.P[0].Z + OZ + OZ1

				/*---------LPによる影を考慮------------------*/
				for j = 0; j < lpn; j++ {
					omote = 0
					for k = 0; k < LP[j].polyd; k++ {
						if s.ps[j][k] > 0.0 {
							omote = 1
							break
						}
					}
					if (g != j) && (omote > 0) {
						KOUTEN(Qx, Qy, Qz, ls, ms, ns, &Px, &Py, &Pz, LP[j].P[0], LP[j].e)
						YOGEN(Qx, Qy, Qz, Px, Py, Pz, &SHITA, lp.e)
						rp2 = LP[j].polyd - 2
						for k = 0; k < rp2; k++ {
							INOROUT(Px, Py, Pz, LP[j].P[0], LP[j].P[k+1], LP[j].P[k+2], &S, &T)
							if ((S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0)) && (SHITA > 0) {
								lp.sum = lp.sum + 1.0
								goto koko
							}
						}
					}
				}

				/*-----------------MPによる影を考慮----------------------*/
				for h = 0; h < mpn; h++ {
					omote = 0
					for k = 0; k < MP[h].polyd; k++ {
						if t.ps[h][k] > 0.0 {
							omote = 1
							break
						}
					}
					if omote > 0 {
						KOUTEN(Qx, Qy, Qz, ls, ms, ns, &Px, &Py, &Pz, MP[h].P[0], MP[h].e)
						YOGEN(Qx, Qy, Qz, Px, Py, Pz, &SHITA, lp.e)
						rp2 = LP[h].polyd - 2
						for k = 0; k < rp2; k++ {
							INOROUT(Px, Py, Pz, LP[h].P[0], LP[h].P[k+1], LP[h].P[k+2], &S, &T)
							if ((S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0)) && (SHITA > 0) {
								lp.sum = lp.sum + 1.0
								goto koko
							}
						}
					}
				}
			koko:
			}
		}
	}
	lp.sum = lp.sum / wap
}
