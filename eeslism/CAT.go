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

	  小数点の切捨て
		   FILE=CAT.c
		   Create Date=1999.6.7
*/

package eeslism

import "math"

func CAT(a, b, c *float64) {
	*a = math.Floor((*a)*10000.0 + 0.5)
	*b = math.Floor((*b)*10000.0 + 0.5)
	*c = math.Floor((*c)*10000.0 + 0.5)

	*a = (*a) / 10000.0
	*b = (*b) / 10000.0
	*c = (*c) / 10000.0

	if math.Signbit(*a) {
		*a = 0.0
	}
	if math.Signbit(*b) {
		*b = 0.0
	}
	if math.Signbit(*c) {
		*c = 0.0
	}
}
