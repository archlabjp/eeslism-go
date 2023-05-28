package main

import (
	"fmt"
	"os"
)

/* 経路に沿ったシステム要素の出力 */

var __Pathprint_id int

func Pathprint(fo *os.File, title string, mon int, day int, time float64, Nmpath int, _Mpath []MPATH) {
	if __Pathprint_id == 0 {
		__Pathprint_id++
		fmt.Fprintf(fo, "%s ;\n", title)
		fmt.Fprintf(fo, "%d\n", Nmpath)

		for i := 0; i < Nmpath; i++ {
			Mpath := &_Mpath[i]
			fmt.Fprintf(fo, "%s %c %c %d\n", Mpath.Name, Mpath.Type, Mpath.Fluid, Mpath.Nlpath)

			if Mpath.Plist[0].Pelm[0].Ci == '>' {
				fmt.Fprint(fo, " >")
			}

			for j := 0; j < Mpath.Nlpath; j++ {
				Pli := &Mpath.Plist[j]

				if Pli.Name != "" {
					fmt.Fprintf(fo, "%s", Pli.Name)
				} else {
					fmt.Fprint(fo, "?")
				}
				fmt.Fprintf(fo, " %c %d\n", Pli.Type, Pli.Nelm)

				if Pli.Pelm[0].Ci == '>' {
					fmt.Fprintf(fo, " %s", Pli.Pelm[0].Cmp.Name)
				}

				for k := 0; k < Pli.Nelm; k++ {
					Pelm := Pli.Pelm[k]
					fmt.Fprintf(fo, " %s", Pelm.Cmp.Name)
				}
				fmt.Fprint(fo, "\n")
			}
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)

	for i := 0; i < Nmpath; i++ {
		Mpath := &_Mpath[i]

		var fm string
		fm = "%4.1f"
		if Mpath.Fluid == AIRx_FLD {
			fm = "%6.4f"
		}

		var c ControlSWType
		if c = Mpath.Control; c == 0 {
			c = '?'
		}
		fmt.Fprintf(fo, "[%c]", c)

		for j := 0; j < Mpath.Nlpath; j++ {
			Pli := &Mpath.Plist[j]

			if c = Pli.Control; c == 0 {
				c = '?'
			}
			fmt.Fprintf(fo, " %5.3g %c: ", Pli.G, c)

			if Pli.Pelm[0].Ci == '>' {
				fmt.Fprintf(fo, fm, Pli.Pelm[0].In.Sysvin)
			}
			for k := 0; k < Pli.Nelm; k++ {
				Pelm := Pli.Pelm[k]
				if Pelm.Out != nil {
					fmt.Fprintf(fo, fm, Pelm.Out.Sysv)
					if c = Pelm.Out.Control; c == 0 {
						c = '?'
					}
					fmt.Fprintf(fo, " %c ", c)
				}

			}
			fmt.Fprint(fo, "\n")
		}
	}
	fmt.Fprintf(fo, " ;\n")
}
