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

   乱数の発生
      FILE=RAND.c
      Create Date=1999.6.7

*/

package eeslism

import (
	"math"
	"math/rand"
)

func RAND(a, v *float64) {
	const RAND_MAX = 2147483647
	*a = 2.0 * math.Pi * (float64(rand.Intn(RAND_MAX)) / float64(RAND_MAX))
	*v = math.Acos(math.Sqrt(1.0 - (float64(rand.Intn(RAND_MAX)) / float64(RAND_MAX))))
}
