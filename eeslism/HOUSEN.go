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
//                         法線ベクトルを求める
//                                        FILE=HOUSEN.c
//                                        Create Date=1998.10.26
//                                        Update 2007.10.11 higuchi
//
//*/

// HOUSEN calculates the normal vector of a polygon.
// LP[i].e に法線ベクトルが入る
func HOUSEN(LP []*P_MENN) {
	for _, _lp := range LP {
		// 多角形のうち2辺のベクトルを求める
		x := _lp.P[1].X - _lp.P[0].X
		y := _lp.P[1].Y - _lp.P[0].Y
		z := _lp.P[1].Z - _lp.P[0].Z
		x1 := _lp.P[2].X - _lp.P[0].X
		y1 := _lp.P[2].Y - _lp.P[0].Y
		z1 := _lp.P[2].Z - _lp.P[0].Z

		// 法線ベクトルを求める
		_lp.e.X = y*z1 - z*y1
		_lp.e.Y = z*x1 - x*z1
		_lp.e.Z = x*y1 - y*x1

		// 法線ベクトルの正規化
		el := math.Sqrt(_lp.e.X*_lp.e.X + _lp.e.Y*_lp.e.Y + _lp.e.Z*_lp.e.Z)
		_lp.e.X = _lp.e.X / el
		_lp.e.Y = _lp.e.Y / el
		_lp.e.Z = _lp.e.Z / el
	}
}

func HOUSEN2(p0, p1, p2, e *XYZ) {
	x := p1.X - p0.X
	y := p1.Y - p0.Y
	z := p1.Z - p0.Z
	x1 := p2.X - p0.X
	y1 := p2.Y - p0.Y
	z1 := p2.Z - p0.Z

	e.X = y*z1 - z*y1
	e.Y = z*x1 - x*z1
	e.Z = x*y1 - y*x1

	el := math.Sqrt(e.X*e.X + e.Y*e.Y + e.Z*e.Z)
	e.X = e.X / el
	e.Y = e.Y / el
	e.Z = e.Z / el
}
