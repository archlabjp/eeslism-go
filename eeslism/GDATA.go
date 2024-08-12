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

	 壁面の中心点の座標を求める
					FILE=GDATA.c
					Create Date=1999.10.26

*/

package eeslism

// GDATA calculates the center of gravity of a polygon.
func GDATA(OP *P_MENN) XYZ {
	var x, y, z float64

	for i := range OP.P {
		x += OP.P[i].X
		y += OP.P[i].Y
		z += OP.P[i].Z
	}

	// the center of gravity of the polygon.
	d := float64(len(OP.P))
	return XYZ{
		X: x / d,
		Y: y / d,
		Z: z / d,
	}
}
