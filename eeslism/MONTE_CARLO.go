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

            形態係数を求める
                       FILE=MONTE_CARLO.c
                       Create Date=1999.6.7
   061110 FFACTOR_LP$B$r=$@5(B  $B;M3Q7A$7$+G'<1$7$F$$$J$+$C$?(B

*/

package eeslism

import (
	"fmt"
	"math"
)

func MONTE_CARLO(
	mpn int,
	lpn int,
	NUM int,
	MP []*P_MENN,
	LP []*P_MENN,
	GP [][]XYZ,
	gpn int,
	nday int,
	startday int,
) {

	var j, h, l, n, i, k, mlpn int
	var suma, sumg float64
	var sumwall []float64
	var G, O, OO, E XYZ
	var mlp []*P_MENN
	var a, va, FF float64
	var gcnt int /*--地面の代表点カウンター--*/

	O.X = 0.0
	O.Y = 0.0
	O.Z = 0.0

	/*--------メモリーの確保-----------------*/
	mlpn = mpn + lpn
	sumwall = make([]float64, mlpn)
	mlp = make([]*P_MENN, mlpn)

	suma = 0
	sumg = 0
	gcnt = 0
	for i = 0; i < mlpn; i++ {
		sumwall[i] = 0.0
	}

	for j = 0; j < mpn; j++ {
		/*-----op面の中心点の座標を求める--------*/
		GDATA(MP[j], &G)
		//printf("MP[%d].opname=%s wlflg=%d\n",j,MP[j].opname,MP[j].wlflg) ;
		/*-----------初期化---------------*/
		for i = 0; i < mlpn; i++ {
			mlp[i] = new(P_MENN)
			mlp[i].e.X = 0.0
			mlp[i].e.Y = 0.0
			mlp[i].e.Z = 0.0
		}
		//MATINIT(mlp,mlpn) ;

		/*---------mlp面の座標変換---------*/
		for h = 0; h < mlpn; h++ {

			if h < mpn {

				mlp[h].sbflg = MP[h].sbflg /*--higuchi add 080915 */

				mlp[h].polyd = MP[h].polyd
				mlp[h].P = make([]XYZ, mlp[h].polyd)
				for i = 0; i < mlp[h].polyd; i++ {
					mlp[h].P[i].X = 0.0
					mlp[h].P[i].Y = 0.0
					mlp[h].P[i].Z = 0.0
				}

				for k = 1; k < 366; k++ {
					mlp[h].shad[k] = 1.0
				}

				for l = 0; l < mlp[h].polyd; l++ {
					ZAHYOU(MP[h].P[l], G, &mlp[h].P[l], MP[j].wa, MP[j].wb)
				}
			} else {
				mlp[h].polyd = LP[h-mpn].polyd

				mlp[h].sbflg = LP[h-mpn].sbflg /*--higuchi add  0809015  --*/

				mlp[h].P = make([]XYZ, mlp[h].polyd)

				for i = 0; i < mlp[h].polyd; i++ {
					mlp[h].P[i].X = 0.0
					mlp[h].P[i].Y = 0.0
					mlp[h].P[i].Z = 0.0
				}

				for k = 1; k < 366; k++ {
					mlp[h].shad[k] = LP[h-mpn].shad[k]
				}

				for l = 0; l < mlp[h].polyd; l++ {
					ZAHYOU(LP[h-mpn].P[l], G, &mlp[h].P[l], MP[j].wa, MP[j].wb)
				}
			}
		}

		const M_rad = math.Pi / 180.0

		ZAHYOU(O, G, &OO, MP[j].wa, MP[j].wb)
		E.X = 0.0
		E.Y = -math.Sin((-MP[j].wb) * M_rad)
		E.Z = math.Cos((-MP[j].wb) * M_rad)

		/*-------op面の法線ベクトルを求める----------*/
		HOUSEN(mlpn, mlp)

		/*-------------点を射出する-------------------*/
		for n = 0; n < NUM; n++ {
			var ls, ms, ns, s float64

			/*----------乱数の発生--------------*/
			RAND(&a, &va)

			ls = math.Sin(va) * math.Cos(a)
			ms = math.Sin(va) * math.Sin(a)
			ns = math.Cos(va)

			URA_M(ls, ms, ns, &s, MP[j].wb)
			KAUNT(mlpn, ls, ms, ns, &suma, &sumg, sumwall, s,
				mlp, GP[j], OO, E, MP[j].wa, MP[j].wb, G, gpn,
				nday, &gcnt, startday, MP[j].wlflg)
		}

		if nday == startday {
			if gcnt >= gpn {
				GP[j][gpn].X = -999
				GP[j][gpn].Y = -999
				GP[j][gpn].Z = -999
			} else {
				GP[j][gcnt].X = -999
				GP[j][gcnt].Y = -999
				GP[j][gcnt].Z = -999
			}
		}

		FF = 0.0
		for i = 0; i < mlpn; i++ {
			MP[j].faiwall[i] = sumwall[i] / float64(NUM)
			FF = FF + MP[j].faiwall[i]
		}
		MP[j].faig = sumg / float64(NUM)
		MP[j].faia = suma / float64(NUM)

		/*--*/
		fmt.Printf("%s faia=%f faig=%f faib=%f\n", MP[j].opname, MP[j].faia, MP[j].faig, FF)
		/*--*/

		FF = FF + MP[j].faia + MP[j].faig

		for i = 0; i < mlpn; i++ {
			sumwall[i] = 0
		}

		suma = 0
		sumg = 0
		gcnt = 0
	}
}

/*-----------地面から見た天空に対する形態係数-----------------------*/
func GR_MONTE_CARLO(mp []*P_MENN, mpn int, lp []*P_MENN, lpn int, monten int, day int) {
	var rp int
	var i, n, l, mlpn, k, h int
	var ls, ms, ns float64
	var a, va float64
	var Px, Py, Pz, U, PX, PY, PZ float64
	var S, T float64
	var shad, sum float64
	var mlp []*P_MENN

	shad = 0.0
	sum = 0.0
	mlpn = mpn + lpn
	mlp = make([]*P_MENN, mlpn)

	/*---MPとLPの結合--*/
	for i = 0; i < mlpn; i++ {
		if i < mpn {
			mlp[i] = mp[i]
			for k = 1; k < 366; k++ {
				mlp[i].shad[k] = 1.0
			}
		} else {
			mlp[i] = lp[i-mpn]
			for k = 1; k < 366; k++ {
				mlp[i].shad[k] = lp[i-mpn].shad[k] //LPの遮蔽率を代入
			}
		}
	}

	for i = 0; i < mpn; i++ {
		for n = 0; n < monten; n++ {
			// ランダムな太陽位置
			// a: 太陽方位角
			// va: 太陽高度
			RAND(&a, &va)

			// ls: 3D空間におけるX軸方向の成分（東西方向）
			// ms: 3D空間におけるY軸方向の成分（南北方向）
			// ns: 3D空間におけるZ軸方向の成分（垂直方向）
			ls = math.Sin(va) * math.Cos(a)
			ms = math.Sin(va) * math.Sin(a)
			ns = math.Cos(va)

			for l = 0; l < mlpn; l++ {
				// 前面地面代表点から太陽方向のベクトルとXXの交点 (Px,Py,Pz)を求める
				KOUTEN(
					mp[i].grp.X, mp[i].grp.Y, mp[i].grp.Z,
					ls, ms, ns,
					&Px, &Py, &Pz,
					mlp[l].P[0], mlp[l].e)

				// 交点(Px,Py,Pz)の移動してベクトル(PX,PY,PZ)を求める
				PX = Px - mp[i].grp.X
				PY = Py - mp[i].grp.Y
				PZ = Pz - mp[i].grp.Z

				// ベクトル(PX,PY,PZ)の
				PRA(&U, ls, ms, ns, PX, PY, PZ)

				rp = mlp[l].polyd - 2
				/*--多角形ループ　三角形：１回、四角形：２回、、、---*/
				for h = 0; h < rp; h++ {
					INOROUT(Px, Py, Pz, mlp[l].P[0], mlp[l].P[h+1], mlp[l].P[h+2], &S, &T)

					if ((S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0)) && (U > 0.0) {
						if mlp[l].shad[day] > 0.0 {
							if shad == 0.0 {
								shad = 1 - mlp[l].shad[day]
							} else {
								shad = shad * (1 - mlp[l].shad[day])
								break
							}
						}
					}
				}

			}
			sum = sum + shad
		}
		mp[i].grpfaia = sum / float64(monten)
		sum = 0
	}
}

/*---------障害物LPから見た天空に対する形態係数-----------*/
func FFACTOR_LP(lpn int, mpn int, NUM int, LP []*P_MENN, MP []*P_MENN) {

	var flg int
	var rp int /*--多角形のループ回数---*/
	var j, h, l, n, i, mlpn, sum int
	var G XYZ // 中心座標
	var mlp []*P_MENN
	var a, va, x, y, z, Px, Py, Pz, U float64
	var S, T float64

	/*---------メモリーの確保-----------------*/
	mlpn = lpn + mpn
	mlp = make([]*P_MENN, mlpn)

	sum = 0

	for j = 0; j < lpn; j++ {
		/*-----lp面の中心点の座標を求める--------*/
		GDATA(LP[j], &G)

		/*-------初期化---------------*/
		for i = 0; i < mlpn; i++ {
			mlp[i] = new(P_MENN)
			mlp[i].e.X = 0.0
			mlp[i].e.Y = 0.0
			mlp[i].e.Z = 0.0
		}
		//MATINIT(mlp,mlpn) ;

		/*---------lp面の座標変換---------*/
		for h = 0; h < mlpn; h++ {
			if h < lpn {

				mlp[h].polyd = LP[h].polyd
				mlp[h].P = make([]XYZ, mlp[h].polyd)
				for i = 0; i < mlp[h].polyd; i++ {
					mlp[h].P[i].X = 0.0
					mlp[h].P[i].Y = 0.0
					mlp[h].P[i].Z = 0.0
				}

				for l = 0; l < mlp[h].polyd; l++ {
					ZAHYOU(LP[h].P[l], G, &mlp[h].P[l], LP[j].wa, LP[j].wb)
				}
			} else {
				mlp[h].polyd = MP[h-lpn].polyd
				mlp[h].P = make([]XYZ, mlp[h].polyd)
				for i = 0; i < mlp[h].polyd; i++ {
					mlp[h].P[i].X = 0.0
					mlp[h].P[i].Y = 0.0
					mlp[h].P[i].Z = 0.0
				}

				for l = 0; l < mlp[h].polyd; l++ {
					ZAHYOU(MP[h-lpn].P[l], G, &mlp[h].P[l], LP[j].wa, LP[j].wb)
				}
			}
		}

		/*-------mlp面の法線ベクトルを求める----------*/
		HOUSEN(mlpn, mlp)

		/*---------点を射出する----------------*/
		for n = 0; n < NUM; n++ {
			var ls, ms, ns, s float64

			/*----------乱数の発生--------------*/
			RAND(&a, &va)

			ls = math.Sin(va) * math.Cos(a)
			ms = math.Sin(va) * math.Sin(a)
			ns = math.Cos(va)

			URA_M(ls, ms, ns, &s, LP[j].wb)

			if s < 0.0 {
				sum = sum + 1
				continue
			}
			flg = 0 /*--内側フラグ初期化--*/
			for i = 0; i < mlpn; i++ {
				KOUTEN(x, y, z, ls, ms, ns, &Px, &Py, &Pz, mlp[i].P[0], mlp[i].e)
				PRA(&U, ls, ms, ns, Px, Py, Pz)
				rp = mlp[i].polyd - 2

				for h = 0; h < rp; h++ {
					INOROUT(Px, Py, Pz, mlp[i].P[0], mlp[i].P[h+1], mlp[i].P[h+2], &S, &T)
					if ((S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0)) && (U > 0.0) {
						sum = sum + 1
						flg = 1
						break
					}
				}
				if flg > 0 {
					break
				}
			}
		}

		LP[j].faia = float64(NUM-sum) / float64(NUM)
	}

	//printf("sum=%d NUM=%d LP[%d].faia=%f\n", sum,NUM,j, LP[j].faia);
	sum = 0
}
