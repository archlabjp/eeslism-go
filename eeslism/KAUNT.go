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

	 障害物に当たった点の合計をカウントする
					FILE=KAUNT.c
					Create Date=1999.6.7

*/

package eeslism

import "math"

func minval(span []float64, u int, min *int, val *float64) {
	*val = 10000.0
	*min = -1

	for i := 0; i < u; i++ {
		if span[i] > 0.0 && *val >= span[i] {
			*val = span[i]
			*min = i
		}
	}
}

func KAUNT(
	mlpn int,
	ls float64,
	ms float64,
	ns float64,
	suma *float64,
	sumg *float64,
	sumwall []float64,
	s float64,
	mlp []P_MENN,
	p []XYZ,
	O XYZ,
	E XYZ,
	wa float64,
	wb float64,
	G XYZ,
	gpn int,
	nday int,
	gcnt *int,
	startday int,
	wlflg int,
) {
	var rp, h int
	var l int
	var sdsum float64

	var i, j int
	U := 0.0
	Px := 0.0
	Py := 0.0
	Pz := 0.0
	var nai float64
	x := 0.0
	y := 0.0
	z := 0.0

	S := 0.0
	T := 0.0
	var span []float64

	var mencnt []int

	var dumy2 int
	var dumy1 float64
	var k int
	var mini float64
	//var minib int

	span = make([]float64, mlpn)
	mencnt = make([]int, mlpn)

	sdsum = 0.0
	i = 0

	for l = 0; l < mlpn; l++ {

		nai = ls*mlp[l].e.X + ms*mlp[l].e.Y + ns*mlp[l].e.Z

		if nai == 0.0 {
			span[l] = -1.0
		} else {
			KOUTEN(x, y, z, ls, ms, ns, &Px, &Py, &Pz, mlp[l].P[0], mlp[l].e)
			CAT(&ls, &ms, &ns)
			PRA(&U, ls, ms, ns, Px, Py, Pz)

			rp = mlp[l].polyd - 2

			/*--多角形ループ　三角形：１回、四角形：２回、、---*/
			for h = 0; h < rp; h++ {
				INOROUT(Px, Py, Pz, mlp[l].P[0], mlp[l].P[h+1], mlp[l].P[h+2], &S, &T)
				if ((S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0)) && (U > 0.0) {
					span[l] = math.Sqrt(Px*Px + Py*Py + Pz*Pz)
					break // 090131 higuchi debug
				} else {
					span[l] = -1.0
				}
			}

		}
	}

	/*--あたった障害物を近い順に並べ替える--*/
	/*--樋口　080920　追加--*/

	for i = 0; i < mlpn; i++ {
		mencnt[i] = i
	}

	for i = 0; i < mlpn-1; i++ {
		mini = span[i]
		//minib = mencnt[i]
		k = i
		for j = i + 1; j < mlpn; j++ {
			if span[j] < mini {
				mini = span[j]
				//minib = mencnt[j]
				k = j
			}
		}
		dumy1 = span[i]
		dumy2 = mencnt[i]
		span[i] = span[k]
		mencnt[i] = mencnt[k]
		span[k] = dumy1
		mencnt[k] = dumy2
	}

	k = 0
	for i = 0; i < mlpn; i++ {
		if span[i] > 0 {
			if k == 0.0 {
				sumwall[mencnt[i]] = sumwall[mencnt[i]] + mlp[mencnt[i]].shad[nday]
				sdsum = 1 - mlp[mencnt[i]].shad[nday] /*--透過分--*/
				k = 1
			} else {
				sumwall[mencnt[i]] = sumwall[mencnt[i]] + sdsum*mlp[mencnt[i]].shad[nday]
				sdsum = sdsum * (1 - mlp[mencnt[i]].shad[nday])
			}
		}
	}

	if k == 0 {
		//どの面にも当たらなかった場合
		if s > 0.0 {
			(*suma) = (*suma) + 1
		} else if s < 0.0 {
			(*sumg) = (*sumg) + 1
		} else {
			(*suma) = (*suma) + 1
			(*sumg) = (*sumg) + 1
		}
	} else {
		//どれかの面にあたった場合
		if s > 0.0 {
			(*suma) = (*suma) + sdsum
		} else {
			(*sumg) = (*sumg) + sdsum
		}
	}

	if (s < 0.0) && (nday == startday) {
		/*--始めの１回のみ地面のポイントを計算する--*/
		if *gcnt < gpn {
			DAINYUU_GP(&p[*gcnt], O, E, ls, ms, ns)
			R_ZAHYOU(p[*gcnt], G, &p[*gcnt], wa, wb)
			(*gcnt) = (*gcnt) + 1
		}
	}
}
