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

   ベクトルの向きを判定する
              FILE=PRA.c
              Create Date=1998.8.15

*/

package eeslism

import (
	"fmt"
	"math"
	"os"
)

// PRA determines the direction of the vector.
// U: direction of the vector [out]
// ls, ms, ns: 太陽方位ベクトル
// x, y, z: coordinates of the vector
func PRA(U *float64, ls, ms, ns, x, y, z float64) {
	epsilon := 1.0e-6

	if math.Abs(ls) > epsilon {
		*U = x / ls
	} else if math.Abs(ms) > epsilon {
		*U = y / ms
	} else if math.Abs(ns) > epsilon {
		*U = z / ns
	} else {
		fmt.Printf("ls=%f ms=%f ns=%f\n", ls, ms, ns)
		fmt.Println("errorPRA")
		os.Exit(1)
	}
}
