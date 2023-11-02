package eeslism

import (
	"fmt"
	"io"
)

/*  システム使用機器についての出力  */

var __Hcmpprint_id int

func Hcmpprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64, Eqsys *EQSYS, Rdpnl []*RDPNL) {
	var j int

	if __Hcmpprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintln(fo, "-cat")
			}

			boiprint(fo, __Hcmpprint_id, Eqsys.Boi)
			refaprint(fo, __Hcmpprint_id, Eqsys.Refa)
			collprint(fo, __Hcmpprint_id, Eqsys.Coll)
			hccprint(fo, __Hcmpprint_id, Eqsys.Hcc)
			pipeprint(fo, __Hcmpprint_id, Eqsys.Pipe)
			hexprint(fo, __Hcmpprint_id, Eqsys.Hex)
			stankcmpprt(fo, __Hcmpprint_id, Eqsys.Stank)
			pumpprint(fo, __Hcmpprint_id, Eqsys.Pump)
			hcldprint(fo, __Hcmpprint_id, Eqsys.Hcload)
			vavprint(fo, __Hcmpprint_id, Eqsys.Vav)
			stheatprint(fo, __Hcmpprint_id, Eqsys.Stheat)
			Thexprint(fo, __Hcmpprint_id, Eqsys.Thex)
			Qmeasprint(fo, __Hcmpprint_id, Eqsys.Qmeas)
			PVprint(fo, __Hcmpprint_id, Eqsys.PVcmp)
			Desiprint(fo, __Hcmpprint_id, Eqsys.Desi)
			Evacprint(fo, __Hcmpprint_id, Eqsys.Evac)

			if j == 0 {
				fmt.Fprintln(fo, "*")
				fmt.Fprintln(fo, "#")
			}

			__Hcmpprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	boiprint(fo, __Hcmpprint_id, Eqsys.Boi)
	refaprint(fo, __Hcmpprint_id, Eqsys.Refa)
	collprint(fo, __Hcmpprint_id, Eqsys.Coll)
	hccprint(fo, __Hcmpprint_id, Eqsys.Hcc)
	pipeprint(fo, __Hcmpprint_id, Eqsys.Pipe)
	hexprint(fo, __Hcmpprint_id, Eqsys.Hex)
	stankcmpprt(fo, __Hcmpprint_id, Eqsys.Stank)
	pumpprint(fo, __Hcmpprint_id, Eqsys.Pump)
	hcldprint(fo, __Hcmpprint_id, Eqsys.Hcload)
	vavprint(fo, __Hcmpprint_id, Eqsys.Vav)
	stheatprint(fo, __Hcmpprint_id, Eqsys.Stheat)
	Thexprint(fo, __Hcmpprint_id, Eqsys.Thex)
	Qmeasprint(fo, __Hcmpprint_id, Eqsys.Qmeas)
	PVprint(fo, __Hcmpprint_id, Eqsys.PVcmp)
	Desiprint(fo, __Hcmpprint_id, Eqsys.Desi)
	Evacprint(fo, __Hcmpprint_id, Eqsys.Evac)

	if SIMUL_BUILDG {
		panelprint(fo, __Hcmpprint_id, Rdpnl)
	}

}

var __Hstkprint_id int = 0

func Hstkprint(fo io.Writer, title string, mon int, day int, time float64, Eqsys *EQSYS) {
	if __Hstkprint_id == 0 {
		fmt.Fprintf(fo, "%s ;\n", title)
		stankivprt(fo, __Hstkprint_id, Eqsys.Stank)
		__Hstkprint_id++
	}
	if len(Eqsys.Stank) > 0 {
		fmt.Fprintf(fo, "%02d %02d %5.2f  ", mon, day, time)
		stankivprt(fo, __Hstkprint_id, Eqsys.Stank)
	}
	fmt.Fprintln(fo, " ;")
}
