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

   OP,LPの座標出力（デバッグ用）
              FILE=ZPRINT.c
              Create Date=1999.10.26

*/

package eeslism

import (
	"fmt"
	"os"
)

func ZPRINT(lp []P_MENN, op []P_MENN, lpn, opn int, name string) {
	name += ".gchi"
	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open")
		os.Exit(1)
	}

	for i := 0; i < opn; i++ {
		fmt.Fprintf(fp, "op[%d] %s\n", i, op[i].opname)
		for j := 0; j < op[i].polyd; j++ {
			fmt.Fprintf(fp, "    P[%d] X=%f Y=%f Z=%f\n", j, op[i].P[j].X, op[i].P[j].Y, op[i].P[j].Z)
		}
		for k := 0; k < op[i].wd; k++ {
			fmt.Fprintf(fp, "op[%d] opw[%d] %s\n ", i, k, op[i].opw[k].opwname)
			for j := 0; j < 4; j++ {
				fmt.Fprintf(fp, "   P[%d] X=%f Y=%f Z=%f\n", j, op[i].opw[k].P[j].X, op[i].opw[k].P[j].Y, op[i].opw[k].P[j].Z)
			}
		}
		fmt.Fprintln(fp)
	}

	for i := 0; i < lpn; i++ {
		fmt.Fprintf(fp, "lp[%d] %s\n ", i, lp[i].opname)
		fmt.Fprintf(fp, "      e.X=%f e.Y=%f e.Z=%f\n", lp[i].e.X, lp[i].e.Y, lp[i].e.Z)
		for j := 0; j < lp[i].polyd; j++ {
			fmt.Fprintf(fp, "   P[%d] X=%f Y=%f Z=%f\n", j, lp[i].P[j].X, lp[i].P[j].Y, lp[i].P[j].Z)
		}
		fmt.Fprintln(fp)
	}

	fp.Close()
}
