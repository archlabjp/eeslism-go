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
	"os"
)

/* -------------------------------- */

func dprdayweek(daywk []int) {
	const dmax = 366

	fmt.Println("---  Day of week -----")
	for d := 0; d < 8; d++ {
		fmt.Printf("  %s=%d  ", DAYweek[d], d)
	}
	fmt.Println()

	k := 1
	for d := 1; d < dmax; d++ {
		if FNNday(k, 1) == d {
			fmt.Printf("\n%2d - ", k)
			k++
		}
		fmt.Printf("%2d", daywk[d-1])
	}
	fmt.Println()
}

/* ----------------------------------------------------------------- */

func dprschtable(Ssn []SEASN, Wkd []WKDY, Dh []DSCH, Dw []DSCW) {

	Ns := len(Ssn)
	Nw := len(Wkd)
	Nsc := len(Dh)
	Nsw := len(Dw)

	if DEBUG {
		fmt.Printf("\n*** dprschtable  ***\n")
		fmt.Printf("\n=== Schtable end  is=%d  iw=%d  sc=%d  sw=%d\n", Ns, Nw, Nsc, Nsw)

		for is := 0; is < Ns; is++ {
			Seasn := &Ssn[is]
			fmt.Printf("\n- %s", Seasn.name)

			for js := 0; js < Seasn.N; js++ {
				sday := &Seasn.sday[js]
				eday := &Seasn.eday[js]
				fmt.Printf("  %4d-%4d", *sday, *eday)
			}
		}

		for iw := 0; iw < Nw; iw++ {
			Wkdy := &Wkd[iw]
			fmt.Printf("\n- %s", Wkdy.name)

			for j := 0; j < 8; j++ {
				wday := &Wkdy.wday[j]
				fmt.Printf("   %d", *wday)
			}
		}

		for sc := 0; sc < Nsc; sc++ {
			Dsch := &Dh[sc]
			fmt.Printf("\n-VL   %10s (%2d) ", Dsch.name, sc)

			for jsc := 0; jsc < Dsch.N; jsc++ {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Printf("  %4d-(%6.2f)-%4d", stime, val, etime)
			}
		}

		for sw := 0; sw < Nsw; sw++ {
			Dscw := &Dw[sw]
			fmt.Printf("\n-SW   %10s (%2d) ", Dscw.name, sw)

			for jsw := 0; jsw < Dscw.N; jsw++ {
				stime := Dscw.stime[jsw]
				mode := Dscw.mode[jsw]
				etime := Dscw.etime[jsw]
				fmt.Printf("  %4d-( %c)-%4d", stime, mode, etime)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprschtable  ***\n")
		fmt.Fprintf(Ferr, "\n=== Schtable end  is=%d  iw=%d  sc=%d  sw=%d\n", Ns, Nw, Nsc, Nsw)

		for is := 0; is < Ns; is++ {
			Seasn := &Ssn[is]
			fmt.Fprintf(Ferr, "\n\t%s", Seasn.name)

			for js := 0; js < Seasn.N; js++ {
				sday := Seasn.sday[js]
				eday := Seasn.eday[js]
				fmt.Fprintf(Ferr, "\t%d-%d", sday, eday)
			}
		}

		for j := 0; j < 8; j++ {
			fmt.Fprintf(Ferr, "\t%s", DAYweek[j])
		}

		for iw := 0; iw < Nw; iw++ {
			Wkdy := &Wkd[iw]
			fmt.Fprintf(Ferr, "\n%s", Wkdy.name)

			for j := 0; j < 8; j++ {
				wday := &Wkdy.wday[j]
				fmt.Fprintf(Ferr, "\t%d", *wday)
			}
		}

		for sc := 0; sc < Nsc; sc++ {
			Dsch := &Dh[sc]
			fmt.Fprintf(Ferr, "\nVL\t%s\t[%d]", Dsch.name, sc)

			for jsc := 0; jsc < Dsch.N; jsc++ {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Fprintf(Ferr, "\t%d-(%.2g)-%d", stime, val, etime)
			}
		}

		for sw := 0; sw < Nsw; sw++ {
			Dscw := &Dw[sw]
			fmt.Fprintf(Ferr, "\nSW\t%s\t[%d]", Dscw.name, sw)

			for jsw := 0; jsw < Dscw.N; jsw++ {
				stime := Dscw.stime[jsw]
				mode := Dscw.mode[jsw]
				etime := Dscw.etime[jsw]
				fmt.Fprintf(Ferr, "\t%d-(%c)-%d", stime, mode, etime)
			}
		}

		fmt.Fprintf(Ferr, "\n\n")
	}
}

/* ----------------------------------------------------------------- */

func dprschdata(Sh []SCH, Sw []SCH) {
	const dmax = 366

	Nsc := len(Sh)
	Nsw := len(Sw)

	if DEBUG {
		fmt.Printf("\n*** dprschdata  ***\n")
		fmt.Printf("\n== Sch.end=%d   Scw.end=%d\n", Nsc, Nsw)

		for i := 0; i < Nsc; i++ {
			Sch := &Sh[i]
			fmt.Printf("\nSCH= %s (%2d) ", Sch.name, i)

			k := 1
			for d := 1; d < dmax; d++ {
				day := Sch.day[d]
				if FNNday(k, 1) == d {
					fmt.Printf("\n%2d - ", k)
					k++
				}
				fmt.Printf("%2d", day)
			}
		}

		for i := 0; i < Nsw; i++ {
			Scw := &Sw[i]
			fmt.Printf("\nSCW= %s (%2d) ", Scw.name, i)
			k := 1
			for d := 1; d < dmax; d++ {
				day := Scw.day[d]
				if FNNday(k, 1) == d {
					fmt.Printf("\n%2d - ", k)
					k++
				}
				fmt.Printf("%2d", day)
			}
		}
		fmt.Printf("\n")
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprschdata  ***\n")
		fmt.Fprintf(Ferr, "\n== Sch.end=%d   Scw.end=%d\n", Nsc, Nsw)

		for i := 0; i < Nsc; i++ {
			Sch := &Sh[i]
			fmt.Fprintf(Ferr, "\nSCH=%s\t[%d]\t", Sch.name, i)

			k := 1
			for d := 1; d < dmax; d++ {
				day := Sch.day[d]
				if FNNday(k, 1) == d {
					fmt.Fprintf(Ferr, "\n%2d - ", k)
					k++
				}
				fmt.Fprintf(Ferr, "%2d", day)
			}
		}

		for i := 0; i < Nsw; i++ {
			Scw := &Sw[i]
			fmt.Fprintf(Ferr, "\nSCW= %s (%2d) ", Scw.name, i)
			k := 1
			for d := 1; d < dmax; d++ {
				day := Scw.day[d]
				if FNNday(k, 1) == d {
					fmt.Fprintf(Ferr, "\n%2d - ", k)
					k++
				}
				fmt.Fprintf(Ferr, "%2d", day)
			}
		}
		fmt.Fprintf(Ferr, "\n")
	}
}

/* ----------------------------------------------------------------- */

func dprachv(Nroom int, Room []ROOM) {

	f := func(s io.Writer) {
		fmt.Fprintln(Ferr, "\n*** dprachv***")

		for i := 0; i < Nroom; i++ {
			Rm := &Room[i]
			fmt.Fprintf(Ferr, "to rm: %-10s   from rms(sch):", Rm.Name)

			for j := 0; j < Rm.Nachr; j++ {
				A := &Rm.achr[j]
				fmt.Fprintf(Ferr, "  %-10s (%3d)", Room[A.rm].Name, A.sch)
			}
			fmt.Fprintln(Ferr)
		}
	}

	if DEBUG {
		f(os.Stdout)
	}

	if Ferr != nil {
		f(Ferr)
	}
}

/* ----------------------------------------------------------------- */

func dprexsf(E []EXSF) {
	if E == nil {
		return
	}

	if DEBUG {
		fmt.Println("\n*** dprexsf ***")
		for i := 0; i < E[0].End; i++ {
			Exs := &E[i]
			fmt.Printf("%2d  %-11s  typ=%c Wa=%6.2f Wb=%5.2f Rg=%4.2f  z=%5.2f edf=%6.2e\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprexsf ***")
		fmt.Fprintln(Ferr, "\tNo.\tName\ttyp\tWa\tWb\tRg\tz\tedf")

		for i := 0; i < E[0].End; i++ {
			Exs := &E[i]
			fmt.Fprintf(Ferr, "\t%d\t%s\t%c\t%.4g\t%.4g\t%.2g\t%.2g\t%.2g\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}
}

/* ----------------------------------------------------------------- */

func dprwwdata(Wa []WALL, Wi []WINDOW) {
	if DEBUG {
		fmt.Printf("\n*** dprwwdata ***\nWALLdata\n")

		for i := 0; i < Wa[0].end; i++ {
			Wall := &Wa[i]
			fmt.Printf("\nWall i=%d %s R=%5.3f IP=%d Ei=%4.2f Eo=%4.2f as=%4.2f\n", i, Wall.name, Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Printf("   %2d  %-10s %5.3f %2d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Printf("\nWINDOWdata\n")

		for i := 0; i < Wi[0].end; i++ {
			Window := &Wi[i]
			fmt.Printf("windows  %s\n", Window.Name)
			fmt.Printf(" R=%f t=%f B=%f  Ei=%f Eo=%f\n", Window.Rwall, Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprwwdata ***\nWALLdata\n")

		for i := 0; i < Wa[0].end; i++ {
			Wall := &Wa[i]
			fmt.Fprintf(Ferr, "\nWall[%d]\t%s\tR=%.3g\tIP=%d\tEi=%.2g\tEo=%.2g\tas=%.2g\n", i, Wall.name, Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			fmt.Fprintf(Ferr, "\tNo.\tcode\tL\tND\n")

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Fprintf(Ferr, "\t%d\t%s\t%.3g\t%d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Fprintf(Ferr, "\nWINDOWdata\n")

		for i := 0; i < Wi[0].end; i++ {
			Window := &Wi[i]
			fmt.Fprintf(Ferr, "windows[%d]\t%s\n", i, Window.Name)
			fmt.Fprintf(Ferr, "\tR=%.3g\tt=%.2g\tB=%.2g\tEi=%.2g\tEo=%.2g\n", Window.Rwall,
				Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}
}

/* ----------------------------------------------------------------- */

func dprroomdata(R []ROOM, S []RMSRF) {
	if DEBUG {
		fmt.Printf("\n*** dprroomdata ***\n")

		for i := 0; i < R[0].end; i++ {
			Room := &R[i]

			fmt.Printf("\n==room=(%d)    %s   N=%d  Ntr=%d Nrp=%d  V=%8.1f   MRM=%10.4e\n",
				i, Room.Name, Room.N, Room.Ntr, Room.Nrp, Room.VRM, Room.MRM)
			fmt.Printf("   Floor area=%6.2f   Total surface area=%6.2f\n", Room.FArea, Room.Area)

			fmt.Printf("   Gve=%f    Gvi=%f\n",
				Room.Gve, Room.Gvi)
			fmt.Printf("   Light=%f  Ltyp=%c  ", Room.Light, Room.Ltyp)
			fmt.Printf("  Nhm=%f\n",
				Room.Nhm)
			fmt.Printf("  Apsc=%f  Apsr=%f   ",
				Room.Apsc, Room.Apsr)
			fmt.Printf("  Apl=%f \n", Room.Apl)

			for j := 0; j < Room.N; j++ {
				Sdd := &S[Room.Brs+j]
				fmt.Printf(" %2d  ble=%c typ=%c name=%8s exs=%2d nxrm=%2d nxn=%2d ",
					Room.Brs+j, Sdd.ble, Sdd.typ, Sdd.Name, Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Printf("wd=%2d Nfn=%2d A=%5.1f mwside=%c mwtype=%c Ei=%.2f Eo=%.2f\n",
					Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprroomdata ***\n")

		for i := 0; i < R[0].end; i++ {
			Room := &R[i]

			fmt.Fprintf(Ferr, "\n==room=(%d)\t%s\tN=%d\tNtr=%d\tNrp=%d\tV=%.3g\tMRM=%.2g\n",
				i, Room.Name, Room.N, Room.Ntr, Room.Nrp, Room.VRM, Room.MRM)
			fmt.Fprintf(Ferr, "\tFloor_area=%.3g\tTotal_surface_area=%.2g\n", Room.FArea, Room.Area)

			fmt.Fprintf(Ferr, "\tGve=%.2g\tGvi=%.2g\n", Room.Gve, Room.Gvi)
			fmt.Fprintf(Ferr, "\tLight=%.2g\tLtyp=%c", Room.Light, Room.Ltyp)
			fmt.Fprintf(Ferr, "\tNhm=%.2g\n", Room.Nhm)
			fmt.Fprintf(Ferr, "\tApsc=%.2g\tApsr=%.2g", Room.Apsc, Room.Apsr)
			fmt.Fprintf(Ferr, "\tApl=%.2g\n", Room.Apl)

			fmt.Fprintf(Ferr, "\tNo.\tble\ttyp\tname\texs\tnxrmd\tnxn\t")
			fmt.Fprintf(Ferr, "wd\tNfn\tA\tmwside\tmwtype\tEi\tEo\n")

			for j := 0; j < Room.N; j++ {
				Sdd := &S[Room.Brs+j]
				fmt.Fprintf(Ferr, "\t%d\t%c\t%c\t%s\t%d\t%d\t%d\t", Room.Brs+j, Sdd.ble, Sdd.typ, Sdd.Name, Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Fprintf(Ferr, "%d\t%d\t%.3g\t%c\t%c\t%.2f\t%.2f\n", Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}
}

/* ----------------------------------------------------------------- */

func dprballoc(M []MWALL, S []RMSRF) {

	if DEBUG {
		fmt.Println("\n*** dprballoc ***")

		N := M[0].end
		for mw := 0; mw < N; mw++ {
			Mw := &M[mw]
			id := S[Mw.ns].wd
			fmt.Printf(" %2d n=%2d  rm=%2d  nxrm=%2d wd=%2d wall=%s M=%2d A=%.2f\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, Mw.wall.name, Mw.M, Mw.sd.A)
		}
	}
	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprballoc ***")
		fmt.Fprintln(Ferr, "\tNo.\tn\trm\tnxrm\twd\twall\tM\tA")

		N := M[0].end
		for mw := 0; mw < N; mw++ {
			Mw := &M[mw]
			id := S[Mw.ns].wd
			fmt.Fprintf(Ferr, "\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%.2g\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, Mw.wall.name, Mw.M, Mw.sd.A)
		}
	}
}
