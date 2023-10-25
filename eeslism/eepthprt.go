package eeslism

import (
	"fmt"
	"io"
)

/* 経路に沿ったシステム要素の出力 */

var __Pathprint_id int

func Pathprint(fo io.Writer, title string, mon int, day int, time float64, Nmpath int, _Mpath []*MPATH) {
	if __Pathprint_id == 0 {
		__Pathprint_id++
		fmt.Fprintf(fo, "%s ;\n", title)
		fmt.Fprintf(fo, "%d\n", Nmpath)

		for _, Mpath := range _Mpath {
			fmt.Fprintf(fo, "%s %c %c %d\n", Mpath.Name, Mpath.Type, Mpath.Fluid, len(Mpath.Plist))

			if Mpath.Plist[0].Pelm[0].Ci == '>' {
				fmt.Fprint(fo, " >")
			}

			for _, Pli := range Mpath.Plist {
				if Pli.Name != "" {
					fmt.Fprintf(fo, "%s", Pli.Name)
				} else {
					fmt.Fprint(fo, "?")
				}
				fmt.Fprintf(fo, " %c %d\n", Pli.Type, len(Pli.Pelm))

				if Pli.Pelm[0].Ci == '>' {
					fmt.Fprintf(fo, " %s", Pli.Pelm[0].Cmp.Name)
				}

				for _, Pelm := range Pli.Pelm {
					fmt.Fprintf(fo, " %s", Pelm.Cmp.Name)
				}
				fmt.Fprint(fo, "\n")
			}
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)

	for _, Mpath := range _Mpath {
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

		for _, Pli := range Mpath.Plist {

			if c = Pli.Control; c == 0 {
				c = '?'
			}
			fmt.Fprintf(fo, " %5.3g %c: ", Pli.G, c)

			if Pli.Pelm[0].Ci == '>' {
				fmt.Fprintf(fo, fm, Pli.Pelm[0].In.Sysvin)
			}
			for _, Pelm := range Pli.Pelm {
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
