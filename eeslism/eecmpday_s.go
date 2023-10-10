package eeslism

import (
	"fmt"
	"io"
)

/* システム要素機器の日集計処理 */
var __Compoday_OldDay int = 0
var __Compoday_OldMon int = 0

func Compoday(Mon, Day, Nday, ttmm int, Eqsys *EQSYS, SimDayend int) {

	// 日集計
	if Nday != __Compoday_OldDay {
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

		__Compoday_OldDay = Nday
	}

	if Mon != __Compoday_OldMon {
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

		__Compoday_OldMon = Mon
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

var __Compodyprt_id int

func Compodyprt(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	if __Compodyprt_id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boidyprt(fo, __Compodyprt_id, Eqsys.Nboi, Eqsys.Boi)
			refadyprt(fo, __Compodyprt_id, Eqsys.Nrefa, Eqsys.Refa)
			colldyprt(fo, __Compodyprt_id, Eqsys.Ncoll, Eqsys.Coll)
			hccdyprt(fo, __Compodyprt_id, Eqsys.Nhcc, Eqsys.Hcc)
			pipedyprt(fo, __Compodyprt_id, Eqsys.Npipe, Eqsys.Pipe)
			hexdyprt(fo, __Compodyprt_id, Eqsys.Nhex, Eqsys.Hex)
			stankdyprt(fo, __Compodyprt_id, Eqsys.Nstank, Eqsys.Stank)
			pumpdyprt(fo, __Compodyprt_id, Eqsys.Npump, Eqsys.Pump)
			hclddyprt(fo, __Compodyprt_id, Eqsys.Nhcload, Eqsys.Hcload)
			stheatdyprt(fo, __Compodyprt_id, Eqsys.Nstheat, Eqsys.Stheat)
			Qmeasdyprt(fo, __Compodyprt_id, Eqsys.Nqmeas, Eqsys.Qmeas)
			Thexdyprt(fo, __Compodyprt_id, Eqsys.Nthex, Eqsys.Thex)
			PVdyprt(fo, __Compodyprt_id, Eqsys.Npv, Eqsys.PVcmp)
			Desidyprt(fo, __Compodyprt_id, Eqsys.Ndesi, Eqsys.Desi)

			paneldyprt(fo, __Compodyprt_id, Nrdpnl, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compodyprt_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boidyprt(fo, __Compodyprt_id, Eqsys.Nboi, Eqsys.Boi)
	refadyprt(fo, __Compodyprt_id, Eqsys.Nrefa, Eqsys.Refa)
	colldyprt(fo, __Compodyprt_id, Eqsys.Ncoll, Eqsys.Coll)
	hccdyprt(fo, __Compodyprt_id, Eqsys.Nhcc, Eqsys.Hcc)
	pipedyprt(fo, __Compodyprt_id, Eqsys.Npipe, Eqsys.Pipe)
	hexdyprt(fo, __Compodyprt_id, Eqsys.Nhex, Eqsys.Hex)
	stankdyprt(fo, __Compodyprt_id, Eqsys.Nstank, Eqsys.Stank)
	pumpdyprt(fo, __Compodyprt_id, Eqsys.Npump, Eqsys.Pump)
	hclddyprt(fo, __Compodyprt_id, Eqsys.Nhcload, Eqsys.Hcload)
	stheatdyprt(fo, __Compodyprt_id, Eqsys.Nstheat, Eqsys.Stheat)
	Qmeasdyprt(fo, __Compodyprt_id, Eqsys.Nqmeas, Eqsys.Qmeas)
	Thexdyprt(fo, __Compodyprt_id, Eqsys.Nthex, Eqsys.Thex)
	PVdyprt(fo, __Compodyprt_id, Eqsys.Npv, Eqsys.PVcmp)
	Desidyprt(fo, __Compodyprt_id, Eqsys.Ndesi, Eqsys.Desi)

	paneldyprt(fo, __Compodyprt_id, Nrdpnl, Rdpnl)

}

/* システム要素機器の月集計結果出力 */

var __Compomonprt_id int

func Compomonprt(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	if __Compomonprt_id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boimonprt(fo, __Compomonprt_id, Eqsys.Nboi, Eqsys.Boi)
			refamonprt(fo, __Compomonprt_id, Eqsys.Nrefa, Eqsys.Refa)
			collmonprt(fo, __Compomonprt_id, Eqsys.Ncoll, Eqsys.Coll)
			hccmonprt(fo, __Compomonprt_id, Eqsys.Nhcc, Eqsys.Hcc)
			pipemonprt(fo, __Compomonprt_id, Eqsys.Npipe, Eqsys.Pipe)
			hexmonprt(fo, __Compomonprt_id, Eqsys.Nhex, Eqsys.Hex)
			stankmonprt(fo, __Compomonprt_id, Eqsys.Nstank, Eqsys.Stank)
			pumpmonprt(fo, __Compomonprt_id, Eqsys.Npump, Eqsys.Pump)
			hcldmonprt(fo, __Compomonprt_id, Eqsys.Nhcload, Eqsys.Hcload)
			stheatmonprt(fo, __Compomonprt_id, Eqsys.Nstheat, Eqsys.Stheat)
			Qmeasmonprt(fo, __Compomonprt_id, Eqsys.Nqmeas, Eqsys.Qmeas)
			Thexmonprt(fo, __Compomonprt_id, Eqsys.Nthex, Eqsys.Thex)
			PVmonprt(fo, __Compomonprt_id, Eqsys.Npv, Eqsys.PVcmp)

			panelmonprt(fo, __Compomonprt_id, Nrdpnl, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compomonprt_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boimonprt(fo, __Compomonprt_id, Eqsys.Nboi, Eqsys.Boi)
	refamonprt(fo, __Compomonprt_id, Eqsys.Nrefa, Eqsys.Refa)
	collmonprt(fo, __Compomonprt_id, Eqsys.Ncoll, Eqsys.Coll)
	hccmonprt(fo, __Compomonprt_id, Eqsys.Nhcc, Eqsys.Hcc)
	pipemonprt(fo, __Compomonprt_id, Eqsys.Npipe, Eqsys.Pipe)
	hexmonprt(fo, __Compomonprt_id, Eqsys.Nhex, Eqsys.Hex)
	stankmonprt(fo, __Compomonprt_id, Eqsys.Nstank, Eqsys.Stank)
	pumpmonprt(fo, __Compomonprt_id, Eqsys.Npump, Eqsys.Pump)
	hcldmonprt(fo, __Compomonprt_id, Eqsys.Nhcload, Eqsys.Hcload)
	stheatmonprt(fo, __Compomonprt_id, Eqsys.Nstheat, Eqsys.Stheat)
	Qmeasmonprt(fo, __Compomonprt_id, Eqsys.Nqmeas, Eqsys.Qmeas)
	Thexmonprt(fo, __Compomonprt_id, Eqsys.Nthex, Eqsys.Thex)
	PVmonprt(fo, __Compomonprt_id, Eqsys.Npv, Eqsys.PVcmp)

	panelmonprt(fo, __Compomonprt_id, Nrdpnl, Rdpnl)
}

/* システム要素機器の年集計結果出力 */

var __Compomtprt_id int = 0

func Compomtprt(fo io.Writer, mrk string, Simc *SIMCONTL, Eqsys *EQSYS, Nrdpnl int, Rdpnl []RDPNL) {
	if __Compomtprt_id == 0 {
		ttlmtprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			stheatmtprt(fo, __Compomtprt_id, Eqsys.Nstheat, Eqsys.Stheat, 0, 0)
			boimtprt(fo, __Compomtprt_id, Eqsys.Nboi, Eqsys.Boi, 0, 0)
			refamtprt(fo, __Compomtprt_id, Eqsys.Nrefa, Eqsys.Refa, 0, 0)
			pumpmtprt(fo, __Compomtprt_id, Eqsys.Npump, Eqsys.Pump, 0, 0)
			PVmtprt(fo, __Compomtprt_id, Eqsys.Npv, Eqsys.PVcmp, 0, 0)
			hcldmtprt(fo, __Compomtprt_id, Eqsys.Nhcload, 0, 0, Eqsys.Hcload)
			panelmtprt(fo, __Compomtprt_id, Nrdpnl, Rdpnl, 0, 0)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compomtprt_id++
		}
	}

	for mo := 1; mo <= 12; mo++ {
		for tt := 1; tt <= 24; tt++ {
			fmt.Fprintf(fo, "%02d %02d\n", mo, tt)
			stheatmtprt(fo, __Compomtprt_id, Eqsys.Nstheat, Eqsys.Stheat, mo, tt)
			boimtprt(fo, __Compomtprt_id, Eqsys.Nboi, Eqsys.Boi, mo, tt)
			refamtprt(fo, __Compomtprt_id, Eqsys.Nrefa, Eqsys.Refa, mo, tt)
			pumpmtprt(fo, __Compomtprt_id, Eqsys.Npump, Eqsys.Pump, mo, tt)
			PVmtprt(fo, __Compomtprt_id, Eqsys.Npv, Eqsys.PVcmp, mo, tt)
			hcldmtprt(fo, __Compomtprt_id, Eqsys.Nhcload, mo, tt, Eqsys.Hcload)
			panelmtprt(fo, __Compomtprt_id, Nrdpnl, Rdpnl, mo, tt)
		}
	}
}
