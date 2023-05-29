package main

import (
	"fmt"
	"io"
)

/*  システム使用機器についての出力  */

var __Hcmpprint_id int

func Hcmpprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	var j int

	if __Hcmpprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintln(fo, "-cat")
			}

			boiprint(fo, __Hcmpprint_id, Eqsys.Nboi, Eqsys.Boi)
			refaprint(fo, __Hcmpprint_id, Eqsys.Nrefa, Eqsys.Refa)
			collprint(fo, __Hcmpprint_id, Eqsys.Ncoll, Eqsys.Coll)
			hccprint(fo, __Hcmpprint_id, Eqsys.Nhcc, Eqsys.Hcc)
			pipeprint(fo, __Hcmpprint_id, Eqsys.Npipe, Eqsys.Pipe)
			hexprint(fo, __Hcmpprint_id, Eqsys.Nhex, Eqsys.Hex)
			stankcmpprt(fo, __Hcmpprint_id, Eqsys.Nstank, Eqsys.Stank)
			pumpprint(fo, __Hcmpprint_id, Eqsys.Npump, Eqsys.Pump)
			hcldprint(fo, __Hcmpprint_id, Eqsys.Nhcload, Eqsys.Hcload)
			vavprint(fo, __Hcmpprint_id, Eqsys.Nvav, Eqsys.Vav)
			stheatprint(fo, __Hcmpprint_id, Eqsys.Nstheat, Eqsys.Stheat)
			Thexprint(fo, __Hcmpprint_id, Eqsys.Nthex, Eqsys.Thex)
			Qmeasprint(fo, __Hcmpprint_id, Eqsys.Nqmeas, Eqsys.Qmeas)
			PVprint(fo, __Hcmpprint_id, Eqsys.Npv, Eqsys.PVcmp)
			Desiprint(fo, __Hcmpprint_id, Eqsys.Ndesi, Eqsys.Desi)
			Evacprint(fo, __Hcmpprint_id, Eqsys.Nevac, Eqsys.Evac)

			if j == 0 {
				fmt.Fprintln(fo, "*")
				fmt.Fprintln(fo, "#")
			}

			__Hcmpprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	boiprint(fo, __Hcmpprint_id, Eqsys.Nboi, Eqsys.Boi)
	refaprint(fo, __Hcmpprint_id, Eqsys.Nrefa, Eqsys.Refa)
	collprint(fo, __Hcmpprint_id, Eqsys.Ncoll, Eqsys.Coll)
	hccprint(fo, __Hcmpprint_id, Eqsys.Nhcc, Eqsys.Hcc)
	pipeprint(fo, __Hcmpprint_id, Eqsys.Npipe, Eqsys.Pipe)
	hexprint(fo, __Hcmpprint_id, Eqsys.Nhex, Eqsys.Hex)
	stankcmpprt(fo, __Hcmpprint_id, Eqsys.Nstank, Eqsys.Stank)
	pumpprint(fo, __Hcmpprint_id, Eqsys.Npump, Eqsys.Pump)
	hcldprint(fo, __Hcmpprint_id, Eqsys.Nhcload, Eqsys.Hcload)
	vavprint(fo, __Hcmpprint_id, Eqsys.Nvav, Eqsys.Vav)
	stheatprint(fo, __Hcmpprint_id, Eqsys.Nstheat, Eqsys.Stheat)
	Thexprint(fo, __Hcmpprint_id, Eqsys.Nthex, Eqsys.Thex)
	Qmeasprint(fo, __Hcmpprint_id, Eqsys.Nqmeas, Eqsys.Qmeas)
	PVprint(fo, __Hcmpprint_id, Eqsys.Npv, Eqsys.PVcmp)
	Desiprint(fo, __Hcmpprint_id, Eqsys.Ndesi, Eqsys.Desi)
	Evacprint(fo, __Hcmpprint_id, Eqsys.Nevac, Eqsys.Evac)

	if SIMUL_BUILDG {
		panelprint(fo, __Hcmpprint_id, Nrdpnl, Rdpnl)
	}

}

func Hstkprint(fo io.Writer, title string, mon int, day int, time float64, Eqsys *EQSYS) {
	staticId := 0
	if staticId == 0 {
		fmt.Fprintf(fo, "%s ;\n", title)
		stankivprt(fo, staticId, Eqsys.Nstank, Eqsys.Stank)
		staticId++
	}
	if Eqsys.Nstank > 0 {
		fmt.Fprintf(fo, "%02d %02d %5.2f  ", mon, day, time)
		stankivprt(fo, staticId, Eqsys.Nstank, Eqsys.Stank)
	}
	fmt.Fprintln(fo, " ;")
}
