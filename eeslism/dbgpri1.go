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
)

/* -------------------------------------------- */

func xprtwallinit(Nmwall int, M []*MWALL) {
	Max := 0
	for j := 0; j < Nmwall; j++ {
		if M[j].M > Max {
			Max = M[j].M
		}
	}

	if DEBUG {
		fmt.Println("--- xprtwallinit")
		for j := 0; j < Nmwall; j++ {
			fmt.Printf("Told  j=%2d", j)
			for m := 0; m < M[j].M; m++ {
				fmt.Printf("  %2d%5.1f", m, M[j].Told[m])
			}
			fmt.Println()
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprtwallinit")
		fmt.Fprint(Ferr, "\tNo.")
		for j := 0; j < Max; j++ {
			fmt.Fprintf(Ferr, "\tT[%d]", j)
		}
		fmt.Fprintln(Ferr)
		for j := 0; j < Nmwall; j++ {
			fmt.Fprintf(Ferr, "\t%d", j)
			for m := 0; m < M[j].M; m++ {
				fmt.Fprintf(Ferr, "\t%.3g", M[j].Told[m])
			}
			fmt.Fprintln(Ferr)
		}
	}
}

/* -------------------------------------------- */

func xprsolrd(Nexs int, E []EXSF) {
	if DEBUG {
		fmt.Println("--- xprsolrd")
		for i := 0; i < Nexs; i++ {
			Exs := &E[i]
			fmt.Printf("EXSF[%2d]=%s  Id=%.0f  Idif=%.0f  Iw=%.0f RN=%.0f cinc=%.3f\n",
				i, Exs.Name, Exs.Idre, Exs.Idf, Exs.Iw, Exs.Rn, Exs.Cinc)
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprsolrd")
		fmt.Fprintln(Ferr, "\tNo.\tName\tId\tIdif\tIw\tRN\tcinc")
		for i := 0; i < Nexs; i++ {
			Exs := &E[i]
			fmt.Fprintf(Ferr, "\t%d\t%s\t%.0f\t%.0f\t%.0f\t%.0f\t%.3f\n",
				i, Exs.Name, Exs.Idre, Exs.Idf, Exs.Iw, Exs.Rn, Exs.Cinc)
		}
	}
}

/* ---------------------------------------------------------- */

func xpralph(Nroom int, _Room []ROOM, S []RMSRF) {
	fmt.Println("--- xpralph")

	for i := 0; i < Nroom; i++ {
		Room := &_Room[i]
		N := Room.N
		brs := Room.Brs

		fmt.Println(" alr(i,j)")
		Matfprint("  %5.1f", N, Room.alr)
		fmt.Println(" alph")

		for n := brs; n < brs+N; n++ {
			Sd := S[n]
			fmt.Printf("  %3d  alo=%5.1f  alir=%5.1f alic=%5.1f  ali=%5.1f\n",
				n, Sd.alo, Sd.alir, Sd.alic, Sd.ali)
		}
	}
}

/* ---------------------------------------------------------- */

func xprxas(Nroom int, R []ROOM, S []RMSRF) {
	if DEBUG {
		fmt.Printf("--- xprxas\n")

		for i := 0; i < Nroom; i++ {
			Room := &R[i]
			N := Room.N
			brs := Room.Brs

			fmt.Printf(" XA(i,j)\n")
			Matprint("%7.4f", N, Room.XA)

			for n := brs; n < brs+N; n++ {
				Sd := &S[n]
				fmt.Printf("%2d  K=%f  alo=%f  FI=%f FO=%f FP=%f  CF=%f\n",
					n, Sd.K, Sd.alo, Sd.FI, Sd.FO, Sd.FP, Sd.CF)
				fmt.Printf("            WSR=%f", Sd.WSR)

				for j := 0; j < Room.Ntr; j++ {
					fmt.Printf(" WSRN=%f", Sd.WSRN[j])
				}

				for j := 0; j < Room.Nrp; j++ {
					fmt.Printf(" WSPL=%f", Sd.WSPL[j])
				}

				fmt.Printf("   WSC=%f\n", Sd.WSC)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "--- xprxas\n")

		for i := 0; i < Nroom; i++ {
			Room := &R[i]
			N := Room.N
			brs := Room.Brs

			fmt.Fprintf(Ferr, "Room=%s\tXA(i,j)\n", Room.Name)
			Matfiprint(Ferr, "\t%.1g", N, Room.XA)

			for n := brs; n < brs+N; n++ {

				Sd := &S[n]
				fmt.Fprintf(Ferr, "\n\n\t%d\tK=%.2g\talo=%.2g\tFI=%.2g\tFO=%.2g\tFP=%.2g\tCF=%.2g\t",
					n, Sd.K, Sd.alo, Sd.FI, Sd.FO, Sd.FP, Sd.CF)
				fmt.Fprintf(Ferr, "\t\tWSR=%.3g\n\t", Sd.WSR)

				for j := 0; j < Room.Ntr; j++ {
					fmt.Fprintf(Ferr, "\tWSRN[%d]=%.3g", j, Sd.WSRN[j])
				}
				fmt.Fprintf(Ferr, "\n\t")

				for j := 0; j < Room.Nrp; j++ {
					fmt.Fprintf(Ferr, "\tWSPL[%d]=%.3g", j, Sd.WSPL[j])
				}
				fmt.Fprintf(Ferr, "\n")

				fmt.Fprintf(Ferr, "\t\tWSC=%.3g\n", Sd.WSC)
			}
		}
	}
}

func xprtwsrf(N int, _Sd []RMSRF) {
	fmt.Println("--- xprtwsrf")

	for n := 0; n < N; n++ {
		Sd := &_Sd[n]
		fmt.Printf("  n=%2d  rm=%d nr=%d  Ts=%6.2f  Tmrt=%6.2f  Te=%6.2f  RS=%7.1f\n",
			n, Sd.rm, Sd.n, Sd.Ts, Sd.Tmrt, Sd.Te, Sd.RS)
	}
}

/* -------------------------------------------------------------------- */

func xprrmsrf(N int, _Sd []RMSRF) {
	fmt.Println("--- xprrmsf")

	for n := 0; n < N; n++ {
		Sd := &_Sd[n]
		fmt.Printf("  [%d]=%6.2f", n, Sd.Ts)
	}
	fmt.Println()
}

/* -------------------------------------------------------------------- */

func xprtwall(Nmwall int, _Mw []MWALL) {
	fmt.Println("--- xprtwall")

	for j := 0; j < Nmwall; j++ {
		Mw := &_Mw[j]
		if Mw.Pc > 0 {
			fmt.Printf("Tw j=%2d", j)

			for m := 0; m < Mw.M; m++ {
				Tw := Mw.Tw[m]
				fmt.Printf("  [%d]=%6.2f", m, Tw)
			}

			fmt.Println()
		}
	}
}
