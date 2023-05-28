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

/****************

 熱損失係数出力ルーチン
 1997.11.18
	FILE=bdhpri.c

****************/

package main

import (
	"fmt"
	"os"
)

func bdhpri(ofile string, rmvls RMVLS, exs *EXSFS) {
	Nroom := rmvls.Nroom
	e := exs.Exs

	file := ofile + "_bdh.es"
	fp, err := os.Create(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fp.Close()

	fmt.Fprintf(fp, "%d\n", Nroom)

	for i := 0; i < Nroom; i++ {
		room := &rmvls.Room[i]
		fmt.Fprintf(fp, "%s\t", room.Name)
		fmt.Fprintf(fp, "%d\t", room.N)
		fmt.Fprintf(fp, "%.3f\t", room.FArea)
		fmt.Fprintf(fp, "%.3f\t#\n", room.VRM)

		for j := 0; j < room.N; j++ {
			r := room.rsrf[j]
			if r.ble == 'E' || r.ble == 'W' || r.ble == 'R' || r.ble == 'F' {
				er := e[r.exs]
				fmt.Fprintf(fp, "\t%s", er.Name)
			} else if r.nextroom != nil && r.nextroom.Name != "" {
				fmt.Fprintf(fp, "\t%s", r.nextroom.Name)
			} else {
				fmt.Fprintf(fp, "\t-")
			}

			if r.Name != "" && len(r.Name) > 0 {
				fmt.Fprintf(fp, "\t%s\t", r.Name)
			} else {
				fmt.Fprintf(fp, "\t-\t")
			}

			fmt.Fprintf(fp, "%c\t", r.ble)
			fmt.Fprintf(fp, "%.3f\t", r.A)

			if r.A < 0.0 {
				E := fmt.Sprintf("RmName=%s Ble=%c A=%.3f\n", room.Name, r.ble, r.A)
				Errprint(1, "<bdhpri>", E)
			}

			if r.Rwall >= 0.0 {
				fmt.Fprintf(fp, "%.3f\t%.2f\t;\n", r.Rwall, r.CAPwall)
			} else {
				//E := fmt.Sprintf("Rmname=%s Ble=%c  Not Defined", room.Name, r.ble)
				//Errprint(1, "<bdhpri>", E)
				fmt.Fprintf(fp, "\t\t;\n")
			}
		}
		fmt.Fprintf(fp, "*\n")
	}
}
