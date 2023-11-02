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
		boidyint(Eqsys.Boi)
		refadyint(Eqsys.Refa)
		colldyint(Eqsys.Coll)
		hccdyint(Eqsys.Hcc)
		pipedyint(Eqsys.Pipe)
		hexdyint(Eqsys.Hex)
		stankdyint(Eqsys.Stank)
		pumpdyint(Eqsys.Pump)
		hclddyint(Eqsys.Hcload)
		stheatdyint(Eqsys.Stheat)
		Thexdyint(Eqsys.Thex)
		Qmeasdyint(Eqsys.Qmeas)
		PVdyint(Eqsys.PVcmp)
		Desidyint(Eqsys.Desi)

		__Compoday_OldDay = Nday
	}

	if Mon != __Compoday_OldMon {
		boimonint(Eqsys.Boi)
		refamonint(Eqsys.Refa)
		collmonint(Eqsys.Coll)
		hccmonint(Eqsys.Hcc)
		pipemonint(Eqsys.Pipe)
		hexmonint(Eqsys.Hex)
		stankmonint(Eqsys.Stank)
		pumpmonint(Eqsys.Pump)
		hcldmonint(Eqsys.Hcload)
		stheatmonint(Eqsys.Stheat)
		Thexmonint(Eqsys.Thex)
		Qmeasmonint(Eqsys.Qmeas)
		PVmonint(Eqsys.PVcmp)

		__Compoday_OldMon = Mon
	}

	// 日集計
	boiday(Mon, Day, ttmm, Eqsys.Boi, Nday, SimDayend)
	refaday(Mon, Day, ttmm, Eqsys.Refa, Nday, SimDayend)
	collday(Mon, Day, ttmm, Eqsys.Coll, Nday, SimDayend)
	hccday(Mon, Day, ttmm, Eqsys.Hcc, Nday, SimDayend)
	pipeday(Mon, Day, ttmm, Eqsys.Pipe, Nday, SimDayend)
	hexday(Mon, Day, ttmm, Eqsys.Hex, Nday, SimDayend)
	stankday(Mon, Day, ttmm, Eqsys.Stank, Nday, SimDayend)
	pumpday(Mon, Day, ttmm, Eqsys.Pump, Nday, SimDayend)
	hcldday(Mon, Day, ttmm, Nday, SimDayend, Eqsys.Hcload)
	stheatday(Mon, Day, ttmm, Eqsys.Stheat, Nday, SimDayend)
	Thexday(Mon, Day, ttmm, Eqsys.Thex, Nday, SimDayend)
	Qmeasday(Mon, Day, ttmm, Eqsys.Qmeas, Nday, SimDayend)
	PVday(Mon, Day, ttmm, Eqsys.PVcmp, Nday, SimDayend)
	Desiday(Mon, Day, ttmm, Eqsys.Desi, Nday, SimDayend)

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

func Compodyprt(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Rdpnl []*RDPNL) {
	if __Compodyprt_id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boidyprt(fo, __Compodyprt_id, Eqsys.Boi)
			refadyprt(fo, __Compodyprt_id, Eqsys.Refa)
			colldyprt(fo, __Compodyprt_id, Eqsys.Coll)
			hccdyprt(fo, __Compodyprt_id, Eqsys.Hcc)
			pipedyprt(fo, __Compodyprt_id, Eqsys.Pipe)
			hexdyprt(fo, __Compodyprt_id, Eqsys.Hex)
			stankdyprt(fo, __Compodyprt_id, Eqsys.Stank)
			pumpdyprt(fo, __Compodyprt_id, Eqsys.Pump)
			hclddyprt(fo, __Compodyprt_id, Eqsys.Hcload)
			stheatdyprt(fo, __Compodyprt_id, Eqsys.Stheat)
			Qmeasdyprt(fo, __Compodyprt_id, Eqsys.Qmeas)
			Thexdyprt(fo, __Compodyprt_id, Eqsys.Thex)
			PVdyprt(fo, __Compodyprt_id, Eqsys.PVcmp)
			Desidyprt(fo, __Compodyprt_id, Eqsys.Desi)

			paneldyprt(fo, __Compodyprt_id, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compodyprt_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boidyprt(fo, __Compodyprt_id, Eqsys.Boi)
	refadyprt(fo, __Compodyprt_id, Eqsys.Refa)
	colldyprt(fo, __Compodyprt_id, Eqsys.Coll)
	hccdyprt(fo, __Compodyprt_id, Eqsys.Hcc)
	pipedyprt(fo, __Compodyprt_id, Eqsys.Pipe)
	hexdyprt(fo, __Compodyprt_id, Eqsys.Hex)
	stankdyprt(fo, __Compodyprt_id, Eqsys.Stank)
	pumpdyprt(fo, __Compodyprt_id, Eqsys.Pump)
	hclddyprt(fo, __Compodyprt_id, Eqsys.Hcload)
	stheatdyprt(fo, __Compodyprt_id, Eqsys.Stheat)
	Qmeasdyprt(fo, __Compodyprt_id, Eqsys.Qmeas)
	Thexdyprt(fo, __Compodyprt_id, Eqsys.Thex)
	PVdyprt(fo, __Compodyprt_id, Eqsys.PVcmp)
	Desidyprt(fo, __Compodyprt_id, Eqsys.Desi)

	paneldyprt(fo, __Compodyprt_id, Rdpnl)

}

/* システム要素機器の月集計結果出力 */

var __Compomonprt_id int

func Compomonprt(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Eqsys *EQSYS, Rdpnl []*RDPNL) {
	if __Compomonprt_id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			boimonprt(fo, __Compomonprt_id, Eqsys.Boi)
			refamonprt(fo, __Compomonprt_id, Eqsys.Refa)
			collmonprt(fo, __Compomonprt_id, Eqsys.Coll)
			hccmonprt(fo, __Compomonprt_id, Eqsys.Hcc)
			pipemonprt(fo, __Compomonprt_id, Eqsys.Pipe)
			hexmonprt(fo, __Compomonprt_id, Eqsys.Hex)
			stankmonprt(fo, __Compomonprt_id, Eqsys.Stank)
			pumpmonprt(fo, __Compomonprt_id, Eqsys.Pump)
			hcldmonprt(fo, __Compomonprt_id, Eqsys.Hcload)
			stheatmonprt(fo, __Compomonprt_id, Eqsys.Stheat)
			Qmeasmonprt(fo, __Compomonprt_id, Eqsys.Qmeas)
			Thexmonprt(fo, __Compomonprt_id, Eqsys.Thex)
			PVmonprt(fo, __Compomonprt_id, Eqsys.PVcmp)

			panelmonprt(fo, __Compomonprt_id, Rdpnl)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compomonprt_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	boimonprt(fo, __Compomonprt_id, Eqsys.Boi)
	refamonprt(fo, __Compomonprt_id, Eqsys.Refa)
	collmonprt(fo, __Compomonprt_id, Eqsys.Coll)
	hccmonprt(fo, __Compomonprt_id, Eqsys.Hcc)
	pipemonprt(fo, __Compomonprt_id, Eqsys.Pipe)
	hexmonprt(fo, __Compomonprt_id, Eqsys.Hex)
	stankmonprt(fo, __Compomonprt_id, Eqsys.Stank)
	pumpmonprt(fo, __Compomonprt_id, Eqsys.Pump)
	hcldmonprt(fo, __Compomonprt_id, Eqsys.Hcload)
	stheatmonprt(fo, __Compomonprt_id, Eqsys.Stheat)
	Qmeasmonprt(fo, __Compomonprt_id, Eqsys.Qmeas)
	Thexmonprt(fo, __Compomonprt_id, Eqsys.Thex)
	PVmonprt(fo, __Compomonprt_id, Eqsys.PVcmp)

	panelmonprt(fo, __Compomonprt_id, Rdpnl)
}

/* システム要素機器の年集計結果出力 */

var __Compomtprt_id int = 0

func Compomtprt(fo io.Writer, mrk string, Simc *SIMCONTL, Eqsys *EQSYS, Rdpnl []*RDPNL) {
	if __Compomtprt_id == 0 {
		ttlmtprint(fo, mrk, Simc)

		for j := 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}

			stheatmtprt(fo, __Compomtprt_id, Eqsys.Stheat, 0, 0)
			boimtprt(fo, __Compomtprt_id, Eqsys.Boi, 0, 0)
			refamtprt(fo, __Compomtprt_id, Eqsys.Refa, 0, 0)
			pumpmtprt(fo, __Compomtprt_id, Eqsys.Pump, 0, 0)
			PVmtprt(fo, __Compomtprt_id, Eqsys.PVcmp, 0, 0)
			hcldmtprt(fo, __Compomtprt_id, 0, 0, Eqsys.Hcload)
			panelmtprt(fo, __Compomtprt_id, Rdpnl, 0, 0)

			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}

			__Compomtprt_id++
		}
	}

	for mo := 1; mo <= 12; mo++ {
		for tt := 1; tt <= 24; tt++ {
			fmt.Fprintf(fo, "%02d %02d\n", mo, tt)
			stheatmtprt(fo, __Compomtprt_id, Eqsys.Stheat, mo, tt)
			boimtprt(fo, __Compomtprt_id, Eqsys.Boi, mo, tt)
			refamtprt(fo, __Compomtprt_id, Eqsys.Refa, mo, tt)
			pumpmtprt(fo, __Compomtprt_id, Eqsys.Pump, mo, tt)
			PVmtprt(fo, __Compomtprt_id, Eqsys.PVcmp, mo, tt)
			hcldmtprt(fo, __Compomtprt_id, mo, tt, Eqsys.Hcload)
			panelmtprt(fo, __Compomtprt_id, Rdpnl, mo, tt)
		}
	}
}
