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

import (
	"fmt"
	"io"
)

/* 室供給熱量、放射パネルについての出力 */
var __Rmpnlprint_id int = 0

func Rmpnlprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64, Room []ROOM) {

	if __Rmpnlprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}
			rmqaprint(fo, __Rmpnlprint_id, Room)
			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}
			__Rmpnlprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	rmqaprint(fo, __Rmpnlprint_id, Room)
}
