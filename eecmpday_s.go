package main

import (
	"fmt"
	"os"
)

/* システム要素機器の日集計処理 */

func Compoday(Mon, Day, Nday, ttmm int, Eqsys *EQSYS, SimDayend int) {
	staticOldDay := 0
	staticOldMon := 0

	// 日集計
	if Nday != staticOldDay {
		boidyint(Eqsys.Nboi, Eqsys.Boi)
		refadyint(Eqsys.Nrefa, Eqsys.Refa)
		colldyint(Eqsys.Ncoll, Eqsys.Coll)
		hccdyint(Eqsys.Nhcc, Eqsys.Hcc)
		pipedyint(Eqsys.Npipe, Eqsys.Pipe)
		hexdyint(Eqsys.Nhex, Eqsys.Hex)
		stankdyint(Eqsys.Nstank, Eqsys.Stank)
		pumpdyint(Eqsys.Npump, Eqsys.Pump)
		hclddyint(Eqsys.Nhcload, Eqsys.Hcload)
		stheatdyint(Eqsys.Nstheat, Eqsys.Stheat)
		Thexdyint(Eqsys.Nthex, Eqsys.Thex)
		Qmeasdyint(Eqsys.Nqmeas, Eqsys.Qmeas)
		PVdyint(Eqsys.Npv, Eqsys.PVcmp)
		Desidyint(Eqsys.Ndesi, Eqsys.Desi)

		staticOldDay = Nday
	}

	if Mon != staticOldMon {
		boimonint(Eqsys.Nboi, Eqsys.Boi)
		refamonint(Eqsys.Nrefa, Eqsys.Refa)
		collmonint(Eqsys.Ncoll, Eqsys.Coll)
		hccmonint(Eqsys.Nhcc, Eqsys.Hcc)
		pipemonint(Eqsys.Npipe, Eqsys.Pipe)
		hexmonint(Eqsys.Nhex, Eqsys.Hex)
		stankmonint(Eqsys.Nstank, Eqsys.Stank)
		pumpmonint(Eqsys.Npump, Eqsys.Pump)
		hcldmonint(Eqsys.Nhcload, Eqsys.Hcload)
		stheatmonint(Eqsys.Nstheat, Eqsys.Stheat)
		Thexmonint(Eqsys.Nthex, Eqsys.Thex)
		Qmeasmonint(Eqsys.Nqmeas, Eqsys.Qmeas)
		PVmonint(Eqsys.Npv, Eqsys.PVcmp)

		staticOldMon = Mon
	}

	// 日集計
	boiday(Mon, Day, ttmm, Eqsys.Nboi, Eqsys.Boi, Nday, SimDayend)
	refaday(Mon, Day, ttmm, Eqsys.Nrefa, Eqsys.Refa, Nday, SimDayend)
	collday(Mon, Day, ttmm, Eqsys.Ncoll, Eqsys.Coll, Nday, SimDayend)
	hccday(Mon, Day, ttmm, Eqsys.Nhcc, Eqsys.Hcc, Nday, SimDayend)
	pipeday(Mon, Day, ttmm, Eqsys.Npipe, Eqsys.Pipe, Nday, SimDayend)
	hexday(Mon, Day, ttmm, Eqsys.Nhex, Eqsys.Hex, Nday, SimDayend)
	stankday(Mon, Day, ttmm, Eqsys.Nstank, Eqsys.Stank, Nday, SimDayend)
	pumpday(Mon, Day, ttmm, Eqsys.Npump, Eqsys.Pump, Nday, SimDayend)
	hcldday(Mon, Day, ttmm, Eqsys.Nhcload, Nday, SimDayend, Eqsys.Hcload)
	stheatday(Mon, Day, ttmm, Eqsys.Nstheat, Eqsys.Stheat, Nday, SimDayend)
	Thexday(Mon, Day, ttmm, Eqsys.Nthex, Eqsys.Thex, Nday, SimDayend)
	Qmeasday(Mon, Day, ttmm, Eqsys.Nqmeas, Eqsys.Qmeas, Nday, SimDayend)
	PVday(Mon, Day, ttmm, Eqsys.Npv, Eqsys.PVcmp, Nday, SimDayend)
	Desiday(Mon, Day, ttmm, Eqsys.Ndesi, Eqsys.Desi, Nday, SimDayend)

	// 月集計
	//boimon(Mon, Day, ttmm, Eqsys.Nboi, Eqsys.Boi);
	//refamon(Mon, Day, ttmm, Eqsys.Nrefa, Eqsys.Refa);
	//collmon(Mon, Day, ttmm, Eqsys.Ncoll, Eqsys.Coll);
	//hccmon(Mon, Day, ttmm, Eqsys.Nhcc, Eqsys.Hcc);
	//pipemon(Mon, Day, ttmm, Eqsys.Npipe, Eqsys.Pipe);
	//hexmon(Mon, Day, ttmm, Eqsys.Nhex, Eqsys.Hex);
	//stankmon(Mon, Day, ttmm, Eqsys.Nstank, Eqsys.Stank);
	//pumpmon(Mon, Day, ttmm, Eqsys.Npump, Eqsys.Pump);
	//hcldmon(Mon, Day, ttmm, Eqsys.Nhcload, Eqsys.Hcload);
	//stheatmon(Mon, Day, ttmm, Eqsys.Nstheat, Eqsys.stheat) ;
	//Thexmon(Mon, Day, ttmm, Eqsys.Nthex, Eqsys.Thex) ;
	//Qmeasmon(Mon, Day, ttmm, Eqsys.Nqmeas, Eqsys.Qmeas) ;
	//PVmon(Mon, Day, ttmm, Eqsys.Npv, Eqsys.PVcmp ) ;
}

/* システム要素機器の日集計結果出力 */

func Compodyprt(fo *os.File, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	var j, id int

	if id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boidyprt(fo, id, Eqsys.Nboi, Eqsys.Boi)
			refadyprt(fo, id, Eqsys.Nrefa, Eqsys.Refa)
			colldyprt(fo, id, Eqsys.Ncoll, Eqsys.Coll)
			hccdyprt(fo, id, Eqsys.Nhcc, Eqsys.Hcc)
			pipedyprt(fo, id, Eqsys.Npipe, Eqsys.Pipe)
			hexdyprt(fo, id, Eqsys.Nhex, Eqsys.Hex)
			stankdyprt(fo, id, Eqsys.Nstank, Eqsys.Stank)
			pumpdyprt(fo, id, Eqsys.Npump, Eqsys.Pump)
			hclddyprt(fo, id, Eqsys.Nhcload, Eqsys.Hcload)
			stheatdyprt(fo, id, Eqsys.Nstheat, Eqsys.Stheat)
			Qmeasdyprt(fo, id, Eqsys.Nqmeas, Eqsys.Qmeas)
			Thexdyprt(fo, id, Eqsys.Nthex, Eqsys.Thex)
			PVdyprt(fo, id, Eqsys.Npv, Eqsys.PVcmp)
			Desidyprt(fo, id, Eqsys.Ndesi, Eqsys.Desi)

			paneldyprt(fo, id, Nrdpnl, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boidyprt(fo, id, Eqsys.Nboi, Eqsys.Boi)
	refadyprt(fo, id, Eqsys.Nrefa, Eqsys.Refa)
	colldyprt(fo, id, Eqsys.Ncoll, Eqsys.Coll)
	hccdyprt(fo, id, Eqsys.Nhcc, Eqsys.Hcc)
	pipedyprt(fo, id, Eqsys.Npipe, Eqsys.Pipe)
	hexdyprt(fo, id, Eqsys.Nhex, Eqsys.Hex)
	stankdyprt(fo, id, Eqsys.Nstank, Eqsys.Stank)
	pumpdyprt(fo, id, Eqsys.Npump, Eqsys.Pump)
	hclddyprt(fo, id, Eqsys.Nhcload, Eqsys.Hcload)
	stheatdyprt(fo, id, Eqsys.Nstheat, Eqsys.Stheat)
	Qmeasdyprt(fo, id, Eqsys.Nqmeas, Eqsys.Qmeas)
	Thexdyprt(fo, id, Eqsys.Nthex, Eqsys.Thex)
	PVdyprt(fo, id, Eqsys.Npv, Eqsys.PVcmp)
	Desidyprt(fo, id, Eqsys.Ndesi, Eqsys.Desi)

	paneldyprt(fo, id, Nrdpnl, Rdpnl)

}

/* システム要素機器の月集計結果出力 */

func Compomonprt(fo *os.File, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	var j int
	staticId := 0

	if staticId == 0 {
		ttldyprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boimonprt(fo, staticId, Eqsys.Nboi, Eqsys.Boi)
			refamonprt(fo, staticId, Eqsys.Nrefa, Eqsys.Refa)
			collmonprt(fo, staticId, Eqsys.Ncoll, Eqsys.Coll)
			hccmonprt(fo, staticId, Eqsys.Nhcc, Eqsys.Hcc)
			pipemonprt(fo, staticId, Eqsys.Npipe, Eqsys.Pipe)
			hexmonprt(fo, staticId, Eqsys.Nhex, Eqsys.Hex)
			stankmonprt(fo, staticId, Eqsys.Nstank, Eqsys.Stank)
			pumpmonprt(fo, staticId, Eqsys.Npump, Eqsys.Pump)
			hcldmonprt(fo, staticId, Eqsys.Nhcload, Eqsys.Hcload)
			stheatmonprt(fo, staticId, Eqsys.Nstheat, Eqsys.Stheat)
			Qmeasmonprt(fo, staticId, Eqsys.Nqmeas, Eqsys.Qmeas)
			Thexmonprt(fo, staticId, Eqsys.Nthex, Eqsys.Thex)
			PVmonprt(fo, staticId, Eqsys.Npv, Eqsys.PVcmp)

			panelmonprt(fo, staticId, Nrdpnl, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			staticId++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boimonprt(fo, staticId, Eqsys.Nboi, Eqsys.Boi)
	refamonprt(fo, staticId, Eqsys.Nrefa, Eqsys.Refa)
	collmonprt(fo, staticId, Eqsys.Ncoll, Eqsys.Coll)
	hccmonprt(fo, staticId, Eqsys.Nhcc, Eqsys.Hcc)
	pipemonprt(fo, staticId, Eqsys.Npipe, Eqsys.Pipe)
	hexmonprt(fo, staticId, Eqsys.Nhex, Eqsys.Hex)
	stankmonprt(fo, staticId, Eqsys.Nstank, Eqsys.Stank)
	pumpmonprt(fo, staticId, Eqsys.Npump, Eqsys.Pump)
	hcldmonprt(fo, staticId, Eqsys.Nhcload, Eqsys.Hcload)
	stheatmonprt(fo, staticId, Eqsys.Nstheat, Eqsys.Stheat)
	Qmeasmonprt(fo, staticId, Eqsys.Nqmeas, Eqsys.Qmeas)
	Thexmonprt(fo, staticId, Eqsys.Nthex, Eqsys.Thex)
	PVmonprt(fo, staticId, Eqsys.Npv, Eqsys.PVcmp)

	panelmonprt(fo, staticId, Nrdpnl, Rdpnl)
}

/* システム要素機器の年集計結果出力 */

func Compomtprt(fo *os.File, mrk string, Simc *SIMCONTL, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	var j int
	var id int = 0
	var mo, tt int

	if id == 0 {
		ttlmtprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			stheatmtprt(fo, id, Eqsys.Nstheat, Eqsys.Stheat, 0, 0)
			boimtprt(fo, id, Eqsys.Nboi, Eqsys.Boi, 0, 0)
			refamtprt(fo, id, Eqsys.Nrefa, Eqsys.Refa, 0, 0)
			pumpmtprt(fo, id, Eqsys.Npump, Eqsys.Pump, 0, 0)
			PVmtprt(fo, id, Eqsys.Npv, Eqsys.PVcmp, 0, 0)
			hcldmtprt(fo, id, Eqsys.Nhcload, 0, 0, Eqsys.Hcload, Cff_kWh)
			panelmtprt(fo, id, Nrdpnl, Rdpnl, 0, 0)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			id++
		}
	}

	for mo = 1; mo <= 12; mo++ {
		for tt = 1; tt <= 24; tt++ {
			fmt.Fprintf(fo, "%02d %02d\n", mo, tt)
			stheatmtprt(fo, id, Eqsys.Nstheat, Eqsys.Stheat, mo, tt)
			boimtprt(fo, id, Eqsys.Nboi, Eqsys.Boi, mo, tt)
			refamtprt(fo, id, Eqsys.Nrefa, Eqsys.Refa, mo, tt)
			pumpmtprt(fo, id, Eqsys.Npump, Eqsys.Pump, mo, tt)
			PVmtprt(fo, id, Eqsys.Npv, Eqsys.PVcmp, mo, tt)
			hcldmtprt(fo, id, Eqsys.Nhcload, mo, tt, Eqsys.Hcload, Cff_kWh)
			panelmtprt(fo, id, Nrdpnl, Rdpnl, mo, tt)
		}
	}
}
