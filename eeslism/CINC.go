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
//						   壁面に対する入射日射角度
//						   FILE=CINC.c
//						   Create Date=1999.6.7
//						   */

func CINC(op *P_MENN, ls, ms, ns float64, co *float64) {
	Wz := math.Cos(op.wb * math.Pi / 180)
	Ww := -math.Sin(op.wb*math.Pi/180) * math.Sin(op.wa*math.Pi/180)
	Ws := -math.Sin(op.wb*math.Pi/180) * math.Cos(op.wa*math.Pi/180)

	*co = ns*Wz + ls*Ww + ms*Ws

	//fmt.Printf("op.wb=%f ns=%f ls=%f ms=%f co=%f\n", op.wb, ns, ls, ms, *co)
}
