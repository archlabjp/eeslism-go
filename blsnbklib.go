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

/*   snbklib.c  */

package main

import "math"

// 入力:
//
// 出力:
//   日よけの影面積率 Fsdw [-]
func FNFsdw(Ksdw, Ksi int, Xazm, Xprf, D, Wr, Hr, Wi1, Hi1, Wi2, Hi2 float64) float64 {
	var Da, Dp, Asdw, Fsdw float64

	Asdw = 0.0

	if Ksdw == 0 {
		return 0.0
	} else {
		Da = D * Xazm
		Dp = D * Xprf
		if Ksdw == 2 || Ksdw == 6 {
			Da = math.Abs(Da)
		}
		if Ksdw == 4 || Ksdw == 8 {
			Da = -Da
		}

		switch Ksdw {
		case 1:
			Asdw = FNAsdw1(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2)
		case 2, 3, 4:
			Asdw = FNAsdw1(Dp, Da, Hr, Wr, Hi1, Wi1, Hi2)
		case 5:
			Asdw = FNAsdw2(Dp, Hr, Wr, Hi1)
		case 6, 7, 8:
			Asdw = FNAsdw2(Da, Wr, Hr, Wi1)
		case 9:
			Asdw = FNAsdw3(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2, Hi2)
		}
		Fsdw = Asdw / (Wr * Hr)
	}

	if Ksi == 1 {
		Fsdw = 1.0 - Fsdw
	}

	return Fsdw
}

/*  -----------------------------------------------------  */

func FNAsdw1(Da, Dp, Wr, Hr, Wi1, Hi, Wi2 float64) float64 {
	var Wi, Daa, Dha, Dhb, Dwa, Dwb, Asdw float64

	if Dp <= 0.0 {
		return 0.0
	} else {
		Wi = Wi1
		if Da < 0.0 {
			Wi = Wi2
		}
		Daa = math.Abs(Da)
		Dha = Wi*Dp/math.Max(Wi, Daa) - Hi
		Dha = math.Min(math.Max(0.0, Dha), Hr)
		Dhb = (Wi+Wr)*Dp/math.Max(Wi+Wr, Daa) - Hi
		Dhb = math.Min(math.Max(0.0, Dhb), Hr)
		if Hi >= Dp {
			Dwa = 0.0
		} else {
			Dwa = (Wi + Wr) - Hi*Daa/Dp
			Dwa = math.Min(math.Max(0.0, Dwa), Wr)
		}
		Dwb = (Wi + Wr) - (Hi+Hr)*Daa/math.Max(Hi+Hr, Dp)
		Dwb = math.Min(math.Max(0.0, Dwb), Wr)
		Asdw = Dwa*Dha + 0.5*(Dwa+Dwb)*(Dhb-Dha)
	}

	return Asdw
}

/*  -----------------------------------------------------  */

func FNAsdw2(Dp, Hr, Wr, Hi float64) float64 {
	var Dh, Asdw float64

	if Dp <= 0.0 {
		return 0.0
	} else {
		Dh = math.Min(math.Max(0.0, Dp-Hi), Hr)
		Asdw = Wr * Dh
	}

	return Asdw
}

/*  -----------------------------------------------------  */

func FNAsdw3(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2, Hi2 float64) float64 {
	var Dw1, Dw2, Dh1, Dh2, Asdw float64

	Dw1 = math.Min(math.Max(0.0, Da-Wi1), Wr)
	Dw2 = math.Min(math.Max(0.0, -Da-Wi2), Wr)
	Dh1 = math.Min(math.Max(0.0, Dp-Hi1), Hr)
	Dh2 = math.Min(math.Max(0.0, -Dp-Hi2), Hr)
	Asdw = Wr*(Dh1+Dh2) + (Dw1+Dw2)*(Hr-Dh1-Dh2)

	return Asdw
}
