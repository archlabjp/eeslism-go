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

	fmt.Print("---  Day of week -----\n   ")
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
		fmt.Printf("%2d", daywk[d])
	}
	fmt.Println()
}

/* ----------------------------------------------------------------- */

// スケジュール設定のデバッグ出力
func (Schdl *SCHDL) dprschtable() {

	Ssn, Wkd, Dh, Dw := Schdl.Seasn, Schdl.Wkdy, Schdl.Dsch, Schdl.Dscw

	Ns := len(Ssn)
	Nw := len(Wkd)
	Nsc := len(Dh)
	Nsw := len(Dw)

	if DEBUG {
		fmt.Printf("\n*** dprschtable  ***\n")
		fmt.Printf("\n=== Schtable end  is=%d  iw=%d  sc=%d  sw=%d\n", Ns, Nw, Nsc, Nsw)

		// 季節設定の出力
		for _, Seasn := range Ssn {
			fmt.Printf("\n- %s", Seasn.name)

			for js := range Seasn.sday {
				sday := Seasn.sday[js]
				eday := Seasn.eday[js]
				fmt.Printf("  %4d-%4d", sday, eday)
			}
		}

		// 曜日設定の出力
		for _, Wkdy := range Wkd {
			fmt.Printf("\n- %s", Wkdy.name)

			for _, wday := range Wkdy.wday {
				if wday {
					fmt.Printf("   1")
				} else {
					fmt.Printf("   0")
				}
			}
		}

		// 1日の設定値スケジュールの出力
		for sc, Dsch := range Dh {
			fmt.Printf("\n-VL   %10s (%2d) ", Dsch.name, sc)

			for jsc := range Dsch.stime {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Printf("  %4d-(%6.2f)-%4d", stime, val, etime)
			}
		}

		// 1日の切替スケジュールの出力
		for sw, Dscw := range Dw {
			fmt.Printf("\n-SW   %10s (%2d) ", Dscw.name, sw)

			for jsw := range Dscw.stime {
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

		// 季節設定の出力
		for _, Seasn := range Ssn {
			fmt.Fprintf(Ferr, "\n\t%s", Seasn.name)

			for js := range Seasn.sday {
				sday := Seasn.sday[js]
				eday := Seasn.eday[js]
				fmt.Fprintf(Ferr, "\t%d-%d", sday, eday)
			}
		}

		// 曜日の出力
		for j := range DAYweek {
			fmt.Fprintf(Ferr, "\t%s", DAYweek[j])
		}

		// 曜日設定の出力
		for _, Wkdy := range Wkd {
			fmt.Fprintf(Ferr, "\n%s", Wkdy.name)

			for _, wday := range Wkdy.wday {
				if wday {
					fmt.Fprintf(Ferr, "\t1")
				} else {
					fmt.Fprintf(Ferr, "\t0")
				}
			}
		}

		// 1日の設定値スケジュールの出力
		for sc, Dsch := range Dh {
			fmt.Fprintf(Ferr, "\nVL\t%s\t[%d]", Dsch.name, sc)

			for jsc := range Dsch.stime {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Fprintf(Ferr, "\t%d-(%.2g)-%d", stime, val, etime)
			}
		}

		// 1日の切替スケジュールの出力
		for sw, Dscw := range Dw {
			fmt.Fprintf(Ferr, "\nSW\t%s\t[%d]", Dscw.name, sw)

			for jsw := range Dscw.stime {
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
		fmt.Printf("\n== len(Sch)=%d   len(Scw)=%d\n", Nsc, Nsw)

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
		fmt.Fprintf(Ferr, "\n== len(Sch)=%d   len(Scw)=%d\n", Nsc, Nsw)

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

func dprachv(Room []ROOM) {

	f := func(s io.Writer) {
		fmt.Fprintln(Ferr, "\n*** dprachv***")

		for i := range Room {
			Rm := Room[i]
			fmt.Fprintf(Ferr, "to rm: %-10s   from rms(sch):", Rm.Name)

			for j := 0; j < Rm.Nachr; j++ {
				A := Rm.achr[j]
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

func (exsfs *EXSFS) dprexsf() {
	if exsfs.Exs == nil {
		return
	}

	if DEBUG {
		fmt.Println("\n*** dprexsf ***")
		for i, Exs := range exsfs.Exs {
			fmt.Printf("%2d  %-11s  typ=%c Wa=%6.2f Wb=%5.2f Rg=%4.2f  z=%5.2f edf=%6.2e\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprexsf ***")
		fmt.Fprintln(Ferr, "\tNo.\tName\ttyp\tWa\tWb\tRg\tz\tedf")

		for i, Exs := range exsfs.Exs {
			fmt.Fprintf(Ferr, "\t%d\t%s\t%c\t%.4g\t%.4g\t%.2g\t%.2g\t%.2g\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}
}

/* ----------------------------------------------------------------- */

func (Rmvls *RMVLS) dprwwdata() {
	if DEBUG {
		fmt.Printf("\n*** dprwwdata ***\nWALLdata\n")

		for i, Wall := range Rmvls.Wall {
			fmt.Printf("\nWall i=%d %s R=%5.3f IP=%d Ei=%4.2f Eo=%4.2f as=%4.2f\n", i, get_string_or_null(Wall.name), Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Printf("   %2d  %-10s %5.3f %2d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Printf("\nWINDOWdata\n")

		for _, Window := range Rmvls.Window {
			fmt.Printf("windows  %s\n", Window.Name)
			fmt.Printf(" R=%f t=%f B=%f  Ei=%f Eo=%f\n", Window.Rwall, Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprwwdata ***\nWALLdata\n")

		for i, Wall := range Rmvls.Wall {
			fmt.Fprintf(Ferr, "\nWall[%d]\t%s\tR=%.3g\tIP=%d\tEi=%.2g\tEo=%.2g\tas=%.2g\n", i, Wall.name, Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			fmt.Fprintf(Ferr, "\tNo.\tcode\tL\tND\n")

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Fprintf(Ferr, "\t%d\t%s\t%.3g\t%d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Fprintf(Ferr, "\nWINDOWdata\n")

		for i, Window := range Rmvls.Window {
			fmt.Fprintf(Ferr, "windows[%d]\t%s\n", i, Window.Name)
			fmt.Fprintf(Ferr, "\tR=%.3g\tt=%.2g\tB=%.2g\tEi=%.2g\tEo=%.2g\n", Window.Rwall,
				Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}
}

/* ----------------------------------------------------------------- */

func (Rmvls *RMVLS) dprroomdata() {
	if DEBUG {
		fmt.Printf("\n*** dprroomdata ***\n")

		for i, Room := range Rmvls.Room {
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
				Sdd := Rmvls.Sd[Room.Brs+j]
				fmt.Printf(" %2d  ble=%c typ=%c name=%8s exs=%2d nxrm=%2d nxn=%2d ",
					Room.Brs+j, Sdd.ble, Sdd.typ, get_string_or_null(Sdd.Name), Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Printf("wd=%2d Nfn=%2d A=%5.1f mwside=%c mwtype=%c Ei=%.2f Eo=%.2f\n",
					Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprroomdata ***\n")

		for i, Room := range Rmvls.Room {
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
				Sdd := Rmvls.Sd[Room.Brs+j]
				fmt.Fprintf(Ferr, "\t%d\t%c\t%c\t%s\t%d\t%d\t%d\t", Room.Brs+j, Sdd.ble, Sdd.typ, Sdd.Name, Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Fprintf(Ferr, "%d\t%d\t%.3g\t%c\t%c\t%.2f\t%.2f\n", Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}
}

/* ----------------------------------------------------------------- */

func (Rmvls *RMVLS) dprballoc() {
	if DEBUG {
		fmt.Println("\n*** dprballoc ***")

		for mw, Mw := range Rmvls.Mw {
			id := Rmvls.Sd[Mw.ns].wd
			fmt.Printf(" %2d n=%2d  rm=%2d  nxrm=%2d wd=%2d wall=%s M=%2d A=%.2f\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, get_string_or_null(Mw.wall.name), Mw.M, Mw.sd.A)
		}
	}
	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprballoc ***")
		fmt.Fprintln(Ferr, "\tNo.\tn\trm\tnxrm\twd\twall\tM\tA")

		for mw, Mw := range Rmvls.Mw {
			id := Rmvls.Sd[Mw.ns].wd
			fmt.Fprintf(Ferr, "\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%.2g\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, Mw.wall.name, Mw.M, Mw.sd.A)
		}
	}
}
