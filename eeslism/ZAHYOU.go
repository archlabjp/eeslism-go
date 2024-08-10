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

   モンテカルロ法を用いる際の座標変換
              FILE=ZAHYOU.c
              Create Date=1999.6.7

*/

package eeslism

import "math"

// ZAHYOU converts coordinates  for Monte Carlo method.
// Op: original coordinates
// G: coordinates of the center of gravity
// op: converted coordinates [out]
// wa: angle of rotation around the Z-axis
// wb: angle of rotation around the X-axis
func ZAHYOU(Op, G XYZ, op *XYZ, wa, wb float64) {
	Cwa := math.Cos(wa * math.Pi / 180)
	Swa := math.Sin(wa * math.Pi / 180)
	Cwb := math.Cos((-wb) * math.Pi / 180)
	Swb := math.Sin((-wb) * math.Pi / 180)

	p := XYZ{
		X: Op.X - G.X,
		Y: Op.Y - G.Y,
		Z: Op.Z - G.Z,
	}

	q := XYZ{
		X: p.X*Cwa - p.Y*Swa,
		Y: p.X*Swa + p.Y*Cwa,
		Z: p.Z,
	}

	op.X = q.X
	op.Y = q.Y*Cwb - q.Z*Swb
	op.Z = q.Y*Swb + q.Z*Cwb

	CAT(&op.X, &op.Y, &op.Z)
}

/*------------------------------------------------------------------*/

func R_ZAHYOU(Op, G XYZ, op *XYZ, wa, wb float64) {
	Cwa := math.Cos((-wa) * math.Pi / 180)
	Swa := math.Sin((-wa) * math.Pi / 180)
	Cwb := math.Cos(wb * math.Pi / 180)
	Swb := math.Sin(wb * math.Pi / 180)

	p := XYZ{
		X: Op.X,
		Y: Op.Y*Cwb - Op.Z*Swb,
		Z: Op.Y*Swb + Op.Z*Cwb,
	}

	q := XYZ{
		X: p.X*Cwa - p.Y*Swa,
		Y: p.X*Swa + p.Y*Cwa,
		Z: p.Z,
	}

	op.X = q.X + G.X
	op.Y = q.Y + G.Y
	op.Z = q.Z + G.Z

	CAT(&op.X, &op.Y, &op.Z)
}
