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

/*    srprint.c            */

package eeslism

/*  時間別PMVの出力  */

import (
	"fmt"
	"io"
)

var __Pmvprint_count = 0

func Pmvprint(fpout io.Writer, title string, Room []*ROOM, Mon, Day int, time float64) {
	var Nr int
	if __Pmvprint_count == 0 && Room != nil {
		for i := range Room {
			Rm := Room[i]
			if Rm.Metsch != nil {
				Nr++
			}
		}

		fmt.Fprintf(fpout, "%s ;\n", title)
		fmt.Fprintf(fpout, "%d ", Nr)

		for i := range Room {
			Rm := Room[i]
			if Rm.Metsch != nil {
				fmt.Fprintf(fpout, "  %s ", Rm.Name)
			}
		}

		fmt.Fprintf(fpout, "\n")

		__Pmvprint_count = 1
	}

	fmt.Fprintf(fpout, "%02d %02d %5.2f ", Mon, Day, time)

	for i := range Room {
		Rm := Room[i]
		if Rm.Metsch != nil {
			fmt.Fprintf(fpout, " %4.3f ", Rm.PMV)
		}
	}
	fmt.Fprintf(fpout, "\n")
}

/* ----------------------------------------------------- */

/*   室内温・湿度、室内表面平均温度の出力
 */

var __Rmevprint_count = 0

func Rmevprint(fpout io.Writer, title string, Room []*ROOM, Mon, Day int, time float64) {
	if __Rmevprint_count == 0 {
		fmt.Fprintf(fpout, "%s ;\n", title)
		fmt.Fprintf(fpout, "%d室\t\t\t", len(Room))

		for i := range Room {
			Rm := Room[i]
			fmt.Fprintf(fpout, "%s\t\t\t\t", Rm.Name)
		}
		fmt.Fprintf(fpout, "\n")

		__Rmevprint_count = 1
	}
	/*======================================= */
	fmt.Fprintf(fpout, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := range Room {
		Rm := Room[i]
		fmt.Fprintf(fpout, "%.1f\t%.4f\t%.1f\t%.0f\t", Rm.Tr, Rm.xr, Rm.Tsav, Rm.RH)
	}

	fmt.Fprintf(fpout, "\n")
}

/* ----------------------------------------------------- */
