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

/*  dbgpri2.c   */

package main

import "fmt"

/* ----------------------------------------- */
func xprroom(Nroom int, R []*ROOM) {
	var i, j int
	var ARN []float64
	var RMP []float64
	var Room *ROOM

	Room = R[0]
	if DEBUG {
		fmt.Println("--- xprroom")
		for i = 0; i < Nroom; i++ {
			Room = R[i]
			fmt.Printf(" Room:  name=%s  MRM=%f  GRM=%f\n", Room.Name, Room.MRM, Room.GRM)
			fmt.Printf("     RMt=%f", Room.RMt)

			ARN = Room.ARN
			for j = 0; j < Room.Ntr; j++ {
				fmt.Printf(" ARN=%f", ARN[j])
			}

			RMP = Room.RMP
			for j = 0; j < Room.Nrp; j++ {
				fmt.Printf(" RMP=%f", RMP[j])
			}

			fmt.Printf(" RMC=%f\n", Room.RMC)
			fmt.Printf("     RMx=%f          RMXC=%f\n", Room.RMx, Room.RMXC)
		}
	}

	Room = R[0]
	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprroom")
		for i = 0; i < Nroom; i++ {
			Room = R[i]
			fmt.Fprintf(Ferr, "Room:\tname=%s\tMRM=%.4g\tGRM=%.4g\n", Room.Name, Room.MRM, Room.GRM)
			fmt.Fprintf(Ferr, "\tRMt=%.4g\n", Room.RMt)

			ARN = Room.ARN
			for j = 0; j < Room.Ntr; j++ {
				fmt.Fprintf(Ferr, "\tARN[%d]=%.4g", j, ARN[j])
			}
			fmt.Fprintln(Ferr)

			RMP = Room.RMP
			for j = 0; j < Room.Nrp; j++ {
				fmt.Fprintf(Ferr, "\tRMP[%d]=%.4g", j, RMP[j])
			}
			fmt.Fprintln(Ferr)

			fmt.Fprintf(Ferr, "\tRMC=%.4g\n", Room.RMC)
			fmt.Fprintf(Ferr, "\tRMx=%.2g\t\tRMXC=%.2g\n", Room.RMx, Room.RMXC)
		}
	}
}

/* ----------------------------------------- */

func xprschval(Nsch int, val []float64, Nscw int, isw []rune) {
	var j int

	fmt.Println("--- xprschval")

	for j = 0; j < Nsch; j++ {
		fmt.Printf("--- val=(%d) %f\n", j, val[j])
	}

	for j = 0; j < Nscw; j++ {
		fmt.Printf("--- isw=(%d) %c\n", j, isw[j])
	}
}

/* --------------------------------------------- */

func xprqin(Nroom int, Room []ROOM) {
	var i int

	fmt.Printf("--- xprqin  Nroom=%d\n", Nroom)

	for i = 0; i < Nroom; i++ {
		r := &Room[i]
		fmt.Printf("[%d] Hc=%f Hr=%f HL=%f Lc=%f Lr=%f Ac=%f Ar=%f AL=%f\n",
			i, r.Hc, r.Hr, r.HL, r.Lc, r.Lr, r.Ac, r.Ar, r.AL)
	}
}

/* --------------------------------------------- */

func xprvent(Nroom int, R []ROOM) {
	var i, j int
	var A *ACHIR
	var Room *ROOM

	if DEBUG {
		fmt.Println("--- xprvent")

		for i = 0; i < Nroom; i++ {
			Room = &R[i]
			fmt.Printf("[%d] %-10s  Gvent=%f  -- Gvr:", i, Room.Name, Room.Gvent)

			for j = 0; j < Room.Nachr; j++ {
				A = &Room.achr[j]
				fmt.Printf(" <%d>=%f", A.rm, A.Gvr)
			}
			fmt.Println()
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n\n--- xprvent")

		for i = 0; i < Nroom; i++ {
			Room = &R[i]
			fmt.Fprintf(Ferr, "\t[%d]\t%s\tGvent=%.3g\n\t\t", i, Room.Name, Room.Gvent)

			for j = 0; j < Room.Nachr; j++ {
				A = &Room.achr[j]
				fmt.Fprintf(Ferr, "\t<%d>=%.2g", A.rm, A.Gvr)
			}
			fmt.Fprintln(Ferr)
		}
	}
}
