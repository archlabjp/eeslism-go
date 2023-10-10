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

 点と面の交点が平面の内か外かを判別する
 FILE=INOROUT.c
 Create Date=1998.10.26

*/

package eeslism

import (
	"fmt"
	"math"
	"os"
)

func INOROUT(Px, Py, Pz float64, P0, P1, P2 XYZ, S, T *float64) {
	var Sx01, Sy01, Sz01, Tx03, Ty03, Tz03 float64
	var AA1, BB1, CC1, aa1, bb1, cc1 float64

	Sx01 = P1.X - P0.X
	Sy01 = P1.Y - P0.Y
	Sz01 = P1.Z - P0.Z
	Tx03 = P2.X - P0.X
	Ty03 = P2.Y - P0.Y
	Tz03 = P2.Z - P0.Z

	AA1 = Tx03*Sy01 - Ty03*Sx01
	BB1 = Ty03*Sz01 - Tz03*Sy01
	CC1 = Tx03*Sz01 - Tz03*Sx01

	CAT(&AA1, &BB1, &CC1)

	if math.Abs(AA1) > 0.0 {
		*T = (Sy01*Px - Sx01*Py - P0.X*Sy01 + P0.Y*Sx01) / (Tx03*Sy01 - Ty03*Sx01)
	} else if math.Abs(BB1) > 0.0 {
		*T = (Sz01*Py - Sy01*Pz - P0.Y*Sz01 + P0.Z*Sy01) / (Ty03*Sz01 - Tz03*Sy01)
	} else if math.Abs(CC1) > 0.0 {
		*T = (Sz01*Px - Sx01*Pz - P0.X*Sz01 + P0.Z*Sx01) / (Tx03*Sz01 - Tz03*Sx01)
	} else {
		fmt.Printf("error inorout\n0X=%f 0Y=%f 0Z=%f\n1X=%f 1Y=%f 1Z=%f\n2X=%f 2Y=%f 2Z=%f\n",
			P0.X, P0.Y, P0.Z, P1.X, P1.Y, P1.Z, P2.X, P2.Y, P2.Z)
		// Exit the program with an error code
	}

	aa1 = Sx01
	bb1 = Sy01
	cc1 = Sz01

	CAT(&aa1, &bb1, &cc1)

	if math.Abs(aa1) > 0.0 {
		*S = (Px - (*T)*Tx03 - P0.X) / Sx01
	} else if math.Abs(bb1) > 0.0 {
		*S = (Py - (*T)*Ty03 - P0.Y) / Sy01
	} else if math.Abs(cc1) > 0.0 {
		*S = (Pz - (*T)*Tz03 - P0.Z) / Sz01
	} else {
		fmt.Println("error inorout2")
		os.Exit(1)
	}
}
