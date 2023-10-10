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

   壁面の法線ベクトルを求める
                 FILE=EOP.c
                 Create Date 1999.6.15

*/

package eeslism

import "math"

func EOP(u int, p []*P_MENN) {
	const M_rad = math.Pi / 180
	for j := 0; j < u; j++ {
		p[j].e.Z = math.Cos(p[j].wb * M_rad)
		p[j].e.Y = -math.Sin(p[j].wb*M_rad) * math.Cos(p[j].wa*M_rad)
		p[j].e.X = -math.Sin(p[j].wb*M_rad) * math.Sin(p[j].wa*M_rad)
		CAT(&p[j].e.Z, &p[j].e.Y, &p[j].e.X)
	}
}
