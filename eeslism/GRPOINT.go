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

// GRGPOINT calculates the representative point of the front ground.
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
