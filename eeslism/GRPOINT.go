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

///*
//
//                      前面地面の代表点を求める
//                      GDATA():MP面の中心点を求める
//                                        FILE=GRGPOINT.c
//                                        Create Date=1999.11.1
//
//*/

/*
GRGPOINT (Ground Representative Point Calculation)

この関数は、各主面（`mp`）の前面地面における代表点（`grp`）を計算します。
これは、地面からの反射日射や、地盤からの熱伝達を考慮する際に用いられる重要な情報です。

建築環境工学的な観点:
- **地面からの反射日射の考慮**: 建物の壁面や窓面は、
  地面からの反射日射を受けることがあります。
  この反射日射は、日射熱取得量を増加させ、
  特に冬季の暖房負荷軽減や、夏季の冷房負荷増加に影響を与えます。
  この関数は、各主面から地面への代表的な光線が当たる点を計算することで、
  地面からの反射日射の影響をモデル化します。
- **地盤からの熱伝達の考慮**: 地盤に接する壁面や床面は、
  地盤からの熱伝達を受けます。
  この関数で計算される代表点は、
  地盤からの熱伝達をモデル化する際の基準点として用いられる可能性があります。
- **重心の利用**: `GDATA(mp[i])`を呼び出して多角形の重心を計算し、
  それを基準として地面の代表点を計算します。
  これにより、多角形の形状を考慮した代表点の算出が可能になります。
- **法線ベクトルの考慮**: `ex`, `ey`, `ez`は、
  多角形の法線ベクトル成分であり、
  多角形が地面に対してどのような向きを向いているかを考慮します。
  `math.Abs(ez) < 1e-6` の条件は、
  多角形がほぼ水平である場合（法線ベクトルがZ軸に平行）に、
  代表点を重心とすることを意味します。

この関数は、建物の日射環境や地盤からの熱伝達を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func GRGPOINT(mp []*P_MENN, mpn int) {
	const M_rad = math.Pi / 180.0

	for i := 0; i < mpn; i++ {
		// Calculate the center of gravity of the polygon.
		mp[i].G = GDATA(mp[i])

		ex := mp[i].P[1].X - mp[i].P[0].X
		ey := mp[i].P[1].Y - mp[i].P[0].Y
		ez := mp[i].P[1].Z - mp[i].P[0].Z

		if math.Abs(ez) < 1e-6 {
			// If the normal vector is parallel to the Z-axis, the representative point is the center of gravity.
			mp[i].grp.X = 0.0
			mp[i].grp.Y = 0.0
			mp[i].grp.Z = 0.0
			continue
		} else {
			// Calculate the representative point of the front ground.
			t := -mp[i].G.Z / ez
			mp[i].grp.X = t*ex + mp[i].G.X - mp[i].grpx*math.Sin(mp[i].wa*M_rad)
			mp[i].grp.Y = t*ey + mp[i].G.Y - mp[i].grpx*math.Cos(mp[i].wa*M_rad)
			mp[i].grp.Z = 0.0
		}
	}
}
